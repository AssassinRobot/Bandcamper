package cmd

import (
	"fmt"
	"log"

	"github.com/AssassinRobot/Bandcamper/downloader"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
	"github.com/spf13/cobra"
)

var (
	urlDownloader downloader.URLDownloader
)

func init() {
	urlFile := utils.NewFileMngmnt()
	urlHttp := utils.NewHttpMngmnt()

	scrapper := scrap.NewScrapper()

	urlDownloader = downloader.NewURLDownloader(urlHttp, urlFile, scrapper)
}

var urlCmd = &cobra.Command{
	Use:   "url [url]",
	Short: "Download album/track by url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		err := urlDownloader.Download(url)
		if err != nil{
			fmt.Println(err)
		}

		log.Println("Done")
	},
}
