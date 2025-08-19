package downloader

import (
	"fmt"
	"log"
	"os"

	// "strconv"
	//
	// "github.com/AssassinRobot/Bandcamper/entities"
	// "github.com/AssassinRobot/Bandcamper/helpers"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type wishlistDownloader struct {
	http      *utils.HttpMngmnt
	file      *utils.FileMngmnt
	downloads []string
	scrapper  scrap.Scrapper
}

// var wg = &sync.WaitGroup{}

func (c *wishlistDownloader) Download(username string) error {
	var errorChan = make(chan error, 500)

	var wishlistUrl = fmt.Sprintf("https://bandcamp.com/%s/wishlist", username)

	// TODO: get cookies from Chrome or Firefox
	var cookies = os.Getenv("BANDCAMP_COOKIES")
	println("Using cookies:", cookies)
	var headers = map[string]string{
		"Cookie":     cookies,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	res, getURLError := c.http.Get(wishlistUrl, headers)

	if getURLError != nil {
		return getURLError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Just print the response body for debugging
	println("Response Status:", res.Status)

	wishlistData, scrapError := c.scrapper.ListWishlist(res.Body)
	if scrapError != nil {
		return scrapError
	}

	// just print wishlist and return
	for _, item := range wishlistData {
		fmt.Printf("Wishlist Item: %s", item.Title)
		fmt.Printf("Album URL: %s\n", item.AlbumURL)
		fmt.Printf("Artwork URL: %s\n", item.ImageURL)
		fmt.Println()
	}
	return nil

	//
	// ticker := helpers.DownloadStatus(&c.downloads)
	//
	// baseFilepath := fmt.Sprintf("./%s%s", helpers.RemoveAlphaNum(trackData.Artist), helpers.RemoveAlphaNum(trackData.Current.Title))
	//
	// trackData.AlbumArtworkFilepath = fmt.Sprintf("%s/%s.jpg", baseFilepath, trackData.Current.Title)
	//
	// createError := c.file.CreateDir(baseFilepath)
	// if createError != nil {
	// 	return createError
	// }
	//
	// imageRes, getImageError := c.http.Get(trackData.ArtworkURL)
	// if getImageError != nil {
	// 	return getImageError
	// }
	//
	// defer func() {
	// 	err := res.Body.Close()
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }()
	//
	// saveImageError := c.file.Save(trackData.AlbumArtworkFilepath, imageRes.Body)
	// if saveImageError != nil {
	// 	return saveImageError
	// }
	//
	// for _, v := range trackData.TrackInfo {
	// 	wg.Add(1)
	//
	// 	currentTrackData := *trackData
	//
	// 	currentTrackData.CurrentTrackNum = strconv.Itoa(v.TrackNum)
	// 	currentTrackData.CurrentTrackTitle = v.Title
	// 	currentTrackData.CurrentTrackURL = v.File.Mp3128
	// 	currentTrackData.CurrentTrackFilepath = baseFilepath +
	// 		"/" + helpers.RemoveAlphaNum(currentTrackData.CurrentTrackNum) +
	// 		"-" + helpers.RemoveAlphaNum(currentTrackData.Artist) +
	// 		"-" + helpers.RemoveAlphaNum(currentTrackData.CurrentTrackTitle) +
	// 		".mp3"
	//
	// 	go func(mp3 entities.TrackData) {
	// 		defer wg.Done()
	//
	// 		c.downloads = append(c.downloads, fmt.Sprintf("%s - %s", mp3.Artist, mp3.CurrentTrackTitle))
	//
	// 		mp3Res, mp3DownloadError := c.http.Get(mp3.CurrentTrackURL)
	// 		if mp3DownloadError != nil {
	// 			errorChan <- mp3DownloadError
	// 			return
	// 		}
	//
	// 		defer func() {
	// 			err := mp3Res.Body.Close()
	// 			if err != nil {
	// 				log.Fatalln(err)
	// 			}
	// 		}()
	//
	// 		saveError := c.file.Save(mp3.CurrentTrackFilepath, mp3Res.Body)
	// 		if saveError != nil {
	// 			errorChan <- saveError
	// 			return
	// 		}
	//
	// 		tagFileError := c.file.TagFile(&mp3)
	// 		if tagFileError != nil {
	// 			errorChan <- tagFileError
	// 			return
	// 		}
	// 	}(currentTrackData)
	// }
	//
	wg.Wait()

	// ticker.Stop()

	close(errorChan)

	return <-errorChan
}

func NewWishlistDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, scrapper scrap.Scrapper) WishlistDownloader {
	return &wishlistDownloader{
		http:     http,
		file:     file,
		scrapper: scrapper,
	}
}
