package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type wishlistDownloader struct {
	http     *utils.HttpMngmnt
	file     *utils.FileMngmnt
	scrapper scrap.Scrapper
}

type errorResponse struct {
	Error bool `json:"error"`
}

func (c *wishlistDownloader) Download(username string, cookies string) error {
	var errorChan = make(chan error, 500)

	var url = fmt.Sprintf("https://bandcamp.com/%s/wishlist", username)
	var headers = map[string]string{
		"Cookie": cookies,
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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Check if body is not { error: true }
	var apiErr errorResponse
	if err := json.Unmarshal(bodyBytes, &apiErr); err == nil && apiErr.Error {
		bodyString := string(bodyBytes)
		return fmt.Errorf("Error response from server: %s", bodyString)
	}

	wishlistData, scrapError := c.scrapper.ListWishlist(res.Body)
	if scrapError != nil {
		return scrapError
	}

	// just print wishlist and return
	var albumUrls []string
	for _, item := range wishlistData {
		fmt.Printf("Wishlist Item: %s", item.Title)
		if item.AlbumURL != "" {
			albumUrls = append(albumUrls, item.AlbumURL)
		}
	}

	urlDownloader := NewURLDownloader(c.http, c.file, c.scrapper)
	err = urlDownloader.DownloadAll(albumUrls)
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
