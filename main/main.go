package main

import (
	"compiladores/Compiler"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func Run() {
	var files []string

	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		matchFile, _ := regexp.MatchString("input(\\d+)?(\\.txt)", file)
		if matchFile {
			// busca por números na string
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			fileNumbers := re.FindAllString(file, -1)

			// cria um novo nome de arquivo "output" com base na string de input
			var outputFileName = "output.txt"
			if len(fileNumbers) >= 1 {
				outputFileName = "output" + fileNumbers[0] + ".txt"
			}

			// remove o arquivo de output anterior, caso exista
			err := os.Remove(outputFileName)
			content, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			Compiler.Lex("", outputFileName, string(content))
		}
	}
}

func main() {
	//	Run()
	content, err := ioutil.ReadFile("./teste.txt")
	if err != nil {
		log.Fatal(err)
	}
	l := Compiler.Lex("", "", string(content))
	Compiler.Analyze(l)
}
