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
	list.Resize(fyne.Size{Height: 434})

	// add := widget.NewButton("Append", func() {
	// 	val := fmt.Sprintf("Item %d", data.Length()+1)
	// 	data.Append(val)
	// })

	// a := app.New()
	// w := a.NewWindow("Dark Age of Camelot - Chat Parser")

	// damageLabel := widget.NewLabel("")
	// createdBy := widget.NewLabel("Created by Theorist")

	// w.SetContent(container.NewVBox(
	// 	damageLabel,
	// 	widget.NewButton("Refresh", func() {
	// 		e := os.Remove(FILE_NAME)
	// 		if e != nil {
	// 			fmt.Println(e)
	// 		}
	// 	}),
	// 	createdBy,
	// ))

	go func() {
		openLogFile(FILE_NAME)
		// fmt.Println(daocLogs.writeLogValues())
		for _, item := range daocLogs.writeLogValues() {
			data.Append(item)
		}
		// damageLabel.SetText(daocLogs.writeLogValues())
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			daocLogs = DaocLogs{}
			data = binding.BindStringList(
				&[]string{},
			)
			openLogFile(FILE_NAME)
			// fmt.Println(daocLogs.writeLogValues())
			// damageLabel.SetText(daocLogs.writeLogValues())
			for _, item := range daocLogs.writeLogValues() {
				data.Append(item)
			}
			data.Reload()
		}
	}()

	myWindow.Resize(fyne.NewSize(600, 300))
	// w.SetContent(scrollContainer)
	// w.ShowAndRun()
	myWindow.SetContent(container.NewBorder(nil, nil, nil, nil, list))
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
