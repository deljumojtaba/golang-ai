package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/StephaneBunel/bresenham"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// constans
const cellSize = 60 // size of each cell in pixels

// variables

var (
	green     = color.RGBA{0, 255, 0, 255}     // color for start cell
	darkGreen = color.RGBA{0, 155, 0, 255}     // color for goal cell
	red       = color.RGBA{255, 0, 0, 255}     // color for walls
	yellow    = color.RGBA{255, 255, 0, 255}   // color for explored cells
	gray      = color.RGBA{200, 200, 200, 255} // color for empty cells
	orange    = color.RGBA{255, 165, 0, 255}   // color for solution path
	blue      = color.RGBA{0, 0, 255, 255}     // color for frontier cells
)

// output image

func (g *Maze) OutputImage(filename ...string) {
	imgWidth := cellSize * (g.Width - 1)
	imgHeight := cellSize * g.Height

	var outFile = "image.png"

	if len(filename) > 0 {
		outFile = filename[0]
	}

	fmt.Printf("Generating image %s...\n", outFile)

	upLeft := image.Point{}
	lowRight := image.Point{X: imgWidth, Y: imgHeight}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	// draw squares on the image
	for i, row := range g.Walls {
		for j, col := range row {
			p := Point{Row: i, Col: j}
			if col.wall {
				g.drawSquare(col, p, img, color.Black, cellSize, j*cellSize, i*cellSize)
			} else if g.inSolution(p) {
				g.drawSquare(col, p, img, green, cellSize, j*cellSize, i*cellSize)
			} else if col.State.Row == g.Start.Row && col.State.Col == g.Start.Col {
				g.drawSquare(col, p, img, darkGreen, cellSize, j*cellSize, i*cellSize)
			} else if col.State.Row == g.Goal.Row && col.State.Col == g.Goal.Col {
				g.drawSquare(col, p, img, red, cellSize, j*cellSize, i*cellSize)
			} else if col.State == g.CurrentNode.State {
				g.drawSquare(col, p, img, orange, cellSize, j*cellSize, i*cellSize)
			} else if inExplored(Point{i, j}, g.Explored) {
				g.drawSquare(col, p, img, yellow, cellSize, j*cellSize, i*cellSize)
			} else {
				g.drawSquare(col, p, img, color.White, cellSize, j*cellSize, i*cellSize)
			}
		}
	}

	// draw a grid
	for i, _ := range g.Walls {
		bresenham.DrawLine(img, 0, i*cellSize, g.Width*cellSize, i*cellSize, gray)
	}

	for i := 0; i <= g.Width; i++ {
		bresenham.DrawLine(img, i*cellSize, 0, i*cellSize, imgHeight, gray)
	}
	// save to file
	f, _ := os.Create(outFile)
	_ = png.Encode(f, img)

}

// draw Square

func (g *Maze) drawSquare(col Wall, p Point, img *image.RGBA, fill color.Color, size, x, y int) {
	patch := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(patch, patch.Bounds(), &image.Uniform{C: fill}, image.Point{}, draw.Src)

	if !col.wall {
		// use package color.Black (not the parameter) for text color
		g.printLocation(p, color.Black, patch)
	}
	draw.Draw(img, image.Rect(x, y, x+size, y+size), patch, image.Point{}, draw.Src)
}

// print location

func (g *Maze) printLocation(p Point, c color.Color, patch *image.RGBA) {
	point := fixed.Point26_6{X: fixed.I(6), Y: fixed.I(40)}

	d := &font.Drawer{
		Dst:  patch,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(fmt.Sprintf("[%d,%d]", p.Row, p.Col))
}
