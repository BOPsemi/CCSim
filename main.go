package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ProgressBar struct {
	max      int
	progress int
	percent  float64
}

func main() {
	// 乱数生成器作成
	rand.Seed(time.Now().UnixNano())
	// プログレスバーを表す構造体
	progressBar := ProgressBar{max: 30, progress: 0}

	for {
		// 加算処理
		randNum := rand.Intn(10)
		progressBar.progress += int(randNum)

		// max値を超えないようにする
		if progressBar.progress > progressBar.max {
			progressBar.progress -= (progressBar.progress - progressBar.max)
		}

		// プログレスのパーセンテージ
		progressBar.percent = (float64(progressBar.progress) / float64(progressBar.max)) * 100
		// いくつ=を表示するかを決めるため、intに変換
		roundPercent := int(progressBar.percent)
		// ===> のように >付きで表示したいので、-1する。
		if roundPercent > 0 {
			roundPercent -= 1
		}
		fmt.Print(fmt.Sprintf("%5.1f%% [%s>%s] %d \r", progressBar.percent, strings.Repeat("=", roundPercent), strings.Repeat(" ", 99-roundPercent), progressBar.progress))

		if progressBar.progress == progressBar.max {
			fmt.Println("")
			break
		}

		time.Sleep(1 * time.Second)
	}
}
