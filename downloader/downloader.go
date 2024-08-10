package downloader

import "github.com/AssassinRobot/Bandcamper/entities"

type BandDownloader interface {
	GetBand(name string) (*entities.Band, error)
	GetAlbum(albumURL string) (*entities.TrackData, error)
	GetTrack(trackURL string) (*entities.TrackData, error)
	DownloadAlbum(albumURL string) error
	DownloadTrack(trackURL string) error
}

type URLDownloader interface {
	Download(url string) error
}
