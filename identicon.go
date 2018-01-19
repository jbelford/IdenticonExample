package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

const imageX = 1024
const imageY = 1024

const centerY = imageY / 2
const scaleY float64 = imageX / 510.0

const spaceX = imageX / 32
const widthX = (imageX - 17*spaceX) / 16

var stringToHash string
var bgCol color.RGBA

func init() {
	if len(os.Args) < 2 {
		panic("Missing arguments: <string>+")
	}
	stringToHash = strings.Join(os.Args[1:], " ")
	bgCol = color.RGBA{20, 0, 50, 255}
}

func main() {
	// Encoding
	userHash := genHash(stringToHash)
	fmt.Printf("Generated hash: %x\n", userHash)
	identicon := genIdenticon(userHash)
	w, _ := os.Create("out.png")
	defer w.Close()
	png.Encode(w, identicon)
	fmt.Println("Finished encoding!")

	// Decoding
	r, _ := os.Open("out.png")
	defer r.Close()
	img, form, _ := image.Decode(r)
	if form != "png" {
		panic("What the haye")
	}
	decodedHash := decodeIdenticon(img)
	fmt.Printf("Decoded hash: %x\n", decodedHash)
}

func genIdenticon(h []byte) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageX, imageY))
	drawBg(img)
	offsetX := spaceX
	ch := make(chan bool, 100)
	for _, b := range h {
		for i := 0; i < widthX; i++ {
			go drawLines(img, offsetX+i, b, ch)
		}
		offsetX += widthX + spaceX
	}
	for i := 0; i < len(h)*widthX; i++ {
		<-ch
	}
	return img
}

func decodeIdenticon(img image.Image) []byte {
	h := make([]byte, 16)
	offsetX := spaceX + widthX/2
	for i := range h {
		_, g, _, _ := img.At(offsetX, centerY).RGBA()
		h[i] = byte(g)
		offsetX += widthX + spaceX
	}
	return h
}

func drawBg(img *image.RGBA) {
	for x := 0; x < imageX; x++ {
		for y := 0; y < imageY; y++ {
			img.Set(x, y, bgCol)
		}
	}
}

func drawLines(img *image.RGBA, offsetX int, b byte, ch chan bool) {
	col := uint8(b)
	val := int(float64(b) * scaleY)
	for i := 0; i < val; i++ {
		img.Set(offsetX, centerY+i, color.RGBA{0, col, 255, 255})
		img.Set(offsetX, centerY-i, color.RGBA{0, col, 255, 255})
	}
	ch <- true
}

func genHash(a string) []byte {
	h := md5.Sum([]byte(a))
	return h[:]
}
