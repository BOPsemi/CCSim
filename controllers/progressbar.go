package controllers

import (
	"fmt"
	"strings"
	"time"
)

/*
ProgressBarCUI :Progress bar
*/
type ProgressBarCUI struct {
	max      int
	progress int
	percent  float64
}

/*
NewProgressBarCUI :Initializer
*/
func NewProgressBarCUI(maxValue int) *ProgressBarCUI {
	return &ProgressBarCUI{
		max:      maxValue,
		progress: 0,
		percent:  0.0,
	}
}

func (pg *ProgressBarCUI) countProgress(count int) {
	for {
		pg.progress++

		if pg.progress > pg.max {
			pg.progress -= (pg.progress - pg.max)
		}

		pg.percent = (float64(pg.progress) / float64(pg.max)) * 100
		roundPercent := int(pg.percent)

		if roundPercent > 0 {
			roundPercent--
		}

		fmt.Print(fmt.Sprintf("%5.1f%% [%s>%s] %d \r", pg.percent, strings.Repeat("=", roundPercent), strings.Repeat(" ", 99-roundPercent), pg.progress))

		if pg.progress == pg.max {
			fmt.Println("")
			break
		}

		time.Sleep(1 * time.Second)

	}
}
