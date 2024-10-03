package main

import (
	"context"
	"flag"
	"log"

	"github.com/Timofey335/auth/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	if err = a.Run(ctx); err != nil {
		log.Fatalf("failed ti run app: %s", err.Error())
	}
}
