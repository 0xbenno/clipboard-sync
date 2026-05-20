package monitor

import (
	"context"

	"golang.design/x/clipboard"
)

func StartMonitor() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	last := ""
	for data := range ch {
		data := string(data)
		if data != last {
			println(data)
			last = data
		}
	}
}

func ReadClipboard() string {
	contents := clipboard.Read(clipboard.FmtText)
	return string(contents)
}

func WriteClipboard(content string) {
	clipboard.Write(clipboard.FmtText, []byte(content))
}
