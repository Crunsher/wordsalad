package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/Crunsher/wordsalad/canvas"
	"github.com/Crunsher/wordsalad/field"
	"os"
	"strconv"
	"strings"
	"text/scanner"
)

func readWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sc := scanner.Scanner{}
	sc.Init(bufio.NewReader(file))

	var words []string

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
	x, err := strconv.Atoi(sizes[0])
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

	newField := field.NewField(fieldX, fieldY, words)
	if err = newField.PositionWords(); err != nil  {
		fmt.Println(err.Error())
		return
	}

	newField.FillWithGarbage()
	fmt.Println(newField.AsciiField())

	pfield, err := canvas.PaintImage(newField. Bytes(), fontPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := canvas.WriteToFile(pfield, "out.png"); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("YAY")
}
