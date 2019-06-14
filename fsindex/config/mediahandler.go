package config

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"tfw.io/Go/fsindex/util"
)

type MediaInfo struct {
	HasImage    bool
	Data        string
	Path        string
	ImageMime   string
	ImageData   string
	Title       string
	Comment     string
	Artist      string
	AlbumArtist string
	Year        int
	Album       string
}

func TagHandler(conf *Configuration, c *gin.Context) {

	action := c.Param("action")
	route := c.Param("route")
	var (
		path string
		err  error
	)
	var mnfo MediaInfo = MediaInfo{
		ImageMime:   "Unkn",
		Path:        "",
		Title:       "",
		Comment:     "",
		AlbumArtist: "",
		Artist:      "",
		Album:       "",
		Year:        0,
		Data:        "",
	}
	if path, err = conf.GetFilePath(route, action); err == nil {
		media, err2 := util.GetMediaFile(path)
		if err2 == nil {
			mpic := media.Picture()
			mnfo = MediaInfo{
				ImageMime:   "Unkn",
				Path:        path,
				Title:       media.Title(),
				Comment:     strings.ReplaceAll(media.Comment(), "\n", "<br />\n"),
				AlbumArtist: media.AlbumArtist(),
				Artist:      media.Artist(),
				Album:       media.Album(),
				Year:        media.Year(),
				Data:        "",
			}
			if mpic != nil {
				mnfo.HasImage = true
				mnfo.ImageData = base64.StdEncoding.EncodeToString(mpic.Data)
				mnfo.ImageMime = mpic.MIMEType
				mnfo.Data += fmt.Sprintf("picture-data<br/>%s\n<br/>", "yep")
			} else {
				mnfo.Data += fmt.Sprintf("picture-none<br/>%s\n", "nope")
			}
			for k, r := range media.Raw() {
				if k != "\xa9cmt" {
					mnfo.Data += fmt.Sprintf("tag: key = %s, %v<br/>\n", k, r)
				}
			}
		}
	} else {
		mnfo.Data += fmt.Sprintf("error: %s\nRoute: %s, Action: %s<br/>\n", err, route, action)
	}
	mediafiletemplate.ExecuteTemplate(c.Writer, "mediafile", mnfo)

}
