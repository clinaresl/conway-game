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
(*Conway's Game*) for the number of generations given in `--generations`:

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

## Examples



# License #

conway-game is free software: you can redistribute it and/or modify it under the
terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

conway-game is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
conway-game. If not, see <http://www.gnu.org/licenses/>.


# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>
