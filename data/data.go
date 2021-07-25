package data

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Entry struct {
	Date    time.Time
	Content string
}

type LogBook struct {
	Entrys []Entry
}

func Load(path string) (logbook LogBook, err error) {
	logbook = LogBook{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	// log.Print(string(data))
	err = json.Unmarshal(data, &logbook)

	return
}

func (logbook LogBook) Write(path string) (err error) {
	data, err := json.Marshal(logbook)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, data, 0644)
	return
}

// func (log LogBook) Write(path, content string) (err error) {
// 	err = ioutil.WriteFile(path, []byte(content), 0644)
// 	return
// }
