# vivo

A very simple tool to get the video url and some other stuff from a [vivo.sx](https://vivo.sx) video.
The alternative domain [vivo.st](https://vivo.st) is also supported.

Only tested on linux, but should work on Mac and Windows too.

## Install

```bash
go get github.com/bytedream/vivo
```

## Usage

Get infos about a video:
```go
package main

import (
    "fmt"
    "github.com/bytedream/vivo"
)


func main() {
    vivoVideo, err := vivo.GetVideo("https://vivo.sx/1234567890")
    // this extract all the infos about the video
    if err != nil {
        panic(err)
    }

    fmt.Println(vivoVideo.VideoURL) // url of the video
}
```

## License

This project is licensed under the Mozilla Public Licence 2.0 (MPL-2.0) licence - see the [LICENSE](LICENCE) file for more details
