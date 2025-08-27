package downloader

import (
	"fmt"
	"log"

	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type collectionDownloader struct {
	http     *utils.HttpMngmnt
	file     *utils.FileMngmnt
	email    *utils.EmailMngmnt
	scrapper scrap.Scrapper
}

func (c *collectionDownloader) Download(username string, cookies string) error {
	var errorChan = make(chan error, 500)

	var url = fmt.Sprintf("https://bandcamp.com/%s", username)
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

	collectionData, scrapError := c.scrapper.ListCollection(res.Body)
	if scrapError != nil {
		return scrapError
	}

	// just print collection and return
	var collectionUrls []string
	for _, item := range collectionData {
		fmt.Printf("Collection Item: %s\n", item.Title)
		if item.DownloadURL != "" {
			collectionUrls = append(collectionUrls, item.DownloadURL)
		}
	}

	collectionDownloader := NewCollectionDownloader(c.http, c.file, c.email, c.scrapper)
	err := collectionDownloader.DownloadAll(collectionUrls)
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	close(errorChan)
	return nil
}

func (c *collectionDownloader) DownloadAll(urls []string) error {
	println("Starting download of all collection items...")
	return nil
}

func NewCollectionDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, email *utils.EmailMngmnt, scrapper scrap.Scrapper) CollectionDownloader {
	return &collectionDownloader{
		http:     http,
		file:     file,
		email:    email,
		scrapper: scrapper,
	}
}
