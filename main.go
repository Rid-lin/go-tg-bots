package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/rid-lin/go-tg-bots/for_Vasiliy/internal/app"
	"github.com/rid-lin/go-tg-bots/for_Vasiliy/internal/config"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("Recovered in f", r)
			main()
		}
	}()
	var cfg *config.Config
	cfg = config.New()
	myApp := app.New(cfg)
	myApp.Configure()

	myApp.Start()
}
