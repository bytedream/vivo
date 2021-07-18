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

	client *http.Client
}

// GetVideo extracts the video url and some other nice information from a vivo.sx page
func GetVideo(URL string) (*Vivo, error) {
	return GetVideoWithClient(URL, http.DefaultClient)
}

// GetVideoWithClient extracts the video url and some other nice information from a vivo.sx page with a pre defined proxy
func GetVideoWithClient(URL string, client *http.Client) (*Vivo, error) {
	var scheme string

	URL = strings.TrimSuffix(URL, "/")

	urlPattern := regexp.MustCompile(`(?m)^(?P<scheme>https?://)?vivo\.(sx|st)/(embed/)?.{10}`)
	urlGroups := urlPattern.SubexpNames()

	if urlMatch := urlPattern.FindAllStringSubmatch(URL, -1); len(urlMatch) == 0 {
		return &Vivo{}, errors.New("Not a valid url")
	} else {
		for _, match := range urlMatch {
			for i, content := range match {
				if urlGroups[i] == "scheme" {
					scheme = content
					break
				}
			}
		}
	}
	if scheme == "" {
		URL = "https://" + URL
	}

	if strings.Contains(URL, "/embed/") {
		URL = strings.ReplaceAll(URL, "/embed/", "/")
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return &Vivo{}, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0")

	response, err := client.Do(req)
	if err != nil {
		return &Vivo{}, err
	}
	defer response.Body.Close()

	bodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Vivo{}, err
	}

	streamPattern := regexp.MustCompile(`(?m)<h1>Watch (?P<title>[\S\s]+)&nbsp|quality:\s*(?P<quality>\S+),|source:\s*'(?P<source>\S+)'`)
	streamGroups := streamPattern.SubexpNames()

	vivo := &Vivo{
		client:  client,
		VivoURL: URL,
		ID:      URL[strings.LastIndex(URL, "/")+1:],
	}

	for _, match := range streamPattern.FindAllSubmatch(bodyAsBytes, -1) {
		for i, content := range match {
			contentAsString := string(content)
			if contentAsString != "" {
				switch streamGroups[i] {
				case "title":
					vivo.Title = contentAsString
				case "quality":
					vivo.Quality = contentAsString + "p"
				case "source":
					decodedURL, err := url.QueryUnescape(contentAsString)
					if err != nil {
						return &Vivo{}, err
					}
					videoURL := rot47(decodedURL)
					vivo.VideoURL = videoURL

					video, err := client.Get(videoURL)
					if err != nil {
						return &Vivo{}, err
					}
					vivo.Mime = video.Header.Get("content-type")
					vivo.Length = video.ContentLength
				}
			}
		}
	}

	if vivo.VideoURL == "" {
		return &Vivo{}, errors.New("Could not find video")
	}

	return vivo, nil
}

// Download downloads the video
func (v Vivo) Download(destination io.Writer) error {
	response, err := v.client.Get(v.VideoURL)
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
