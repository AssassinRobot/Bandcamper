package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
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

func (c *collectionDownloader) DebugEmail() error {
	inbox, err := c.email.ReadInbox()
	if err != nil {
		return err
	}
	if len(inbox) == 0 {
		return fmt.Errorf("no email found")
	}
	// just print email subjects
	fmt.Println()
	for _, email := range inbox {
		fmt.Printf("Email From: %s\n", email.From)
		fmt.Printf("%s\n\n", email.Subject)
	}
	return nil
}

func (c *collectionDownloader) getCollectionData(username string, cookies string) (*entities.CollectionPage, error) {
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

	collectionPageData, scrapError := c.scrapper.CollectionPage(res.Body)
	if scrapError != nil {
		return nil, scrapError
	}
	return collectionPageData, nil
}

func (c *collectionDownloader) getCollectionItems(username string, cookies string, email string) ([]*entities.CollectionItem, error) {
	collectionPageData, err := c.getCollectionData(username, cookies)
	if err != nil {
		return nil, err
	}

	url := "https://bandcamp.com/api/fancollection/1/collection_items"
	headers := map[string]string{
		"Cookie":       cookies,
		"Content-Type": "application/json",
		"Referer":      fmt.Sprintf("https://bandcamp.com/%s", username),
	}

	fanId := collectionPageData.FanData.FanId
	lastToken := collectionPageData.CollectionData.LastToken

	if fanId == 0 || strings.TrimSpace(lastToken) == "" {
		return nil, fmt.Errorf("invalid fan ID or last token")
	}

	var body = map[string]any{
		"fan_id":           fanId,
		"older_than_token": lastToken,
		"count":            100,
	}

	res, err := c.http.Post(url, headers, body)
	if err != nil {
		fmt.Printf("Error making POST request: %v\n", err)
		return nil, err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}
	// fmt.Println(string(bodyBytes))

	// parse json response
	var data *entities.CollectionItemsData
	err = json.Unmarshal(bodyBytes, &data)

	if err != nil {
		fmt.Printf("Error decoding JSON response: %v\n", err)
		return nil, err
	}

	if (len(data.Items) == 0) || (len(data.RedownloadURLs) == 0) {
		return nil, fmt.Errorf("no collection items found")
	}

	var items []*entities.CollectionItem
	for _, item := range data.Items {
		index := "p" + strconv.Itoa(item.SaleID)
		redownloadURL := data.RedownloadURLs[index]
		collectionItem := &entities.CollectionItem{
			Title:       item.Title,
			ArtURL:      item.ArtURL,
			DownloadURL: redownloadURL,
		}
		items = append(items, collectionItem)
	}
	return items, nil
}

func (c *collectionDownloader) DownloadAll(username string, cookies string, email string) error {
	var errorChan = make(chan error, 500)

	collectionData, err := c.getCollectionItems(username, cookies, email)
	if err != nil {
		return err
	}
	fmt.Printf("Total collection items fetched: %d\n", len(collectionData))

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

	err = c.downloadAll(collectionUrls, cookies, email)
	if err != nil {
		return err
	}

	close(errorChan)
	return nil
}

func (c *collectionDownloader) Download(downloadURL string, cookies string, email string) error {
	var errorChan = make(chan error, 500)

	// downloadURL example: https://bandcamp.com/download?payment_id=xxxx&sitem_id=yyyyy
	// check that it looks like that with a regex
	if !strings.Contains(downloadURL, "payment_id=") {
		return fmt.Errorf("invalid download URL: %s", downloadURL)
	}
	if !strings.Contains(downloadURL, "sitem_id=") {
		return fmt.Errorf("invalid download URL: %s", downloadURL)
	}

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

	var downloadLinks []string
	for _, email := range inbox {
		subject := strings.TrimSpace(email.Subject)
		subject = strings.ToLower(subject)
		subjectCheck := strings.Contains(subject, "download link")
		subjectCheck = subjectCheck || strings.Contains(subject, "bandcamp")

		if !subjectCheck {
			continue
		}

		bodyLines := strings.Split(email.Body, "\n")
		for _, line := range bodyLines {
			line = strings.TrimSpace(line)
			lineCheck := strings.HasPrefix(line, "http://")
			lineCheck = lineCheck || strings.HasPrefix(line, "https://")
			lineCheck = lineCheck && strings.Contains(line, "bandcamp.com/download/")
			lineCheck = lineCheck && strings.Contains(line, "payment_id=")
			lineCheck = lineCheck && strings.Contains(line, "sitem_id=")

			if lineCheck {
				downloadLinks = append(downloadLinks, line)
			}
		}
	}

	fmt.Printf("Total download links found in email: %d\n", len(downloadLinks))
	for _, link := range downloadLinks {
		fmt.Printf("Download Link: %s\n", link)
	}

	close(errorChan)
	return nil
}

func (c *collectionDownloader) downloadAll(urls []string, cookies string, email string) error {
	println("Starting download of all collection items...")
	for _, url := range urls {
		println("Downloading from URL:", url)
		// Here you would implement the actual download logic
		// For demonstration, we'll just simulate a download with a print statement
		err := c.Download(url, cookies, email)
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
