package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/internal/assets"

	"github.com/duke-git/lancet/v2/slice"
)

// Options controls initialization behavior.
type Options struct {
	UpdateBasicInfo  bool
	NoStockData      bool
	AutoMigrateAsync bool
}

// Init initializes database, sentiment, migrations, and optional base data.
func Init(dbPath string, stockAssets *assets.StockAssets, opts Options) error {
	ensureDataDir(dbPath)
	db.Init(dbPath)
	data.InitAnalyzeSentiment()

	if opts.AutoMigrateAsync {
		go AutoMigrate()
	} else {
		AutoMigrate()
	}

	if !opts.NoStockData && stockAssets != nil {
		if err := ImportStockData(stockAssets); err != nil {
			logger.SugaredLogger.Errorf("import stock data error: %v", err)
		}
	}

	if opts.UpdateBasicInfo {
		UpdateBasicInfo()
	}
	return nil
}

func ensureDataDir(dbPath string) {
	cleanPath := cleanDbPath(dbPath)
	if cleanPath == "" || cleanPath == ":memory:" {
		return
	}
	dir := filepath.Dir(cleanPath)
	if dir == "." || dir == "" {
		return
	}
	_ = os.MkdirAll(dir, os.ModePerm)
}

func cleanDbPath(dbPath string) string {
	if dbPath == "" {
		return "data/stock.db"
	}
	path := strings.SplitN(dbPath, "?", 2)[0]
	return path
}

// AutoMigrate migrates all DB tables.
func AutoMigrate() {
	db.Dao.AutoMigrate(&data.StockInfo{})
	db.Dao.AutoMigrate(&data.StockBasic{})
	db.Dao.AutoMigrate(&data.FollowedStock{})
	db.Dao.AutoMigrate(&data.IndexBasic{})
	db.Dao.AutoMigrate(&data.Settings{})
	db.Dao.AutoMigrate(&models.AIResponseResult{})
	db.Dao.AutoMigrate(&models.StockInfoHK{})
	db.Dao.AutoMigrate(&models.StockInfoUS{})
	db.Dao.AutoMigrate(&data.FollowedFund{})
	db.Dao.AutoMigrate(&data.FundBasic{})
	db.Dao.AutoMigrate(&models.PromptTemplate{})
	db.Dao.AutoMigrate(&data.Group{})
	db.Dao.AutoMigrate(&data.GroupStock{})
	db.Dao.AutoMigrate(&models.Tags{})
	db.Dao.AutoMigrate(&models.Telegraph{})
	db.Dao.AutoMigrate(&models.TelegraphTags{})
	db.Dao.AutoMigrate(&models.LongTigerRankData{})
	db.Dao.AutoMigrate(&data.AIConfig{})
	db.Dao.AutoMigrate(&models.BKDict{})
	db.Dao.AutoMigrate(&models.WordAnalyze{})
	db.Dao.AutoMigrate(&models.SentimentResultAnalyze{})
	db.Dao.AutoMigrate(&models.AiRecommendStocks{})
}

// UpdateBasicInfo updates stock base info from Tushare.
func UpdateBasicInfo() {
	config := data.GetSettingConfig()
	if config.UpdateBasicInfoOnStart {
		go data.NewStockDataApi().GetStockBaseInfo()
		go data.NewStockDataApi().GetIndexBasic()
	}
}

// ImportStockData imports embedded stock base data into DB.
func ImportStockData(stockAssets *assets.StockAssets) error {
	if stockAssets == nil {
		return errors.New("stock assets is nil")
	}
	if len(stockAssets.StockBasic) > 0 {
		if err := importStockData(stockAssets.StockBasic); err != nil {
			return err
		}
	}
	if len(stockAssets.StockHK) > 0 {
		if err := importStockDataHK(stockAssets.StockHK); err != nil {
			return err
		}
	}
	if len(stockAssets.StockUS) > 0 {
		if err := importStockDataUS(stockAssets.StockUS); err != nil {
			return err
		}
	}
	return nil
}

func importStockDataUS(bin []byte) error {
	var v []models.StockInfoUS
	if err := json.Unmarshal(bin, &v); err != nil {
		return err
	}
	logger.SugaredLogger.Infof("init stock data us %d", len(v))
	var total int64
	db.Dao.Model(&models.StockInfoUS{}).Count(&total)
	if total != int64(len(v)) {
		for _, item := range v {
			var count int64
			db.Dao.Model(&models.StockInfoUS{}).Where("code = ?", item.Code).Count(&count)
			if count > 0 {
				continue
			}
			db.Dao.Model(&models.StockInfoUS{}).Create(&item)
		}
	}
	return nil
}

func importStockDataHK(bin []byte) error {
	var v []models.StockInfoHK
	if err := json.Unmarshal(bin, &v); err != nil {
		return err
	}
	logger.SugaredLogger.Infof("init stock data hk %d", len(v))
	var total int64
	db.Dao.Model(&models.StockInfoHK{}).Count(&total)
	if total != int64(len(v)) {
		for _, item := range v {
			var count int64
			db.Dao.Model(&models.StockInfoHK{}).Where("code = ?", item.Code).Count(&count)
			if count > 0 {
				continue
			}
			db.Dao.Model(&models.StockInfoHK{}).Create(&item)
		}
	}
	return nil
}

func importStockData(bin []byte) error {
	fields := "ts_code,symbol,name,area,industry,cnspell,market,list_date,act_name,act_ent_type,fullname,exchange,list_status,curr_type,enname,delist_date,is_hs"
	res := &data.TushareStockBasicResponse{}
	if err := json.Unmarshal(bin, res); err != nil {
		return err
	}
	for _, item := range res.Data.Items {
		stock := &data.StockBasic{}
		stockData := map[string]any{}
		for _, field := range strings.Split(fields, ",") {
			idx := slice.IndexOf(res.Data.Fields, field)
			if idx == -1 {
				continue
			}
			stockData[field] = item[idx]
		}
		jsonData, _ := json.Marshal(stockData)
		if err := json.Unmarshal(jsonData, stock); err != nil {
			continue
		}
		stock.ID = 0
		var count int64
		db.Dao.Model(&data.StockBasic{}).Where("ts_code = ?", stock.TsCode).Count(&count)
		if count > 0 {
			continue
		}
		db.Dao.Create(stock)
	}
	return nil
}

// InitStockDataWithContext is kept for compatibility when callers need a context.
func InitStockDataWithContext(ctx context.Context, stockAssets *assets.StockAssets) error {
	_ = ctx
	return ImportStockData(stockAssets)
}
