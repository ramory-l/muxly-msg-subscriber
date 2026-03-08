package main

import (
	sharedapp "github.com/Muxly-Corp/muxly-shared/app"
	"github.com/Muxly-Corp/muxly-msg-subscriber/config"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/app"
)

func main() {
	sharedapp.Run(
		sharedapp.NewConfig[config.Config],
		app.InitializeApp,
		sharedapp.WithQueue(),
	)
}
