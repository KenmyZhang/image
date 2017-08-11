package app

import (
	"io"
	"fmt"
	"errors"
	"bytes"
	"os"
	"io/ioutil"
	"image"
	"image/jpeg"
	"path/filepath"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/disintegration/imaging"	
)

const (
	/*
	  EXIF Image Orientations
	  1        2       3      4         5            6           7          8

	  888888  888888      88  88      8888888888  88                  88  8888888888
	  88          88      88  88      88  88      88  88          88  88      88  88
	  8888      8888    8888  8888    88          8888888888  8888888888          88
	  88          88      88  88
	  88          88  888888  888888
	*/
	Upright            = 1
	UprightMirrored    = 2
	UpsideDown         = 3
	UpsideDownMirrored = 4
	RotatedCWMirrored  = 5
	RotatedCCW         = 6
	RotatedCCWMirrored = 7
	RotatedCW          = 8

	MaxImgSize         = 6048 * 4032 // 24 megapixels, roughly 36MB as a raw image
)

func getImageOrientation(input io.Reader) (int, error) {
	if exifData, err := exif.Decode(input); err != nil {
		return Upright, err
	} else {
		if tag, err := exifData.Get("Orientation"); err != nil {
			return Upright, err
		} else {
			orientation, err := tag.Int(0)
			if err != nil {
				return Upright, err
			} else {
				return orientation, nil
			}
		}
	}
}

func makeImageUpright(img image.Image, orientation int) image.Image {
	switch orientation {
	case UprightMirrored:
		return imaging.FlipH(img)
	case UpsideDown:
		return imaging.Rotate180(img)
	case UpsideDownMirrored:
		return imaging.FlipV(img)
	case RotatedCWMirrored:
		return imaging.Transpose(img)
	case RotatedCCW:
		return imaging.Rotate270(img)
	case RotatedCCWMirrored:
		return imaging.Transverse(img)
	case RotatedCW:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func SetScaleImage(body []byte, width, height, option int) (*bytes.Buffer, error) {
	input := bytes.NewReader(body)

    if config, err := jpeg.DecodeConfig(input); err !=nil {
		return nil, errors.New("jpeg.DecodeConfig:" + err.Error())
	} else if config.Width * config.Height > MaxImgSize {
	   fmt.Println("image is too large",)
    }

    input = bytes.NewReader(body)
	orientation, _ := getImageOrientation(input)


    input = bytes.NewReader(body)	
	img, err := jpeg.Decode(input); 
	if err != nil {
		return nil, errors.New("decode err=" + err.Error())
	}

   	img = makeImageUpright(img, orientation)

	img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{option})
	if err != nil {
		return nil, errors.New("jpeg encode err=" + err.Error())
	}

	return buf, nil 
} 

func SaveImage(data []byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0774); err != nil {
		directory, _ := filepath.Abs(filepath.Dir(path))
		return errors.New("directory=" + directory + ", err=" + err.Error())
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return errors.New("path=" + path + ", err=" + err.Error())
	}

	return nil
}