package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	App        *tview.Application
	LogView    *tview.TextView
	InputField *tview.InputField
	mu         sync.Mutex
	root       tview.Primitive
	inputChan  chan string
}

func NewUI() *UI {
	app := tview.NewApplication()

	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	logView.
		SetBorder(true).
		SetTitle(" Messages / Events ")

	inputChan := make(chan string, 100)
	input := tview.NewInputField().
		SetLabel(">>> ").
		SetFieldWidth(0)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := input.GetText()
			if text != "" {
				inputChan <- text
				input.SetText("")
			}
		}
	})

	input.
		SetBorder(true).
		SetTitle(" Input ")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logView, 0, 1, false).
		AddItem(input, 3, 0, true)

	app.SetRoot(flex, true)

	return &UI{
		App:        app,
		LogView:    logView,
		InputField: input,
		root:       flex,
		inputChan:  inputChan,
	}
}

func (u *UI) Root() tview.Primitive {
	return u.root
}

func (ui *UI) Log(format string, args ...any) {
	ui.App.QueueUpdateDraw(func() {
		ui.LogView.Write([]byte(
			fmt.Sprintf(format+"\n", args...),
		))
	})
}

func (ui *UI) GetInputChannel() <-chan string {
	return ui.inputChan
}
