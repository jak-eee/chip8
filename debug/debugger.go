package debug

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// DebugProvider is what the debugger expects.
type DebugProvider interface {
	GetDelayTimer() byte
	GetSoundTimer() byte
	GetReg() [16]byte
	GetKeyPressed() [16]bool
	GetMem() [4096]byte
	GetFramesI() uint64
	GetFrameTimes() []time.Duration
	GetFramesAvg() float64
	GetOpCodes() []uint16
	GetOpCodesI() uint64
	GetOpCodesAvg() float64
}

// Debugger collects and displays system information on in the terminal.
type Debugger struct {
	Chip8 DebugProvider
}

// Start displays debug info on screen until quit<-true.
func (d *Debugger) Start(quit chan bool) {
	go func() {
		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		l1 := widgets.NewList()
		l1.Title = "Registers"
		l1.Rows = d.registerRows()
		l1.TextStyle = ui.NewStyle(ui.ColorYellow)
		l1.WrapText = false
		l1.SetRect(0, 0, 20, 18)

		l2 := widgets.NewList()
		l2.Title = "Last OpCodes"
		l2.Rows = d.opCodesRows()
		l2.TextStyle = ui.NewStyle(ui.ColorYellow)
		l2.WrapText = false
		l2.SetRect(20, 0, 40, 18)

		l3 := widgets.NewList()
		l3.Title = "OpCodes I"
		l3.Rows = d.opCodesIRows()
		l3.TextStyle = ui.NewStyle(ui.ColorYellow)
		l3.WrapText = false
		l3.SetRect(40, 0, 55, 3)

		l4 := widgets.NewList()
		l4.Title = "OpCodes / S"
		l4.Rows = d.opCodesAvgRows()
		l4.TextStyle = ui.NewStyle(ui.ColorYellow)
		l4.WrapText = false
		l4.SetRect(40, 3, 55, 6)

		l5 := widgets.NewList()
		l5.Title = "Frames I"
		l5.Rows = d.framesIRows()
		l5.TextStyle = ui.NewStyle(ui.ColorYellow)
		l5.WrapText = false
		l5.SetRect(40, 6, 55, 9)

		l6 := widgets.NewList()
		l6.Title = "Frames / S"
		l6.Rows = d.framesAvgRows()
		l6.TextStyle = ui.NewStyle(ui.ColorYellow)
		l6.WrapText = false
		l6.SetRect(40, 9, 55, 12)

		l7 := widgets.NewList()
		l7.Title = "Delay Timer"
		l7.Rows = d.delayTimerRows()
		l7.TextStyle = ui.NewStyle(ui.ColorYellow)
		l7.WrapText = false
		l7.SetRect(40, 12, 55, 15)

		l8 := widgets.NewList()
		l8.Title = "Sound Timer"
		l8.Rows = d.soundTimerRows()
		l8.TextStyle = ui.NewStyle(ui.ColorYellow)
		l8.WrapText = false
		l8.SetRect(40, 15, 55, 18)

		l9 := widgets.NewList()
		l9.Title = "Keys Pressed"
		l9.Rows = d.keyPressedRows()
		l9.TextStyle = ui.NewStyle(ui.ColorYellow)
		l9.WrapText = false
		l9.SetRect(55, 0, 75, 18)

		p1 := widgets.NewPlot()
		p1.Title = "Cycle Time ms"
		p1.Marker = widgets.MarkerBraille
		p1.Data = d.frameTimeRows()
		p1.SetRect(0, 18, 75, 28)
		p1.DotMarkerRune = '+'
		p1.AxesColor = ui.ColorWhite
		p1.LineColors[0] = ui.ColorYellow
		p1.DrawDirection = widgets.DrawLeft

		ticker := time.NewTicker(time.Millisecond * 333)
		for {
			select {
			case <-quit:
				{
					return
				}
			case <-ticker.C:
				{
					l1.Rows = d.registerRows()

					l2.Rows = d.opCodesRows()

					l3.Rows = d.opCodesIRows()

					l4.Rows = d.opCodesAvgRows()

					l5.Rows = d.framesIRows()

					l6.Rows = d.framesAvgRows()

					l7.Rows = d.delayTimerRows()

					l8.Rows = d.soundTimerRows()

					l9.Rows = d.keyPressedRows()

					p1.Data = d.frameTimeRows()

					ui.Render(l1, l2, l3, l4, l5, l6, l7, l8, p1, l9)
				}
			}
		}
	}()
}

func (d *Debugger) delayTimerRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetDelayTimer())}
}

func (d *Debugger) soundTimerRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetSoundTimer())}
}

func (d *Debugger) framesIRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetFramesI())}
}

func (d *Debugger) framesAvgRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetFramesAvg())}
}

func (d *Debugger) opCodesRows() []string {
	var rows []string
	for i, v := range d.Chip8.GetOpCodes() {
		rows = append(rows, fmt.Sprintf("%d - 0x%x", i, v))
	}
	return rows
}

func (d *Debugger) opCodesIRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetOpCodesI())}
}

func (d *Debugger) opCodesAvgRows() []string {
	return []string{fmt.Sprint(d.Chip8.GetOpCodesAvg())}
}

func (d *Debugger) keyPressedRows() []string {
	var rows []string
	for i, v := range d.Chip8.GetKeyPressed() {
		rows = append(rows, fmt.Sprintf("Key[0x%x] = %t", i, v))
	}
	return rows
}

func (d *Debugger) registerRows() []string {
	var rows []string
	for i, v := range d.Chip8.GetReg() {
		rows = append(rows, fmt.Sprintf("Reg[0x%x] = 0x%x", i, v))
	}
	return rows
}

func (d *Debugger) frameTimeRows() [][]float64 {
	data := make([][]float64, 1)
	data[0] = make([]float64, 64)
	for i, t := range d.Chip8.GetFrameTimes() {
		data[0][i] = float64(t.Milliseconds())
	}
	return data
}
