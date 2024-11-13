package modals

import (
	"gopasskeeper/pkg/tui/styles"

	"github.com/rivo/tview"
)

type SecretModal struct {
	tview *tview.Flex
}

func NewSecretModal(text string, callback func()) *SecretModal {
	const (
		inputFieldWidth int  = 34
		formWidth       int  = 11
		formHeight      int  = 47
		fieldWidth      int  = 47
		fieldHight      int  = 5
		resizable       int  = 0
		oneWeight       int  = 1
		twoWeight       int  = 2
		dynamicColors   bool = false
		scrollable      bool = false
		focused         bool = true
		unfocused       bool = false
	)

	frame := tview.NewFlex()

	form := tview.NewForm()
	form.SetTitle(" Secret ")
	styles.ApplyFormStyle(form)

	form.AddTextView(
		"Secret: ",
		text,
		fieldWidth,
		fieldHight,
		dynamicColors,
		scrollable,
	)
	form.AddButton("Back", func() {
		callback()
	}).
		SetButtonsAlign(tview.AlignCenter)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	innerFlex.
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused).
		AddItem(form, formWidth, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused)

	frame.
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused).
		AddItem(innerFlex, formHeight, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused)

	return &SecretModal{tview: frame}
}

func (s *SecretModal) Flex() *tview.Flex {
	return s.tview
}

type ErrorModal struct {
	tview *tview.Flex
}

func NewErrorModal(text string, callback func()) *ErrorModal {
	const (
		inputFieldWidth int  = 34
		formWidth       int  = 11
		formHeight      int  = 47
		fieldWidth      int  = 47
		fieldHight      int  = 5
		resizable       int  = 0
		oneWeight       int  = 1
		twoWeight       int  = 2
		dynamicColors   bool = false
		scrollable      bool = false
		focused         bool = true
		unfocused       bool = false
	)

	frame := tview.NewFlex()

	form := tview.NewForm()
	form.SetTitle(" Error ")
	styles.ApplyFormStyle(form)

	form.AddTextView(
		"Message: ",
		text,
		fieldWidth,
		fieldHight,
		dynamicColors,
		scrollable,
	)
	form.AddButton("Ok", func() {
		callback()
	}).
		SetButtonsAlign(tview.AlignCenter)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	innerFlex.
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused).
		AddItem(form, formWidth, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused)

	frame.
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused).
		AddItem(innerFlex, formHeight, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused)

	return &ErrorModal{tview: frame}
}

func (s *ErrorModal) Flex() *tview.Flex {
	return s.tview
}
