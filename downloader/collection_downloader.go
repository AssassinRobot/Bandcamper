package downloader

import (
	"fmt"
	"log"

	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type collectionDownloader struct {
	http     *utils.HttpMngmnt
	file     *utils.FileMngmnt
	email    *utils.EmailMngmnt
	scrapper scrap.Scrapper
}
