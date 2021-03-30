package main

import (
    "io"
    "net/http"
    "os"
    "strconv"
    "time"
)

func main() {
        start := time.Now().Unix()
	for i, arg := range os.Args[1:] {
		res, err := http.Get(arg)
			if err != nil {
				panic(err)
			}
		f, err := os.Create(strconv.Itoa(int(start) + i) + ".mp4")
			if err != nil {
				panic(err)
			}
		io.Copy(f, res.Body)	
	}
	
}
