package main

import (
	"flag"
	"github.com/crunsher/wordsalad/canvas"
)

func readWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
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
	var fp string
	var fs string
	flag.StringVar(&fp, "words", "", "Path to csv with words")
	flag.StringVar(&fs, "size", "30:30", "Size of field. Default \"30:30\"")
	flag.Parse()
	
	fieldX, fieldY, err := parseSize(fs)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	words, err := readWordsFromFile(fp)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err.Error())
		return
	}

	finalField := field{}
	yay := false
	for	i := 0 ; i < MaxTries; i++ {
		f, success := canvas.fillField(words, fieldX, fieldY)
		if success {
			finalField = f
			yay = true
			break;
		}
	}

	if yay {
		finalField.printField()
		fmt.Printf("\n%s\n\n", strings.Repeat("-", finalField.width*2))
		finalField.solveField()
	} else {
		fmt.Printf("Gave up after %d tries with %d tries per word.\n", MaxTries, TriesPerWord)
	}
}
