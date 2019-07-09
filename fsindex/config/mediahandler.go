package config

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/tfwio/sekhm/util"
)

// MediaInfo struct is just what it says.
type MediaInfo struct {
	HasImage    bool   `json:"hasimage,omitempty"`
	Data        string `json:"dbg,omitempty"`
	Path        string `json:"-"`
	ImageMime   string `json:"mime"`
	ImageData   string `json:"pic"`
	Title       string `json:"title"`
	Comment     string `json:"cmmt"`
	Artist      string `json:"artist"`
	AlbumArtist string `json:"albumartist"`
	Year        int    `json:"year"`
	Album       string `json:"album"`
}

func getTagData(conf *Configuration, c *gin.Context) MediaInfo {
	action := c.Param("action")
	route := c.Param("route")
	var (
		path string
		err  error
	)
	var mnfo = MediaInfo{
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
				ImageMime:   "unknown",
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
	return mnfo
}

// TagHandler wraps JSON data into an in memory html.Template.
func TagHandler(conf *Configuration, c *gin.Context) {
	mnfo := getTagData(conf, c)
	mediafiletemplate.ExecuteTemplate(c.Writer, "mediafile", mnfo)
}

// TagHandlerJSON serves JSON data of course.
func TagHandlerJSON(conf *Configuration, c *gin.Context) {
	mnfo := getTagData(conf, c)
	mnfo.Path = ""
	if mnfo.HasImage {
		mnfo.ImageData = fmt.Sprintf("data:%s;base64,%s", mnfo.ImageMime, mnfo.ImageData)
	}
	c.JSON(http.StatusOK, mnfo)
}
