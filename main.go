package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	// "fyne.io/fyne"
	// "fyne.io/fyne/app"
	// "fyne.io/fyne/container"
	// "fyne.io/fyne/v2/data/binding"
	// "fyne.io/fyne/widget"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	LOG_STREAM_TIME = 3
	FILE_NAME       = "chat.log"
)

var (
	daocLogs DaocLogs
)

func readChatFile(r fyne.URIReadCloser, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.URI())
}

type ListItem struct {
	Text string
}

func main() {
	daocLogs = DaocLogs{}

	myApp := app.New()
	myWindow := myApp.NewWindow("Dark Age of Camelot - Chat Parser")

	data := binding.BindStringList(
		&[]string{},
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	createdBy := widget.NewLabel("Created by Theorist")

	refresh := widget.NewButton("Refresh", func() {
		e := os.Remove(FILE_NAME)
		if e != nil {
			fmt.Println(e)
		} else {
			data = binding.BindStringList(
				&[]string{},
			)
			data.Reload()
		}

	})

	go func() {
		openLogFile(FILE_NAME)
		for _, item := range daocLogs.writeLogValues() {
			data.Append(item)
		}
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			daocLogs = DaocLogs{}
			data = binding.BindStringList(
				&[]string{},
			)
			openLogFile(FILE_NAME)
			for _, item := range daocLogs.writeLogValues() {
				data.Append(item)
			}
			data.Reload()
		}
	}()

	myWindow.Resize(fyne.NewSize(600, 300))
	myWindow.SetContent(container.NewBorder(createdBy, refresh, nil, nil, list))
	myWindow.ShowAndRun()
}

func openLogFile(logPath string) {
	f, err := os.OpenFile(logPath, os.O_RDONLY|os.O_EXCL, 0666)
	defer f.Close()
	if err == nil {
		iterateLogFile(f)
	}
}

func iterateLogFile(f *os.File) {
	reader := bufio.NewReader(f)
	style := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		style = daocLogs.regexOffensive(line, style)
		daocLogs.regexDefensives(line)
		daocLogs.regexSupport(line)
		daocLogs.regexPets(line)
		daocLogs.regexMisc(line)
		daocLogs.regexEnemy(line)
		daocLogs.regexTime(line)
	}
	daocLogs.writeLogValues()
}
