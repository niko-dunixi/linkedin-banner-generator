package main

import (
	"bufio"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const (
	width  = 646
	height = 220
)

func main() {
	log.Println("Generating banner of dimensions:", width, height)
	//image_mask := image.NewRGBA(image.Rect(0, 0, width, height))
	//ioutil.ReadFile()

	generateImageMask()
}
func generateImageMask() {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalln(err)
	}
	foreground, background := image.Black, image.White
	imageMask := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(imageMask, imageMask.Bounds(), background, image.ZP, draw.Src)
	context := freetype.NewContext()
	context.SetDPI(float64(72))
	context.SetFont(font)
	context.SetFontSize(float64(16))
	context.SetClip(imageMask.Bounds())
	context.SetDst(imageMask)
	context.SetSrc(foreground)
	options := &truetype.Options{
		Size: 125.0,
	}
	face := truetype.NewFace(font, options)
	for i, currentChar := range "paul.nelson.baker@gmail.com" {
		currentCharAsString := string(currentChar)
		runeWidth, ok := face.GlyphAdvance(rune(currentChar))
		if ok != true {
			log.Fatalln("Something happened with the rune")
		}
		widthAsInt := int(float64(runeWidth) / 64)
		pt := freetype.Pt(i*250+(125-widthAsInt/2), 128)
		context.DrawString(currentCharAsString, pt)
	}
	outputFile, err := os.Create("/tmp/out.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()
	buffer := bufio.NewWriter(outputFile)
	err = png.Encode(buffer, imageMask)
	if err != nil {
		log.Fatalln(err)
	}
	err = buffer.Flush()
	if err != nil {
		log.Fatalln(err)
	}
}
