package main

import (
	"log"
	"path/filepath"
	"timekeeper/data"
	"timekeeper/helper"
	"timekeeper/ui"

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

	p := tea.NewProgram(ui.InitialModel(logBook, dataFilePath))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
