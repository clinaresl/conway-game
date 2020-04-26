// This package provides all means for creating the Conway's Game for arbitrary
// sizes and number of generations
package conway

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
	"math"
)

// Functions
// ----------------------------------------------------------------------------

// FarestPoint
//
// Return the farest corner from a given point which has to be circumscribed in
// a rectangle. In case the point does not belong to the rectangle, it returns
// an error
func farestPoint(p image.Point, r image.Rectangle) (f image.Point, err error) {

	// if the point does not belong to the rectangle, immediately exit with an
	// error
	// if !p.In(r) {
	// 	return f, errors.New("The point is outside the rectangle")
	// }

	// Now, if the point falls in the upper half of the rectangle
	if p.Y >= (r.Max.Y-r.Min.Y)/2 {

		// Now, if it falls in the right half
		if p.X >= (r.Max.X-r.Min.X)/2 {

			// then return the lower-left corner
			return image.Point{r.Min.X, r.Min.Y}, nil
		} else {

			// then return the lower-right corner
			return image.Point{r.Max.X, r.Min.Y}, nil
		}
	}

	// At this point, it is known to fall in the lower half of the rectangle

	// If it falls in the right half
	if p.X >= (r.Max.X-r.Min.X)/2 {

		// then return the upper-left corner
		return image.Point{r.Min.X, r.Max.Y}, nil
	}

	// We know that the point falls in the lower-left quarter, so then directly
	// return the upper-right corner
	return image.Point{r.Max.X, r.Max.Y}, nil
}

// EuclideanDistance
//
// Return the euclidean distance between two points
func EuclideanDistance(p1, p2 image.Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) +
		math.Pow(float64(p1.Y-p2.Y), 2))
}

// Generation
// ----------------------------------------------------------------------------

// type

// Frames can be generated with arbitrary aspect ratios, whose definition is
// given below. They represent a magnification of the underlying matrix by 1:X
// and 1:Y
type AspectRatio struct {
	X, Y int
}

// Generation
// ----------------------------------------------------------------------------

// type

// A generation, consists of a specification of those cells that are alive and
// those that are dead over a bidimensional matrix which is subjected to an
// aspect ratio. Each generation has an index running in the range [1,
// nbgenerations]
//
// Since this implementation acknowledges colors these are represented as
// indexes to a palette. Additionally, the dimensions of the rectangle that
// circumscribes the generation is stored also. This all comes to define
// generations as paletted images along with information of the aspect ratio
//
// Pixels are coloured according to the given color model:
//
//    * Gradient: all living cells are coloured with a different color in each
//    generation
//
//    * Radial: living cells are coloured with an RGB combination according to
//    its distance to the farest corner from a center point
//
// In all cases, dead cells are coloured always with the same RGB combination
//
// Because the radial color model computes distances from a corner, this is
// stored in each generation as well
type generation struct {
	img                         image.Paletted
	ratio                       AspectRatio
	model                       string
	nbgeneration, nbgenerations int
	center                      image.Point
}

// methods

// return a new generation which is initially empty. Since this implementation
// honours colors, the contents are stored as indexes to a color palette and a
// colour model has to be given. The new generation is given index nbgeneration
// among nbgenerations
func NewGeneration(rectangle image.Rectangle,
	palette color.Palette,
	ratio AspectRatio,
	model string,
	nbgeneration, nbgenerations int) *generation {

	// compute the number of pixels to use
	nbpixels := (1 + rectangle.Max.Y) * ratio.Y *
		(1 + rectangle.Max.X) * ratio.X

	// note that in creation, room is allocated for storing the contents but
	// these are all empty. The true size of the underling image is:
	//
	// (1 + rectangle.Max.X) * ratio.X wide
	//
	// (1 + rectangle.Max.Y) * ratio.Y tall
	//
	// where the minimum values of the rectangle are just ignored
	return &generation{
		img: image.Paletted{
			Pix:    make([]uint8, nbpixels),
			Stride: 1 + (rectangle.Max.X-rectangle.Min.X)*ratio.X,
			Rect: image.Rectangle{
				Min: image.Point{
					X: rectangle.Min.X * ratio.X,
					Y: rectangle.Min.Y * ratio.Y},
				Max: image.Point{
					X: rectangle.Max.X * ratio.X,
					Y: rectangle.Max.Y * ratio.Y}},
			Palette: palette},
		ratio:         ratio,
		model:         model,
		nbgeneration:  nbgeneration,
		nbgenerations: nbgenerations}
}

