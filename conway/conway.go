// This package provides all means for creating the Conway's Game for arbitrary
// sizes and number of generations
package conway

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
)

// Generation
// ----------------------------------------------------------------------------

// type

// A generation, consists of a specification of those cells that are alive and
// those that are dead. Since this implementation acknowledges colors these are
// represented as indexes to a palette. Additionally, the dimensions of the
// rectangle that circumscribes the generation is stored also. This all comes to
// define generations as paletted images
type generation image.Paletted

// methods

// return a new generation which is initially empty. Since this implementation
// honours colors, the contents are stored as indexes to a color palette
func NewGeneration(rectangle image.Rectangle, palette color.Palette) *generation {

	// compute the number of pixels to use
	nbpixels := (1 + rectangle.Max.Y - rectangle.Min.Y) *
		(1 + rectangle.Max.X - rectangle.Min.X)

	// note that in creation, room is allocated for storing the contents but
	// these are all empty
	return &generation{Pix: make([]uint8, nbpixels),
		Stride:  rectangle.Max.X,
		Rect:    rectangle,
		Palette: palette}
}

// return the number of cells alive around the given position
func (g *generation) nbalive(x, y int) (result int) {

	// cast the generation into an image paletted to access its methods
	img := image.Paletted(*g)

	// if (x, y) is not at the top row
	if y < img.Rect.Max.Y {

		if img.ColorIndexAt(x, y+1) != 0 {
			result += 1
		}

		// if this is not the leftmost column
		if x > 0 {
			if img.ColorIndexAt(x-1, y+1) != 0 {
				result += 1
			}
		}

		// if this is not the rightmost column
		if x < img.Rect.Max.X {
			if img.ColorIndexAt(x+1, y+1) != 0 {
				result += 1
			}
		}
	}

	// if (x, y) is not at the bottom row
	if y > 0 {

		if img.ColorIndexAt(x, y-1) != 0 {
			result += 1
		}

		// if this is not the leftmost column
		if x > 0 {
			if img.ColorIndexAt(x-1, y-1) != 0 {
				result += 1
			}
		}

		// if this is not the rightmost column
		if x < img.Rect.Max.X {
			if img.ColorIndexAt(x+1, y-1) != 0 {
				result += 1
			}
		}
	}

	// if (x,y) is not at the leftmost column
	if x > 0 {
		if img.ColorIndexAt(x-1, y) != 0 {
			result += 1
		}
	}

	// if (x,y) is not at the rightmost column
	if x < img.Rect.Max.X {
		if img.ColorIndexAt(x+1, y) != 0 {
			result += 1
		}
	}

	// and return the number of alive cells around (x, y)
	return
}

// Return the next generation, i.e., apply the rules of the Conway's Game
func (g *generation) Next() *generation {

	// create a new generation with the same dimensions and palette than this
	// one
	next := NewGeneration(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: g.Rect.Max.X, Y: g.Rect.Max.Y}},
		g.Palette)

	// cast both generations into an image paletted to access its methods
	img := image.Paletted(*g)
	nxt := image.Paletted(*next)

	// for all cells in this generation
	for x := 0; x <= g.Rect.Max.X; x++ {
		for y := 0; y <= g.Rect.Max.Y; y++ {

			// get the number of cells alive around cell (x, y)
			alive := g.nbalive(x, y)

			// by default, the next generation is empty, i.e., all of them are
			// dead and thus, the only rules considered are those that make some
			// cells take birth or survive

			// -- survival: Any live cell with two or three live neighbors
			// survives
			if img.ColorIndexAt(x, y) != 0 && (alive == 2 || alive == 3) {
				nxt.SetColorIndex(x, y, 1)
			}

			// -- birth: Any dead cell with three live neighbors becomes a live
			// cell
			if img.ColorIndexAt(x, y) == 0 && alive == 3 {
				nxt.SetColorIndex(x, y, 1)
			}
		}
	}

	// and return the next generation
	return next
}

// Set the contents of a generation to those given in contents. In case the
// given slice and the length of the contents do not match an error is returned
func (g *generation) Set(contents []bool) error {

	if len(contents) != (1+g.Rect.Max.Y-g.Rect.Min.Y)*
		(1+g.Rect.Max.X-g.Rect.Min.X) {
		return errors.New("Mismatched dimensions")
	}

	// cast the generation into an image paletted to access its methods
	img := (*image.Paletted)(g)

	// otherwise, just set the contents of the generation to those given in the
	// slice
	for x := 0; x <= img.Rect.Max.X; x++ {
		for y := 0; y <= img.Rect.Max.Y; y++ {
			if contents[y*(img.Rect.Max.X)+x] {
				img.SetColorIndex(x, y, 1)
			}
		}
	}

	// and return no error
	return nil
}

// Generations are stringers also so that they can be displayed on an output
// stream
func (g generation) String() (output string) {

	// cast the generation into an image paletted to access its methods
	img := image.Paletted(g)

	for irow := 0; irow <= g.Rect.Max.Y; irow++ {
		for icol := 0; icol <= g.Rect.Max.X; icol++ {

			if img.ColorIndexAt(icol, irow) != 0 {
				output += "█"
			} else {
				output += "░"
			}
		}

		// and move to the next line
		output += "\n"
	}

	return output
}

// Conway
// ----------------------------------------------------------------------------

// type

// The Conway's Game consists of a slice with a number of generations each with
// a given width and height
type Conway struct {
	width, height int
	nbgenerations int
	generations   []*generation
}

// methods

// Return a new Conway's Game. Note that it is necessary to specify the first
// generation
func NewConway(width, height, generations int, contents *generation) Conway {

	// when creating a new instance of the Conway's Game, note that space is
	// allocated for all generations, but these are not initialized as a matter
	// of fact
	conway := Conway{
		width:         width,
		height:        height,
		nbgenerations: generations,
		generations:   make([]*generation, generations)}

	// set the initial contents
	conway.generations[0] = contents

	// and return the new instance
	return conway
}

// Run the entire game and generate all generations from the initial population
// in the given instance of the Conway's Game
func (game *Conway) Run() {

	// for all generations but the first one
	for igeneration := 1; igeneration < game.nbgenerations; igeneration++ {

		// compute the generation next to the previous one
		game.generations[igeneration] = game.generations[igeneration-1].Next()
	}
}

// return a gif animation of the Conway's Game with the given delay in 100th of
// a second between frames
func (game *Conway) GetGIF(delay int) gif.GIF {

	// create an array of images and delays between successive frames
	var delays []int = make([]int, game.nbgenerations)
	var images []*image.Paletted = make([]*image.Paletted, game.nbgenerations)

	// transform each generation of the game into a paletted image
	for index, generation := range game.generations {
		delays[index] = delay
		images[index] = (*image.Paletted)(generation)
	}

	// and now return the GIF image
	return gif.GIF{Delay: delays, Image: images}
}
