package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func OverlayImageWithColor(input_filepath string, hex_color string) {
	// Replaces file in-place
	reader, err := os.Open(input_filepath)
	if err != nil {
		log.Fatal(err)
	}

	original, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	original_bounds := original.Bounds()

	println("Original dimensions:")
	println(original_bounds.String())

	destination := image.NewRGBA(image.Rect(0,0, original_bounds.Dx(), original_bounds.Dy()))
	// overlay with a color
	//"#5094D0"

	backfill_color, _ := parseHexColor(hex_color)
	draw.Draw(destination, destination.Bounds(), image.NewUniform(backfill_color), image.ZP, draw.Src)
	draw.Draw(destination, original_bounds, original, original_bounds.Min, draw.Over)

	f, err := os.Create(input_filepath)
	err = png.Encode(f, destination)
	if err != nil {
		log.Fatal(err)
	}
}


func ConvertPossibleJpegToPNG(input_filepath string) {
	// Replaces file in-place
	fmt.Println("ConvertPossibleJpegToPNG 'input_filepath':", input_filepath)
	fmt.Println("ConvertPossibleJpegToPNG : opening file...")
	reader, err := os.Open(input_filepath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ConvertPossibleJpegToPNG : decoding file...")
	original, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	original_bounds := original.Bounds()

	destination := image.NewRGBA(image.Rect(0,0, original_bounds.Dx(), original_bounds.Dy()))
	// overlay with a color
	//"#5094D0"
	draw.Draw(destination, original_bounds, original, original_bounds.Min, draw.Src)

	fmt.Println("ConvertPossibleJpegToPNG : creating new file...")
	f, err := os.Create(input_filepath)
	err = png.Encode(f, destination)
	if err != nil {
		log.Fatal(err)
	}
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}