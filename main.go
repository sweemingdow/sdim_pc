package main

import (
	"embed"
	"flag"
	"sdim_pc/backend/binder/syncbinder"
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
		Title:     "Sdim For PC",
		Width:     1024,
		Height:    820,
		MinWidth:  650,
		MinHeight: 550,
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
			syncbinder.NewSyncBinder(app.ci),
		),
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
