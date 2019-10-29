package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var symbols string
	var inputPath string
	var outPath string
	flag.StringVar(&symbols, "symbols", "EURUSD,USDJPY", "Symbol name")
	flag.StringVar(&inputPath, "inputPath", "file.log", "Input File")
	flag.StringVar(&outPath, "outputPath", "result", "Out File")
	flag.Parse()
	symbolsArray := strings.Split(symbols, ",")
	filesArray := strings.Split(inputPath, ",")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		os.Mkdir(outPath, os.ModeDir)
	}

	nubbersChanels := len(filesArray) * len(symbolsArray)
	fmt.Println(nubbersChanels)
	// responses := make(chan string, nubbersChanels)
	start := time.Now()
	// defer close(responses)
	for _, file := range filesArray {

		for _, symbol := range symbolsArray {
			wg.Add(1)
			fmt.Println(symbol)
			go ParseFile(symbol, file, outPath, &wg)

		}
	}

	// for response := range <-responses {

	// 	fmt.Println("from channel", response)
	// }

	wg.Wait()
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println(elapsed)
}

//ParseFile is parsing file
func ParseFile(symbol string, inputPath string, outPath string, wg *sync.WaitGroup) {

	pattern := "\\d*\\sGMT\\s{1,3}%s\\s.*"
	pattern = fmt.Sprintf(pattern, symbol)
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	outfilename := outPath + "/" + strings.Split(inputPath, ".")[0] + "_" + symbol + ".log"
	outfile, err := os.Create(outfilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outfile)
	defer wg.Done()
	for scanner.Scan() {
		res, err := regexp.MatchString(pattern, scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		if res {
			// fmt.Println(scanner.Text())
			fmt.Fprintln(writer, scanner.Text())
			scanner.Scan()
			pattern1 := "\\s\\s\\s\\w{3}.*|-----"
			for res1, err := regexp.MatchString(pattern1, scanner.Text()); res1; scanner.Scan() {
				if err != nil {
					log.Fatal(err)
				}
				// fmt.Println(scanner.Text())
				fmt.Fprintln(writer, scanner.Text())
				res1, _ = regexp.MatchString(pattern1, scanner.Text())
			}
		}
		writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// chanel <- outfilename
}
