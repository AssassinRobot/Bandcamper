package downloader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type collectionDownloader struct {
	http     *utils.HttpMngmnt
	file     *utils.FileMngmnt
	email    *utils.EmailMngmnt
	scrapper scrap.Scrapper
}

func (c *collectionDownloader) getCollectionData(username string, cookies string) (*entities.CollectionData, error) {
	var url = fmt.Sprintf("https://bandcamp.com/%s", username)
	var headers = map[string]string{
		"Cookie": cookies,
	}
	res, getURLError := c.http.Get(url, headers)
	if getURLError != nil {
		return nil, getURLError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	collectionData, scrapError := c.scrapper.CollectionData(res.Body)
	if scrapError != nil {
		return nil, scrapError
	}
	return collectionData, nil
}

func (c *collectionDownloader) getCollectionItems(username string, cookies string) ([]*entities.CollectionItem, error) {
	collectionData, err := c.getCollectionData(username, cookies)
	if err != nil {
		return nil, err
	}

	url := "https://bandcamp.com/api/fancollection/1/collection_items"
	headers := map[string]string{
		"Cookie":       cookies,
		"Content-Type": "application/json",
		"Referer":      fmt.Sprintf("https://bandcamp.com/%s", username),
	}

	var body = map[string]any{
		"fan_id":           collectionData.FanData.FanId,
		"older_than_token": collectionData.LastToken,
		"count":            100,
	}

	res, err := c.http.Post(url, headers, body)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// parse json response
	var data *entities.CollectionItemsData
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	var items []*entities.CollectionItem
	for _, item := range data.Items {
		redownloadURL := data.RedownloadURLs[item.ID]
		if redownloadURL == "" {
			continue
		}
		collectionItem := &entities.CollectionItem{
			Title:       item.Title,
			ArtURL:      item.ArtURL,
			DownloadURL: redownloadURL,
		}
		items = append(items, collectionItem)
	}
	return items, nil
}

func (c *collectionDownloader) Download(username string, cookies string) error {
	var errorChan = make(chan error, 500)

	collectionData, err := c.getCollectionItems(username, cookies)
	if err != nil {
		return err
	}

	// just print collection and return
	var collectionUrls []string
	for _, item := range collectionData {
		fmt.Printf("Collection Item: %s\n", item.Title)
		if item.DownloadURL != "" {
			collectionUrls = append(collectionUrls, item.DownloadURL)
		}
	}
	fmt.Printf("Total collection items with download URLs: %d\n", len(collectionUrls))

	if len(collectionUrls) == 0 {
		fmt.Println("No downloadable items found in the collection.")
		close(errorChan)
		return nil
	}

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

	inbox, err := c.email.ReadInbox()
	if err != nil {
		return err
	}
	if len(inbox) == 0 {
		return fmt.Errorf("no email received for reauth")
	}
	// just print email subjects
	for _, email := range inbox {
		fmt.Printf("Email Subject: %s\n", email)
	}
	// panic("not implemented")

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
