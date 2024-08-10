package entities

type(
	File struct {
		Mp3128 string `json:"mp3-128"`
	}

	TrackInfo struct {
		File File `json:"file"`
		TrackNum int    `json:"track_num"`
		TrackID  int    `json:"track_id"`
		Title    string `json:"title"`
		Duration          float64     `json:"duration"`
		HasLyrics         bool        `json:"has_lyrics"`  
		Lyrics            string      `json:"lyrics"`
	}

	Current struct {
		ReleaseDate string      `json:"release_date"`
		Artist      interface{} `json:"artist"`
		Title       string      `json:"title"`
		ID          int64       `json:"id"`
		ArtID       int         `json:"art_id"`
	} 

	TrackData struct {
		Current Current `json:"current"`
		TrackInfo []TrackInfo `json:"trackinfo"`
		ItemType                   string      `json:"item_type"`     
		Artist               string `json:"artist"`
		AlbumReleaseDate     string `json:"album_release_date"`
		ArtID                int    `json:"art_id"`
		ArtworkURL string
		AlbumArtworkFilepath string
		CurrentTrackTitle    string
		CurrentTrackURL      string
		CurrentTrackFilepath string
		CurrentTrackNum string
	}
)
