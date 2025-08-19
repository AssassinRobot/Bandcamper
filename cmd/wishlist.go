package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AssassinRobot/Bandcamper/downloader"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
	"github.com/spf13/cobra"
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

		var cookies string
		if cookieEnv := os.Getenv("BANDCAMP_COOKIES"); cookieEnv != "" {
			cookies = cookieEnv
		} else {
			fmt.Println("Please set the BANDCAMP_COOKIES environment variable with your Bandcamp cookies.")
			return
		}

		err := wishlistDownloader.Download(username, cookies)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Done")
	},
}
