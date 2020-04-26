# Conway's Game

John Horton Conway sadly passed away on April, 11, 2020. Among so many
contributions he also conceived the Game of Life ---here denoted as Conway's
Game

This Go package is a personal tribute to his memory


## Installation

To install the `conway-game` package execute:

``` sh
$ go get github.com/clinaresl/conway-game
```

and a directory called `github.com/clinaresl/conway-game` will be created under
the `src/` folder of your `$GOPATH`.

To compile the program go to `$GOPATH/src/github.com/clinaresl/conway-game` and
type:

``` sh
$ go build
```

and the executable `conway-game` will be generated. You can test it with:

``` sh
$ ./conway-game --version
```


## Usage

`conway-game` provides various functionalities for generating animated GIF
images which are stored in the file specified with `--filename`. It randomly
locates an arbitrary number of living cells (specified with `--population`) over
a grid of dimensions *width* and *height* (which are specified with the flags
`--width` and `--height` respectively) and applies the rules of the Game of Life
(*Conway's Game*) for the number of **generations** given in `--generations`,
described in the
[Wikipedia](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) as follows:

1. Any live cell with fewer than two live neighbours dies, as if by
   *underpopulation*.
2. Any live cell with two or three live neighbours lives on to the next
   generation.
3. Any live cell with more than three live neighbours dies, as if by
   *overpopulation*.
4. Any dead cell with exactly three live neighbours becomes a live cell, as if
   by *reproduction*.

It also provides the following functionalities:

* It is possible to specify the delay between frames (with `--delay`), and also
  the delay of the first frame (`--delay0`), so that the first one can become
  visible any amount of time.
  
* By default, each cell takes a pixel of the GIF image. It is possible, however,
  to apply any *x*/*y* aspect ratio to the image with `--xratio`/`--yratio`,
  which are not expected to be necessarily the same.

* It acknowledges various *color models* through `--model`. To get a complete
  overview of the different colour models use `--help-model`.

* Finally, to give a sense of *evolution*, it is also feasible to compute the
  average of colors of the same cell over an arbitrary number of consecutive
  generations with `--average`.


## Examples

The most basic of all examples is given below:

``` sh
./conway-game --filename test.gif --generations 100 --width 50 --height 50 
              --population 720 --xratio 10 --yratio 10 
              --model "gradient #000000:#00ff00:#00ff00"
```

which produces an old-fashioned image with big pixels, live cells shown in green
and the background in black:

![Example 0](pics/example-0.gif)

As an example of the possibility of modifying the aspect ratio, the following
command:

```sh
./conway-game --filename test.gif --generations 300 --width 150 --height 150 
              --population 7500 --xratio 8 --yratio 2 
              --model "gradient #000000:#ffff00:#ffff00"
```

modifies the aspect ratio so that cells are four times larger than wider than
tall:

![Example 1](pics/example-1.gif)

# License #

`conway-game` is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

`conway-game` is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
`conway-game`. If not, see <http://www.gnu.org/licenses/>.


# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>
