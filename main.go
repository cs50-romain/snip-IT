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
	Code string
}

var datamap map[string]map[string][]CodeSnip

func checkFileExistence(ln string, fn string) bool {
	data := datamap[ln]["Files"]

	for _, value := range data {
		if value.Filename == fn {
			return true
		}
	}
	return false
}

func updateMap(ln string, fn string, data []string) {
	codesnip := CodeSnip {
		Filename: fn,
		Code: strings.Join(data, " "),
	}

	if datamap[ln] == nil {
		datamap[ln] = make(map[string][]CodeSnip)	
	}

	if checkFileExistence(ln, fn) {
		fmt.Println("File already exists")
		return
	}

	datamap[ln]["Files"] = append(datamap[ln]["Files"], codesnip)
}

func deserializeJSON(filename string){
	jsonContents := readJSON(filename)
	
	err := json.Unmarshal(jsonContents, &datamap)
	if err != nil {
		fmt.Println("error: ", err)
	}
}

func readJSON(filename string) []byte {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func serializeJSON(data map[string]map[string][]CodeSnip) {
	b, err := json.Marshal(datamap)
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

func getInput(scanner *bufio.Scanner) []string{
	var lines []string
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text() + "\n"

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

func getContents(ln string, fn string) string {
	array := datamap[ln]["Files"]
	for _, data := range array {
		if data.Filename == fn {
			return data.Code
		}
	}
	return "File Not Found" 
}

func main() {
	datamap = make(map[string]map[string][]CodeSnip)
	
	var lanFlag = flag.String("l", "Go", "Choose language to load/save code snippet")
	var sFlag = flag.Bool("s", false, "Save file")
	var fFlag = flag.String("f", "default.txt", "Specify filename")
	flag.Parse()

	deserializeJSON("data.json")
	err := os.Truncate("data.json", 0)
	if err != nil {
		fmt.Println(err)
	}

	if *sFlag == true {
		fmt.Println("Input text (q to quit):")
		code := getInput(bufio.NewScanner(os.Stdin))
		updateMap(*lanFlag, *fFlag, code)
	}else {
		fmt.Println(getContents(*lanFlag, *fFlag))
	}

	serializeJSON(datamap)
}
