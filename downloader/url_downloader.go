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

func (c *urlDownloader) Download(url string) error {
	var errorChan = make(chan error, 500)

	res, getURLError := c.http.Get(url)
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

	var artwork string
	switch trackData.ItemType {
	case "track":
		artwork = fmt.Sprintf("https://f4.bcbits.com/img/a%d_10.jpg", trackData.ArtID)
	case "album":
		artwork = fmt.Sprintf("https://f4.bcbits.com/img/a%d_16.jpg", trackData.Current.ArtID)
	default:
		return fmt.Errorf("error get image:%d", trackData.Current.ID)
	}

	trackData.AlbumArtworkFilepath = fmt.Sprintf("%s/%s.jpg", baseFilepath, trackData.Current.Title)

	createError := c.file.CreateDir(baseFilepath)
	if createError != nil {
		return createError
	}

	imageRes, getImageError := c.http.Get(artwork)
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

		go func(mp3 entities.TrackData) {
			defer wg.Done()

			c.downloads = append(c.downloads, fmt.Sprintf("%s - %s", mp3.Artist, mp3.CurrentTrackTitle))

			mp3Res, mp3DownloadError := c.http.Get(mp3.CurrentTrackURL)
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

func NewURLDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, scrapper scrap.Scrapper) URLDownloader {
	return &urlDownloader{
		http:     http,
		file:     file,
		scrapper: scrapper,
	}
}
