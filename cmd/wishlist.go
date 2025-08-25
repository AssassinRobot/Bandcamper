package cmd

import (
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
	Use:   "wishlist [username]",
	Short: "Download wishlist albums",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		err = wishlistDownloader.Download(username, cookies, force)
		if err != nil {
			return err
		}

		log.Println("Done")
		return nil
	},
}
