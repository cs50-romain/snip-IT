package main

import (
	//"encode/json"
	"bufio"
	"fmt"
	"flag"
	"log"
	"os"
	"strings"
)

type CodeSnip struct {
	language string
	filename string
	content string
}
var languages []string
var files = map[string][]string{}

func saveLang(lang string) {
	count := 0
	for _, val := range languages {
		if lang == val {
			count++
		}
	}

	if count == 0 {
		languages = append(languages, lang)
	}
}

func updateMap(fkey string, contents []string) {
	files[fkey] = contents
}

func serializeJSON(langs []string) {
	
}

func gatherContents(scanner *bufio.Scanner) []string{
	var lines []string
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()

		if strings.ToLower(strings.TrimSpace(line)) == "q" {
			break
		}
		lines = append(lines, line)
	}

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return lines
}

func getContent(key string) []string{
	return files[key]
}

func main() {
	var lanFlag = flag.String("l", "Go", "Choose language to load/save code snippet")
	var sFlag = flag.Bool("s", false, "Save file")
	var fFlag = flag.String("f", "default.txt", "Specify filename")
	flag.Parse()

	if *sFlag == true {
		saveLang(*lanFlag)
		fmt.Println("Input text:")
		updateMap(*fFlag, gatherContents(bufio.NewScanner(os.Stdin)))
	}
	
	// Will most likely need changed to also include the language specified
	content := getContent(*fFlag)

	for _, lines := range content {
		fmt.Println(lines)
	}
}

