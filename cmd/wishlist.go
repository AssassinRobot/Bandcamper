package cmd

import (
	"fmt"
	"github.com/AssassinRobot/Bandcamper/downloader"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
	"github.com/spf13/cobra"
	"log"
)

var (
	wishlistDownloader downloader.WishlistDownloader
)

func init() {
	urlFile := utils.NewFileMngmnt()
	urlHttp := utils.NewHttpMngmnt()

	scrapper := scrap.NewScrapper()

	wishlistDownloader = downloader.NewWishlistDownloader(urlHttp, urlFile, scrapper)
}

var wishlistCmd = &cobra.Command{
	Use:   "wishlist",
	Short: "Download wishlist albums",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		err := wishlistDownloader.Download(username)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Done")
	},
}
