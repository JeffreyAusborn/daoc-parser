package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	fyne "fyne.io/fyne"
	fa "fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
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

func main() {
	daocLogs = DaocLogs{}
	a := fa.New()
	w := a.NewWindow("Dark Age of Camelot - Chat Parser")
	damageLabel := widget.NewLabel("")
	createdBy := widget.NewLabel("Created by Theorist")

	w.SetContent(container.NewVBox(
		damageLabel,
		widget.NewButton("Refresh", func() {
			e := os.Remove(FILE_NAME)
			if e != nil {
				log.Fatal(e)
			}
		}),
		createdBy,
	))

	go func() {
		openLogFile(FILE_NAME)
		damageLabel.SetText(daocLogs.writeLogValues())
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			daocLogs = DaocLogs{}
			openLogFile(FILE_NAME)
			damageLabel.SetText(daocLogs.writeLogValues())
		}
	}()

	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
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
