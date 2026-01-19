package main

import (
	"embed"
	"log"
	"os"

	"github.com/The-Robin-Hood/bunnymail/internal/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	loggerCfg := logger.Config{
		Level:      logger.LogLevel("debug"),
		OutputFile: home + "/.bunnymail/logs/bunnymail.log",
		UseJSON:    false,
	}
	if err := logger.Init(loggerCfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	cmdlineArg := os.Args
	if len(cmdlineArg) > 1 && cmdlineArg[1] == "sandbox" {
		sandbox()
		return
	}

	app := BunnyMailApp()

	err = wails.Run(&options.App{
		Title: "Bunny Mail",

		MinWidth:  1200,
		MinHeight: 800,

		Width:  1200,
		Height: 800,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
