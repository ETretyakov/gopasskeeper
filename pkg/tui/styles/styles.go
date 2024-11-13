package styles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	BackgroundColor   tcell.Color = tcell.NewRGBColor(30, 33, 39)
	MainAccentColor   tcell.Color = tcell.NewRGBColor(184, 93, 214)
	SecondAccentColor tcell.Color = tcell.NewRGBColor(135, 183, 101)
	SecondTextColor   tcell.Color = tcell.NewRGBColor(152, 161, 177)
)

func ApplyFormStyleNoBorder(form *tview.Form) {
	form.
		SetTitleAlign(tview.AlignCenter)
	form.
		SetLabelColor(MainAccentColor).
		SetFieldBackgroundColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorBlack).
		SetButtonBackgroundColor(BackgroundColor).
		SetButtonTextColor(tcell.ColorBlack).
		SetButtonActivatedStyle(tcell.StyleDefault.Background(MainAccentColor)).
		SetBackgroundColor(BackgroundColor)
}

func ApplyFormStyle(form *tview.Form) {
	ApplyFormStyleNoBorder(form)
	form.SetBorder(true)
}

func ApplyFrameStyle(flex *tview.Flex) {
	flex.
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)

	flex.
		SetBackgroundColor(BackgroundColor)
}
