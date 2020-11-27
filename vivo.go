package vivo

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

//Vivo is the base struct where all information about a vivo.sx video are saved
type Vivo struct {
	VivoURL  string
	VideoURL string
	ID       string
	Title    string
	Mime     string
	Quality  string
	Length   int64
}

// GetVideo extracts the video url and some other nice information from a vivo.sx page
func GetVideo(URL string) (*Vivo, error) {
	if !regexp.MustCompile("(vivo\\.sx/)(.*)(.{10}$)").MatchString(URL) {
		return &Vivo{}, errors.New("Not a valid vivo.sx url")
	}

	if strings.Contains(URL, "/embed/") {
		URL = strings.ReplaceAll(URL, "/embed/", "/")
	}

	response, err := http.Get(URL)
	if err != nil {
		return &Vivo{}, err
	}
	defer response.Body.Close()

	bodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Vivo{}, err
	}
	bodyAsString := string(bodyAsBytes)

	parameter := regexp.MustCompile("(?s)InitializeStream\\s*\\(\\s*({.+?})\\s*\\)\\s*;").FindString(bodyAsString)
	parameter = strings.NewReplacer("\n", "", "\t", "", "InitializeStream ({", "", "});", "", "'", "\"").Replace(strings.TrimSpace(parameter))

	vivo := &Vivo{VivoURL: URL,
		ID:    URL[strings.LastIndex(URL, "/")+1:],
		Title: strings.TrimPrefix(strings.TrimSuffix(regexp.MustCompile(`<h1>(.*?)<strong>`).FindString(bodyAsString), "&nbsp;<strong>"), "<h1>Watch ")}

	for _, info := range strings.Split(parameter, ",") {
		keyValue := strings.Split(info, ": ")
		if len(keyValue) <= 1 {
			continue
		}
		key := keyValue[0]
		value := strings.ReplaceAll(keyValue[1], "\"", "")

		switch key {
		case "quality":
			vivo.Quality = value + "p"
		case "source":
			decodedURL, err := url.QueryUnescape(value)
			if err != nil {
				return &Vivo{}, err
			}
			videoURL := rot47(decodedURL)
			vivo.VideoURL = videoURL

			video, err := http.Get(videoURL)
			if err != nil {
				return vivo, nil
			}
			vivo.Mime = video.Header.Get("content-type")
			vivo.Length = video.ContentLength
		}
	}

	return vivo, nil
}

// Download downloads the video
func (v Vivo) Download(destination io.Writer) error {
	response, err := http.Get(v.VideoURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(destination, response.Body)
	return err
}

// rot47 decrypts a input string with the ROT47 algorithm.
// This is needed because the vivo.sx video url is encrypted in ROT47
func rot47(input string) string {
	var result []string
	for i := range input[:] {
		j := int(input[i])
		if (j >= 33) && (j <= 126) {
			result = append(result, string(rune(33+((j+14)%94))))
		} else {
			result = append(result, string(input[i]))
		}

	}
	return strings.Join(result, "")
}
