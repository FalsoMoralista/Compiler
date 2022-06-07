package Compiler

import "fmt"

type syntacticAnalyzer struct {
	lex       *lexer
	tokenList []token
	index     int
}

func (s *syntacticAnalyzer) nextToken() token {
	if s.index < len(s.tokenList) {
		t := s.tokenList[s.index]
		s.index++
		return t
	} else {
		fmt.Println("Fim de lista de tokens.")
		return token{tokenEOF, "", 0} // todo: remove (unnecessary)
	}
}

func Analyze(l *lexer) *syntacticAnalyzer {
	tokens := make([]token, 0)
	for t := range l.Tokens() {
		tokens = append(tokens, t)
	}
	fmt.Println(tokens[0])
	s := &syntacticAnalyzer{lex: l, tokenList: tokens, index: 0}
	s.start()
	return s
}

func (s *syntacticAnalyzer) start() {
	t := s.nextToken()
	if t.Val == ProgramKeyword {
		t = s.nextToken()
		if t.Typ == tokenIdentifier { // TODO: whenever find an identifier, markup for symbol table.
			t = s.nextToken()
			if t.Val == ";" {
				s.globalStatement()
			}
		}
	}
}

func (s *syntacticAnalyzer) globalStatement() {
	for {
		switch t := s.nextToken(); {
		case t.Val == VarKeyword:
			s.varStatement()
		case t.Val == ConstKeyword:
			s.constStatement()
		case t.Val == RegisterKeyword:
		case t.Val == ProcedureKeyword:
		case t.Val == MainKeyword:
		case t.Typ == tokenEOF:
			fmt.Println("Reconhecido com sucesso")
			return
		}
	}
}

func (s *syntacticAnalyzer) varStatement() {
	t := s.nextToken()
	if t.Val == "{" {
		s.varList()
	} else {
		if t.Val == "}" {
			//ok
		} else {
			fmt.Println("Erro 1")
			// todo: throw error
		}
	}
}

func (s *syntacticAnalyzer) varList() {
	t := s.nextToken()
	if t.Val == IntegerKeyword || t.Val == StringKeyword || t.Val == RealKeyword || t.Val == BooleanKeyword || t.Val == CharKeyword || t.Typ == tokenIdentifier {
		s.varDeclaration()
		s.varList1()
	} else {
		if t.Val == "}" {
			// ok: return
		} else {
			fmt.Println("Erro 2 ")
			// todo: throw error
		}
	}
}

func (s *syntacticAnalyzer) varList1() {
	t := s.nextToken()
	if t.Val == IntegerKeyword || t.Val == StringKeyword || t.Val == RealKeyword || t.Val == BooleanKeyword || t.Val == CharKeyword || t.Typ == tokenIdentifier {
		s.varDeclaration()
		s.varList1()
	} else {
		if t.Val == "}" {
			// ok: return
		} else {
			fmt.Println("Erro 3")
			// todo: throw error
		}
	}
}

func (s *syntacticAnalyzer) varDeclaration() {
	t := s.nextToken()
	if t.Typ == tokenIdentifier {
		s.varDeclaration1()
	} else {
		fmt.Println("Erro 4")
		// todo: throw error
	}

}

func (s *syntacticAnalyzer) varDeclaration1() {
	t := s.nextToken()
	if t.Val == "," {
		t = s.nextToken()
		if t.Typ == tokenIdentifier {
			s.varDeclaration1()
		} else {
			fmt.Println("Erro 5")
			// todo: throw error
		}
	} else {
		if t.Val == ";" {
			// ok:  return
		} else {
			fmt.Println("Erro 5")
			// todo: throw error
		}
	}
}

func (s *syntacticAnalyzer) constStatement() {
	// todo
}
