package main

import "lexical_analysis"

func main() {
	//l := lexical_analysis.Lex("teste", "a2a1_a_ if -6 +7.4 725 -72.5")
	l := lexical_analysis.Lex("teste", "2 if a a2a1_a_ 2 +3 -44 -4.5")
	l.Debug()

}
