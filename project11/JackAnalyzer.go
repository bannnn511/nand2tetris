package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// os.Args = []string{"", "test/Seven/Main.jack"}
	if len(os.Args) < 2 {
		printErr("invalid number of arguments")
	}

	jackFiles := make([]string, 0)

	// Open file
	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		printErr(fmt.Sprintf("%s file not exists\n", os.Args[1]))
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		printErr(err.Error())
	}

	isDir := false
	if fileInfo.IsDir() {
		isDir = true
		jackFiles = getJackFiles(file.Name())
	} else {
		jackFiles = append(jackFiles, file.Name())
	}

	for _, jack := range jackFiles {
		var parser Parser
		readFile := jack
		if isDir {
			readFile = os.Args[1] + "/" + jack
		}
		src, err := os.ReadFile(readFile)
		if err != nil {
			printErr(err.Error())
		}
		parser.Init(jack, src)
		parser.ParseFile()

		fileName := ""
		if isDir {
			fileName = os.Args[1] + "/" + jack[:strings.Index(jack, ".jack")]
		} else {
			fileName = os.Args[1][:strings.Index(os.Args[1], ".jack")]
		}
		writeErr := os.WriteFile(fileName+".an", []byte(parser.VmOut()), 0644)
		if writeErr != nil {
			printErr(err.Error())
		}
	}
}

func printErr(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

func getJackFiles(dir string) []string {
	files, err := os.ReadDir(dir)
	jackFiles := make([]string, 0, len(files))
	if err != nil {
		printErr(err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jack") {
			if err != nil {
				printErr(err.Error())
			}
			jackFiles = append(jackFiles, file.Name())
		}
	}

	return jackFiles
}