// Get the color index at location (x, y) taking into account the aspect ratio
// of this generation
func (g *generation) ColorIndexAt(x, y int) uint8 {

	// note that all pixels in the underlying image corresponding to the same
	// cell are assumed to be coloured with the same index of the palette!
	return g.img.ColorIndexAt(x*g.ratio.X, y*g.ratio.Y)
}

// Set the color index at location (x, y) taking into account the aspect ratio
// of this generation
func (g *generation) SetColorIndex(x, y int, c uint8) {

	for xoffset := 0; xoffset < g.ratio.X; xoffset++ {
		for yoffset := 0; yoffset < g.ratio.Y; yoffset++ {
			g.img.SetColorIndex(x*g.ratio.X+xoffset, y*g.ratio.Y+yoffset, c)
		}
	}
}

// Set the center to use for deciding the colour to use
func (g *generation) SetCenter(p image.Point) {
	g.center = p
}

// return the number of cells alive around the given position
func (g *generation) nbalive(x, y int) (result int) {

	// transform the given coordinates according to the aspect ratio
	xt, yt := x*g.ratio.X, y*g.ratio.Y

	// if (x, y) is not at the top row
	if yt < g.img.Rect.Max.Y {

		if g.ColorIndexAt(x, y+1) != 0 {
			result += 1
		}

		// if this is not the leftmost column
		if xt > 0 {
			if g.ColorIndexAt(x-1, y+1) != 0 {
				result += 1
			}
		}

		// if this is not the rightmost column
		if xt < g.img.Rect.Max.X {
			if g.ColorIndexAt(x+1, y+1) != 0 {
				result += 1
			}
		}
	}

	// if (x, y) is not at the bottom row
	if yt > 0 {

		if g.ColorIndexAt(x, y-1) != 0 {
			result += 1
		}

		// if this is not the leftmost column
		if xt > 0 {
			if g.ColorIndexAt(x-1, y-1) != 0 {
				result += 1
			}
		}

		// if this is not the rightmost column
		if xt < g.img.Rect.Max.X {
			if g.ColorIndexAt(x+1, y-1) != 0 {
				result += 1
			}
		}
	}

	// if (x,y) is not at the leftmost column
	if xt > 0 {
		if g.ColorIndexAt(x-1, y) != 0 {
			result += 1
		}
	}

	// if (x,y) is not at the rightmost column
	if xt < g.img.Rect.Max.X {
		if g.ColorIndexAt(x+1, y) != 0 {
			result += 1
		}
	}

	// and return the number of alive cells around (x, y)
	return
}

