package auth

import (
	"fmt"
	"gopasskeeper/pkg/tui/modals"
	"gopasskeeper/pkg/tui/registry"
	"gopasskeeper/pkg/tui/styles"

	"github.com/rivo/tview"
)

type API interface {
	Login(endpoint, login, password string) error
}

type AuthWidget struct {
	api   API
	pages *tview.Pages
	tview *tview.Flex
}

func New(
	api API,
	pages *tview.Pages,
	returnCallback func(),
	version string,
) *AuthWidget {
	const (
		inputFieldWidth int  = 33
		formHeight      int  = 11
		formWidth       int  = 47
		resizable       int  = 0
		oneWeight       int  = 1
		twoWeight       int  = 2
		focused         bool = true
		unfocused       bool = false
	)

	serverLogin := LoginForm{}

	flex := tview.NewFlex()

	form := tview.NewForm()
	form.SetTitle(fmt.Sprintf(" Server (%s) ", version))
	styles.ApplyFormStyle(form)

	form.
		AddInputField("Host:", "", inputFieldWidth, nil, func(server string) {
			serverLogin.server = server
		}).
		AddInputField("User:", "", inputFieldWidth, nil, func(login string) {
			serverLogin.login = login
		}).
		AddPasswordField("Password:", "", inputFieldWidth, '*', func(password string) {
			serverLogin.password = password
		}).
		AddButton("Login", func() {
			if err := api.Login(
				serverLogin.server,
				serverLogin.login,
				serverLogin.password,
			); err != nil {
				callback := func() { pages.SwitchToPage(registry.AuthWidgetPage) }

				pages.RemovePage(registry.ErrorModalPage)
				modal := modals.NewErrorModal(
					"failed to get account id",
					callback,
				)
				pages.AddAndSwitchToPage(
					registry.ErrorModalPage,
					modal.Flex(),
					true,
				)
			} else {
				returnCallback()
			}
		}).
		SetButtonsAlign(tview.AlignCenter)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	innerFlex.
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused).
		AddItem(form, formHeight, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused)

	flex.
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused).
		AddItem(innerFlex, formWidth, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused)

	return &AuthWidget{
		api:   api,
		pages: pages,
		tview: flex,
	}
}

func (a *AuthWidget) Flex() *tview.Flex {
	return a.tview
}
