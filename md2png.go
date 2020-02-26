package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

func url2html(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	baseUrl := path.Dir(u.Path)
	if baseUrl != "." && baseUrl != "" {
		u.Path = baseUrl
	}
	base := u.String()

	rv := ""
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		rv += fmt.Sprintf("<head><base href=\"%v\"></head>", base)
		rv += string(contents)
	}
	return rv, nil
}

type ImageRender struct {
	BinaryPath *string
}

func (r *ImageRender) generateImage(html, format, output string, width, quality int) ([]byte, error) {
	c := ImageOptions{
		BinaryPath: *r.BinaryPath,
		Input:      "-",
		Html:       html,
		Format:     format,
		Width:      width,
		Quality:    quality,
		Output:     output,
	}
	out, err := GenerateImage(&c)
	return out, err
}
