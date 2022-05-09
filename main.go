package main

import (
	_ "net/http/pprof"

	"github.com/rid-lin/go-tg-bots/for_Vasiliy/internal/app"
	"github.com/rid-lin/go-tg-bots/for_Vasiliy/internal/config"
)

func main() {
	var cfg *config.Config
	cfg = config.New()
	App := app.New(cfg)
	App.Configure()
	App.Start()
}
