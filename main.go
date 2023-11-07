package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type fn func(int)

const (
	LOG_STREAM_TIME = 2
	FILE_NAME       = "chat.log"
)

var (
	daocLogs DaocLogs
	mu       sync.Mutex
)

func openLogFile() {
	f, err := os.OpenFile(FILE_NAME, os.O_RDONLY|os.O_EXCL, 0666)
	defer f.Close()
	if err == nil {
		iterateLogFile(f)
	}
}

var (
	tempLines []string
)

func iterateLogFile(f *os.File) {
	reader := bufio.NewReader(f)
	style := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		daocLogs.regexOffensive(line, style)
		daocLogs.regexDefensives(line)
		daocLogs.regexSupport(line)
		daocLogs.regexPets(line)
		daocLogs.regexMisc(line)
		daocLogs.regexEnemy(line)
		daocLogs.regexTime(line)
	}
}

// TODO: Remove duplication on each render
func main() {
	daocLogs = DaocLogs{}
	myApp := app.New()
	myWindow := myApp.NewWindow("Dark Age of Camelot - Chat Parser - Theorist")
	openLogFile()
	spellLogs, _ := renderSpells(myWindow)
	killLog, _ := renderKills(myWindow)
	damageArmorLogs, _ := renderArmorDamage(myWindow)
	styleLogs, _ := renderStyles(myWindow)
	healLogs, _ := renderHeals(myWindow)
	defensiveLogs, _ := renderDefensives(myWindow)
	combativeLogs, _ := renderCombatives(myWindow)

	go func() {
		for range time.Tick(time.Second * LOG_STREAM_TIME) {
			mu.Lock()
			daocLogs = DaocLogs{}
			openLogFile()
			mu.Unlock()
		}
	}()

	tabs := container.NewAppTabs(
		container.NewTabItem("Spells", spellLogs),
		container.NewTabItem("Styles", styleLogs),
		container.NewTabItem("Healing and Absorb", healLogs),
		container.NewTabItem("Defensive", defensiveLogs),
		container.NewTabItem("Armor Damaged", damageArmorLogs),
		container.NewTabItem("Enemy", combativeLogs),
		container.NewTabItem("Kills", killLog),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.Size{Width: 700, Height: 500})
	myWindow.ShowAndRun()
}

func renderSpells(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateSpells()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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
			chatLogs = daocLogs.calculateSpells()
			l.Refresh()
			mu.Unlock()
		}
	}()

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderHeals(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateHeal()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderDefensives(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateDensives()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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
			chatLogs = daocLogs.calculateDensives()
			l.Refresh()
			mu.Unlock()
		}
	}()

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderKills(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.getKills()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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
			chatLogs = daocLogs.getKills()
			l.Refresh()
			mu.Unlock()
		}
	}()

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderArmorDamage(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateArmorhits()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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
			chatLogs = daocLogs.calculateArmorhits()
			l.Refresh()
			mu.Unlock()
		}
	}()

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderStyles(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.calculateStyles()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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
			chatLogs = daocLogs.calculateStyles()
			l.Refresh()
			mu.Unlock()
		}
	}()

	resetLabel := widget.NewLabel("")
	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), resetLabel, resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}

func renderCombatives(w fyne.Window) (fyne.CanvasObject, error) {
	chatLogs := daocLogs.getCombativeUsers()
	l := widget.NewList(func() int {
		return len(chatLogs)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
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

	resetbtn := widget.NewButton("Reset Logs - This will delete your log file.", func() {
		mu.Lock()
		e := os.Remove(FILE_NAME)
		mu.Unlock()
		if e != nil {
			fmt.Println(e)
		}
	})

	grid := container.New(layout.NewFormLayout(), widget.NewLabel(""), resetbtn)
	tab := container.New(layout.NewBorderLayout(grid, nil, nil, nil), grid, l)
	return tab, nil
}
