package cards

import (
	"fmt"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"
	"gopasskeeper/pkg/tui/modals"
	"gopasskeeper/pkg/tui/registry"
	"gopasskeeper/pkg/tui/styles"
	"strconv"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type API interface {
	SearchCard(substring string, offset uint64, limit uint32) (*cardsv1.CardSearchResponse, error)
	GetCard(uuid string) (string, error)
	AddCard(name, number string, month, year int32, ccv, pin string) error
	RemoveCard(secredID string) error
}

type CardsWidget struct {
	mu          *sync.RWMutex
	api         API
	table       *CardsTable
	pages       *tview.Pages
	tview       *tview.Flex
	searchInput *SearchInput
}

func New(api API, pages *tview.Pages) *CardsWidget {
	cardWidget := &CardsWidget{
		mu:    &sync.RWMutex{},
		api:   api,
		pages: pages,
		searchInput: &SearchInput{
			value:  "",
			offset: 0,
			limit:  20,
			step:   20,
		},
	}

	cardWidget.draw()

	return cardWidget
}

func (a *CardsWidget) Flex() *tview.Flex {
	return a.tview
}

func (a *CardsWidget) draw() {
	const (
		leftPaneWidth   int  = 20
		leftPaneWeight  int  = 1
		rightPaneWidth  int  = 0
		rightPaneWeight int  = 10
		unfocused       bool = false
		focused         bool = true
	)

	mainFrame := tview.NewFlex()

	leftPaneFrame := a.drawLeftPaneFrame()
	rightPaneFrame := a.drawRightPaneFrame()

	mainFrame.
		AddItem(leftPaneFrame, leftPaneWidth, leftPaneWeight, unfocused).
		AddItem(rightPaneFrame, rightPaneWidth, rightPaneWeight, focused)

	a.tview = mainFrame
}

func (a *CardsWidget) drawRightPaneFrame() *tview.Flex {
	const (
		unfocused bool = false
		focused   bool = true
	)

	rightPaneFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPaneFrame.SetTitle(" Cards ")
	styles.ApplyFrameStyle(rightPaneFrame)

	topMenu := a.drawTopMenu()
	searchForm := a.drawSearchViewForm()
	table := NewTable()

	table.table.SetSelectionChangedFunc(
		func(row int, column int) {
			if row == a.table.table.GetRowCount()-1 {
				a.Paginate()
			}
		},
	)

	table.table.SetSelectedFunc(func(row int, column int) {
		cell := a.table.table.GetCell(row, 0)
		ref := cell.GetReference()
		if val, ok := ref.(string); ok {
			a.ShowPass(val)
		} else {
			callback := func() { a.pages.SwitchToPage(registry.CardsWidgetPage) }

			a.pages.RemovePage(registry.ErrorModalPage)
			modal := modals.NewErrorModal(
				"empty selection",
				callback,
			)
			a.pages.AddAndSwitchToPage(
				registry.ErrorModalPage,
				modal.Flex(),
				true,
			)
		}
	})

	a.table = table

	rightPaneFrame.AddItem(topMenu, 3, 1, unfocused)
	rightPaneFrame.AddItem(searchForm, 3, 1, unfocused)
	rightPaneFrame.AddItem(table.table, 0, 12, focused)

	return rightPaneFrame
}

func (a *CardsWidget) drawSearchViewForm() *tview.Form {
	searchViewFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	searchViewFrame.SetTitle(registry.CardsMainFrameTitle)
	styles.ApplyFrameStyle(searchViewFrame)

	return a.drawSearchInput()
}

func (a *CardsWidget) drawTopMenu() *tview.Form {
	menuForm := tview.NewForm()
	styles.ApplyFormStyleNoBorder(menuForm)

	menuForm.AddButton("Add", func() {
		modal := NewAddCard(
			a.api.AddCard,
			a.Refresh,
			func() { a.pages.SwitchToPage(registry.CardsWidgetPage) },
		)

		a.pages.AddAndSwitchToPage(
			registry.AddCardsWidgetPage,
			modal,
			true,
		)
	})

	menuForm.AddButton("Remove", func() {
		row, column := a.table.table.GetSelection()
		cell := a.table.table.GetCell(row, column)
		ref := cell.GetReference()
		if val, ok := ref.(string); ok {
			a.Remove(val)
		} else {
			callback := func() { a.pages.SwitchToPage(registry.CardsWidgetPage) }

			a.pages.RemovePage(registry.ErrorModalPage)
			modal := modals.NewErrorModal(
				"failed to get card id",
				callback,
			)
			a.pages.AddAndSwitchToPage(
				registry.ErrorModalPage,
				modal.Flex(),
				true,
			)
		}
	})

	menuForm.SetButtonsAlign(tview.AlignRight)

	return menuForm
}

func (a *CardsWidget) drawSearchInput() *tview.Form {
	const (
		searchInputWidth int = 0
	)

	searchForm := tview.NewForm()
	styles.ApplyFormStyleNoBorder(searchForm)

	searchForm.AddInputField(
		"Search: ", "", searchInputWidth, nil,
		func(searchInputValue string) {
			a.mu.Lock()
			a.searchInput.value = searchInputValue
			a.mu.Unlock()

			a.Search()
		},
	)

	searchForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			a.Search()
			a.pages.SwitchToPage(registry.CardsWidgetPage)
		case tcell.KeyDown:
			a.pages.SwitchToPage(registry.CardsWidgetPage)
		}

		return event
	})

	return searchForm
}

