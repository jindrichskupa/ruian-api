package main

import (
	"github.com/jindrichskupa/ruian-api/app"
	"github.com/jindrichskupa/ruian-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
