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
	URLDownloader
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

	data.BasePath = helpers.GetBasePath(albumURL)
	return data, nil
}
func (c *bandDownloader) GetTrack(trackURL string) (*entities.TrackData, error) {
	res, getURLError := c.http.Get(trackURL)
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
	data.BasePath = helpers.GetBasePath(trackURL)
	log.Println(data.BasePath )

	return data, nil
}

func (c *bandDownloader) DownloadAlbum(albumURL string) error {
	err := c.URLDownloader.Download(albumURL)
	if err != nil{
		return err
	}
	return nil
}

func (c *bandDownloader) DownloadTrack(trackURL string) error {
	err := c.URLDownloader.Download(trackURL)
	if err != nil{
		return err
	}
	return nil
}

func NewBandDownloader(http *utils.HttpMngmnt, file *utils.FileMngmnt, bandScrapper scrap.BandScrapper,downloader URLDownloader) BandDownloader {
	return &bandDownloader{
		http:         http,
		file:         file,
		bandScrapper: bandScrapper,
		URLDownloader: downloader,
	}
}
