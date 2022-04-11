package main

import "compiladores/Compiler"

func main() {
	//l := lexical_analysis.Lex("teste", "a2a1_a_ if -6 +7.4 725 -72.5")
	//l := Compiler.Lex("teste", "2 if a a2a1_a_ 2 +3 -44 -4.5 ++ -- * / */  /# abc abc abc abc abc abc #/ /# aaa #/ #/")
	l := Compiler.Lex("teste", "a != 2 && b >= 3 || c == 8")
	l.Debug()
}
