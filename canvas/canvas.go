package canvas

import (
	"fmt"
	"sort"
	"math/rand"
	"time"
	. "github.com/logrusorgru/aurora"
)

type orientation int
const none  orientation = 0
const right orientation = 1
const down  orientation = 2

func (o orientation) String() string {
	switch o {
	case none:
		return "none"
	case right:
		return "right"
	case down:
		return "down"
	default:
		return "ERROR"
	}
}

type wordPosition struct {
	word string
	x, y int
	orient orientation
}

type element struct {
	c byte
	isWord bool
}

type field struct {
	width, height int
	field [][]element
	wordPositions []wordPosition
}

// newField returns a new field with the given height and width.
func newField (width int, height int) field {
	intField := make([][]element, height)
	for i := range intField {
		intField[i] = make([]element, width)
	}
	return field{width, height, intField, nil}
}

func (f field) printField() {
	for _, row := range f.field {
		for _, c := range row {
			if c.c != 0 {
				fmt.Printf(" %v", string(c.c))
			} else {
				fmt.Printf(" /")
			}
		}
		fmt.Println()
	}
}

func (f field) solveField() {
	for _, row := range f.field {
		for _, elem := range row {
			if elem.c != 0 {
				if elem.isWord {
					fmt.Printf(" %v", Magenta(string(elem.c)))
				} else {
					fmt.Printf(" %v",string(elem.c))
				}
			} else {
				fmt.Printf(" /")
			}
		}
		fmt.Println()
	}
}

// fillWithGarbage fills a fields empty spaces with random lower case alphabetical characters.
func (f field) fillWithGarbage() {
	rand.Seed(time.Now().Unix())

	for _, row := range f.field {
		for x, c := range row {
			if c.c == 0 {
				row[x] = element{ byte(rand.Int()%26 + 97), false }
			}
		}
	}
}

// placeWord tries to place a single string word on a field with the given starting coordinates and orientation.
// Returns true if successful and false if not.
func (f field) placeWord(word string, x int, y int, orient orientation) bool {
	if orient == 0 {
		if x + len(word) > f.width {
			return false
		}
		for i, c := range word {
			fc := f.field[y][x+i]
			if fc.c != 0 && fc.c != byte(c) {
				return false
			}
		}
		for i, c := range word {
			f.field[y][x+i] = element{byte(c), true}
		}
	} else {
		if y+len(word) > f.height {
			return false
		}
		for i, c := range word {
			fc := f.field[y+i][x]
			if fc.c != 0 && fc.c != byte(c) {
				return false
			}
		}
		for i, c := range word {
			f.field[y+i][x] = element{byte(c), true}
		}
	}

	return true
}

// TryPositionWords takes a slice of words and tries to position them on a field f. Should it fail to place all the words
// on the board it returns false, otherwise returns true. Words which are placed on the board are also stored in the
// fields wordPosition with their value, starting coordinates and orientation.
func (f field) tryPositionWords(words []string, maxTries int) bool {
	if len(words) == 0 {
		return true
	}

	// Sort the words by descending length
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) > len(words[j])
	})

	if len(words[0]) > f.width || len(words[0]) > f.height {
		return false
	}

	f.wordPositions = make([]wordPosition, len(words))

	rand.Seed(time.Now().Unix())
	for i, word := range words {
		tries := 0
		for ; tries < maxTries; tries++ {
			x := rand.Int() % f.width
			y := rand.Int() % f.height
			orient := orientation(rand.Int() % 2)
			if f.placeWord(word, x, y, orient) {
				f.wordPositions[i] = wordPosition{word, x, y, orient}
				break
			}
		}

		if tries == maxTries {
			fmt.Printf("Gave up after %d tries. (%s)\n", maxTries, word)
			return false
		}
	}

	return true
}

// fillField creates a new field with width and height and then tries to place the words onto it. Should it
// succeed it returns it and true, in a failure case empty and false
func fillField(words []string, width int, height int, TriesPerWord int) (field, bool) {
	arr := newField(width, height)
	if !arr.tryPositionWords(words, TriesPerWord) {
		return field{}, false
	}
	arr.fillWithGarbage()

	return arr, true
}
