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
	scrapper := scrap.NewScrapper()

	var emailAddress, emailPassword, emailServer string
	if emailAddressEnv := os.Getenv("EMAIL_ADDRESS"); emailAddressEnv != "" {
		emailAddress = emailAddressEnv
		log.Printf("using EMAIL_ADDRESS env var")
	} else {
		panic("EMAIL_ADDRESS env var is required for collection command")
	}
	if emailPasswordEnv := os.Getenv("EMAIL_PASSWORD"); emailPasswordEnv != "" {
		emailPassword = emailPasswordEnv
		log.Printf("using EMAIL_PASSWORD env var")
	} else {
		panic("EMAIL_PASSWORD env var is required for collection command")
	}
	if emailServerEnv := os.Getenv("EMAIL_SERVER"); emailServerEnv != "" {
		emailServer = emailServerEnv
		log.Printf("using EMAIL_SERVER env var")
	} else {
		panic("EMAIL_SERVER env var is required for collection command")
	}

	urlEmail := utils.NewEmailMngmnt(emailAddress, emailPassword, emailServer)

	collectionDownloader = downloader.NewCollectionDownloader(urlHttp, urlFile, urlEmail, scrapper)
}

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Download collection albums",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var cookies string
		if cookieEnv := os.Getenv("BANDCAMP_COOKIES"); cookieEnv != "" {
			cookies = cookieEnv
			log.Printf("using BANDCAMP_COOKIES env var")
		} else {
			return fmt.Errorf("BANDCAMP_COOKIES env var is required for collection command")
		}

		var username string
		if usernameEnv := os.Getenv("BANDCAMP_USERNAME"); usernameEnv != "" {
			username = usernameEnv
			log.Printf("using BANDCAMP_USERNAME env var")
		} else {
			username = args[0]
		}

		err := collectionDownloader.DownloadAll(username, cookies)
		if err != nil {
			return err
		}

		log.Println("Done")
		return nil
	},
}
