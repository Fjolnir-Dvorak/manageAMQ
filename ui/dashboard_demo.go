//// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
//// Use of this source code is governed by a MIT license that can
//// be found in the LICENSE file.
//
// +build ignore
//
package main
//
//import (
//	"math"
//
//	ui "github.com/gizak/termui"
//)
//
//func main() {
//	if err := ui.Init(); err != nil {
//		panic(err)
//	}
//	defer ui.Close()
//
//	simpleText := ui.NewPar(":PRESS q TO QUIT DEMO:")
//	simpleText.Height = 3
//	simpleText.Width = len(simpleText.Text) + 2
//	simpleText.TextFgColor = ui.ColorWhite
//	//simpleText.BorderLabel = "Text Box"
//	simpleText.BorderFg = ui.ColorCyan
//	//simpleText.Handle("/timer/1s", func(e ui.Event) {
//	//	cnt := e.Data.(ui.EvtTimer)
//	//	if cnt.Count%2 == 0 {
//	//		simpleText.TextFgColor = ui.ColorRed
//	//	} else {
//	//		simpleText.TextFgColor = ui.ColorWhite
//	//	}
//	//})
//
//	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] interrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
//	list := ui.NewList()
//	list.Items = strs
//	list.ItemFgColor = ui.ColorYellow
//	list.BorderLabel = "List"
//	list.Height = 7
//	list.Width = 25
//	list.Y = 4
//
//	gauge := ui.NewGauge()
//	gauge.Percent = 50
//	gauge.Width = 50
//	gauge.Height = 3
//	gauge.Y = 11
//	gauge.BorderLabel = "Gauge"
//	gauge.BarColor = ui.ColorRed
//	gauge.BorderFg = ui.ColorWhite
//	gauge.BorderLabelFg = ui.ColorCyan
//
//	spark_srv0 := ui.Sparkline{}
//	spark_srv0.Height = 1
//	spark_srv0.Title = "srv 0:"
//	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
//	spark_srv0.Data = spdata
//	spark_srv0.LineColor = ui.ColorCyan
//	spark_srv0.TitleColor = ui.ColorWhite
//
//	spark_srv1 := ui.Sparkline{}
//	spark_srv1.Height = 1
//	spark_srv1.Title = "srv 1:"
//	spark_srv1.Data = spdata
//	spark_srv1.TitleColor = ui.ColorWhite
//	spark_srv1.LineColor = ui.ColorRed
//
//	sparkline := ui.NewSparklines(spark_srv0, spark_srv1)
//	sparkline.Width = 25
//	sparkline.Height = 7
//	sparkline.BorderLabel = "Sparkline"
//	sparkline.Y = 4
//	sparkline.X = 25
//
//	sinps := (func() []float64 {
//		n := 220
//		ps := make([]float64, n)
//		for i := range ps {
//			ps[i] = 1 + math.Sin(float64(i)/5)
//		}
//		return ps
//	})()
//
//	dotLineChart := ui.NewLineChart()
//	dotLineChart.BorderLabel = "dot-mode Line Chart"
//	dotLineChart.Data["default"] = sinps
//	dotLineChart.Width = 50
//	dotLineChart.Height = 11
//	dotLineChart.X = 0
//	dotLineChart.Y = 14
//	dotLineChart.AxesColor = ui.ColorWhite
//	dotLineChart.LineColor["default"] = ui.ColorRed | ui.AttrBold
//	dotLineChart.Mode = "dot"
//
//	barChart := ui.NewBarChart()
//	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
//	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
//	barChart.BorderLabel = "Bar Chart"
//	barChart.Width = 26
//	barChart.Height = 10
//	barChart.X = 51
//	barChart.Y = 0
//	barChart.DataLabels = bclabels
//	barChart.BarColor = ui.ColorGreen
//	barChart.NumColor = ui.ColorBlack
//
//	brailleLineChart := ui.NewLineChart()
//	brailleLineChart.BorderLabel = "braille-mode Line Chart"
//	brailleLineChart.Data["default"] = sinps
//	brailleLineChart.Width = 26
//	brailleLineChart.Height = 11
//	brailleLineChart.X = 51
//	brailleLineChart.Y = 14
//	brailleLineChart.AxesColor = ui.ColorWhite
//	brailleLineChart.LineColor["default"] = ui.ColorYellow | ui.AttrBold
//
//	p1 := ui.NewPar("Hey!\nI am a borderless block!")
//	p1.Border = false
//	p1.Width = 26
//	p1.Height = 2
//	p1.TextFgColor = ui.ColorMagenta
//	p1.X = 52
//	p1.Y = 11
//
//	draw := func(t int) {
//		gauge.Percent = t % 101
//		list.Items = strs[t%9:]
//		sparkline.Lines[0].Data = spdata[:30+t%50]
//		sparkline.Lines[1].Data = spdata[:35+t%50]
//		dotLineChart.Data["default"] = sinps[:t/2%220]
//		brailleLineChart.Data["default"] = sinps[:2*t%220]
//		barChart.Data = bcdata[t/2%10:]
//		ui.Render(simpleText, list, gauge, sparkline, dotLineChart, barChart, brailleLineChart, p1)
//	}
//	ui.Handle("/sys/kbd/q", func(ui.Event) {
//		ui.StopLoop()
//	})
//	ui.Handle("/timer/1s", func(e ui.Event) {
//		t := e.Data.(ui.EvtTimer)
//		draw(int(t.Count))
//	})
//	ui.Loop()
//}
