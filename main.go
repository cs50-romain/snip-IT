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

type JSON struct {
	Language string
	Content CodeSnip
}

var languages map[string]map[string][]string
var jsonArray []JSON

func saveLang(lang string) {
	languages[lang] = make(map[string][]string)
}

func updateMap(lkey string, fkey string, contents []string) {
	languages[lkey][fkey] = contents
}

func deserializeJSON(){
	
}

func serializeJSON(langs map[string]map[string][]string) {
	for language, nestmap:= range languages{
		for nestkey, value := range nestmap {
			codesnip := CodeSnip{
				Filename: nestkey,
				Code: value,
			}

			json := JSON{
				Language: language,
				Content: codesnip,
			}
			fmt.Println(json.Content)
			jsonArray = append(jsonArray, json)
		}
	}

	for _, jsonObj := range jsonArray {
		fmt.Println("JSON Object: ", jsonObj)
		b, err := json.Marshal(jsonObj)

		if err != nil {
			log.Fatal(err)
		}
		saveToJSONFile(b)
	}
}

func saveToJSONFile(data []byte){
	fmt.Println(string(data))
	f, _ := os.Create("data.json")
	defer f.Close()
	_, err := f.Write(data)

	if err != nil {
		log.Fatal(err)
	}
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

func getContent(lkey string, fkey string) []string{
	return languages[lkey][fkey]
}

func main() {
	languages = make(map[string]map[string][]string)

	var lanFlag = flag.String("l", "Go", "Choose language to load/save code snippet")
	var sFlag = flag.Bool("s", false, "Save file")
	var fFlag = flag.String("f", "default.txt", "Specify filename")
	flag.Parse()

	if *sFlag == true {
		saveLang(*lanFlag)
		fmt.Println("Input text (q to quit):")
		updateMap(*lanFlag, *fFlag, gatherContents(bufio.NewScanner(os.Stdin)))
	}
	
	// Will most likely need changed to also include the language specified
	content := getContent(*lanFlag, *fFlag)

	for _, lines := range content {
		fmt.Println(lines)
	}
	serializeJSON(languages)
}

