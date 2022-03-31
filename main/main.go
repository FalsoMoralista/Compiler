package main

import "lexical_analysis"

func main() {
	l := lexical_analysis.Lex("teste", " program B ll var \n consta_a1 a2a1_a_ -6 +7.4 725 -72.5")
	l.Debug()
}
