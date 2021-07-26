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

	if logBook.Entrys == nil {
		logBook.Entrys = []data.Entry{}
	}

	p := tea.NewProgram(initialModel(logBook))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type tickMsg struct{}
type errMsg error

type model struct {
	cursor    int
	logBook   data.LogBook
	textInput textinput.Model
	err       error
}

func initialModel(logBook data.LogBook) model {
	ti := textinput.NewModel()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	return model{
		cursor:    -1,
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
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.logBook.Entrys) - 1
			}
			break
		case tea.KeyDown:
			if m.cursor < len(m.logBook.Entrys)-1 {
				m.cursor++
			} else {
				m.cursor = -1
			}
			break
		case tea.KeyEnter:
			if m.cursor > -1 && len(m.logBook.Entrys) > m.cursor {
				entry := m.logBook.Entrys[m.cursor]
				entry.Content = m.textInput.Value()
				m.logBook.Entrys[m.cursor] = entry
			} else {
				m.logBook.Entrys = append(m.logBook.Entrys, data.Entry{Date: time.Now(), Content: m.textInput.Value()})
			}
			err := m.logBook.Write(dataFilePath)
			if err != nil {
				log.Fatal(err)
			}
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	lowerbound, upperbound := calculateBoundaries(len(m.logBook.Entrys), m.cursor)
	for i, v := range m.logBook.Entrys {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		if i <= upperbound && i >= lowerbound {
			s += fmt.Sprintf("%s [%s] %s\r\n", cursor, v.Date.Format(time.RFC822), v.Content)
		}
	}

	return s + fmt.Sprintf(
		"What are you working on?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func calculateBoundaries(count, index int) (lowerbound, upperbound int) {
	if count <= 10 {
		lowerbound = 0
		upperbound = 10
		return
	}
	maxIndex := count - 1
	if index == -1 {
		upperbound = maxIndex
		lowerbound = maxIndex - 8
		if lowerbound < 0 {
			lowerbound = 0
		}
		return
	}
	lowerbound = index - 4
	upperbound = index + 4
	if upperbound > maxIndex {
		lowerbound -= upperbound % maxIndex
		upperbound = maxIndex
	} else if lowerbound < 0 {
		upperbound -= lowerbound
		lowerbound = 0
	}

	return
}