func (a *CardsWidget) drawLeftPaneFrame() *tview.Flex {
	const (
		paneHeight int  = 0
		paneWeight int  = 1
		unfocused  bool = false
	)

	leftPaneMenu := tview.NewForm().SetHorizontal(true)
	leftPaneMenu.SetTitle(" Sections ")
	styles.ApplyFormStyle(leftPaneMenu)

	leftPaneMenu.
		AddButton(
			registry.LeftPaneMenuAccount,
			func() { a.pages.SwitchToPage(registry.AccountWidgetPage) },
		).
		AddButton(
			registry.LeftPaneMenuCardsActive,
			func() { a.pages.SwitchToPage(registry.CardsWidgetPage) },
		).
		AddButton(
			registry.LeftPaneMenuNotes,
			func() { a.pages.SwitchToPage(registry.NotesWidgetPage) },
		).
		AddButton(
			registry.LeftPaneMenuFiles,
			func() { a.pages.SwitchToPage(registry.FilesWidgetPage) },
		)

	leftPaneMenu.SetButtonsAlign(tview.AlignCenter)

	leftPaneFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPaneFrame.AddItem(leftPaneMenu, paneHeight, paneWeight, unfocused)

	return leftPaneFrame
}

func (a *CardsWidget) Update() {
	a.table.Clean()

	resp, _ := a.api.SearchCard(
		a.searchInput.value,
		a.searchInput.offset,
		a.searchInput.limit,
	)

	a.table.Fill(resp)
	a.table.table.ScrollToBeginning()
}

func (a *CardsWidget) Refresh() {
	a.Search()
	a.pages.SwitchToPage(registry.CardsWidgetPage)
}

func (a *CardsWidget) Init() {
	a.Refresh()
}

func (a *CardsWidget) ShowPass(secretID string) {
	callback := func() { a.pages.SwitchToPage(registry.CardsWidgetPage) }

	secret, err := a.api.GetCard(secretID)
	if err != nil {
		a.pages.RemovePage(registry.ErrorModalPage)
		modal := modals.NewErrorModal(
			fmt.Sprintf("failed to get card: %s", err),
			callback,
		)
		a.pages.AddAndSwitchToPage(
			registry.ErrorModalPage,
			modal.Flex(),
			true,
		)

		return
	}

	a.pages.RemovePage(registry.SecretModalPage)
	modal := modals.NewSecretModal(secret, callback)
	a.pages.AddAndSwitchToPage(
		registry.SecretModalPage,
		modal.Flex(),
		true,
	)
}

func (a *CardsWidget) Remove(secretID string) {
	callback := func() { a.pages.SwitchToPage(registry.CardsWidgetPage) }

	if err := a.api.RemoveCard(secretID); err != nil {
		a.pages.RemovePage(registry.ErrorModalPage)
		modal := modals.NewErrorModal(
			fmt.Sprintf("failed to get card: %s", err),
			callback,
		)
		a.pages.AddAndSwitchToPage(
			registry.ErrorModalPage,
			modal.Flex(),
			true,
		)

		return
	}

	a.Refresh()
}

func (a *CardsWidget) resetPagination() {
	a.mu.Lock()
	a.searchInput.offset = 0
	a.searchInput.limit = uint32(a.searchInput.step)
	a.mu.Unlock()
}

func (a *CardsWidget) Search() {
	a.resetPagination()
	a.table.Clean()

	resp, err := a.api.SearchCard(
		a.searchInput.value,
		a.searchInput.offset,
		a.searchInput.limit,
	)
	if err != nil {
		a.pages.SwitchToPage(registry.ErrorModalPage)
		return
	}

	a.table.Fill(resp)
}

func (a *CardsWidget) Paginate() {
	a.mu.Lock()
	a.searchInput.limit += uint32(a.searchInput.step)
	a.searchInput.offset += uint64(a.searchInput.step)
	a.mu.Unlock()

	resp, err := a.api.SearchCard(
		a.searchInput.value,
		a.searchInput.offset,
		a.searchInput.limit,
	)
	if err != nil {
		a.pages.SwitchToPage(registry.ErrorModalPage)
		return
	}

	a.table.Fill(resp)
}

type CardsTable struct {
	table *tview.Table
}

func NewTable() *CardsTable {
	const (
		rows    = true
		columns = false
		row     = 1
		column  = 0
	)

	table := tview.NewTable()
	table.SetBackgroundColor(styles.BackgroundColor)
	table.SetSelectable(rows, columns)
	table.SetOffset(row, column)

	cardsTable := &CardsTable{table: table}
	cardsTable.FillHeader()

	return cardsTable
}

