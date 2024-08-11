package cmd

import (
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
			helpers.PrintErrorAndExit("Error get band:", err)
		} else {
			helpers.GetBandInfo(band)

			switch helpers.GetScan("\n\n Do you want continue or quit? (c/q): ") {
			case "c":

				switch helpers.GetScan("\nDo you want to download (album/single) or get info (album/single)? (d/i/q): ") {
				case "d":

					switch helpers.GetScan("\nDo you want to download a album or single? (a/s/q): ") {
					case "a":
						albumNumber := helpers.GetScan("\nEnter Album number: ")

						album := helpers.GetByNumber[entities.Album](albumNumber, band.Albums)
						if album == nil {
							helpers.PrintErrorAndExit("Invalid number")
						}

						err := bandDownloader.DownloadAlbum(album.AlbumURL)

						if err != nil {
							helpers.PrintErrorAndExit("Error download album:", err)
						}
					case "s":
						trackNumber := helpers.GetScan("\nEnter Track number: ")

						single := helpers.GetByNumber[entities.Single](trackNumber, band.Singles)
						if single == nil {
							helpers.PrintErrorAndExit("Invalid number")
						}

						err := bandDownloader.DownloadTrack(single.SingleURL)

						if err != nil {
							helpers.PrintErrorAndExit("Error download single track:", err)
						}
					case "q":
						helpers.Exit()
					default:
						helpers.InvalidOption()
					}
				case "i":
					switch helpers.GetScan("\nDo you want get a album info or single info? (a/s/q): ") {
					case "a":
						albumNUmber := helpers.GetScan("\nEnter Album number: ")

						album := helpers.GetByNumber[entities.Album](albumNUmber, band.Albums)
						if album == nil {
							helpers.PrintErrorAndExit("Invalid number")
						}

						albumData, err := bandDownloader.GetAlbum(album.AlbumURL)

						if err != nil {
							helpers.PrintErrorAndExit("Error get album info:", err)
						}

						helpers.GetAlbumInfo(albumData)

						switch helpers.GetScan("\n\nDo you want download specific track or album? (t/a/q): ") {
						case "a":
							err := bandDownloader.DownloadAlbum(album.AlbumURL)

							if err != nil {
								helpers.PrintErrorAndExit("Error download album:", err)
							}
						case "t":
							trackNumber := helpers.GetScan("\nEnter Track number: ")

							track := helpers.GetByNumber[entities.TrackInfo](trackNumber, albumData.TrackInfo)
							if track == nil {
								helpers.PrintErrorAndExit("Invalid number")
							}

							err := bandDownloader.DownloadTrack(helpers.GetSpecificTrackURL(albumData.BasePath, track.TitleLink))
							if err != nil {
								helpers.PrintErrorAndExit("Error download track:", err)
							}
						case "q":
							helpers.Exit()
						default:
							helpers.InvalidOption()
						}
					case "s":
						singleTrackNumber := helpers.GetScan("\nEnter Track number: ")

						single := helpers.GetByNumber[entities.Single](singleTrackNumber, band.Singles)
						if single == nil {
							helpers.PrintErrorAndExit("Invalid number")
						}

						singleTrackData, err := bandDownloader.GetTrack(single.SingleURL)
						if err != nil {
							helpers.PrintErrorAndExit("Error get single track info:", err)
						}

						helpers.GetSingleTrackInfo(singleTrackData)

						switch helpers.GetScan("\nDo you want download it? (y/n/q): ") {
						case "y":
							err := bandDownloader.DownloadTrack(single.SingleURL)
							if err != nil {
								helpers.PrintErrorAndExit("Error get single track info:", err)
							}
						case "n":
							helpers.Exit()
						case "q":
							helpers.Exit()
						default:
							helpers.InvalidOption()
						}
					case "q":
						helpers.Exit()
					default:
						helpers.InvalidOption()
					}
				case "q":
					helpers.Exit()
				default:
					helpers.InvalidOption()
				}
			case "q":
				helpers.Exit()
			default:
				helpers.InvalidOption()
			}

		}
	},
}
