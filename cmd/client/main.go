package main

import (
	"fmt"
	"gopasskeeper/pkg/tui"
	"gopasskeeper/pkg/tui/api"
)

type DefaultClientConfig struct{}

func (c *DefaultClientConfig) CertPath() string {
	return "certs/cert.pem"
}

func main() {
	api := api.New(&DefaultClientConfig{})
	tuiAPI := &tui.API{
		AuthAPI:     api,
		AccountsAPI: api,
		CardsAPI:    api,
		NotesAPI:    api,
		FilesAPI:    api,
		SyncAPI:     api,
	}
	app := tui.New(tuiAPI)

	if err := app.Run(); err != nil {
		fmt.Printf("got error: %+v", err)
	}
}
