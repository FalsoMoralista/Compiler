package main

import "lexical_analysis"

func main() {
	l := lexical_analysis.Lex("teste", "a2a1_a_ if -6 +7.4 725 -72.5")
	l.Debug()
}
