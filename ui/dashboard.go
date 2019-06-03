// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// // +build ignore

package ui

import (
	"bufio"
	"fmt"
	"github.com/Fjolnir-Dvorak/manageAMQ/queue"
	ui "github.com/gizak/termui"
	"os"
)

var (
	list *queue.FileList
	filenames []string
	runner *queue.ActiveRunner
)

const(
	FmtFileProcess = "File %d / %d:"
	fmtLineProcess = "Line %d / %d:"
)

func StartUI(activeRunner *queue.ActiveRunner) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	runner = activeRunner
	list = runner.FileList

	nodes := make([]ui.Bufferer, 0, 5)

	helpQuit := createTextBox(":PRESS q TO QUIT:")
	nodes = append(nodes, helpQuit)
	helpPause := createTextBox(":PRESS \u2423 TO PAUSE:")
	helpPause.X = helpQuit.X + helpQuit.Width
	nodes = append(nodes, helpPause)

	fileListProgress := createGauge(fmt.Sprintf(FmtFileProcess, list.CurrentFile, list.TotalFiles))
	fileListProgress.Y = helpQuit.Y + helpQuit.Height
	nodes = append(nodes, fileListProgress)

	currentFileProgress := createGauge(fmt.Sprintf(fmtLineProcess, list.CurrentOrFirst().ReadingPosition, list.CurrentOrFirst().TotalLines))
	currentFileProgress.Y = helpQuit.Y + helpQuit.Height
	currentFileProgress.X = fileListProgress.X + fileListProgress.Width
	nodes = append(nodes, currentFileProgress)

	listView := createFileListView(list)
	listView.Y = fileListProgress.Y + fileListProgress.Height
	nodes = append(nodes, listView)

	draw := func() {
		fileListProgress.Percent = int((float64(list.CurrentFile) / float64(list.TotalFiles)) * 100.0)
		fileListProgress.BorderLabel = fmt.Sprintf(FmtFileProcess, list.CurrentFile, list.TotalFiles)

		currentFile := list.CurrentOrFirst()
		currentFileProgress.Percent = int((float64(currentFile.ReadingPosition) / float64(currentFile.TotalLines)) * 100.0)
		currentFileProgress.BorderLabel = fmt.Sprintf(fmtLineProcess, currentFile.ReadingPosition, currentFile.TotalLines)

		currentFileIndex := list.CurrentFile
		if currentFileIndex > 0 {
			currentFileIndex = currentFileIndex - 1
		}
		listView.Items = filenames[(currentFileIndex):]
		ui.Render(nodes...)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
		runner.Lock()
		runner.Running = false
		runner.Paused = false
		runner.Unlock()
		runner.Channel <- struct{}{}
		runner.Waiter.Wait()
		return
	})
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		runner.Lock()
		runner.Paused = !runner.Paused
		runner.Unlock()
		runner.Channel <- struct{}{}
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	if !waitForUserInput() {
		runner.Lock()
		runner.Running = false
		runner.Paused = false
		runner.Unlock()
		runner.Channel <- struct{}{}
		runner.Waiter.Wait()
		return
	}

	runner.Lock()
	runner.Paused = false
	runner.Unlock()
	runner.Channel <- struct{}{}
	ui.Loop()
}

func waitForUserInput() bool {
	if true {
		return true
	}
	fmt.Println("Weitermachen?")
	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadRune()
	if err != nil {
		return false
	}
	switch input {
	case 0x000A, 'y', 'Y':
		fmt.Println("... yes")
		return true
	case 'a', 'A':
		fmt.Println("... Aborting")
		return false
	default:
		fmt.Println("... That was no valid character. Perhaps you meant to say No...")
		return false
	}
}

func createGauge(baseText string) *ui.Gauge{
	gauge := ui.NewGauge()
	gauge.Percent = 50
	gauge.Width = 50
	gauge.Height = 3
	gauge.BorderLabel = baseText
	gauge.BarColor = ui.ColorMagenta
	gauge.PercentColor = ui.ColorBlue
	gauge.PercentColorHighlighted = ui.ColorBlack
	gauge.BorderFg = ui.ColorWhite
	gauge.BorderLabelFg = ui.ColorBlue
	return gauge
}
func createTextBox(text string) *ui.Par {
	simpleText := ui.NewPar(text)
	simpleText.Height = 3
	simpleText.Width = len(simpleText.Text) + 2
	simpleText.TextFgColor = ui.ColorWhite
	simpleText.BorderFg = ui.ColorBlue
	return simpleText
}
func createFileListView(list *queue.FileList) *ui.List {
	names := make([]string, 0, list.TotalFiles)
	for index, single := range list.Files {
		names = append(names, fmt.Sprintf("[%d] %s", index + 1, single.ReprenstativeName))
	}
	filenames = names

	listView := ui.NewList()
	listView.Items = filenames
	listView.ItemFgColor = ui.ColorMagenta
	listView.BorderLabelFg = ui.ColorBlue
	listView.BorderLabel = "List"
	listView.Height = 7
	listView.Width = 50
	return listView
}