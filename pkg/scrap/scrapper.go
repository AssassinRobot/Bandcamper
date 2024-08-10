package scrap

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/PuerkitoBio/goquery"
)

type Scrapper interface {
	ListInfos(reader io.Reader) (*entities.TrackData, error)
}

func NewScrapper() Scrapper {
	return &dataScrapper{}
}

type dataScrapper struct{}

func (s *dataScrapper) ListInfos(reader io.Reader) (*entities.TrackData, error) {
	var listInfoError error

	doc, newDocError := goquery.NewDocumentFromReader(reader)
	if newDocError != nil {
		listInfoError = newDocError
	}

	var trackData = &entities.TrackData{}

	doc.Find("script[data-tralbum]").Each(func(i int, s *goquery.Selection) {
		trackDataString, _ := s.Attr("data-tralbum")

		log.Printf("track data received:%s\n", trackDataString)

		jsonUnmarshalError := json.Unmarshal([]byte(trackDataString), trackData)
		if jsonUnmarshalError != nil {
			listInfoError = jsonUnmarshalError
		}
	})

	if listInfoError != nil {
		return nil, listInfoError
	}

	var artwork string
	switch trackData.ItemType {
	case "track":
		artwork = fmt.Sprintf("https://f4.bcbits.com/img/a%d_10.jpg", trackData.ArtID)
	case "album":
		artwork = fmt.Sprintf("https://f4.bcbits.com/img/a%d_16.jpg", trackData.Current.ArtID)
	default:
		return nil,fmt.Errorf("error get image:%d", trackData.Current.ID)
	}
	trackData.ArtworkURL = artwork

	return trackData, nil
}
