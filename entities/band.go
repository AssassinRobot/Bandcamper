package entities

type (
	Band struct {
		ImageURL  string
		Title    string
		Location string
		Bio      string
		Albums    []Album
		Singles []Single
	}
	Album struct {
		Title     string
		AlbumURL string
		ImageURL   string
	}
	Single struct {
		Title     string
		SingleURL string
		ImageURL   string
	}
)
