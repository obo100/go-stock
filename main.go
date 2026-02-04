package main

import (
	"embed"
	"fmt"
	"go-stock/backend/bootstrap"
	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/internal/assets"
	"os"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed build/app.ico
var icon2 []byte

//go:embed build/screenshot/alipay.jpg
var alipay []byte

//go:embed build/screenshot/wxpay.jpg
var wxpay []byte

//go:embed build/screenshot/扫码_搜索联合传播样式-白色版.png
var wxgzh []byte

//go:generate cp -R ./data ./build/bin

var Version string
var VersionCommit string
var OFFICIAL_STATEMENT string
var BuildKey string

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.SugaredLogger.Error("panic: ", r)
			log.SugaredLogger.Error("stack: ", string(debug.Stack()))
		}
	}()

	checkDir("data")
	_ = bootstrap.Init("", assets.StockData(), bootstrap.Options{
		AutoMigrateAsync: true,
	})

	//db.Dao.Model(&data.Group{}).Where("id = ?", 0).FirstOrCreate(&data.Group{
	//	Name: "默认分组",
	//	Sort: 0,
	//})

	log.SugaredLogger.Info("starting...")
	log.SugaredLogger.Infof("version: %s  commit: %s", Version, VersionCommit)
	//log.SugaredLogger.Infof("build key: %s", BuildKey)

	// Create an instance of the app structure
	app := NewApp()
	AppMenu := menu.NewMenu()
	if IsMacOS() {
		AppMenu.Append(menu.EditMenu())
	}
	//FileMenu := AppMenu.AddSubmenu("设置")
	//FileMenu.AddText("窗口全屏", keys.CmdOrCtrl("f"), func(callback *menu.CallbackData) {
	//	runtime.WindowFullscreen(app.ctx)
	//})
	//FileMenu.AddText("窗口还原", keys.Key("Esc"), func(callback *menu.CallbackData) {
	//	runtime.WindowUnfullscreen(app.ctx)
	//})
	//FileMenu.AddText("显示搜索框", keys.CmdOrCtrl("s"), func(callbackData *menu.CallbackData) {
	//	runtime.EventsEmit(app.ctx, "showSearch", 1)
	//})
	//FileMenu.AddText("隐藏搜索框", keys.CmdOrCtrl("d"), func(callbackData *menu.CallbackData) {
	//	runtime.EventsEmit(app.ctx, "showSearch", 0)
	//})
	//FileMenu.AddText("刷新数据", keys.CmdOrCtrl("r"), func(callbackData *menu.CallbackData) {
	//	//runtime.EventsEmit(app.ctx, "refresh", "setting-"+time.Now().Format("2006-01-02 15:04:05"))
	//	runtime.EventsEmit(app.ctx, "refreshFollowList", "refresh-"+time.Now().Format("2006-01-02 15:04:05"))
	//})
	//FileMenu.AddSeparator()

	//if goruntime.GOOS == "windows" {
	//	FileMenu.AddText("隐藏到托盘区", keys.CmdOrCtrl("z"), func(_ *menu.CallbackData) {
	//		runtime.WindowHide(app.ctx)
	//	})
	//}

	//FileMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
	//	runtime.Quit(app.ctx)
	//})
	log.SugaredLogger.Info("version: " + Version)
	log.SugaredLogger.Info("commit: " + VersionCommit)
	// Create application with options
	//var width, height int
	//var err error
	//
	width, _, minWidth, minHeight, err := getScreenResolution()
	if err != nil {
		log.SugaredLogger.Error("get screen resolution error")
		width = 1456
		//height = 768
	}

	darkTheme := data.GetSettingConfig().DarkTheme
	backgroundColour := &options.RGBA{R: 255, G: 255, B: 255, A: 1}
	if darkTheme {
		backgroundColour = &options.RGBA{R: 27, G: 38, B: 54, A: 1}
	}

	//frameless := getFrameless()

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "go-stock：AI赋能股票分析✨ " + OFFICIAL_STATEMENT,
		Width:     width * 4 / 5,
		Height:    920,
		MinWidth:  minWidth,
		MinHeight: minHeight,
		//MaxWidth:                 width,
		//MaxHeight:                height,
		DisableResize:            false,
		Fullscreen:               false,
		Frameless:                false,
		StartHidden:              false,
		HideWindowOnClose:        false,
		EnableDefaultContextMenu: true,
		BackgroundColour:         backgroundColour,
		Assets:                   assets,
		Menu:                     AppMenu,
		Logger:                   logger.NewFileLogger("./logs/wails.log"),
		LogLevel:                 logger.DEBUG,
		LogLevelProduction:       logger.INFO,
		OnStartup:                app.startup,
		OnDomReady:               app.domReady,
		OnBeforeClose:            app.beforeClose,
		OnShutdown:               app.shutdown,
		WindowStartState:         options.Normal,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "go-stock",
			OnSecondInstanceLaunch: OnSecondInstanceLaunch,
		},
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "go-stock",
				Message: "go-stock：AI赋能股票分析✨ ",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.SugaredLogger.Fatal(err)
	}

}

func checkDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
		log.SugaredLogger.Info("create dir: " + dir)
	}
	if BuildKey == "" {
		BuildKey = "cc1e0d684e32f176c56ff1fcf384dcd9"
	}
}

// PanicHandler 捕获 panic 的包装函数
func PanicHandler() {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from panic: %v\n", r)
		debug.PrintStack()
	}
}
