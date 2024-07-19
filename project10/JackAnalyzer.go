package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// os.Args = []string{"", "test/ExpressionLessSquare/Main.jack"}
	// os.Args[1] = "test/ExpressionLessSquare/Main.jack"
	if len(os.Args) < 2 {
		printErr("invalid number of arguments")
	}

	jackFiles := make([]*os.File, 0)

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

	if fileInfo.IsDir() {
		jackFiles = getJackFiles(file.Name())
	} else {
		jackFiles = append(jackFiles, file)
	}

	defer func(files []*os.File) {
		for _, f := range jackFiles {
			f.Close()
		}
	}(jackFiles)

	var parser Parser
	for _, jack := range jackFiles {
		src, err := os.ReadFile(jack.Name())
		if err != nil {
			printErr(err.Error())
		}
		parser.init(jack.Name(), src)
		parser.parseFile()
		parser.printTree()

		fileXML := parser.GetXML()
		writeErr := os.WriteFile(jack.Name()+".xml", []byte(fileXML), 0644)
		if writeErr != nil {
			printErr(err.Error())
		}
	}

}

func printErr(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

func getJackFiles(dir string) []*os.File {
	files, err := os.ReadDir(dir)
	jackFiles := make([]*os.File, 0, len(files))
	if err != nil {
		printErr(err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jack") {
			f, err := os.OpenFile(dir+"/"+file.Name(), os.O_RDONLY, 0)
			if err != nil {
				printErr(err.Error())
			}
			jackFiles = append(jackFiles, f)
		}
	}

	return jackFiles
}
