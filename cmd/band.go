package cmd

import (
	"fmt"
	"os"

	"github.com/AssassinRobot/Bandcamper/downloader"
	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/helpers"
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

	generalDownloader := downloader.NewURLDownloader(bandHttp, bandFile, generalScrapper)
	bandDownloader = downloader.NewBandDownloader(bandHttp, bandFile, bandScrapper, generalDownloader)
}

var bandCmd = &cobra.Command{
	Use:   "band [band name]",
	Short: "Get Band Information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bandName := args[0]

		band, err := bandDownloader.GetBand(bandName)
		if err != nil {
			fmt.Println("Error get band:", err)
			os.Exit(1)
		} else {
			var input string

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

			fmt.Print("\n\n Do you want continue or quit? (c/q):")
			fmt.Scanln(&input)

			switch helpers.ToLower(input) {
			case "c":
				fmt.Print("\nDo you want to download (a album or single) or get info (album/single)? (d/i/q): ")
				fmt.Scanln(&input)

				switch helpers.ToLower(input) {
				case "d":
					fmt.Print("\nDo you want to download a album or single? (a/s/q): ")
					fmt.Scanln(&input)

					switch helpers.ToLower(input) {
					case "a":
						fmt.Print("\nEnter Album number: ")
						fmt.Scanln(&input)

						album := helpers.GetByNumber[entities.Album](input, band.Albums)
						if album == nil {
							fmt.Println("Invalid number")
							os.Exit(1)
						}

						err := bandDownloader.DownloadAlbum(album.AlbumURL)

						if err != nil {
							fmt.Println("Error download album:", err)
							os.Exit(1)
						}

						fmt.Println("Done")
					case "s":
						fmt.Print("\nEnter Track number: ")
						fmt.Scanln(&input)

						single := helpers.GetByNumber[entities.Single](input, band.Singles)
						if single == nil {
							fmt.Println("Invalid number")
							os.Exit(1)
						}

						err := bandDownloader.DownloadTrack(single.SingleURL)

						if err != nil {
							fmt.Println("Error download single track:", err)
							os.Exit(1)
						}

						fmt.Println("Done")
					case "q":
						exit()
					default:
						invalidOption()
					}
				case "i":
					fmt.Print("\nDo you want get a album info or single info? (a/s/q): ")
					fmt.Scanln(&input)

					switch helpers.ToLower(input) {
					case "a":
						fmt.Print("\nEnter Album number: ")
						fmt.Scanln(&input)

						album := helpers.GetByNumber[entities.Album](input, band.Albums)
						if album == nil {
							fmt.Println("Invalid number")
							os.Exit(1)
						}

						albumData, err := bandDownloader.GetAlbum(album.AlbumURL)

						if err != nil {
							fmt.Println("Error get album info:", err)
							os.Exit(1)
						}

						fmt.Println("Title: ", albumData.Current.Title)
						fmt.Println("Artist: ", albumData.Artist)
						fmt.Println("Release Date: ", albumData.Current.ReleaseDate)
						fmt.Println("Tracks: ")
						for _, v := range albumData.TrackInfo {
							fmt.Println("\tNumber: ", v.TrackNum)
							fmt.Println("\tTitle: ", v.Title)
							fmt.Printf("\tDuration: %.2f\n", v.Duration/60)
							fmt.Println("\tHas Lyrics: ", v.HasLyrics)
							fmt.Println("\tAudio url: ", v.File.Mp3128)
							fmt.Println("_____________________")
						}
						fmt.Print("\n\nDo you want download specific track or album? (t/a/q): ")
						fmt.Scanln(&input)

						switch helpers.ToLower(input) {
						case "a":
							err := bandDownloader.DownloadAlbum(album.AlbumURL)
							
							if err != nil {
								fmt.Println("Error download album:", err)
								os.Exit(1)
							}

							fmt.Println("Done")
						case "t":
							fmt.Print("\nEnter Track number: ")
							fmt.Scanln(&input)

							track := helpers.GetByNumber[entities.TrackInfo](input,albumData.TrackInfo)
							if track == nil{
								fmt.Println("Invalid number")
								os.Exit(1)
							}

							err := bandDownloader.DownloadTrack(track.File.Mp3128)
							if err != nil {
								fmt.Println("Error download track:", err)
								os.Exit(1)
							}

							fmt.Println("Done")
						case "q":
							exit()
						default:
							invalidOption()
						}
					case "s":
						fmt.Print("\nEnter Track number: ")
						fmt.Scanln(&input)

						single := helpers.GetByNumber[entities.Single](input, band.Singles)
						if single == nil {
							fmt.Println("Invalid number")
							os.Exit(1)
						}

						data, err := bandDownloader.GetTrack(single.SingleURL)
						if err != nil {
							fmt.Println("Error get single track info:", err)
							os.Exit(1)
						}

						fmt.Println("Title: ", data.Current.Title)
						fmt.Println("Artist: ", data.Artist)
						fmt.Println("About: ",data.About)
						fmt.Println("Release Date: ", data.Current.ReleaseDate)
						for _, v := range data.TrackInfo {
							fmt.Printf("\tDuration: %.2f\n", v.Duration/60)
							fmt.Println("\tAudio url: ", v.File.Mp3128)
							fmt.Println("\tHas Lyrics: ", v.HasLyrics)
							fmt.Println("_____________________")
						}

						fmt.Print("\nDo you want download it? (y/n/q): ")
						fmt.Scanln(&input)

						switch helpers.ToLower(input) {
						case "y":
							err := bandDownloader.DownloadTrack(single.SingleURL)
							if err != nil {
								fmt.Println("Error get single track info:", err)
								os.Exit(1)
							}
							fmt.Println("Done")
						case "n":
							exit()
						case "q":
							exit()
						default:
							invalidOption()
						}
					case "q":
						exit()
					default:
						invalidOption()
					}
				case "q":
					exit()
				default:
					invalidOption()
				}

				fmt.Print("\n\nEnter number of album: ")

			case "q":
				exit()
			default:
				invalidOption()
			}

		}
	},
}

func exit() {
	fmt.Println("goodbye")

	os.Exit(0)
}

func invalidOption() {
	fmt.Println("Invalid Options")

	os.Exit(1)
}
