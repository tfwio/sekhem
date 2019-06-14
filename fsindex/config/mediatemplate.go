package config

import (
	"fmt"
	"html/template"
)

func noDown(args ...interface{}) template.HTML {
	bytes := []byte(fmt.Sprintf("%s", args...))
	return template.HTML(bytes)
}

var (
	mediafiletemplate = template.Must(template.New("mediafile").Funcs(template.FuncMap{"noDown": noDown}).Parse(`
{{define "mediafile"}}<html>
	<head>
	<meta charset="utf-8">
	<meta http-equiv="content-language" content="en-US" />
	<style>
	body { font-family: FreeSans; }
	</style>
	</head>
	<body>
	
	Path: {{.Path}}<br/>
	{{if .HasImage}}<img width="256" src="data:{{.ImageMime}};base64,{{.ImageData}}">
	{{else}}[no image]
	{{end}}<br/>

	ImageMime: {{.ImageMime}}<br/>
	Title: <b>{{.Title}}</b><br/>
	Album-Artist: {{.AlbumArtist}}<br/>
	Artist: {{.Artist}}<br/>
	Year: {{.Year}}<br/>
	Album: {{.Album}}<br/>
	Comment: {{.Comment | noDown}}<br/>
	<!--// <hr/> //-->
	<!--// {{.Data | noDown}} //-->
	</body>
</html>
{{end}}`))
)
