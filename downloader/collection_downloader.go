package downloader

import (
	"fmt"
	"log"
	"net/url"
	"strings"

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
	err := collectionDownloader.DownloadAll(collectionUrls, cookies)
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	close(errorChan)
	return nil
}

func (c *collectionDownloader) downloadItem(downloadURL string, cookies string) error {
	var errorChan = make(chan error, 500)

	var reauthURL = "https://bandcamp.com/api/downloadsreauth/1/reauth"
	var headers = map[string]string{
		"Cookie": cookies,
	}

	// get ?payment_id=xxxx from url
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return err
	}
	queryParams := parsedURL.Query()
	paymentID := queryParams.Get("payment_id")
	if strings.TrimSpace(paymentID) == "" {
		return fmt.Errorf("payment_id not found in download URL")
	}

	var body = map[string]any{
		"payment_id":   paymentID,
		"reauth_email": c.email.Address,
	}
	res, getURLError := c.http.Post(reauthURL, headers, body)

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

	// collectionData, scrapError := c.scrapper.ListCollection(res.Body)
	// if scrapError != nil {
	// 	return scrapError
	// }

	close(errorChan)
	return nil
}

func (c *collectionDownloader) DownloadAll(urls []string, cookies string) error {
	println("Starting download of all collection items...")
	for _, url := range urls {
		println("Downloading from URL:", url)
		// Here you would implement the actual download logic
		// For demonstration, we'll just simulate a download with a print statement
		err := c.downloadItem(url, cookies)
		if err != nil {
			log.Printf("Error downloading %s: %v", url, err)
			continue
		}
		println("Successfully downloaded from URL:", url)
	}
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
