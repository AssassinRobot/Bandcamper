package scrap

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/helpers"
	"github.com/PuerkitoBio/goquery"
)

type BandScrapper interface {
	ListBandInfo(reader io.Reader) (*entities.Band, error)
	Scrapper
}

type bandScrapper struct {
	Scrapper
}

func NewBandScrapper(scrapper Scrapper) BandScrapper {
	return &bandScrapper{
		Scrapper: scrapper,
	}
}

func (s *bandScrapper) ListBandInfo(reader io.Reader) (*entities.Band, error) {
	doc, newDocError := goquery.NewDocumentFromReader(reader)

	if newDocError != nil {
		return nil, newDocError
	}

	var bandData = &entities.Band{}
	var listBandInfoError error

	doc.Find("#band-name-location .title").Each(func(i int, s *goquery.Selection) {
		bandData.Title = s.Text()
	})

	doc.Find("#band-name-location .location").Each(func(i int, s *goquery.Selection) {
		bandData.Location = s.Text()
	})

	doc.Find(".artists-bio-pic .bio-pic a").Each(func(i int, s *goquery.Selection) {
		imageURL, exists := s.Attr("href")
		if !exists {
			listBandInfoError = errors.New("failed getting image link")
			return
		}
		bandData.ImageURL = imageURL
	})

	doc.Find("#bio-text").Each(func(i int, s *goquery.Selection) {
		bandData.Bio = helpers.RemoveSpaces(helpers.Remove(s.Text(), "...\u00a0more"))
	})

	doc.Find("ol#music-grid > li").Each(func(index int, s *goquery.Selection) {
		title := helpers.RemoveSpaces(s.Find("p.title").Text())

		url, exists := s.Find("a").Attr("href")
		if !exists {
			listBandInfoError = errors.New("failed getting album link")
			return
		}

		fullURL := fmt.Sprintf("https://%s.bandcamp.com%s", helpers.Remove(strings.ToLower(bandData.Title)," "), url)

		imageURL, exists := s.Find("img").Attr("src")
		if !exists {
			listBandInfoError = errors.New("failed getting image link")
			return
		}

		if imageURL == "/img/0.gif" {
			imageURL, exists = s.Find("img.lazy").Attr("data-original")
			if !exists {
				listBandInfoError = errors.New("failed getting image link")
				return
			}
		}

		kind := helpers.GetKind(url)

		switch kind {

		case "track":
			single := &entities.Single{}
			single.Title = title
			single.SingleURL = fullURL
			single.ImageURL = imageURL
			bandData.Singles = append(bandData.Singles, *single)
		
		case "album":
			album := &entities.Album{}
			album.Title = title
			album.AlbumURL = fullURL
			album.ImageURL = imageURL
			bandData.Albums = append(bandData.Albums, *album)
		
		default:
			listBandInfoError = errors.New("failed getting data")
			return
		}
	})

	if listBandInfoError != nil {
		return nil, listBandInfoError
	}

	return bandData, nil
}
