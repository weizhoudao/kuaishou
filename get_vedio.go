package main

import (
   kuaishou "github.com/RogerLiNing/kuaishou_watermark_remover"
   "fmt"
   "os"
)

func main() {
   for _, arg := range os.Args[1:] {
       url, _ := kuaishou.WatermarkRemover(arg)
       fmt.Println(url.VideoLink)
    }
}
