package field

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type orientation int
const right orientation = 0
const down orientation = 1

type element struct {
	char byte
	isWord bool
}

type Field struct {
	height, width int
	Words         []string
	matrix        [][]element
}

func NewField(height int, width int, words []string) *Field {
	byteField := make([][]element, height)
	for x, _ := range byteField {
		byteField[x] = make([]element, width)
	}

	// Sort in descending order
	sort.Slice(words, func(i, k int) bool {
		return len(words[i]) > len(words[k])
	})
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}

	return &Field{
		height: height,
		width:  width,
		Words:  words,
		matrix: byteField,
	}
}

func (f Field) wordFits(word string, posX int, posY int, orientation orientation) bool {
	if orientation == right {
		if posY + len(word) > len(f.matrix[posX]) {
			return false
		}
		for i, c := range []byte(word) {
			if f.matrix[posX][posY+i].isWord && f.matrix[posX][posY+i].char != c {
				return false
			}
		}
	} else {
		if posX + len(word) > len(f.matrix) {
			return false
		}
		for i, c := range []byte(word) {
			if f.matrix[posX+i][posY].isWord && f.matrix[posX+i][posY].char != c {
				return false
			}
		}
	}

	return true
}

func (f Field) positionWord(word string, posX int, posY int, orientation orientation) {
	if orientation == right {
		for i, c := range []byte(word) {
			f.matrix[posX][posY+i] = element{c, true}
		}
	} else {
		for i, c := range []byte(word) {
			f.matrix[posX+i][posY] = element{c, true}
		}
	}
}

func (f Field) PositionWords() error {
	if len(f.Words[0]) > f.width || len(f.Words[0]) > f.height {
		return errors.New("word larger than field")
	}

	rand.Seed(time.Now().Unix())
	var orient orientation

	for _, word := range f.Words {
		tries := 0
		for tries < 100 {
			orient = orientation(rand.Int() % 2)
			posX, posY := rand.Int() % f.width, rand.Int() % f.height
			if (orient == right && len(word)+posX > f.width) ||
				(orient == down && len(word)+posY > f.height) {
				tries++
				continue
			}

			if !f.wordFits(word, posX, posY, orient) {
				tries++
				continue
			}

			fmt.Printf("%d:%d:%d %s (%d)\n", posX, posY, orient, word, tries)
			f.positionWord(word, posX, posY, orient)
			break
		}

		if tries == 100 {
			return errors.New(fmt.Sprintf("gave up after 100 tries for word \"%s\"", word))
		}
	}

	fmt.Println()
	return nil
}

func (f Field) AsciiField() string {
	var ascii bytes.Buffer

	for y := range f.matrix {
		for _, e := range f.matrix[y] {
			ascii.WriteByte(e.char)
		}
		ascii.WriteByte('\n')
	}

	return ascii.String()
}

func (f Field) FillWithGarbage() {
	rand.Seed(time.Now().Unix())

	for y := range f.matrix {
		for x, e := range f.matrix[y] {
			if e.isWord {
				continue
			}
			f.matrix[y][x] = element{byte(rand.Int()%26 + 65), false}
		}
	}
}

func (f Field) Bytes() [][]byte {
	bites := make([][]byte, f.width)


	for x := range f.matrix {
		bites[x] = make([]byte, f.height)
		for y, e := range f.matrix[x] {
			bites[x][y] = e.char
		}
	}

	return bites
}