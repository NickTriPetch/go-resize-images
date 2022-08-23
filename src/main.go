package main

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func loadImageFromURL(URL string) (image.Image, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/image", func(c *gin.Context) {
		defaultUrl := "https://static.truck2hand.com/public/upload/item/98260/f2903f43-795c-4444-9ebe-0aaedc7877e1.jpeg"
		url := c.DefaultQuery("url", defaultUrl)
		rawWidth := c.DefaultQuery("w", "240")
		width, _ := strconv.Atoi(rawWidth)
		rawHeight := c.DefaultQuery("h", "0")
		height, _ := strconv.Atoi(rawHeight)

		image, _ := loadImageFromURL(url)

		resizeImage := resize.Resize(0, uint(height), image, resize.Lanczos3)

		newImage, _ := cutter.Crop(resizeImage, cutter.Config{
			Width:  width,
			Height: height,
			Mode:   cutter.Centered,
		})

		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, newImage, nil); err != nil {
			log.Println("unable to encode image.")
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		c.Writer.Write(buffer.Bytes())
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
