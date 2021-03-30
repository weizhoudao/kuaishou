package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "bytes"
    "io"
    "os"
    "strings"
    "bufio"
)

type FeedPhoto struct {
    Id string `json:"id"`
    Duration int `json:"duration"`
    CoverUrl string `json:"coverUrl"`
    PhotoUrl string `json:"photoUrl"`
}
type FeedMsg struct {
    Type int `json:"type"`
    Photo FeedPhoto `json:"photo"`
    CurrentPcursor string `json:"currentPcursor"`
}

type PhotoList struct {
    Result int `json:"result"`
    Llsid string `json:"llsid"`
    HostName string `json:"hostName"`
    Pcursor string `json:"pcursor"`
    TypeName string `json:"__typename"`
    Feeds []FeedMsg `json:"feeds"`
}
type ProfileData struct {
    VisionProfilePhotoList PhotoList `json:"visionProfilePhotoList"`
}
type ProfileMsg struct {
    Data ProfileData `json:"data"`
}

var cookie = ""
var url = "https://video.kuaishou.com/graphql"
var query = "query visionProfilePhotoList($pcursor: String, $userId: String, $page: String) {\n  visionProfilePhotoList(pcursor: $pcursor, userId: $userId, page: $page) {\n    result\n    llsid\n    feeds {\n      type\n      author {\n        id\n        name\n        following\n        headerUrl\n        headerUrls {\n          cdn\n          url\n          __typename\n        }\n        __typename\n      }\n      tags {\n        type\n        name\n        __typename\n      }\n      photo {\n        id\n        duration\n        caption\n        likeCount\n        realLikeCount\n        coverUrl\n        coverUrls {\n          cdn\n          url\n          __typename\n        }\n        photoUrls {\n          cdn\n          url\n          __typename\n        }\n        photoUrl\n        liked\n        timestamp\n        expTag\n        __typename\n      }\n      canAddComment\n      currentPcursor\n      llsid\n      status\n      __typename\n    }\n    hostName\n    pcursor\n    __typename\n  }\n}\n"
var userid = ""

func InitCookieAndUserid(){
    fi, err := os.Open("/Users/jianghuiqiang/Desktop/workspace/download.conf")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    defer fi.Close()

    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
        }
        line := string(a)
        idx := strings.Index(line, "=")
        if idx == -1 {
            break
        }
        if line[0:idx] == "userid" {
            userid = line[idx + 1:]
        }
        if line[0:idx] == "cookie" {
            cookie = line[idx + 1:]
        }
    }
}

func DoPost(userid, pcursor string) string {
    info := make(map[string]interface{})
    info["operationName"] = "visionProfilePhotoList"
    info["query"] = query 
    subinfo := make(map[string]string)
    subinfo["page"] = "profile"
    subinfo["pcursor"] = pcursor 
    subinfo["userId"] = userid
    info["variables"] = subinfo
    data, err := json.Marshal(info)
    if err != nil {
        return err.Error()
    }
    reader := bytes.NewReader(data)
    req, err := http.NewRequest("POST", url, reader)
    if err != nil {
        return err.Error()
    }
    req.Header.Set("cookie", cookie)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
    req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("content-type", "application/json")
    req.Header.Set("Host", "video.kuaishou.com")
    req.Header.Set("Origin", "https://video.kuaishou.com")
    req.Header.Set("Referer", "https://video.kuaishou.com/profile/" + userid)

    client := http.Client{}
    rsp, err := client.Do(req)
    if err != nil {
        return err.Error()
    }
    rspData, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        return err.Error()
    }
    return string(rspData)
}

func main() {
    InitCookieAndUserid()
    cursor := ""
    for cursor != "no_more" {
        rspData := DoPost(userid, cursor)
        msg := ProfileMsg{}
	err := json.Unmarshal([]byte(rspData), &msg)
	if err != nil {
	    fmt.Println(err.Error())
	    return
	}
	cursor = msg.Data.VisionProfilePhotoList.Pcursor
	for _, val := range msg.Data.VisionProfilePhotoList.Feeds {
	    fmt.Println(val.Photo.PhotoUrl)
	}
    }

}
