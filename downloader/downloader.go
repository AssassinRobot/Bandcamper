package downloader

import "github.com/AssassinRobot/Bandcamper/entities"

type BandDownloader interface {
	GetBand(name string) (*entities.Band, error)
	GetAlbum(albumNumber string) (*entities.TrackData, error)
	DownloadAlbum(albumNumber string) error
	DownloadTrack(trackNumber string) error
}

type URLDownloader interface {
	Download(url string) error
}