// Return the next generation, i.e., apply the rules of the Conway's Game
func (g *generation) Next() *generation {

	// color of living cells
	var c uint8

	// create a new generation with the same dimensions and palette than this
	// one following also the same colour model
	next := NewGeneration(image.Rectangle{
		Min: image.Point{X: g.img.Rect.Min.X / g.ratio.X, Y: g.img.Rect.Min.Y / g.ratio.Y},
		Max: image.Point{X: g.img.Rect.Max.X / g.ratio.X, Y: g.img.Rect.Max.Y / g.ratio.Y}},
		g.img.Palette,
		AspectRatio{X: g.ratio.X, Y: g.ratio.Y},
		g.model,
		1+g.nbgeneration,
		g.nbgenerations)

	// and also reuse the same center
	next.SetCenter(g.center)

	// compute the color to use for the living cells in this generation in case
	// this generation uses the gradient color model
	if g.model == "gradient" {
		c = uint8(g.nbgeneration * 255.0 / g.nbgenerations)
	}

	// for all cells in this generation
	for x := 0; x <= g.img.Rect.Max.X/g.ratio.X; x++ {
		for y := 0; y <= g.img.Rect.Max.Y/g.ratio.Y; y++ {

			// get the number of cells alive around cell (x, y)
			alive := g.nbalive(x, y)

			// compute the color of this cell in case this generation follows
			// the radial color model, and make sure that the maximum index is
			// used
			if g.model == "radial" {

				// get the farest corner from the center used in this
				// generation, and also the distance from this cell to the same
				// corner
				farest, _ := farestPoint(g.center,
					image.Rectangle{Min: image.Point{
						X: g.img.Rect.Min.X / g.ratio.X,
						Y: g.img.Rect.Min.Y / g.ratio.Y},
						Max: image.Point{
							X: g.img.Rect.Max.X / g.ratio.X,
							Y: g.img.Rect.Max.Y / g.ratio.Y}})
				c = uint8(255.0 * EuclideanDistance(g.center, image.Point{X: x, Y: y}) /
					EuclideanDistance(g.center, farest))
			}

			// by default, the next generation is empty, i.e., all of them are
			// dead and thus, the only rules considered are those that make some
			// cells take birth or survive

			// -- survival: Any live cell with two or three live neighbors
			// survives
			if g.ColorIndexAt(x, y) != 0 && (alive == 2 || alive == 3) {
				next.SetColorIndex(x, y, c)
			}

			// -- birth: Any dead cell with three live neighbors becomes a live
			// cell
			if g.ColorIndexAt(x, y) == 0 && alive == 3 {
				next.SetColorIndex(x, y, c)
			}
		}
	}

	// and return the next generation
	return next
}

// Set the contents of a generation to those given in contents. In case the
// given slice and the length of the contents do not match an error is returned
func (g *generation) Set(contents []bool) error {

	// color of living cells
	var c uint8

	if len(contents) != (1+g.img.Rect.Max.Y/g.ratio.Y-g.img.Rect.Min.Y/g.ratio.Y)*
		(1+g.img.Rect.Max.X/g.ratio.X-g.img.Rect.Min.X/g.ratio.X) {
		return errors.New("Mismatched dimensions")
	}

	// compute the color to use for the living cells in this generation in case
	// this generation uses the gradient color model
	if g.model == "gradient" {
		c = uint8(g.nbgeneration * 255.0 / g.nbgenerations)
	}

	// otherwise, just set the contents of the generation to those given in the
	// slice
	for x := 0; x <= g.img.Rect.Max.X/g.ratio.X; x++ {
		for y := 0; y <= g.img.Rect.Max.Y/g.ratio.Y; y++ {
			if contents[y*(g.img.Rect.Max.X/g.ratio.X)+x] {

				// compute the color of this cell in case this generation
				// follows the radial color model, and make sure that the
				// maximum index is used
				if g.model == "radial" {

					// get the farest corner from the center used in this
					// generation, and also the distance from this cell to the
					// same corner
					farest, _ := farestPoint(g.center,
						image.Rectangle{Min: image.Point{
							X: g.img.Rect.Min.X / g.ratio.X,
							Y: g.img.Rect.Min.Y / g.ratio.Y},
							Max: image.Point{
								X: g.img.Rect.Max.X / g.ratio.X,
								Y: g.img.Rect.Max.Y / g.ratio.Y}})
					c = uint8(255.0 * EuclideanDistance(g.center, image.Point{X: x, Y: y}) /
						EuclideanDistance(g.center, farest))
				}
				g.SetColorIndex(x, y, c)
			}
		}
	}

	// and return no error
	return nil
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
// a second between frames, and an initial delay equal to delay0 100th of a
// second
func (game *Conway) GetGIF(delay0, delay int) gif.GIF {

	// create an array of images and delays between successive frames
	var delays []int = make([]int, game.nbgenerations)
	var images []*image.Paletted = make([]*image.Paletted, game.nbgenerations)

	// transform each generation of the game into a paletted image
	for index, generation := range game.generations {
		if index == 0 {
			delays[index] = delay0
		} else {
			delays[index] = delay
		}
		images[index] = (*image.Paletted)(&generation.img)
	}

	// and now return the GIF image
	return gif.GIF{Delay: delays, Image: images}
}
