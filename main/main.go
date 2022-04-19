package main

import (
	"compiladores/Compiler"
	"io/ioutil"
	"log"
	"os"
)

func Read(filename string) string {
	err := os.Remove("output.txt")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func main() {
	if len(os.Args) >= 2 {
		Compiler.Lex("test", Read(os.Args[1]))
	} else {
		Compiler.Lex("test", Read("example.txt"))
	}

	//l := lexical_analysis.Lex("teste", "a2a1_a_ if -6 +7.4 725 -72.5")
	//l := Compiler.Lex("teste", "2 if a a2a1_a_ 2 +3 -44 -4.5 ++ -- * / */  /# abc abc abc abc abc abc #/ /# aaa #/ #/")
	//l := Compiler.Lex("teste", "a != 2 && b >= 3 || c == 8")
	//l := Compiler.Lex("teste", "b = 'a' =, + '1fsdf'")
	//Compiler.Lex("teste", "!&= \n != == <= >= \n\n\n\n- b = 'a' =, + '1' \"asdfds fdfsd dfsd d\" 123 01.49 /#abc ab aba aa ")
	//l.Debug()
}
