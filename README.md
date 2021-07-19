# vivo

A very simple tool to get the video url and some other stuff from a [vivo.sx](https://vivo.sx) video.
The alternative domain [vivo.st](https://vivo.st) is also supported.
This package can be used as a **[CLI](#cli-usage)** and as **[library](#library-usage)**.

Only tested on Linux, but should work on Mac and Windows too.

## Installation

### Binary

Get the [latest binaries](https://github.com/ByteDream/vivo/releases/latest) if you just want the command line version.
- [Linux (amd64)](https://github.com/ByteDream/vivo/releases/download/v1.0/vivo-v1.0_linux)
- [Windows (amd64)](https://github.com/ByteDream/vivo/releases/download/v1.0/vivo-v1.0_windows)
- [Mac (amd64)](https://github.com/ByteDream/vivo/releases/download/v1.0/vivo-v1.0_darwin)

Or install and run it directly from source
```
$ git clone https://github.com/bytedream/vivo
$ cd vivo
$ go run ./cmd/vivo
```

##### For help how to use the cli, see [here](#cli-usage)

### Library

If you want to install this package as library, use
```
$ go get github.com/bytedream/vivo
```

##### For a example how to use the libary, see [here](#library-usage).

## CLI usage

For general help, use the `-h` flag when executing the binary.
The cli has multi video support, so you can safely specify multiple urls when executing it. 

#### Download a video
Multiple video download is supported, so you can safely specify more than one url / video to download.
```
$ vivo https://vivo.sx/1234567890
```

#### Get only infos about a video without downloading it
The `-i` flag shows only information about the video
```
$ vivo -i https://vivo.sx/1234567890
```

#### Specify output
With the `-o` flag, a custom output path / file can be specified. By default the file will be downloaded in the current path.
```
$ vivo -o OwO.mp4 https://vivo.sx/1234567890
```

#### Use a proxy
If you want to hide your real location or something other which requires a proxy, you can use the `-q` flag to specify it.
```
$ vivo -p https://0.0.0.0:0000
```

##### Other useful options / flags:
- `-c` - **clean**
  
  *Disable colors and the separator between multiple video downloads.*
- `-s` - **source**
  
  *Shows only the source video url.*
- `-q` - **quiet**
  
  *Disable the complete output.*


## Library usage

Get infos about a video
```go
package main

import (
    "fmt"
    
    "github.com/bytedream/vivo"
)


func main() {
  // this extract all the infos about the video
    vivoVideo, err := vivo.GetVideo("https://vivo.sx/cf2137f496")
    if err != nil {
        panic(err)
    }

  // url of the video
    fmt.Println(vivoVideo.VideoURL)
}
```

## License

This project is licensed under the Mozilla Public License 2.0 (MPL-2.0) - see the [LICENSE](LICENSE) file for more details.
