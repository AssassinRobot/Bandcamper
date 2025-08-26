package cmd

import (
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
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		err = urlDownloader.Download(url, force)
		if err != nil {
			return err
		}

		log.Println("Done")
		return nil
	},
}
