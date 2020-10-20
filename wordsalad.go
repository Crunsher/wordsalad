package main

import (
	"flag"
	"image"
	"os"
	"text/scanner"
	"bufio"
	"strings"
	"errors"
	"fmt"
	"strconv"
	"github.com/Crunsher/wordsalad/canvas"
)

func readWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sc := scanner.Scanner{}
	sc.Init(bufio.NewReader(file))

	var words= []string{}

	for tok := sc.Scan(); tok != scanner.EOF; tok = sc.Scan() {
		words = append(words, strings.ToLower(sc.TokenText()))
	}

	return words, nil
}

func parseSize(fs string) (int, int, error) {
	var errPrefix = "Failed to parse canvas size:"
	sizes := strings.Split(fs, ":")
	if len(sizes) != 2 {
		return 0, 0, errors.New(fmt.Sprintf("%s \"%v\" is not in the correct format.", errPrefix, fs))
	}
	x, err := strconv.Atoi((sizes[0]))
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("%s %v", errPrefix, err))
	}

	y, err := strconv.Atoi(sizes[1])
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("%s %v", errPrefix, err))
	}

	if x <= 0 || y <= 0 {
		return 0, 0, errors.New(fmt.Sprintf("%s Canvas size can't be 0 or less.", errPrefix, err))
	}

	return x, y, nil
}

var MaxTries = 2000
var TriesPerWord = 200

func main() {
	var wordList string
	var fieldSize string
	var fontPath string
	flag.StringVar(&wordList, "words", "words", "Path to csv word list. Default \"words\"")
	flag.StringVar(&fieldSize, "size", "30:30", "Size of field. Default \"30:30\"")
	flag.StringVar(&fontPath, "font", "RobotoMono-Medium.ttf", "Path to font file. Default \"RobotoMono-Medium.ttf\"")
	flag.Parse()
	
	fieldX, fieldY, err := parseSize(fieldSize)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	words, err := readWordsFromFile(wordList)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err.Error())
		return
	}

	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		fmt.Printf("Font file \"%s\" does not exist", fontPath)
		return
	}

	finalField := canvas.Field{}
	yay := false
	for	i := 0 ; i < MaxTries; i++ {
		f, success := canvas.FillField(words, fieldX, fieldY, TriesPerWord)
		if success {
			finalField = f
			yay = true
			break
		}
	}

	if yay {
		finalField.PrintField()
		fmt.Printf("\n%s\n\n", strings.Repeat("-", fieldX*2))
		finalField.SolveField()
	} else {
		fmt.Printf("Gave up after %d tries with %d tries per word.\n", MaxTries, TriesPerWord)
	}

	var cv *image.RGBA
	if cv, err = canvas.PaintImage(finalField, fontPath); err != nil {
		fmt.Printf(err.Error())
	}

	canvas.WriteToFile(cv, "out.png")
}
