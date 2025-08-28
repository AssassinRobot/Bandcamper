package scrap

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/PuerkitoBio/goquery"
)

type Scrapper interface {
	ListInfos(reader io.Reader) (*entities.TrackData, error)
	ListWishlist(reader io.Reader) ([]*entities.Album, error)
	ListCollection(reader io.Reader) ([]*entities.CollectionItem, error)
	CollectionData(reader io.Reader) (*entities.CollectionData, error)
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
		return nil, fmt.Errorf("error get image:%d", trackData.Current.ID)
	}
	trackData.ArtworkURL = artwork

	return trackData, nil
}

func (s *dataScrapper) ListWishlist(reader io.Reader) ([]*entities.Album, error) {
	doc, newDocError := goquery.NewDocumentFromReader(reader)
	if newDocError != nil {
		return nil, newDocError
	}

	var albums []*entities.Album

	// When user sets the wishlist to private, bandcamp does a rewrite of the page
	// to the user collection, which cannot be made private.
	// So we will try to find the wishlist items in the collection items container.

	var wishlistItems = doc.Find("#wishlist-items-container .collection-items")

	if wishlistItems.Length() == 0 {
		hasCookies := os.Getenv("BANDCAMP_COOKIES") != ""
		var message string

		if hasCookies {
			message = "Could not find wishlist using the provided cookies. Please verify that your username and BANDCAMP_COOKIES value are correct and not expired."
		} else {
			message = "Could not find a public wishlist. If your wishlist is private, set the BANDCAMP_COOKIES environment variable with your cookies."
		}

		return nil, fmt.Errorf("error getting wishlist: %s", message)
	}

	wishlistItems.Each(func(i int, s *goquery.Selection) {
		var album = &entities.Album{}
		album.Title = s.Find(".collection-item-title").Text()
		album.AlbumURL = s.Find(".item-link").AttrOr("href", "")
		album.ImageURL = s.Find(".collection-item-art").AttrOr("src", "")
		albums = append(albums, album)
	})

	return albums, nil
}

func (s *dataScrapper) ListCollection(reader io.Reader) ([]*entities.CollectionItem, error) {
	doc, newDocError := goquery.NewDocumentFromReader(reader)
	if newDocError != nil {
		return nil, newDocError
	}

	// Create debug file
	// html, _ := doc.Html()
	// os.WriteFile("collection.html", []byte(html), 0644)

	var collectionItems []*entities.CollectionItem

	var items = doc.Find("#collection-items .collection-item-container")

	items.Each(func(i int, s *goquery.Selection) {
		var item = &entities.CollectionItem{}

		title := s.Find(".collection-item-title").Text()
		title = strings.Join(strings.Fields(title), " ")
		title = strings.ReplaceAll(title, "(gift given)", "")

		item.Title = title
		item.ArtURL = s.Find(".collection-item-art").AttrOr("src", "")
		item.DownloadURL = s.Find(".redownload-item a").AttrOr("href", "")
		collectionItems = append(collectionItems, item)
	})

	return collectionItems, nil
}

func (s *dataScrapper) CollectionData(reader io.Reader) (*entities.CollectionData, error) {
	doc, newDocError := goquery.NewDocumentFromReader(reader)
	if newDocError != nil {
		return nil, newDocError
	}

	var pageData = &entities.CollectionData{}

	p := doc.Find("#pagedata")
	blob, exists := p.Attr("data-blob")
	if !exists {
		return nil, fmt.Errorf("error getting collection data")
	}

	jsonUnmarshalError := json.Unmarshal([]byte(blob), pageData)
	if jsonUnmarshalError != nil {
		return nil, jsonUnmarshalError
	}

	return pageData, nil
}

func (s *dataScrapper) WishlistData(reader io.Reader) (*entities.WishlistData, error) {
	doc, newDocError := goquery.NewDocumentFromReader(reader)
	if newDocError != nil {
		return nil, newDocError
	}

	var pageData = &entities.WishlistData{}

	p := doc.Find("#pagedata")
	blob, exists := p.Attr("data-blob")
	if !exists {
		return nil, fmt.Errorf("error getting collection data")
	}

	jsonUnmarshalError := json.Unmarshal([]byte(blob), pageData)
	if jsonUnmarshalError != nil {
		return nil, jsonUnmarshalError
	}

	return pageData, nil
}
