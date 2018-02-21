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
	"github.com/disintegration/imaging"
	"golang.org/x/oauth2"
	"net/http"
	"fmt"
	"image/png"
	"path/filepath"
)

var (
	topic  = "sky"
	width  = 1584
	height = 396
)

var (
	utf8FontFile = "/usr/share/fonts/truetype/ubuntu-font-family/Ubuntu-B.ttf"
	utf8FontSize = 55.0
	spacing      = float64(1.5)
	dpi          = float64(72)
	opaque       = color.Alpha{255}
	transparent  = color.Alpha{0}
)

func main() {
	log.Println("Generating banner of dimensions:", width, height)
	maskImage := generateImageMask(width, height, []string{
		"",
		"                       Paul Baker - AWS Certified Developer",
		"                             github.com/paul-nelson-baker",
	})
	backgroundImage := getRandomUnsplashNaive(width, height, topic)
	generateFinalImage(width, height, backgroundImage, maskImage)
}
func generateImageMask(width int, height int, text []string) image.Image {
	// Based off the implementation here: https://socketloop.com/tutorials/golang-print-utf-8-fonts-on-image-example
	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		log.Fatalln(err)
	}
	utf8Font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalln(err)
	}
	fontForegroundColor, fontBackgroundColor := image.NewUniform(opaque), image.NewUniform(transparent)
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
	return imageMask
}

// Naive method of getting an unsplash image.
func getRandomUnsplashNaive(width int, height int, topic string) image.Image {
	randomImageUrl := fmt.Sprintf("https://source.unsplash.com/random/%dx%d?%s", width, height, topic)
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

func generateFinalImage(width int, height int, backgroundImage image.Image, maskImage image.Image) {
	bounds := image.Rect(0, 0, width, height)
	debugPath := filepath.Join(".", "debug")
	os.MkdirAll(debugPath, os.ModePerm)

	saveImgImg(backgroundImage, "./debug/background.png")
	saveImgImg(maskImage, "./debug/mask.png")
	invertedImage := imaging.Invert(backgroundImage)
	saveDrwImg(invertedImage, "./debug/inverted.png")

	//outlineOfMask := image.NewRGBA(bounds)
	//offset := 4
	//draw.Draw(outlineOfMask, image.Rect(bounds.Min.X-offset, bounds.Min.Y-offset, bounds.Max.X-offset, bounds.Max.Y-offset), maskImage, image.ZP, draw.Over)
	//draw.Draw(outlineOfMask, image.Rect(bounds.Min.X-offset, bounds.Min.Y+offset, bounds.Max.X-offset, bounds.Max.Y+offset), maskImage, image.ZP, draw.Over)
	//draw.Draw(outlineOfMask, image.Rect(bounds.Min.X+offset, bounds.Min.Y-offset, bounds.Max.X+offset, bounds.Max.Y-offset), maskImage, image.ZP, draw.Over)
	//draw.Draw(outlineOfMask, image.Rect(bounds.Min.X+offset, bounds.Min.Y+offset, bounds.Max.X+offset, bounds.Max.Y+offset), maskImage, image.ZP, draw.Over)
	//draw.DrawMask(outlineOfMask, bounds, outlineOfMask, image.ZP, maskImage, image.ZP, draw.Src)
	//saveDrwImg(outlineOfMask, "/tmp/debug/mask-outline.png")

	inversionWithMask := image.NewRGBA(bounds)
	draw.DrawMask(inversionWithMask, bounds, invertedImage, image.ZP, maskImage, image.ZP, draw.Src)
	saveDrwImg(inversionWithMask, "./debug/inverted-with-mask.png")

	finalDestination := image.NewRGBA(bounds)
	draw.Draw(finalDestination, bounds, backgroundImage, image.ZP, draw.Over)
	//draw.Draw(finalDestination, bounds, outlineOfMask, image.ZP, draw.Over)
	draw.Draw(finalDestination, bounds, inversionWithMask, image.ZP, draw.Over)
	saveImgImg(finalDestination, fmt.Sprintf("out-final_%dx%d.png", width, height))
}

func saveDrwImg(image draw.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	png.Encode(f, image)
}

func saveImgImg(image image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	png.Encode(f, image)
}
