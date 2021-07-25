package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
	"timekeeper/data"
	"timekeeper/helper"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var path = helper.CorrectPath("~/AppData/Roaming/timekeeper")
var dataFilePath = filepath.Join(path, "data.json")

func main() {
	path = filepath.Clean(path)
	helper.MakeDirectoryIfNotExists(path)
	logBook, err := data.Load(dataFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	p := tea.NewProgram(initialModel(logBook))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type tickMsg struct{}
type errMsg error

type model struct {
	logBook   data.LogBook
	textInput textinput.Model
	err       error
}

func initialModel(logBook data.LogBook) model {
	ti := textinput.NewModel()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		logBook:   logBook,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.logBook.Entrys == nil {
				m.logBook.Entrys = []data.Entry{}
			}
			m.logBook.Entrys = append(m.logBook.Entrys, data.Entry{Date: time.Now(), Content: m.textInput.Value()})
			err := m.logBook.Write(dataFilePath)
			if err != nil {
				log.Fatal(err)
			}
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var entrys string
	if m.logBook.Entrys != nil {
		for _, v := range m.logBook.Entrys {
			entrys += fmt.Sprintf("[%s] %s\r\n", v.Date.Format(time.RFC822), v.Content)
		}
		entrys += "\r\n"
	}
	return entrys + fmt.Sprintf(
		"What are you working on?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
