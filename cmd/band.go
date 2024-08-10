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
	bandDownloader downloader.BandDownloader
)

func init() {
	bandFile := utils.NewFileMngmnt()
	bandHttp := utils.NewHttpMngmnt()

	generalScrapper := scrap.NewScrapper()
	bandScrapper := scrap.NewBandScrapper(generalScrapper)

	bandDownloader = downloader.NewBandDownloader(bandHttp, bandFile, bandScrapper)

}

var bandCmd = &cobra.Command{
	Use:   "band [band name]",
	Short: "Get Band Information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bandName := args[0]

		band, err := bandDownloader.GetBand(bandName)
		if err != nil {
			log.Println("Error get band:", err)
		} else {
			fmt.Println("Image Link: ", band.ImageURL)
			fmt.Println("Title: ", band.Title)
			fmt.Println("Location: ", band.Location)
			fmt.Println("Bio: ", band.Bio)
			fmt.Println("Albums: ")
			for num, album := range band.Albums {
				num++
				fmt.Println("\tImage Link: ", album.ImageURL)
				fmt.Println("\tNumber: ", num)
				fmt.Println("\tTitle: ", album.Title)
				fmt.Println("\tLink: ", album.AlbumURL)
				fmt.Println("_____________________")
			}
			fmt.Println("\n\nSingles: ")
			for num, Singles := range band.Singles {
				num++
				fmt.Println("\tImage Link: ", Singles.ImageURL)
				fmt.Println("\tNumber: ", num)
				fmt.Println("\tTitle: ", Singles.Title)
				fmt.Println("\tLink: ", Singles.SingleURL)
				fmt.Println("_____________________")
			}
		}
	},
}
