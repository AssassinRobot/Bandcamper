package downloader

import "github.com/AssassinRobot/Bandcamper/entities"

type BandDownloader interface {
	GetBand(name string) (*entities.Band, error)
	GetAlbum(albumURL string) (*entities.TrackData, error)
	GetTrack(trackURL string) (*entities.TrackData, error)
	DownloadAlbum(albumURL string, force bool) error
	DownloadTrack(trackURL string, force bool) error
}

type URLDownloader interface {
	Download(url string, force bool) error
	DownloadAll(urls []string, force bool) error
}

type WishlistDownloader interface {
	Download(username string, cookies string, force bool) error
}
