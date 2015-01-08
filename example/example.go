package main

import (
	"flag"
	"fmt"
	"github.com/enobrev/go-mediainfo/mediainfo"
)

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Not enough arguments.")
		return
	}

	/* Load the shared library. */
	mediainfo.Init()

	/* Open and parse the file. */
	info, err := mediainfo.Open(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer info.Close()

	if len(args) > 1 {
		/* Get the info. */
		val, err := info.Get(args[1], 0, mediainfo.Video)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(val)
	} else {
		/* Get all the info. */
		completeInfo, err := info.Info(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(completeInfo)
	}
}
