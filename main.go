package main

import (
	"bufio"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"io/ioutil"
	"image/color"
	"github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"
	"net/http"
	"io"
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
	// more color at https://github.com/golang/image/blob/master/colornames/table.go
)

func main() {
	log.Println("Generating banner of dimensions:", width, height)
	maskImage := generateImageMask([]string{
		"Paul Baker - AWS Certified Developer",
		"Email: paul.nelson.baker@gmail.com",
		"Github: github.com/paul-nelson-baker",
	})
	backgroundImage := getRandomUnsplashURL()

	generateFinalImage(backgroundImage, maskImage)
	//log.Println(maskImage, backgroundImage)
	//getRandomUnsplashAPI()
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
func getRandomUnsplashURL() image.Image {
	randomImageUrl := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	http := http.DefaultClient
	resp, err := http.Get(randomImageUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	outfile, err := os.Create("/tmp/out-background.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	backgroundImage, _, err := image.Decode(outfile)
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

// Based off the imple}mentation here: https://socketloop.com/tutorials/golang-print-utf-8-fonts-on-image-example
func generateImageMask(text []string) image.Image {
	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		log.Fatalln(err)
	}
	utf8Font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalln(err)
	}
	fontForegroundColor, fontBackgroundColor := image.NewUniform(black), image.NewUniform(white)
	imageMask := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(imageMask, imageMask.Bounds(), fontBackgroundColor, image.ZP, draw.Src)
	context := freetype.NewContext()
	context.SetDPI(dpi)
	context.SetFont(utf8Font)
	context.SetFontSize(utf8FontSize)
	context.SetClip(imageMask.Bounds())
	context.SetDst(imageMask)
	context.SetSrc(fontForegroundColor)
	//var text = []string{
	//	"paul.nelson.baker@gmail.com",
	//	"github.com/paul-nelson-baker",
	//}
	pt := freetype.Pt(10, 10+int(context.PointToFixed(utf8FontSize)>>6))
	for _, str := range text {
		_, err := context.DrawString(str, pt)
		if err != nil {
			log.Fatalln(err)
		}
		pt.Y += context.PointToFixed(utf8FontSize * spacing)
	}
	outFile, err := os.Create("/tmp/out-mask.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outFile.Close()
	buff := bufio.NewWriter(outFile)
	err = png.Encode(buff, imageMask)
	if err != nil {
		log.Fatalln(err)
	}
	err = buff.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	return imageMask
}

//func generateImageMask() {
//	font, err := truetype.Parse(goregular.TTF)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	foreground, background := image.Black, image.White
//	imageMask := image.NewRGBA(image.Rect(0, 0, width, height))
//	draw.Draw(imageMask, imageMask.Bounds(), background, image.ZP, draw.Src)
//	context := freetype.NewContext()
//	context.SetDPI(float64(72))
//	context.SetFont(font)
//	context.SetFontSize(float64(16))
//	context.SetClip(imageMask.Bounds())
//	context.SetDst(imageMask)
//	context.SetSrc(foreground)
//	options := &truetype.Options{
//		Size: 125.0,
//	}
//	face := truetype.NewFace(font, options)
//	for i, currentChar := range "📧: paul.nelson.baker@gmail.com" {
//		currentCharAsString := string(currentChar)
//		runeWidth, ok := face.GlyphAdvance(rune(currentChar))
//		if ok != true {
//			log.Fatalln("Something happened with the rune")
//		}
//		widthAsInt := int(float64(runeWidth) / 64)
//		pt := freetype.Pt(i*250+(125-widthAsInt/2), 128)
//		context.DrawString(currentCharAsString, pt)
//	}
//	outputFile, err := os.Create("/tmp/out.png")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer outputFile.Close()
//	buffer := bufio.NewWriter(outputFile)
//	err = png.Encode(buffer, imageMask)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	err = buffer.Flush()
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
