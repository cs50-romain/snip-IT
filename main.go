package main

import (
	"encoding/json"
	"bufio"
	"fmt"
	"flag"
	"log"
	"os"
	"strings"
)

type CodeSnip struct {
	Filename string
	Code []string
}

var languages map[string]map[string][]string
var jsonmap map[string][]CodeSnip

func saveLang(lang string) {
	languages[lang] = make(map[string][]string)
}

func updateMap(lkey string, fkey string, contents []string) {
	languages[lkey][fkey] = contents
}

func deserializeJSON(filename string){
	fmt.Println("Got here")
	jsonContents := readJSON(filename)
	
	err := json.Unmarshal(jsonContents, &jsonmap)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(jsonmap)
}

func readJSON(filename string) []byte {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func serializeJSON(langs map[string]map[string][]string) {
	for language, nestmap:= range languages{
		for nestkey, value := range nestmap {
			codesnip := CodeSnip{
				Filename: nestkey,
				Code: value,
			}

			jsonmap[language] = append(jsonmap[language], codesnip)
		}
	}

	b, err := json.Marshal(jsonmap)
	if err != nil {
		log.Fatal(err)
	}
	saveToJSONFile(b)
}

func saveToJSONFile(data []byte){
	f, _ := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
	defer f.Close()
	_, err := f.Write(data)

	if err != nil {
		log.Fatal(err)
	}
}

func gatherFiless(scanner *bufio.Scanner) []string{
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

func getFiles(lkey string, fkey string) []string{
	return languages[lkey][fkey]
}

func main() {
	languages = make(map[string]map[string][]string)
	jsonmap = make(map[string][]CodeSnip)

	var lanFlag = flag.String("l", "Go", "Choose language to load/save code snippet")
	var sFlag = flag.Bool("s", false, "Save file")
	var fFlag = flag.String("f", "default.txt", "Specify filename")
	flag.Parse()

	if *sFlag == true {
		saveLang(*lanFlag)
		fmt.Println("Input text (q to quit):")
		updateMap(*lanFlag, *fFlag, gatherFiless(bufio.NewScanner(os.Stdin)))
		serializeJSON(languages)
	} else {
		// Populate map fro json data
		deserializeJSON("data.json")

		// Get desired Code Snippet
		content := getFiles(*lanFlag, *fFlag)
	
		for _, lines := range content {
			fmt.Println(lines)
		}
	}
}

