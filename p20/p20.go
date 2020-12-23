package p20

import (
	"errors"
	"fmt"
	"github.com/fritzr/advent2020/util"
	"gonum.org/v1/gonum/mat" // yeah it might be overkill, but I want to learn
	"math"
	"strconv"
	"strings"
)

type Tile struct {
	*mat.Dense
	id int
}

func validateDimensions(tileStr string, ncols int) (nrows int, err error) {
	// Each tile should be square, and we should have a square number of them.
	return
}

func BytesToFloat(data []byte) []float64 {
	float := make([]float64, len(data))
	for index, c := range data {
		float[index] = float64(c)
	}
	return float
}

func NewTile(id int, nrows int, ncols int, data []byte) Tile {
	return Tile{mat.NewDense(nrows, ncols, BytesToFloat(data)), id}
}

// parseTiles creates tiles from string representations.
//
// Each tile string describes a square matrix by line-delimited rows.
func parseTiles(tileStrings []string) (tiles map[int]Tile, err error) {
	tiles = make(map[int]Tile, len(tileStrings))
	for _, tileStr := range tileStrings {
		// First line is "Tile <id>:\n"
		idStart := strings.IndexRune(tileStr, ' ') + 1
		firstLine := strings.IndexRune(tileStr, '\n')
		var id int
		id, err = strconv.Atoi(tileStr[idStart : firstLine-1])
		if err != nil {
			break
		}
		ncols := strings.IndexRune(tileStr[firstLine+1:], '\n')
		tileData := strings.Replace(tileStr[firstLine+1:], "\n", "", -1)
		nrows := len(tileData) / ncols
		if ncols != nrows || len(tileData) != nrows*ncols {
			err = errors.New(fmt.Sprintf(
				"tile has invalid dimensions %d x %d (%d elements)",
				nrows, ncols, len(tileData)))
			break
		}
		tiles[id] = NewTile(id, nrows, ncols, []byte(tileData))
	}
	return
}

// Directions.
const (
	NONE = iota
	RIGHT
	UP
	LEFT
	DOWN
)

const HOW = 0
const NUM_DIRECTIONS = 1 + DOWN

// Reconstitute organizes tiles into a square grid of width `sqrt(len(tiles))`
// such that adjacent tiles are equal on their borders.
//
// Tiles are square and may be reflected vertically or horizontally and rotated
// in intervals of 90 degrees to obtain a valid arrangement.
//
// For example, the following is a valid arrangement of four 4x4 tiles:
//
//   #.#.  .###
//   #...  .#..
//   #..#  #.##
//   ..##  #.#.
//
//   ..##  #.#.
//   #...  .##.
//   .#.#  ##.#
//   ..#.  ..#.
//
// Note the adjacent vertical and horizontal edges:
//
//   #.#|.  .|###
//   #..|.  .|#..
//   #..|#  #|.##
//   ---+----+---
//   ..#|#  #|.#.
//
//   ..#|#  #|.#.
//   ---+----+---
//   #..|.  .|##.
//   .#.|#  #|#.#
//   ..#|.  .|.#.
//
func Reconstitute(tiles map[int]Tile) (order []int, err error) {
	// Map tile IDs to their adjacent tiles (by ID).
	// The indexes of the value are the constants RIGHT, LEFT, UP, DOWN, etc...
	// A zero value (we assume no tile has ID zero) means unknown adjancency.
	// A value of -1 means not adjacent to any tile.
	type adjacencyList = [NUM_DIRECTIONS]int
	type adjacencyMap = map[int]adjacencyList
	adjacent := make(adjacencyMap, len(tiles))

	// We continually visit a random tile, checking all other tiles.
	//
	// Each tile is itself square, and thus can be rotated four times
	// and flipped in two dimensions, providing twelve possible arrangements.
	// Naively, every arrangement of each tile must be checked with every other
	// arrangement of every other tile.
	//
	// As a first attempt, we'll try a greedy algorithm which accepts any adjacent
	// pairs which are found immediately. We may need to keep a list of possible
	// adjacencies and filter them.
Tile1:
	for id1, tile1 := range tiles {
		for direction, adjacentId := range adjacent[id1] {
			if adjacentId == -1 {
				for _, arrangement1 := range Arrangements(tile1) {
					for id2, tile2 := range tiles {
						for _, arrangement2 := range Arrangements(tile2) {
							if which := arrangement1.Borders(arrangement2); which != NONE {
								adjacent[id1][HOW] = arrangement1
								adjacent[id1][which] = id2
								break Done
							}
						}
					}
				}
			}
		}
	}
Done:

	var _ = adjacent // XXX suppress unused variables for build
	return
}

func Main(input_path string, verbose bool, args []string) error {
	tileStrings, err := util.ReadLineGroupsFromFile(input_path)
	if err != nil || len(tileStrings) == 0 {
		return err
	}

	// Convert tiles from string to matrix representation.
	tiles, err := parseTiles(tileStrings)
	if err != nil {
		return err
	}

	// Find the square dimensions of the output.
	fSqrt := math.Sqrt(float64(len(tiles)))
	if fSqrt != math.Floor(fSqrt) {
		return errors.New("number of tiles is not square")
	}
	dim := int(fSqrt)

	var nrows, ncols int
	for _, tile := range tiles {
		nrows, ncols = tile.Dims()
		fmt.Printf("Read %dx%d=%d tiles sized %dx%d.\n",
			dim, dim, len(tiles), nrows, ncols)
		break
	}

	return nil
}
