package assets

import _ "embed"

// StockAssets holds embedded stock base data.
type StockAssets struct {
	StockBasic []byte
	StockHK    []byte
	StockUS    []byte
}

//go:embed data/stock_basic.json
var stockBasic []byte

//go:embed data/stock_base_info_hk.json
var stockHK []byte

//go:embed data/stock_base_info_us.json
var stockUS []byte

// StockData returns embedded stock assets.
func StockData() *StockAssets {
	return &StockAssets{
		StockBasic: stockBasic,
		StockHK:    stockHK,
		StockUS:    stockUS,
	}
}
