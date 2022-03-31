package main

import "lexical_analysis"

func main() {
	l := lexical_analysis.Lex("teste", " B A C ")
	l.Debug()
}
