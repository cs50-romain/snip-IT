package main

import (
	"encoding/json"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CodeSnip struct {
	Filename string
	Code string
}

type MyError struct{}

func (e *MyError) Error() string {
	return "Invalid command"
}

var datamap map[string]map[string][]CodeSnip
var jsonfile string

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
	f, _ := os.OpenFile(jsonfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
	defer f.Close()
	_, err := f.Write(data)

	if err != nil {
		log.Fatal(err)
	}
}

func getInput(scanner *bufio.Scanner) []string{
	var lines []string
	for {
		fmt.Print(">> ")
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

func getClosestMatch(arr []CodeSnip, fn string) []string {
	var pfiles []string
	for _, string := range arr {
		pmatch := string.Filename[:4]
		if len(fn) > 3 {
			if fn[:4] == pmatch {
				pfiles = append(pfiles, string.Filename)
			}
		}
	}
	return pfiles
}

func getContents(ln string, fn string) string {
	array := datamap[ln]["Files"]
	for _, data := range array {
		if data.Filename == fn {
			return data.Code
		}
	}

	pfiles := getClosestMatch(array, fn)
	if pfiles != nil {
		return "Here are possible matches: " + strings.Join(pfiles, " ")
	}

	return "File Not Found" 
}

func parseCmd(input string) ([]string, error) {
	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")
	cmd := args[0]

	if cmd != "get" && cmd != "save" && cmd != "exit" {
		return nil, &MyError{}
	}

	return args, nil 
}

func isExist(fn string) bool {
	f, err := os.Open(fn)
	defer f.Close()

	if err != nil {
		return false
	}
	return true
}

func printASCIIArt(){
	asciiArt := ` ________  ________   ___  ________                ___  _________   
|\   ____\|\   ___  \|\  \|\   __  \              |\  \|\___   ___\ 
\ \  \___|\ \  \\ \  \ \  \ \  \|\  \ ____________\ \  \|___ \  \_| 
 \ \_____  \ \  \\ \  \ \  \ \   ____\\____________\ \  \   \ \  \  
  \|____|\  \ \  \\ \  \ \  \ \  \___\|____________|\ \  \   \ \  \ 
    ____\_\  \ \__\\ \__\ \__\ \__\                  \ \__\   \ \__\
   |\_________\|__| \|__|\|__|\|__|                   \|__|    \|__|
   \|_________|                                                     
                               `
	fmt.Println(asciiArt)
}

func main() {
	printASCIIArt()

	datamap = make(map[string]map[string][]CodeSnip)
	jsonfile = "/home/lettuce/snip-IT/data.json"
	
	if isExist(jsonfile) {
		deserializeJSON(jsonfile)
		err := os.Truncate(jsonfile, 0)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		f, error := os.Create("data.json")
		defer f.Close()
		if error != nil {
			fmt.Println(error)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		args, error := parseCmd(input)
		if error != nil {
			fmt.Println(error)
		} else {
			cmd := args[0]

			if cmd == "exit" {
				fmt.Println("[!] Quitting...")
				serializeJSON(datamap)
				os.Exit(0)
			}

			ln := args[1]
			fn := args[2]
	
			if cmd == "save" {
				fmt.Println("Input text (q to quit):")
				code := getInput(bufio.NewScanner(os.Stdin))
				updateMap(ln, fn, code)
			} else if cmd == "get" {
				fmt.Println(getContents(ln, fn))
			} 
		}
	}
	serializeJSON(datamap)
}
