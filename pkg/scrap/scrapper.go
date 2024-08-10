package scrap

import (
	"encoding/json"
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

	return trackData, nil
}
