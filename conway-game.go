// Conway's Game
//
// John Horton Conway sadly passed away on April, 11, 2020 as another victim of
// COVID-19. Among so many contributions he also conceived the Game of Life
// ---here denoted as Conway's Game
//
// This Go package is a personal tribute to his memory
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/clinaresl/conway-game/conway"
)

// globals
// ----------------------------------------------------------------------------
const EXIT_SUCCESS = 0
const EXIT_FAILURE = 1

const version = "0.1"

// flag parameters
var (
	filename        string
	width, height   int
	xratio, yratio  int
	delay, delay0   int
	population      int
	nbgenerations   int
	model           string
	want_model_help bool
	want_version    bool
)

// functions
// ----------------------------------------------------------------------------

// init module
//
// setup the flag environment for the on-line help
func init() {

	// command line arguments for parsing the name of the gif file
	flag.StringVar(&filename, "filename", "conway.gif", "name of the GIF file")

	// command line arguments for parsing the dimensions of the grid
	flag.IntVar(&width, "width", 100, "Width of the grid")
	flag.IntVar(&height, "height", 100, "Height of the grid")

	// command line arguments for parsing the aspect ratio
	flag.IntVar(&xratio, "xratio", 1, "x aspect ratio")
	flag.IntVar(&yratio, "yratio", 1, "y aspect ratio")

	// command line argument for parsing the delays between frames
	flag.IntVar(&delay0, "delay0", 100, "delay of the first frame")
	flag.IntVar(&delay, "delay", 1, "delay between frames in 100th of a second")

	// command line argument to determine the initial number of alive cells
	flag.IntVar(&population, "population", 100, "initial population")

	// command line argument for getting the desired number of generations
	flag.IntVar(&nbgenerations, "generations", 100, "number of generations")

	// command line argument for parsing the color model
	flag.StringVar(&model, "model", "", "color model. Type --help-model to show additional help")

	// whether additional help on color models was requested
	flag.BoolVar(&want_model_help, "help-model", false, "shows additional information on color models")

	// also, create an additional flag for showing the version
	flag.BoolVar(&want_version, "version", false, "shows version info and exits")
}

// showModelHelp
//
// show additional information on color models
func showModelHelp(signal int) {

	fmt.Println(`
 In all cases colors are given in the format #RRGGBB in hexadecimal format:

   -model "bichrome COLOR[:COLOR]"
		If only one color is given, dead cells are shown in black, and living cells
		with the specified color. If two colors are given, then they are used for
		dead and living cells respectively

   -model "gradient COLOR:COLOR[:COLOR]"
		In case two colors are given, the first is used for dead cells and the living
		cells are used with a gradient of color from black to the given color; if three
		colors are given, then the gradient is computed from the second to the third
		color
`)
	os.Exit(signal)
}

// showVersion
//
// show the current version of this program and exits with the given signal
func showVersion(signal int) {

	fmt.Printf(" %v %v\n", os.Args[0], version)
	os.Exit(signal)
}

// parseHex
//
// return the decimal representation of a number in hexadecimal notation
func parseHex(hexnum string) uint8 {
	var err error
	var result int64
	if result, err = strconv.ParseInt(hexnum, 16, 0); err != nil {
		log.Fatalf("It was not possible to convert the hexadecimal number '%v'", hexnum)
	}
	return uint8(result)
}

// getRGB
//
// return the RGB components of a RGB color as uint8
func getRGB(c color.Color) (r, g, b uint8) {

	// get the rgb components as uint32
	r2, g2, b2, _ := c.RGBA()

	// and transform them now
	return uint8(r2), uint8(g2), uint8(b2)
}

// getColor
//
// return a color from an hexadecimal representation #RRGGBB
func getColor(hexcolor string) color.Color {

	// parse all the hexadecimal components
	re := regexp.MustCompile(`([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})`)
	match := re.FindStringSubmatch(hexcolor)

	// and return a color under the RGB model
	return color.RGBA{parseHex(match[1]), parseHex(match[2]), parseHex(match[3]), 255}
}

