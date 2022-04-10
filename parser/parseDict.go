package parser

import (
	"github.com/BEN00262/simpleLang/lexer"
)

func (parser *Parser) _parseDict() DictNode {
	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		lexer.CURLY_BRACES, "{",
		"Expected '{'",
	)

	dict_elements := []DictElementNode{}

	// blockStatements := parser._parseBlockStatements()
	for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(parser.CurrentToken(), lexer.CURLY_BRACES, "}") {
		dict := DictElementNode{}

		if parser.CurrentToken().Type == lexer.VARIABLE {
			dict.Key = VariableNode{
				Value: parser.CurrentToken().Value.(string),
			}

			parser.eatToken()
		}

		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			lexer.COLON, ":",
			"Expected ':'",
		)

		dict.Value = parser._parseExpression()
		dict_elements = append(dict_elements, dict)

		if parser.CurrentToken().Type == lexer.COMMA {
			parser.eatToken()
		}
	}

	parser.IsExpectedEatElsePanic(
		parser.CurrentToken(),
		lexer.CURLY_BRACES, "}",
		"Expected '}'",
	)

	return DictNode{
		Elements: dict_elements,
	}
}
