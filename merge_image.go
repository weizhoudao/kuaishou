package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"github.com/nfnt/resize"
        "errors"
)

func GetImageObj(filePath string) (img image.Image, err error) {
	f1Src, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f1Src.Close()

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = f1Src.Read(buff)
	if err != nil {
		return nil, err
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(filetype)
	fSrc, err := os.Open(filePath)
	defer fSrc.Close()

	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(fSrc)
		if err != nil {
			fmt.Println("jpeg error")
			return nil, err
		}
	case "image/gif":
		img, err = gif.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}

func MergeImage(worksDir string, file1 string, file2 string, newName string) (string, error) {
	src, err := GetImageObj(worksDir + file1)
	if err != nil {
		return "", err
	}
	srcB := src.Bounds().Max

	src1, err := GetImageObj(worksDir + file2)
	if err != nil {
		return "", err
	}
	src1B := src.Bounds().Max

	newWidth := srcB.X + src1B.X
	newHeight := srcB.Y
	if src1B.Y > newHeight {
		newHeight = src1B.Y
	}

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板

	draw.Draw(des, des.Bounds(), src, src.Bounds().Min, draw.Over)                      //首先将一个图片信息存入jpg
	draw.Draw(des, image.Rect(srcB.X, 0, newWidth, src1B.Y), src1, image.ZP, draw.Over) //将另外一张图片信息存入jpg

	fSave, err := os.Create(worksDir + newName)
	if err != nil {
		return "", err
	}

	defer fSave.Close()

	var opt jpeg.Options
	opt.Quality = 100

	newImage := resize.Resize(1024, 0, des, resize.Lanczos3)

	err = jpeg.Encode(fSave, newImage, &opt) // put quality to 80%
	if err != nil {
		return "", err
	}
	return newName, nil
}

func GetImages(str_images ...string) []image.Image {
    var images []image.Image
    for _, img := range str_images {
        src, err := GetImageObj(img)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        images = append(images, src)
    }
    return images
}

func MergeMultiImage(new_name string, images ...image.Image) error {
    if len(images) < 2 {
        return errors.New("at least two images")
    }
    max_height := 0
    total_width := 0
    for _, img := range images {
        total_width += img.Bounds().Max.X
        if img.Bounds().Max.Y > max_height {
            max_height = img.Bounds().Max.Y
        }
    }
    cur_x := 0
    des := image.NewRGBA(image.Rect(0, 0, total_width, max_height)) // 底板
    for _, img := range images {
        draw.Draw(des, image.Rect(cur_x, 0, total_width, img.Bounds().Max.Y), img, image.ZP, draw.Over)
        cur_x += img.Bounds().Max.X
    }
    fsave, err := os.Create(new_name)
    if err != nil {
        return err
    }
    defer fsave.Close()
    var opt jpeg.Options
    opt.Quality = 100
    
    new_img := resize.Resize(1024, 0, des, resize.Lanczos3)
    err = jpeg.Encode(fsave, new_img, &opt)
    if err != nil {
        return err
    }
    return nil
}

func main() {
    images := GetImages(os.Args[1:]...)
    err := MergeMultiImage("newpic.jpeg", images...)
    if err != nil {
        fmt.Println(err.Error())
    }
}
