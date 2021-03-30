package main

import "os"
import "strings"
import "fmt"

func main() {
    var url string
    for _, arg := range os.Args[1:] {
        idx := strings.Index(arg, "http")
        url = arg[idx:]
        fmt.Println("raw=", arg)
        fmt.Println(url)
    }
}
