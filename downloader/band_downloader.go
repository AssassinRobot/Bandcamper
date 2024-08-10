package downloader

import (
	"fmt"
	"log"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/helpers"
	"github.com/AssassinRobot/Bandcamper/pkg/scrap"
	"github.com/AssassinRobot/Bandcamper/utils"
)

type bandDownloader struct {
	http         *utils.HttpMngmnt
	file         *utils.FileMngmnt
	bandScrapper scrap.BandScrapper
}

func (c *bandDownloader) GetBand(name string) (*entities.Band, error) {
	url := fmt.Sprintf("https://%s.bandcamp.com/music", helpers.GetValidName(name))

	res, getURLError := c.http.Get(url)
	if getURLError != nil {
		return nil, getURLError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	bandData, scrapError := c.bandScrapper.ListBandInfo(res.Body)
	if scrapError != nil {
		return nil, scrapError
	}

	return bandData, nil
}

func (c *bandDownloader) GetAlbum(albumURL string) (*entities.TrackData, error) {
	res, getURLError := c.http.Get(albumURL)
	if getURLError != nil {
		return nil, getURLError
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	data, scrapError := c.bandScrapper.ListInfos(res.Body)
	if scrapError != nil {
		return nil, scrapError
	}

	return data, nil
}
func (c *bandDownloader) GetTrack(trackURL string) (*entities.TrackData, error) {
	panic("uninplemented")
}

func (c *bandDownloader) DownloadAlbum(albumURL string) error {
	panic("uninplemented")
}

func (c *bandDownloader) DownloadTrack(trackURL string) error {
	panic("uninplemented")
}

func NewBandDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, bandScrapper scrap.BandScrapper) BandDownloader {
	return &bandDownloader{
		http:         http,
		file:         file,
		bandScrapper: bandScrapper,
	}
}
