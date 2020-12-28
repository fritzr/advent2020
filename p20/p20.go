package p20

import (
	"errors"
	"fmt"
	"github.com/fritzr/advent2020/util"
	"gonum.org/v1/gonum/mat" // yeah it might be overkill, but I want to learn
	"math"
	"sort"
	"strconv"
	"strings"
)

var gVerbose = false

// Directions.
const (
	RIGHT = iota
	UP
	LEFT
	DOWN
)

var directionString = [...]string{"Right", "Up", "Left", "Down"}

type Direction = int

const NUM_DIRECTIONS = 1 + DOWN

const (
	UNFLIPPED = iota
	FLIPPED
)

type Flip = int

var flipString = [...]string{"flipped", "unflipped"}

const NUM_FLIPS = 1 + FLIPPED

// Some random prime numbers.
// Specifically, 10 of them, since our input consists of 10x10 tiles.
// If you want to support larger input, use a larger set of primes.
var primes = [...]int{
	2267, 63103, 30713, 8861, 43441, 1753, 5023, 14173, 46601, 35491}

// *Very* simple hash routine.
func floatHash(values mat.Vector) (hash float64) {
	for index := 0; index < values.Len(); index++ {
		hash += float64(primes[index%len(primes)]) * values.AtVec(index)
	}
	return
}

type edgeHashList = [NUM_DIRECTIONS][NUM_FLIPS]float64

type Tile struct {
	*mat.Dense
	// Hash of each edge.
	//
	// `edge[d][f]` is the hash of the edge in the given direction d
	// after applying the flip f. Some examples:
	//
	//   right := edge[RIGHT][UNFLIPPED]
	//   left := edge[LEFT][UNFLIPPED]
	//   upFlipped := edge[UP][FLIPPED]
	//
	id   int
	edge edgeHashList
}

type reverseVector struct {
	mat.Vector
}

func (r *reverseVector) AtVec(index int) float64 {
	return r.Vector.AtVec(r.Len() - index - 1)
}

func (t *Tile) ReverseRowView(row int) mat.Vector {
	return &reverseVector{t.RowView(row)}
}

func (t *Tile) ReverseColView(row int) mat.Vector {
	return &reverseVector{t.ColView(row)}
}

func NewTile(id int, nrows int, ncols int, data []float64) *Tile {
	t := Tile{mat.NewDense(nrows, ncols, data), id, edgeHashList{}}
	// The matrix is oriented with directional edges as follows:
	//
	//        0  ...  NC-1
	//      +-------------
	// 0    | LU  U   RU
	// ...  | L  ...  R
	// NR-1 | LD  D   RD
	//
	// Compute the hashes of each edge now.
	// An arrangement of tile T1 is adjacent to an arrangement of another tile T2
	// when T1.edge[d1][f1] = T2.edge[d2][f2] for some e1,e2,f1,f2.
	// By pre-computing all edge hashes we can query this quickly (assuming our
	// hashes have no collisions for the input).
	t.edge[UP][UNFLIPPED] = floatHash(t.RowView(0))
	t.edge[UP][FLIPPED] = floatHash(t.ReverseRowView(0))
	t.edge[DOWN][UNFLIPPED] = floatHash(t.RowView(nrows - 1))
	t.edge[DOWN][FLIPPED] = floatHash(t.ReverseRowView(nrows - 1))
	t.edge[LEFT][UNFLIPPED] = floatHash(t.ColView(0))
	t.edge[LEFT][FLIPPED] = floatHash(t.ReverseColView(0))
	t.edge[RIGHT][UNFLIPPED] = floatHash(t.ColView(ncols - 1))
	t.edge[RIGHT][FLIPPED] = floatHash(t.ReverseColView(ncols - 1))
	return &t
}

func BytesToFloat(data []byte) []float64 {
	float := make([]float64, len(data))
	for index, c := range data {
		float[index] = float64(c)
	}
	return float
}

// parseTiles creates tiles from string representations.
//
// Each tile string describes a square matrix by line-delimited rows.
func parseTiles(tileStrings []string) (tiles map[int]*Tile, err error) {
	tiles = make(map[int]*Tile, len(tileStrings))
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
		tiles[id] = NewTile(id, nrows, ncols, BytesToFloat([]byte(tileData)))
	}
	return
}

type TileEdge struct {
	id        int // tile ID
	direction int // UP, RIGHT, DOWN, or LEFT
	flipped   int // FLIPPED or UNFLIPPED
}

func (e TileEdge) String() string {
	return fmt.Sprintf("%d.%s(%c)",
		e.id, directionString[e.direction], flipString[e.flipped][0])
}

