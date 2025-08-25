package downloader

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/helpers"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type urlDownloader struct {
	http      *utils.HttpMngmnt
	file      *utils.FileMngmnt
	downloads []string
	scrapper  scrap.Scrapper
}

var wg = &sync.WaitGroup{}

func (c *urlDownloader) Download(url string, force bool) error {
	var errorChan = make(chan error, 500)

	res, getURLError := c.http.Get(url, nil)
	if getURLError != nil {
		return getURLError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	trackData, scrapError := c.scrapper.ListInfos(res.Body)
	if scrapError != nil {
		return scrapError
	}

	ticker := helpers.DownloadStatus(&c.downloads)

	baseFilepath := fmt.Sprintf("./%s%s", helpers.RemoveAlphaNum(trackData.Artist), helpers.RemoveAlphaNum(trackData.Current.Title))

	trackData.AlbumArtworkFilepath = fmt.Sprintf("%s/%s.jpg", baseFilepath, trackData.Current.Title)

	createError := c.file.CreateDir(baseFilepath)
	if createError != nil {
		return createError
	}

	imageRes, getImageError := c.http.Get(trackData.ArtworkURL, nil)
	if getImageError != nil {
		return getImageError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	saveImageError := c.file.Save(trackData.AlbumArtworkFilepath, imageRes.Body)
	if saveImageError != nil {
		return saveImageError
	}

	for _, v := range trackData.TrackInfo {
		wg.Add(1)

		currentTrackData := *trackData

		currentTrackData.CurrentTrackNum = strconv.Itoa(v.TrackNum)
		currentTrackData.CurrentTrackTitle = v.Title
		currentTrackData.CurrentTrackURL = v.File.Mp3128
		currentTrackData.CurrentTrackFilepath = baseFilepath +
			"/" + helpers.RemoveAlphaNum(currentTrackData.CurrentTrackNum) +
			"-" + helpers.RemoveAlphaNum(currentTrackData.Artist) +
			"-" + helpers.RemoveAlphaNum(currentTrackData.CurrentTrackTitle) +
			".mp3"

		if !force {
			exists, statError := c.file.Exists(currentTrackData.CurrentTrackFilepath)
			if statError != nil {
				return statError
			}
			if exists {
				println("Skipping " + currentTrackData.CurrentTrackFilepath + " (file exists)")
				wg.Done()
				continue
			}
		}

		go func(mp3 entities.TrackData) {
			defer wg.Done()

			c.downloads = append(c.downloads, fmt.Sprintf("%s - %s", mp3.Artist, mp3.CurrentTrackTitle))

			mp3Res, mp3DownloadError := c.http.Get(mp3.CurrentTrackURL, nil)
			if mp3DownloadError != nil {
				errorChan <- mp3DownloadError
				return
			}

			defer func() {
				err := mp3Res.Body.Close()
				if err != nil {
					log.Fatalln(err)
				}
			}()

			saveError := c.file.Save(mp3.CurrentTrackFilepath, mp3Res.Body)
			if saveError != nil {
				errorChan <- saveError
				return
			}

			tagFileError := c.file.TagFile(&mp3)
			if tagFileError != nil {
				errorChan <- tagFileError
				return
			}
		}(currentTrackData)
	}

	wg.Wait()

	ticker.Stop()

	close(errorChan)

	return <-errorChan
}

func (c *urlDownloader) DownloadAll(urls []string) error {
	var wg = &sync.WaitGroup{}
	var errorChan = make(chan error, len(urls))

	// Add the number of tasks (URLs) to wait on
	wg.Add(len(urls))

	// Loop over URLs and create a goroutine for each URL
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			// Call the Download function for each URL
			err := c.Download(url)
			if err != nil {
				errorChan <- err
			}
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the error channel and return the first error encountered
	close(errorChan)

	// Return the first error from the channel (if any)
	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func NewURLDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, scrapper scrap.Scrapper) URLDownloader {
	return &urlDownloader{
		http:     http,
		file:     file,
		scrapper: scrapper,
	}
}
