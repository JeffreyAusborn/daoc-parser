package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	// "fyne.io/fyne"
	// "fyne.io/fyne/app"
	// "fyne.io/fyne/container"
	// "fyne.io/fyne/v2/data/binding"
	// "fyne.io/fyne/widget"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type fn func(int)

const (
	LOG_STREAM_TIME = 3
	FILE_NAME       = "chat.log"
)

var (
	daocLogs DaocLogs
	mu       sync.Mutex
)

func readChatFile(r fyne.URIReadCloser, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.URI())
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

// TODO: Remove duplication on each render

func main() {
	daocLogs = DaocLogs{}
	myApp := app.New()
	myWindow := myApp.NewWindow("Dark Age of Camelot - Chat Parser - Theorist")
	openLogFile(FILE_NAME)
	allLogs, _ := renderAll(myWindow)
	damageInLogs, _ := renderDamageIn(myWindow)
	damageInOut, _ := renderDamagOut(myWindow)
	healLogs, _ := renderHeals(myWindow)
	combativeLogs, _ := renderCombatives(myWindow)

	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			daocLogs = DaocLogs{}
			openLogFile(FILE_NAME)
			mu.Unlock()
		}
	}()

	tabs := container.NewAppTabs(
		container.NewTabItem("All", allLogs),
		container.NewTabItem("Damage Out", damageInOut),
		container.NewTabItem("Damage In", damageInLogs),
		container.NewTabItem("Healing", healLogs),
		container.NewTabItem("Combatives", combativeLogs),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.Size{Width: 700, Height: 500})
	myWindow.ShowAndRun()
}

func renderAll(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.writeLogValues()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Combat Logs")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		if lii >= len(chatLogs) {
			return
		}
		val := chatLogs[lii]
		label := co.(*widget.Label)
		label.SetText(val)
	})
	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			chatLogs = daocLogs.writeLogValues()
			l.Refresh()
			mu.Unlock()
		}
	}()
	grid := container.New(layout.NewFormLayout())
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderHeals(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateHeal()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Combat Logs")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		if lii >= len(chatLogs) {
			return
		}
		val := chatLogs[lii]
		label := co.(*widget.Label)
		label.SetText(val)
	})
	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			chatLogs = daocLogs.calculateHeal()
			l.Refresh()
			mu.Unlock()
		}
	}()

	grid := container.New(layout.NewFormLayout())
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderDamageIn(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateDamageIn()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Combat Logs")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		if lii >= len(chatLogs) {
			return
		}
		val := chatLogs[lii]
		label := co.(*widget.Label)
		label.SetText(val)
	})
	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			chatLogs = daocLogs.calculateDamageIn()
			l.Refresh()
			mu.Unlock()
		}
	}()

	grid := container.New(layout.NewFormLayout())
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderDamagOut(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateDamageOut()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Combat Logs")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		if lii >= len(chatLogs) {
			return
		}
		val := chatLogs[lii]
		label := co.(*widget.Label)
		label.SetText(val)
	})
	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			chatLogs = daocLogs.calculateDamageOut()
			l.Refresh()
			mu.Unlock()
		}
	}()

	grid := container.New(layout.NewFormLayout())
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderCombatives(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.getCombativeUsers()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Combat Logs")
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		if lii >= len(chatLogs) {
			return
		}
		val := chatLogs[lii]
		label := co.(*widget.Label)
		label.SetText(val)
	})
	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			chatLogs = daocLogs.getCombativeUsers()
			l.Refresh()
			mu.Unlock()
		}
	}()

	grid := container.New(layout.NewFormLayout())
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}
