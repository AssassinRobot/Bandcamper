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
	collectionDownloader downloader.CollectionDownloader
)

func init() {
	urlFile := utils.NewFileMngmnt()
	urlHttp := utils.NewHttpMngmnt()
	urlEmail := utils.NewEmailMngmnt()

	scrapper := scrap.NewScrapper()

	collectionDownloader = downloader.NewCollectionDownloader(urlHttp, urlFile, urlEmail, scrapper)
}

var collectionCmd = &cobra.Command{
	Use:   "collection [username]",
	Short: "Download collection albums",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var cookies string
		if cookieEnv := os.Getenv("BANDCAMP_COOKIES"); cookieEnv != "" {
			cookies = cookieEnv
			log.Printf("using BANDCAMP_COOKIES env var")
		} else {
			cookies = ""
		}

		var username string
		if usernameEnv := os.Getenv("BANDCAMP_USERNAME"); usernameEnv != "" {
			username = usernameEnv
			log.Printf("using BANDCAMP_USERNAME env var")
		} else {
			username = args[0]
		}

		err := collectionDownloader.Download(username, cookies)
		if err != nil {
			fmt.Println(err)
		}

		log.Println("Done")
	},
}