func IsAdjacent(t1 *Tile, e1 *TileEdge, t2 *Tile, e2 *TileEdge) bool {
	return t1.edge[e1.direction][e1.flipped] == t2.edge[e2.direction][e2.flipped]
}

type Adjacency = [2]TileEdge
type AdjacencyList = []Adjacency
type AdjacencyMap = map[int]AdjacencyList

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
func Adjacencies(tiles map[int]*Tile) AdjacencyMap {
	// Map tile IDs to their adjacent tiles (by ID).
	// The indexes of the value are the constants RIGHT, LEFT, UP, DOWN, etc...
	// A zero value (we assume no tile has ID zero) means unknown adjancency.
	// A value of -1 means not adjacent to any tile.
	adjacent := make(AdjacencyMap, len(tiles))

	// We continually visit a random tile, checking all other tiles.
	//
	// Each tile is itself square, and thus can be rotated four times
	// and flipped in two dimensions, providing twelve possible arrangements.
	// Naively, every arrangement of each tile must be checked with every other
	// arrangement of every other tile.

	ids := make([]int, len(tiles))
	index := 0
	for id := range tiles {
		ids[index] = id
		index++
	}
	// Sort the list of IDs for repeatable runs in verbose mode.
	if gVerbose {
		sort.Ints(ids)
	}

	// Naively compute every possible adjacency.
	for index1, id1 := range ids {
		t1 := tiles[id1]
		for d1 := 0; d1 < len(t1.edge); d1++ {
			for f1 := 0; f1 < len(t1.edge[d1]); f1++ {
				for index2 := index1 + 1; index2 < len(ids); index2++ {
					id2 := ids[index2]
					t2 := tiles[id2]
					for d2 := 0; d2 < len(t2.edge); d2++ {
						for f2 := 0; f2 < len(t2.edge[d2]); f2++ {
							e1, e2 := TileEdge{id1, d1, f1}, TileEdge{id2, d2, f2}
							if IsAdjacent(t1, &e1, t2, &e2) {
								adjacency := Adjacency{e1, e2}
								if gVerbose {
									fmt.Printf("%v adjacent to %v\n", e1, e2)
								}
								if _, ok := adjacent[id1]; !ok {
									adjacent[id1] = make(AdjacencyList, 0, len(tiles)-1)
								}
								adjacent[id1] = append(adjacent[id1], adjacency)
								if _, ok := adjacent[id2]; !ok {
									adjacent[id2] = make(AdjacencyList, 0, len(tiles)-1)
								}
								adjacent[id2] = append(adjacent[id2], adjacency)
							} /* else if gVerbose {
								e1, e2 := TileEdge{id1, d1, f1}, TileEdge{id2, d2, f2}
								fmt.Printf("%v not adjacent to %v\n", e1, e2)
							} */
						}
					}
				}
			}
		}
	}

	return adjacent
}

func Corners(tiles map[int]*Tile) (corners [4]int, err error) {
	// There should be only four tiles with exactly two adjacent tiles.
	// XXX Currently we count each adjacency twice, so look for four adjacents.
	adjacencies := Adjacencies(tiles)

	cornerNum := 0
	for id, adjacency := range adjacencies {
		if 2 == len(adjacency)/2 /* XXX fix div 2 */ {
			if gVerbose {
				fmt.Printf("%d has 2 adjacencies: %v\n", id, adjacency)
			}
			if cornerNum == 4 {
				err = errors.New("too many corners!")
				break
			}
			corners[cornerNum] = id
			cornerNum++
		}
	}

	if cornerNum != 4 {
		err = errors.New(fmt.Sprintf("not enough corners (%d)!", cornerNum))
	}

	return
}

func Main(input_path string, verbose bool, args []string) error {
	gVerbose = verbose
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

	/*
		grid, err := Reconstitute(tiles)
		if err != nil {
			return err
		}

		if len(grid) != nrows*ncols {
			return errors.New("result is incomplete")
		}

		// Print the reconstruction.
		for row := 0; row != nrows; row++ {
			for col := 0; col != ncols; col++ {
				if col != 0 {
					fmt.Printf("    ")
				}
				fmt.Printf("%4d", grid[row*ncols+col])
			}
			fmt.Printf("\n")
		}
		corners := [...]int{grid[0], grid[ncols-1], grid[(nrows-1)*ncols],
			grid[nrows*ncols-1]}
	*/

	corners, err := Corners(tiles)
	if err != nil {
		return err
	}

	// Print the product of the corner IDs.
	result := 1
	for index, value := range corners {
		if index != 0 {
			fmt.Printf(" * ")
		}
		fmt.Printf("%d", value)
		result *= value
	}
	fmt.Printf(" = %d\n", result)

	return nil
}
