package helpers

import (
	"fmt"
	"os"

	"github.com/AssassinRobot/Bandcamper/entities"
)

func GetBandInfo(band *entities.Band) {
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

func GetAlbumInfo(albumData *entities.TrackData) {

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
}

func GetSingleTrackInfo(singleTrack *entities.TrackData) {
	fmt.Println("Title: ", singleTrack.Current.Title)
	fmt.Println("Artist: ", singleTrack.Artist)
	fmt.Println("About: ", singleTrack.About)
	fmt.Println("Release Date: ", singleTrack.Current.ReleaseDate)
	for _, v := range singleTrack.TrackInfo {
		fmt.Printf("\tDuration: %.2f\n", v.Duration/60)
		fmt.Println("\tAudio url: ", v.File.Mp3128)
		fmt.Println("\tHas Lyrics: ", v.HasLyrics)
		fmt.Println("_____________________")
	}
}

func GetScan(scanText string) string {
	var input string

	fmt.Print(scanText)

	fmt.Scanln(&input)

	return ToLower(input)
}

func PrintErrorAndExit(a ...any) {
	fmt.Println(a...)
	os.Exit(1)
}

func Exit() {
	fmt.Println("goodbye")

	os.Exit(0)
}

func InvalidOption() {
	fmt.Println("Invalid Option")

	os.Exit(1)
}
