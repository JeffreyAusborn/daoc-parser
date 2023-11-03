package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	LOG_STREAM_TIME = 3
)

var (
	daocLogs DaocLogs
)

func main() {
	daocLogs = DaocLogs{}
	a := app.New()
	w := a.NewWindow("Dark Age of Camelot - Chat Parser\nWritten by: Theorist\nIf you have any feedback, feel free to DM in Discord.\n\n")
	w.Resize(fyne.NewSize(300.0, 300.0))
	damageLabel := widget.NewLabel("")

	go func() {
		logPath := flag.String("file", "", "Path to chat.log")
		streamLogs := flag.Bool("stream", false, "")
		flag.Parse()
		if *logPath != "" {
			openLogFile(*logPath)
			damageLabel.SetText(daocLogs.writeLogValues())
			if *streamLogs {
				for range time.Tick(time.Second * LOG_STREAM_TIME) {
					daocLogs = DaocLogs{}
					openLogFile(*logPath)
					damageLabel.SetText(daocLogs.writeLogValues())
				}
			}
		} else {
			fmt.Println("Requied argument missing: --file path/to/chat.log")
		}
	}()

	w.SetContent(container.NewVBox(damageLabel))
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
