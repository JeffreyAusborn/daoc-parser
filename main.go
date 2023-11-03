package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	fyne "fyne.io/fyne"
	fa "fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

const (
	LOG_STREAM_TIME = 3
)

var (
	daocLogs    DaocLogs
	chatLogFile string
)

func readChatFile(r fyne.URIReadCloser, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.URI())
}

func main() {
	logPath := flag.String("file", "", "Path to chat.log")
	streamLogs := flag.Bool("stream", false, "")
	flag.Parse()
	if *logPath == "" {
		log.Fatal("Requied argument missing: --file path/to/chat.log")
	}
	daocLogs = DaocLogs{}
	a := fa.New()
	w := a.NewWindow("Dark Age of Camelot - Chat Parser - By: Theorist")
	damageLabel := widget.NewLabel("")

	w.SetContent(container.NewVBox(
		damageLabel,
		widget.NewButton("Load Chat Log", func() {
			fd := dialog.NewFileOpen(func(fyne.URIReadCloser, error) {
			}, w)
			fd.Show()
		}),
		widget.NewButton("Refresh", func() {
			e := os.Remove(*logPath)
			if e != nil {
				log.Fatal(e)
			}
		}),
	))

	go func() {
		openLogFile(*logPath)
		damageLabel.SetText(daocLogs.writeLogValues())
		if *streamLogs {
			for range time.Tick(time.Second * LOG_STREAM_TIME) {
				daocLogs = DaocLogs{}
				openLogFile(*logPath)
				damageLabel.SetText(daocLogs.writeLogValues())
			}
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