func (t *CardsTable) FillHeader() {
	const (
		numberColWidth int  = 1
		serverColWidth int  = 20
		loginColWidth  int  = 6
		notSelectable  bool = false
	)

	t.table.SetFixed(1, 0)
	t.table.SetCell(0, 0, tview.NewTableCell(registry.CardsColumnNumberTitle).
		SetSelectable(notSelectable).
		SetExpansion(numberColWidth).
		SetTextColor(styles.SecondAccentColor).
		SetAlign(tview.AlignCenter))

	t.table.SetCell(0, 1, tview.NewTableCell(registry.CardsColumnNumTitle).
		SetSelectable(notSelectable).
		SetExpansion(serverColWidth).
		SetTextColor(styles.SecondAccentColor).
		SetAlign(tview.AlignCenter))

	t.table.SetCell(0, 2, tview.NewTableCell(registry.CardsColumnNameTitle).
		SetSelectable(notSelectable).
		SetExpansion(loginColWidth).
		SetTextColor(styles.SecondAccentColor).
		SetAlign(tview.AlignCenter))
}

func (t *CardsTable) Fill(searchResponse *cardsv1.CardSearchResponse) {
	const (
		tableColumns   int  = 3
		selectable     bool = true
		numberColWidth int  = 1
		serverColWidth int  = 20
		loginColWidth  int  = 6
	)

	currentRows := t.table.GetRowCount()

	rows := searchResponse.GetItems()

	for r := currentRows; r < len(rows)+currentRows; r++ {
		for c := 0; c < tableColumns; c++ {
			cell := tview.NewTableCell("")

			number := strconv.Itoa(currentRows + r - currentRows)

			switch c {
			case 0:
				cell.SetText(number).
					SetReference(rows[r-currentRows].GetId()).
					SetTextColor(styles.SecondAccentColor).
					SetAlign(tview.AlignCenter).
					SetSelectable(selectable).
					SetExpansion(numberColWidth)
			case 1:
				cell.SetText(rows[r-currentRows].GetName()).
					SetReference(rows[r-currentRows].Id).
					SetTextColor(styles.SecondTextColor).
					SetAlign(tview.AlignLeft).
					SetSelectable(selectable).
					SetExpansion(serverColWidth)
			case 2:
				cell.SetText(rows[r-currentRows].GetMask()).
					SetReference(rows[r-currentRows].Id).
					SetTextColor(styles.SecondTextColor).
					SetAlign(tview.AlignCenter).
					SetSelectable(selectable).
					SetExpansion(loginColWidth)
			}

			t.table.SetCell(r, c, cell)
		}
	}
}

func (t *CardsTable) Clean() {
	t.table.Clear()
	t.FillHeader()
}

func NewAddCard(
	callbackAdd func(name, number string, month, year int32, ccv, pin string) error,
	callbackRefresh func(),
	callbackReturn func(),
) *tview.Flex {
	const (
		inputFieldWidth int  = 34
		formWidth       int  = 17
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

	cardAdd := CardAdd{}

	frame := tview.NewFlex()

	infoLine := tview.NewTextView()

	form := tview.NewForm()
	form.SetTitle(" Add card ")
	styles.ApplyFormStyle(form)

	form.
		AddInputField("Name:", "", inputFieldWidth, nil, func(name string) {
			cardAdd.name = name
		}).
		AddInputField("Number:", "", inputFieldWidth, nil, func(number string) {
			cardAdd.number = number
		}).
		AddInputField("Month:", "", inputFieldWidth, nil, func(month string) {
			i, err := strconv.ParseInt(month, 10, 32)
			if err != nil {
				i = 1
			}
			cardAdd.month = int32(i)
		}).
		AddInputField("Year:", "", inputFieldWidth, nil, func(year string) {
			i, err := strconv.ParseInt(year, 10, 32)
			if err != nil {
				i = 1
			}
			cardAdd.year = int32(i)
		}).
		AddInputField("CVC:", "", inputFieldWidth, nil, func(cvc string) {
			cardAdd.cvc = cvc
		}).
		AddInputField("PIN:", "", inputFieldWidth, nil, func(pin string) {
			cardAdd.pin = pin
		}).
		AddButton("Add", func() {
			if err := callbackAdd(
				cardAdd.name,
				cardAdd.number,
				cardAdd.month,
				cardAdd.year,
				cardAdd.cvc,
				cardAdd.pin,
			); err != nil {
				infoLine.SetText(fmt.Sprintf("%s", err))
			} else {
				callbackRefresh()
			}
		}).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Back", func() {
			callbackReturn()
		}).
		SetButtonsAlign(tview.AlignCenter)

	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	innerFlex.
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused).
		AddItem(form, formWidth, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, oneWeight, unfocused).
		AddItem(infoLine, resizable, oneWeight, unfocused)

	frame.
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused).
		AddItem(innerFlex, formHeight, twoWeight, focused).
		AddItem(tview.NewBox(), resizable, twoWeight, unfocused)

	return frame
}
