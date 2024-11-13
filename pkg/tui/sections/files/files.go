package files

import (
	"fmt"
	filesv1 "gopasskeeper/internal/grpc/secretstore/files/gen/files"
	"gopasskeeper/pkg/tui/modals"
	"gopasskeeper/pkg/tui/registry"
	"gopasskeeper/pkg/tui/styles"
	"strconv"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type API interface {
	SearchFile(substring string, offset uint64, limit uint32) (*filesv1.FileSearchResponse, error)
	GetFile(uuid, filePath string) (string, error)
	AddFile(name, content string) error
	RemoveFile(secredID string) error
}

type FilesWidget struct {
	mu          *sync.RWMutex
	api         API
	table       *FilesTable
	pages       *tview.Pages
	tview       *tview.Flex
	searchInput *SearchInput
}

func New(api API, pages *tview.Pages) *FilesWidget {
	fileWidget := &FilesWidget{
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

	fileWidget.draw()

	return fileWidget
}

func (a *FilesWidget) Flex() *tview.Flex {
	return a.tview
}

func (a *FilesWidget) draw() {
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

func (a *FilesWidget) drawRightPaneFrame() *tview.Flex {
	const (
		unfocused bool = false
		focused   bool = true
	)

	rightPaneFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPaneFrame.SetTitle(" Files ")
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
			callback := func() { a.pages.SwitchToPage(registry.FilesWidgetPage) }

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

func (a *FilesWidget) drawSearchViewForm() *tview.Form {
	searchViewFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	searchViewFrame.SetTitle(registry.FilesMainFrameTitle)
	styles.ApplyFrameStyle(searchViewFrame)

	return a.drawSearchInput()
}

func (a *FilesWidget) drawTopMenu() *tview.Form {
	menuForm := tview.NewForm()
	styles.ApplyFormStyleNoBorder(menuForm)

	menuForm.AddButton("Add", func() {
		modal := NewAddFile(
			a.api.AddFile,
			a.Refresh,
			func() { a.pages.SwitchToPage(registry.FilesWidgetPage) },
		)

		a.pages.AddAndSwitchToPage(
			registry.AddFilesWidgetPage,
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
			callback := func() { a.pages.SwitchToPage(registry.FilesWidgetPage) }

			a.pages.RemovePage(registry.ErrorModalPage)
			modal := modals.NewErrorModal(
				"failed to get file id",
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

func (a *FilesWidget) drawSearchInput() *tview.Form {
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
			a.pages.SwitchToPage(registry.FilesWidgetPage)
		case tcell.KeyDown:
			a.pages.SwitchToPage(registry.FilesWidgetPage)
		}

		return event
	})

	return searchForm
}

func (a *FilesWidget) drawLeftPaneFrame() *tview.Flex {
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
			registry.LeftPaneMenuCards,
			func() { a.pages.SwitchToPage(registry.CardsWidgetPage) },
		).
		AddButton(
			registry.LeftPaneMenuNotes,
			func() { a.pages.SwitchToPage(registry.NotesWidgetPage) },
		).
		AddButton(
			registry.LeftPaneMenuFilesActive,
			func() { a.pages.SwitchToPage(registry.FilesWidgetPage) },
		)

	leftPaneMenu.SetButtonsAlign(tview.AlignCenter)

	leftPaneFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPaneFrame.AddItem(leftPaneMenu, paneHeight, paneWeight, unfocused)

	return leftPaneFrame
}

func (a *FilesWidget) Update() {
	a.table.Clean()

	resp, _ := a.api.SearchFile(
		a.searchInput.value,
		a.searchInput.offset,
		a.searchInput.limit,
	)

	a.table.Fill(resp)
	a.table.table.ScrollToBeginning()
}

func (a *FilesWidget) Refresh() {
	a.Search()
	a.pages.SwitchToPage(registry.FilesWidgetPage)
}

func (a *FilesWidget) Init() {
	a.Refresh()
}

func (a *FilesWidget) ShowPass(secretID string) {
	callbackReturn := func() { a.pages.SwitchToPage(registry.FilesWidgetPage) }

	callbackSelected := func(secretID, filePath string, callbackErr func()) {
		_, err := a.api.GetFile(secretID, filePath)
		if err != nil {
			a.pages.RemovePage(registry.ErrorModalPage)
			modal := modals.NewErrorModal(
				fmt.Sprintf("failed to get file: %s", err),
				callbackErr,
			)
			a.pages.AddAndSwitchToPage(
				registry.ErrorModalPage,
				modal.Flex(),
				true,
			)

			return
		}

		a.pages.SwitchToPage(registry.FilesWidgetPage)
	}

	callbackErr := func() { a.pages.SwitchToPage(registry.FilesWidgetPage) }

	modal := SelectOutputFile(
		secretID,
		callbackSelected,
		callbackReturn,
		callbackErr,
	)

	a.pages.AddAndSwitchToPage(
		registry.AddFilesWidgetPage,
		modal,
		true,
	)
}

func (a *FilesWidget) Remove(secretID string) {
	callback := func() { a.pages.SwitchToPage(registry.FilesWidgetPage) }

	if err := a.api.RemoveFile(secretID); err != nil {
		a.pages.RemovePage(registry.ErrorModalPage)
		modal := modals.NewErrorModal(
			fmt.Sprintf("failed to get file: %s", err),
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

func (a *FilesWidget) resetPagination() {
	a.mu.Lock()
	a.searchInput.offset = 0
	a.searchInput.limit = uint32(a.searchInput.step)
	a.mu.Unlock()
}

func (a *FilesWidget) Search() {
	a.resetPagination()
	a.table.Clean()

	resp, err := a.api.SearchFile(
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

func (a *FilesWidget) Paginate() {
	a.mu.Lock()
	a.searchInput.limit += uint32(a.searchInput.step)
	a.searchInput.offset += uint64(a.searchInput.step)
	a.mu.Unlock()

	resp, err := a.api.SearchFile(
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

type FilesTable struct {
	table *tview.Table
}

func NewTable() *FilesTable {
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

	filesTable := &FilesTable{table: table}
	filesTable.FillHeader()

	return filesTable
}

func (t *FilesTable) FillHeader() {
	const (
		numberColWidth int  = 1
		serverColWidth int  = 20
		loginColWidth  int  = 6
		notSelectable  bool = false
	)

	t.table.SetFixed(1, 0)
	t.table.SetCell(0, 0, tview.NewTableCell(registry.FilesColumnNumberTitle).
		SetSelectable(notSelectable).
		SetExpansion(numberColWidth).
		SetTextColor(styles.SecondAccentColor).
		SetAlign(tview.AlignCenter))

	t.table.SetCell(0, 1, tview.NewTableCell(registry.FilesColumnNameTitle).
		SetSelectable(notSelectable).
		SetExpansion(serverColWidth).
		SetTextColor(styles.SecondAccentColor).
		SetAlign(tview.AlignCenter))
}

func (t *FilesTable) Fill(searchResponse *filesv1.FileSearchResponse) {
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
			}

			t.table.SetCell(r, c, cell)
		}
	}
}

func (t *FilesTable) Clean() {
	t.table.Clear()
	t.FillHeader()
}

func NewAddFile(
	callbackAdd func(name, filePath string) error,
	callbackRefresh func(),
	callbackReturn func(),
) *tview.Flex {
	const (
		inputFieldWidth int  = 34
		formWidth       int  = 10
		formHeight      int  = 47
		fieldHight      int  = 5
		resizable       int  = 0
		oneWeight       int  = 1
		twoWeight       int  = 2
		dynamicColors   bool = false
		scrollable      bool = false
		focused         bool = true
		unfocused       bool = false
	)

	fileAdd := FileAdd{}

	frame := tview.NewFlex()

	infoLine := tview.NewTextView()

	form := tview.NewForm()
	form.SetTitle(" Add file ")
	styles.ApplyFormStyle(form)

	form.
		AddInputField("Name:", "", inputFieldWidth, nil, func(name string) {
			fileAdd.name = name
		}).
		AddInputField("FilePath:", "", inputFieldWidth, nil, func(filePath string) {
			fileAdd.filePath = filePath
		}).
		AddButton("Add", func() {
			if err := callbackAdd(
				fileAdd.name,
				fileAdd.filePath,
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

func SelectOutputFile(
	secretID string,
	callbackSelected func(secretID, filePath string, callbackErr func()),
	callbackReturn func(),
	callbackErr func(),
) *tview.Flex {
	const (
		inputFieldWidth int  = 34
		formWidth       int  = 17
		formHeight      int  = 47
		fieldHight      int  = 5
		resizable       int  = 0
		oneWeight       int  = 1
		twoWeight       int  = 2
		dynamicColors   bool = false
		scrollable      bool = false
		focused         bool = true
		unfocused       bool = false
	)

	fileAdd := FileAdd{}

	frame := tview.NewFlex()

	infoLine := tview.NewTextView()

	form := tview.NewForm()
	form.SetTitle(" Output file ")
	styles.ApplyFormStyle(form)

	form.
		AddInputField("FilePath:", "", inputFieldWidth, nil, func(filePath string) {
			fileAdd.filePath = filePath
		}).
		AddButton("Select", func() {
			callbackSelected(
				secretID,
				fileAdd.filePath,
				callbackErr,
			)
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
