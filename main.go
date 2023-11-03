package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

const (
	LOG_STREAM_TIME = 3
)

func main() {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	logPath := flag.String("file", "", "Path to chat.log")
	streamLogs := flag.Bool("stream", false, "")
	flag.Parse()
	if *logPath != "" {
		openLogFile(*logPath, writer)
		if *streamLogs {
			for range time.Tick(time.Second * LOG_STREAM_TIME) {
				openLogFile(*logPath, writer)
			}
		}
	} else {
		fmt.Println("Requied argument missing: --file path/to/chat.log")
	}
}

func openLogFile(logPath string, writer *uilive.Writer) {
	f, err := os.OpenFile(logPath, os.O_RDONLY|os.O_EXCL, 0666)
	defer f.Close()
	if err == nil {
		iterateLogFile(f, writer)
	}
}

func iterateLogFile(f *os.File, writer *uilive.Writer) {
	var daocLogs DaocLogs
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
	daocLogs.writeLogValues(writer)
}
