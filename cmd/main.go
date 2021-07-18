package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"vivo"
)

type Color int

var (
	infos  bool
	source bool
	output string
	proxy  string
	clean  bool
	quiet  bool
)

const (
	RED Color = iota + 1
	GREEN
	MAGENTA
	CYAN
)

func colorize(color Color, content string) string {
	var ansiColor, ansiReset string

	if runtime.GOOS != "windows" && !clean {
		switch color {
		case RED:
			ansiColor = "\033[31m"
		case GREEN:
			ansiColor = "\033[32m"
		case MAGENTA:
			ansiColor = "\033[95m"
		case CYAN:
			ansiColor = "\033[96m"
		}
		ansiReset = "\033[0m"
	}

	return fmt.Sprint(ansiColor + content + ansiReset)
}

func main() {
	flag.BoolVar(&infos, "i", false, "Print information about a vivo video without downloading it")
	flag.BoolVar(&source, "s", false, "Print the source path to the video file without downloading it")
	flag.StringVar(&output, "o", ".", "Destination of the file")
	flag.StringVar(&proxy, "p", "", "Proxy to use")
	flag.BoolVar(&clean, "c", false, "Show clean output / disable all additions (colors and separator between multiple downloads)")
	flag.BoolVar(&quiet, "q", false, "Disable the output")

	flag.Parse()

	var o, e *log.Logger

	if quiet {
		o = log.New(ioutil.Discard, "", 0)
		e = log.New(ioutil.Discard, "", 0)
	} else {
		o = log.New(os.Stdout, "", 0)
		e = log.New(os.Stderr, "", 0)
	}

	if flag.NArg() == 0 {
		e.Fatalln(colorize(RED, "At least one url must be specified"))
	}

	var client *http.Client
	if proxy == "" {
		client = http.DefaultClient
	} else {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			e.Fatalln(colorize(RED, err.Error()))
		}
		client = &http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
				Proxy:              http.ProxyURL(proxyURL),
			},
		}
	}

	for i, URL := range flag.Args() {
		video, err := vivo.GetVideoWithClient(URL, client)
		if err != nil {
			e.Println(colorize(RED, err.Error()))
			continue
		}

		videoInformation := fmt.Sprintf("Vivo URL: %s\nSource video URL: %s\nTitle: %s\nQuality: %s\nSize: %s\nMime: %s",
			colorize(MAGENTA, video.VivoURL),
			colorize(MAGENTA, video.VideoURL),
			colorize(MAGENTA, video.Title),
			colorize(MAGENTA, video.Quality),
			colorize(MAGENTA, fmt.Sprintf("%.2fMB", float64(video.Length)/1024/1024)),
			colorize(MAGENTA, video.Mime))

		if infos {
			o.Println(videoInformation)
		}
		if source {
			if infos {
				o.Println()
			}
			o.Println(video.VideoURL)
		}

		if !(infos || source) {
			o.Printf("%s\n\n", videoInformation)

			fileInfo, err := os.Stat(output)
			if fileInfo != nil && fileInfo.IsDir() {
				output = path.Join(output, video.Title)
			}

			file, err := os.Create(output)
			if err != nil {
				if os.IsPermission(err) {
					e.Fatalf(colorize(RED, "Permissions denied: Cannot create file '%s'\n"+
						"You may want to run this program as root again or change the output directory via the `-o` flag\n"), output)
				} else {
					e.Fatalln(colorize(RED, err.Error()))
				}
			}

			o.Printf(colorize(CYAN, "Downloading %s to '%s'..."), video.VivoURL, output)
			if err := video.Download(file); err != nil {
				e.Fatalln(colorize(RED, err.Error()))
			}

			o.Printf(colorize(GREEN, " finished\n"))
		}

		if i != flag.NArg()-1 && !clean {
			o.Println("\n--------------------------------------------------")
		}
	}
}
