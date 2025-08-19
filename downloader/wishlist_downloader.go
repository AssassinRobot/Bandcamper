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

	var url = fmt.Sprintf("https://bandcamp.com/%s/wishlist", username)
	var cookies = os.Getenv("BANDCAMP_COOKIES")
	var headers = map[string]string{
		"Cookie":     cookies,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	res, getURLError := c.http.Get(url, headers)

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
	var album_urls []string
	for _, item := range wishlistData {
		fmt.Printf("Wishlist Item: %s", item.Title)
		if item.AlbumURL != "" {
			album_urls = append(album_urls, item.AlbumURL)
		}
	}

	urlDownloader := NewURLDownloader(c.http, c.file, c.scrapper)
	err := urlDownloader.DownloadAll(album_urls)
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

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
