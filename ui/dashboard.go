// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// // +build ignore

package main

import (
	"fmt"
	"github.com/Fjolnir-Dvorak/manageAMQ/queue"
	ui "github.com/gizak/termui"
)

var (
	list *queue.FileList
	filenames []string
	runner queue.ActiveRunner
)

const(
	FmtFileProcess = "File %d / %d:"
	fmtLineProcess = "Line %d / %d:"
)

func startUI(activeRunner queue.ActiveRunner) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	runner = activeRunner
	list = &runner.FileList

	nodes := make([]ui.Bufferer, 0, 5)

	helpQuit := createTextBox(":PRESS q TO QUIT:")
	nodes = append(nodes, helpQuit)
	helpPause := createTextBox(":PRESS \u2423 TO PAUSE:")
	helpPause.X = helpQuit.X + helpQuit.Width
	nodes = append(nodes, helpPause)

	fileListProgress := createGauge(fmt.Sprintf(FmtFileProcess, list.CurrentFile, list.TotalFiles))
	fileListProgress.Y = helpQuit.Y + helpQuit.Height
	nodes = append(nodes, fileListProgress)

	currentFileProgress := createGauge(fmt.Sprintf(fmtLineProcess, list.Current().ReadingPosition, list.Current().TotalLines))
	currentFileProgress.Y = helpQuit.Y + helpQuit.Height
	currentFileProgress.X = fileListProgress.X + fileListProgress.Width
	nodes = append(nodes, currentFileProgress)

	listView := createFileListView(list)
	listView.Y = fileListProgress.Y + fileListProgress.Height
	nodes = append(nodes, listView)

	draw := func() {
		fileListProgress.Percent = list.CurrentFile / list.TotalFiles
		fileListProgress.Label = fmt.Sprintf(FmtFileProcess, list.CurrentFile, list.TotalFiles)

		currentFile := list.Files[list.CurrentFile]
		currentFileProgress.Percent = currentFile.ReadingPosition / currentFile.TotalLines
		currentFileProgress.Label = fmt.Sprintf(fmtLineProcess, list.Current().ReadingPosition, list.Current().TotalLines)
		ui.Render(nodes...)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		runner.Lock()
		runner.Running = !runner.Running
		runner.Unlock()
		runner.Cond.Broadcast()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})
	runner.Lock()
	runner.Running = true
	runner.Unlock()
	runner.Cond.Broadcast()
	ui.Loop()
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
	names := make([]string, list.TotalFiles)
	for index, single := range list.Files {
		names = append(names, fmt.Sprintf("[%d] %s", index + 1, single.ReprenstativeName))
	}
	filenames = names

	listView := ui.NewList()
	listView.Items = filenames[list.CurrentFile - 1:]
	listView.ItemFgColor = ui.ColorYellow
	listView.BorderLabel = "List"
	listView.Height = 7
	listView.Width = 50
	return listView
}