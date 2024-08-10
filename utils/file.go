package utils

import (
	"io"
	"log"
	"os"

	"github.com/AssassinRobot/Bandcamper/entities"
	"github.com/AssassinRobot/Bandcamper/helpers"
	"github.com/bogem/id3v2"
)

type FileMngmnt struct {
	file *os.File
}

func NewFileMngmnt() *FileMngmnt {
	return &FileMngmnt{}
}

func (f *FileMngmnt) Save(path string, content io.Reader) error {
	out, osCreateError := os.Create(path)
	if osCreateError != nil {
		return osCreateError
	}

	_, ioCopyError := io.Copy(out, content)
	if ioCopyError != nil {
		return ioCopyError
	}

	f.file = out

	defer f.close()

	return nil
}

func (f *FileMngmnt) CreateDir(path string) error {
	osMkdirError := os.MkdirAll(path, 0o700)

	if osMkdirError != nil {
		return osMkdirError
	}

	return nil
}

func (f *FileMngmnt) TagFile(mp3 *entities.TrackData) error {
	tag, mp3OpenError := id3v2.Open(mp3.CurrentTrackFilepath, id3v2.Options{Parse: true, ParseFrames: nil})
	if mp3OpenError != nil {
		return mp3OpenError
	}

	artwork, readFileError := os.ReadFile(mp3.AlbumArtworkFilepath)
	if readFileError != nil {
		return readFileError
	}

	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     artwork,
	}

	tag.AddAttachedPicture(pic)

	artist := helpers.RemoveAlphaNum(mp3.Artist)
	tag.SetArtist(artist)

	title := helpers.RemoveAlphaNum(mp3.CurrentTrackTitle)
	tag.SetTitle(title)

	album := helpers.RemoveAlphaNum(mp3.Current.Title)
	tag.SetAlbum(album)

	if saveTagError := tag.Save(); saveTagError != nil {
		return saveTagError
	}

	if tagCloseError := tag.Close(); tagCloseError != nil {
		return tagCloseError
	}

	return nil
}

func (f *FileMngmnt) close() {
	if f.file != nil {
		err := f.file.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
		return
	}

	log.Fatalln("missing file")
}
