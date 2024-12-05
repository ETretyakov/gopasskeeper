package tui

import (
	"context"
	"gopasskeeper/pkg/tui/config"
	"gopasskeeper/pkg/tui/registry"
	"gopasskeeper/pkg/tui/sections/accounts"
	"gopasskeeper/pkg/tui/sections/auth"
	"gopasskeeper/pkg/tui/sections/cards"
	"gopasskeeper/pkg/tui/sections/files"
	"gopasskeeper/pkg/tui/sections/notes"
	"time"

	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type SyncAPI interface {
	Outdated() (bool, error)
}

type API struct {
	AuthAPI     auth.API
	AccountsAPI accounts.API
	CardsAPI    cards.API
	NotesAPI    notes.API
	FilesAPI    files.API
	SyncAPI     SyncAPI
}

type App struct {
	app   *tview.Application
	pages *tview.Pages
	api   *API
}

var appliaction *tview.Application

func New(version string, api *API) *App {
	app := tview.NewApplication()
	appliaction = app
	pages := tview.NewPages()

	mainFrame := tview.NewFlex()
	pages.AddPage(registry.MainWidgetPage, mainFrame, true, false)

	// accounts
	accountsPage := accounts.New(api.AccountsAPI, pages)
	pages.AddPage(registry.AccountWidgetPage, accountsPage.Flex(), true, false)

	// cards
	cardsPage := cards.New(api.CardsAPI, pages)
	pages.AddPage(registry.CardsWidgetPage, cardsPage.Flex(), true, false)

	// notes
	notesPage := notes.New(api.NotesAPI, pages)
	pages.AddPage(registry.NotesWidgetPage, notesPage.Flex(), true, false)

	// files
	filesPage := files.New(api.FilesAPI, pages)
	pages.AddPage(registry.FilesWidgetPage, filesPage.Flex(), true, false)

	// auth
	authPage := auth.New(
		api.AuthAPI,
		pages,
		func() { initWidgetsContent(accountsPage, cardsPage, notesPage, filesPage, api.SyncAPI) },
		version,
	)
	pages.AddPage(registry.AuthWidgetPage, authPage.Flex(), true, true)

	return &App{
		app:   app,
		pages: pages,
		api:   api,
	}
}

func initWidgetsContent(
	a *accounts.AccountsWidget,
	c *cards.CardsWidget,
	n *notes.NotesWidget,
	f *files.FilesWidget,
	syncAPI SyncAPI,
) {
	f.Init()
	n.Init()
	c.Init()
	a.Init()

	go func() {
		ctx := context.Background()

		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				if outdated, _ := syncAPI.Outdated(); outdated {
					f.Update()
					n.Update()
					c.Update()
					a.Update()
					if appliaction != nil {
						appliaction.Draw()
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *App) Run() error {
	cfg := config.NewAppConfig()

	if err := a.app.
		SetRoot(a.pages, cfg.Fullscreen).
		EnableMouse(cfg.EnableMouse).
		Run(); err != nil {
		return errors.Wrap(err, "failed to start TUI")
	}

	return nil
}
