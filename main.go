package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "DadBot",
		Width:  1280,
		Height: 820,
		MinWidth: 980,
		MinHeight: 640,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup: app.startup,
		Bind: []interface{}{app},
	})
	if err != nil {
		log.Fatal(err)
	}
}