// getGradientPalette
//
// return the palette of colors to use for gradient palettes. It receives a
// slice of strings which is the output of the regexp matching the color model
// with the user specification
func getGradientPalette(match []string) (gpalette []color.Color) {

	// extract all colors
	first, second, third := getColor(match[2]), getColor(match[3]), getColor(match[4])

	// get the rgb components of the second and third color
	r2, g2, b2 := getRGB(second)
	r3, g3, b3 := getRGB(third)

	// insert the color used for dead cells as first in the palette
	gpalette = append(gpalette, first)

	// now, create a gradient of colors from the second to the third
	for i := 1.0; i <= 255.0; i++ {

		// and add an intermediate color
		r, g, b := uint8(float64(r2)+(i*(float64(r3)-float64(r2))/255.0)),
			uint8(float64(g2)+(i*(float64(g3)-float64(g2))/255.0)),
			uint8(float64(b2)+(i*(float64(b3)-float64(b2))/255.0))
		gpalette = append(gpalette, color.RGBA{r, g, b, 255})
	}

	// and return the palette of colors
	return
}

// getPalette
//
// return the colour model chosen by the user, the center given (if any, by
// default the point 0,0) and a palette of colours, along with an error if any
// is found
func getPalette(model string) (string, image.Point, []color.Color, error) {

	// set up a regular expression to match the color model specifications
	re := regexp.MustCompile(`\s*(gradient|radial)\s+(\#[a-fA-F0-9]{6}):(\#[a-fA-F0-9]{6}):(\#[a-fA-F0-9]{6})(;(\d+),\s*(\d+))?$`)

	// and match the given color model
	match := re.FindStringSubmatch(model)
	if len(match) == 0 {
		return "", image.Point{},
			[]color.Color{},
			errors.New("Syntax error in the specification of the color model")
	}

	// get the center provided by the user, and if not is given, then use the
	// default values 0, 0
	var xcenter, ycenter int64
	if match[6] != "" {
		xcenter, _ = strconv.ParseInt(match[6], 10, 0)
	}
	if match[7] != "" {
		ycenter, _ = strconv.ParseInt(match[7], 10, 0)
	}

	// and apply the given color model
	switch {

	// gradient color model
	case match[1] == "gradient" || match[1] == "radial":
		return match[1], image.Point{X: int(xcenter), Y: int(ycenter)}, getGradientPalette(match), nil
	}

	// in case the previous switch did not return a palette then an error
	// occurred
	return "", image.Point{}, []color.Color{}, errors.New("Unknown model specification")
}

// main function
//
// given a number decide whether it is divisible by 7 or not
func main() {

	// first things first, parse the flags
	flag.Parse()

	// if additional information has been requested on color models show it and
	// then gracefully exit
	if want_model_help {
		showModelHelp(EXIT_SUCCESS)
	}

	// if the current version is requested, then show it on the standard output
	// and exit
	if want_version {
		showVersion(EXIT_SUCCESS)
	}

	// initialize the first generation randomly
	if population > (1+width)*(1+height) {
		log.Printf(" Pruning the initial population to %v individuals", (1+width)*(1+height))
	}

	contents := make([]bool, (1+width)*(1+height))
	for i := 0; i < population; i++ {
		contents[i] = true
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(width*height, func(i, j int) {
		contents[i], contents[j] = contents[j], contents[i]
	})

	// get a palette according to the user's specification along with the colour
	// model and the center used in the radial model
	var ok error
	var usermodel string
	var center image.Point
	var palette []color.Color
	if usermodel, center, palette, ok = getPalette(model); ok != nil {
		log.Fatalf(" Unknown color model: %v", ok)
	}

	// create the first generation and set its contents
	initial := conway.NewGeneration(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height}},
		palette,
		conway.AspectRatio{X: xratio, Y: yratio},
		usermodel,
		1, nbgenerations)

	// and set the center. Note that the conway package will use it only in case
	// the colour model requested by the user is radial
	initial.SetCenter(center)

	if ok := initial.Set(contents); ok != nil {
		log.Fatalf(" It was not possible to initialize the first generation: %v", ok)
	}

	// Create a Conway's Game with this phase
	game := conway.NewConway(width, height, nbgenerations, initial)

	// and run the Conway's Game over this initial generation
	game.Run()

	// get the image of the entire Conway's game
	anim := game.GetGIF(delay0, delay)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	gif.EncodeAll(f, &anim)
}
