package helpers

import (
	"log"
	"time"

	"github.com/inancgumus/screen"
)

func DownloadStatus(downloads *[]string) *time.Ticker {
	const refreshTime = 200

	rot := [4]string{"|", "/", "—", "\\"}
	rotations := len(rot) - 1
	ticker := time.NewTicker(refreshTime * time.Millisecond)
	pos := 1

	go func() {
		for range ticker.C {
			if pos > rotations {
				pos = 0
			}

			screen.Clear()
			screen.MoveTopLeft()

			for _, v := range *downloads {
				log.Println(rot[pos], " ", v)
			}
			pos++
		}
	}()

	return ticker
}
