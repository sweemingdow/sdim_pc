package main

import (
	"embed"
	"flag"
	"sdim_pc/backend/binder/convbinder"
	"sdim_pc/backend/binder/msgbinder"
	"sdim_pc/backend/binder/syncbinder"
	"sdim_pc/backend/binder/userbinder"
	"sdim_pc/backend/config"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/utils"
	"sdim_pc/backend/utils/parser/yaml"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfgPath := flag.String("cfg-path", "./configs/sdim_pc.yaml", "-cfg-path")
	flag.Parse()

	contents, err := utils.ReadAll(*cfgPath)
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	err = yaml.Parse(contents, &cfg)
	if err != nil {
		panic(err)
	}

	err = mylog.InitLogger(cfg.LogCfg, cfg.AppCfg.Profile)
	if err != nil {
		panic(err)
	}

	// Create an instance of the app structure
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:         "Sdim For PC",
		Width:         960,
		Height:        768,
		MinWidth:      960,
		MinHeight:     768,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		Bind: append(
			[]interface{}{
				app,
			},
			syncbinder.NewSyncBinder(app.ci, app.cm),
			userbinder.NewUserBinder(app.ui),
			msgbinder.NewMsgBinder(app.mi, app.cm),
			convbinder.NewConvBinder(app.ci, app.cm),
		),
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
