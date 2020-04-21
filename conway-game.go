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

   -model "bichrome COLOR[COLOR]"
		If only one color is given, dead cells are shown in black, and living cells
		with the specified color. If two colors are given, then they are used for
		dead and living cells respectively
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

// getPalette
//
// gets a palette from the colour model provided by the user
func getPalette(model string) ([]color.Color, error) {

	fmt.Printf(" model: %v\n", model)

	// set up a regular expression to match color model specifications
	re := regexp.MustCompile(`\s*(bichrome)\s+(\#[a-fA-F0-9]{6})(:\#[a-fA-F0-9]{6})?$`)

	// and match the given color model
	match := re.FindStringSubmatch(model)
	fmt.Printf("%q\n", match)
	if len(match) != 4 {
		return []color.Color{},
			errors.New("Syntax error in the specification of the color model")
	}

	// if the color model was successfully parsed, then extract the first color
	first := getColor(match[2])

	// if a second color was given then extract it as well
	second := color.RGBA{0, 0, 0, 255}
	if match[3] != "" {

		// then get the second color after removing the heading colon
		second := getColor(match[3][1:])

		// and return a palette with these two colors using the first for dead
		// cells and the second for living cells
		return []color.Color{first, second}, nil
	}

	// otherwise, return a palette using black (which is the second one in case
	// it was not given) as the color for dead cells
	return []color.Color{second, first}, nil
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

	// get a palette according to the user's specification
	var ok error
	var palette []color.Color
	if palette, ok = getPalette(model); ok != nil {
		log.Fatalf(" Unknown color model: %v", ok)
	}

	// create the first generation and set its contents
	initial := conway.NewGeneration(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height}},
		palette,
		conway.AspectRatio{X: xratio, Y: yratio})
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
