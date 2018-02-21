package main

import (
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"log"
	"os"
	"io/ioutil"
	"image/color"
	"github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"
	"net/http"
	"fmt"
	"image/jpeg"
)

var (
	width  = 646
	height = 220
)

var (
	utf8FontFile = "/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-B.ttf" //"wqy-zenhei.ttf"
	//utf8FontSize     = float64(15.0)
	utf8FontSize = float64(25.0)
	spacing      = float64(1.5)
	dpi          = float64(72)
	red          = color.RGBA{255, 0, 0, 255}
	blue         = color.RGBA{0, 0, 255, 255}
	white        = color.RGBA{255, 255, 255, 255}
	black        = color.RGBA{0, 0, 0, 255}
	transparent  = color.RGBA{0, 0, 0, 0}
	// more color at https://github.com/golang/image/blob/master/colornames/table.go
)

func main() {
	log.Println("Generating banner of dimensions:", width, height)
	maskImage := generateImageMask([]string{
		"Paul Baker - AWS Certified Developer",
		"github.com/paul-nelson-baker",
		"paul.nelson.baker@gmail.com",
	})
	backgroundImage := getRandomUnsplashURL()

	generateFinalImage(backgroundImage, maskImage)
	//log.Println(maskImage, backgroundImage)
	//getRandomUnsplashAPI()
}
func generateImageMask(text []string) image.Image {
	// Based off the implementation here: https://socketloop.com/tutorials/golang-print-utf-8-fonts-on-image-example
	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		log.Fatalln(err)
	}
	utf8Font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalln(err)
	}
	fontForegroundColor, fontBackgroundColor := image.NewUniform(black), image.NewUniform(transparent)
	imageMask := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(imageMask, imageMask.Bounds(), fontBackgroundColor, image.ZP, draw.Src)
	context := freetype.NewContext()
	context.SetDPI(dpi)
	context.SetFont(utf8Font)
	context.SetFontSize(utf8FontSize)
	context.SetClip(imageMask.Bounds())
	context.SetDst(imageMask)
	context.SetSrc(fontForegroundColor)
	pt := freetype.Pt(10, 10+int(context.PointToFixed(utf8FontSize)>>6))
	for _, str := range text {
		_, err := context.DrawString(str, pt)
		if err != nil {
			log.Fatalln(err)
		}
		pt.Y += context.PointToFixed(utf8FontSize * spacing)
	}

	//finalMask := image.NewRGBA(imageMask.Bounds())
	//draw.ApproxBiLinear.Scale(finalMask, imageMask.Bounds(), imageMask, imageMask.Bounds(), draw.Over, nil)

	return imageMask
}
func getRandomUnsplashURL() image.Image {
	randomImageUrl := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	http := http.DefaultClient
	resp, err := http.Get(randomImageUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	backgroundImage, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return backgroundImage
}
func getRandomUnsplashAPI() {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TODO: WE NEED TO GET AN ACCESS TOKEN VIA OAUTH2"},
	)
	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	usplsh := unsplash.New(client)
	photos, _, err := usplsh.Photos.Random(&unsplash.RandomPhotoOpt{
		Height: height,
		Width:  width,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for photo := range *photos {
		log.Println(photo)
	}
}

func generateFinalImage(backgroundImage image.Image, maskImage image.Image) {
	finalDestination := image.NewRGBA(backgroundImage.Bounds())
	draw.Draw(finalDestination, finalDestination.Bounds(), backgroundImage, finalDestination.Bounds().Min, draw.Src)
	draw.DrawMask(finalDestination, finalDestination.Bounds(), maskImage, image.ZP, maskImage, image.ZP, draw.Over)
	output, err := os.Create("/tmp/out-final.jpg")
	defer output.Close()
	if err != nil {
		log.Fatalln(err)
	}
	err = jpeg.Encode(output, finalDestination, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
