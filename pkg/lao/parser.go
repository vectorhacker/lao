package lao

import (
	"fmt"
	"strings"
)

// Parser parses tokens into an AST
type Parser interface {
	Parse() ([]Node, error)
}

type Visitor interface {
	Visit(Node)
}

type VisitorFunc func(Node)

type parser struct {
	tokenizer Tokenizer
}

func NewParser(tokenizer Tokenizer) Parser {
	return parser{tokenizer: tokenizer}
}

func (p parser) Parse() ([]Node, error) {

	nodes := []Node{}

	for {
		switch p.tokenizer.Current().Kind {
		case KindIdentifier:
			node, err := p.parseAssignmentStatement()
			if err != nil {
				return nodes, err
			}

			nodes = append(nodes, node)
			continue
		case KindKeyword:
			node, err := p.parseKeywordStatement()
			if err != nil {
				return nodes, err
			}

			nodes = append(nodes, node)
			continue
		case KindEnd:
			return nodes, nil
		}

		p.tokenizer.Next()
	}
}

func (p parser) parseAssignmentStatement() (Node, error) {
	variable, err := p.parseVariable()
	if err != nil {
		return nil, err
	}

	tokens := append(variable.Tokens(), p.tokenizer.Current())
	if p.tokenizer.Current().Kind != KindAssignment {
		return nil, fmt.Errorf(
			"Not proper variable statment line %d column %d",
			p.tokenizer.Current().Line,
			p.tokenizer.Current().Column,
		)
	}
	p.tokenizer.Next() // eat assignemt token

	var (
		left,
		exp Node
	)
	run(func() error {
		left, err = p.parseAtom(p.tokenizer.Current().Line)
		return err
	}, func() error {

		exp, err = p.parseArithmeticExpression(left, 0)
		return err
	})
	if err != nil {
		return nil, err
	}

	return AssignmentStatement{
		Variable:             variable.(Variable),
		ArithmeticExpression: exp.(ArithmeticExpression),
		tokens:               append(tokens, exp.Tokens()...),
	}, nil
}

var arithmeticPrecedence = map[ArithmeticOperator]int{
	ArithmeticAdd:            6,
	ArithmeticSubtract:       6,
	ArithmeticDivision:       7,
	ArithmeticMultiplication: 7,
}

func (p parser) getArithmeticOperator() ArithmeticOperator {
	switch strings.ToLower(p.tokenizer.Current().Value) {
	case ".add.":
		return ArithmeticAdd
	case ".sub.":
		return ArithmeticSubtract
	case ".div.":
		return ArithmeticDivision
	case ".mul.":
		return ArithmeticMultiplication
	}
	return 0
}

func (p parser) parseArithmeticExpression(
	left Node,
	prec int,
) (Node, error) {

	current := p.tokenizer.Current()
	if p.tokenizer.Current().Kind == KindArithmeticOperator {
		op := p.getArithmeticOperator()
		nextPrec := arithmeticPrecedence[op]
		if nextPrec > prec {
			p.tokenizer.Next()

			var (
				atom,
				right Node
			)
			err := run(func() error {
				var err error
				atom, err = p.parseAtom(current.Line)
				return err
			}, func() error {
				var err error
				right, err = p.parseArithmeticExpression(atom, nextPrec)
				return err
			})

			if err != nil {
				return nil, err
			}

			return p.parseArithmeticExpression(ArithmeticExpression{
				Left:     left,
				Right:    right,
				Operator: op,
				tokens:   append(append(left.Tokens(), current), right.Tokens()...),
			}, prec)
		}
	}

	return left, nil
}

func (p parser) parseReadStatement() (Node, error) {
	current := p.tokenizer.Current()
	p.tokenizer.Next()
	variable, err := p.parseVariable()
	if err != nil {
		return nil, err
	}

	return ReadStatement{
		Variable: variable.(Variable),
		tokens:   append([]Token{current}, variable.Tokens()...),
	}, nil
}

func (p parser) parseKeywordStatement() (Node, error) {
	switch strings.ToLower(p.tokenizer.Current().Value) {
	case "if":
		return p.parseIfStatement()
	case "read":
		return p.parseReadStatement()
	case "print":
		return p.parsePrintStatement()
	case "rem":
		return p.parseRemStatement()
	case "end":
		return p.parseEndStatement()
	}
	return nil, nil
}

func (p parser) getBinaryOperator() BinaryOperator {
	switch strings.ToLower(p.tokenizer.Current().Value) {
	case ".and.":
		return And
	case ".or.":
		return Or
	case ".not.":
		return Not
	case ".lt.":
		return LessThan
	case ".le.":
		return LessThanEqual
	case ".gt.":
		return GreaterThan
	case ".ge.":
		return GreaterThanEqual
	case ".eq.":
		return Equal
	case ".ne.":
		return NotEqual
	}

	return 0
}

func (p parser) parseString() (Node, error) {
	current := p.tokenizer.Current()
	return String{Value: current.Value, tokens: []Token{current}}, nil
}

func (p parser) parseNumber() (Node, error) {
	current := p.tokenizer.Current()
	switch current.Kind {
	case KindInteger:
		return IntegerNumber{Value: current.Value, tokens: []Token{current}}, nil
	case KindReal:
		return RealNumber{Value: current.Value, tokens: []Token{current}}, nil
	}
	return nil, fmt.Errorf("Not a number")
}

// exp is a hack
func (p parser) exp() (Node, error) {
	op := p.getBinaryOperator()
	current := p.tokenizer.Current()
	p.tokenizer.Next()
	right, err := p.parseAtom(current.Line)
	if err != nil {
		return nil, err
	}
	if op == Not {
		return p.parseExpresion(ConditionalExpression{
			Right:    right,
			Operator: op,
			tokens:   append([]Token{current}, right.Tokens()...),
		}, precedence[op])
	}
	return nil, fmt.Errorf("unxpected token %s at line: %d column: %d", current.Value, current.Line, current.Column)
}

func (p parser) parseAtom(line int) (Node, error) {

	current := p.tokenizer.Current()
	switch p.tokenizer.Current().Kind {
	case KindInteger, KindReal:
		defer p.tokenizer.Next()
		return p.parseNumber()
	case KindString:
		defer p.tokenizer.Next()
		return p.parseString()
	case KindIdentifier:
		return p.parseVariable()
	case KindLogicalOperator:
		return p.exp()
	}
	return nil, fmt.Errorf("unxpected token %s at line: %d column: %d", current.Value, current.Line, current.Column)
}

func (p parser) parseExpresion(left Node, prec int) (Node, error) {
	current := p.tokenizer.Current()
	if current.Kind == KindLogicalOperator ||
		current.Kind == KindRelationalOperator {

		operator := p.getBinaryOperator()
		nextPrec := precedence[operator]
		if nextPrec > prec {
			p.tokenizer.Next()

			atom, err := p.parseAtom(current.Line)
			if err != nil {
				return nil, err
			}
			right, err := p.parseExpresion(atom, nextPrec)
			if err != nil {
				return nil, err
			}

			return p.parseExpresion(ConditionalExpression{
				Left:     left,
				Right:    right,
				Operator: operator,
				tokens:   append(append(left.Tokens(), current), right.Tokens()...),
			}, prec)
		}

	}

	return left, nil
}

func run(steps ...func() error) error {
	var err error

	for _, step := range steps {
		if err != nil {
			return err
		}
		err = step()
	}

	return err
}

func (p parser) parseIfStatement() (Node, error) {

	current := p.tokenizer.Current()

	var (
		left,
		condition,
		statement Node
	)
	err := run(func() error {
		var err error
		p.tokenizer.Next() // eat if token

		left, err = p.parseAtom(current.Line)
		return err
	}, func() error {
		var err error
		condition, err = p.parseExpresion(left, 0)
		return err
	}, func() error {

		if p.tokenizer.Current().Kind != KindKeyword &&
			strings.ToLower(p.tokenizer.Current().Value) == "then" {
			return fmt.Errorf(
				"Not valid if then statement line: %d coulumn: %d",
				p.tokenizer.Current().Line,
				p.tokenizer.Current().Column,
			)
		}
		p.tokenizer.Next() // Eat then token

		var err error
		// statement, err = p.parseAtom(current.Line, true)
		statement, err = p.parseKeywordStatement()

		return err
	})
	if err != nil {
		return nil, err
	}
	return IfStatement{
		Condition:     condition.(ConditionalExpression),
		ThenStatement: statement,
		tokens:        append(append([]Token{current}, condition.Tokens()...), statement.Tokens()...),
	}, nil
}

func (p parser) parseEndStatement() (Node, error) {
	current := p.tokenizer.Current()

	p.tokenizer.Next()
	if p.tokenizer.Current().Kind != KindPeriod {
		return nil, fmt.Errorf("Not valid end statement. line: %d column %d", current.Line, current.Column)
	}

	return EndStatement{
		tokens: []Token{current, p.tokenizer.Current()},
	}, nil
}

func (p parser) parseRemStatement() (Node, error) {

	tokens := []Token{}
	line := p.tokenizer.Current().Line
	for p.tokenizer.Current().Line == line &&
		p.tokenizer.Current().Kind != KindEnd {
		tokens = append(tokens, p.tokenizer.Current())

		p.tokenizer.Next()
	}

	return RemStatement{
		tokens: tokens,
	}, nil
}

func (p parser) parsePrintStatement() (Node, error) {
	current := p.tokenizer.Current()
	p.tokenizer.Next()
	argument, err := p.parsePrintArguments(current.Line)
	if err != nil {
		return nil, err
	}

	tokens := []Token{current}
	if argument != nil {
		tokens = append(tokens, argument.Tokens()...)
	}

	return PrintStatement{
		tokens:     tokens,
		Argumenent: argument,
	}, nil
}

func (p parser) parseVariable() (Node, error) {
	if p.tokenizer.Current().Kind != KindIdentifier {
		return nil, fmt.Errorf("Not variable")
	}

	name := strings.ToLower(p.tokenizer.Current().Value)
	current := p.tokenizer.Current()
	tokens := []Token{current}

	switch {
	case name[0] >= 'a' && name[0] <= 'f':
		p.tokenizer.Next()
		return Variable{
			Type:   VariableInteger,
			Name:   name,
			tokens: tokens,
		}, nil
	case name[0] >= 'g' && name[0] <= 'n':
		p.tokenizer.Next()
		return Variable{
			Type:   VariableReal,
			Name:   name,
			tokens: tokens,
		}, nil
	case name[0] >= '0' && name[0] <= 'z':
		p.tokenizer.Next()
		return Variable{
			Type:   VariableString,
			Name:   name,
			tokens: tokens,
		}, nil
	}

	return nil, fmt.Errorf("Invalid identifier used as variable line: %d column: %d", p.tokenizer.Current().Line, p.tokenizer.Current().Column)
}

func (p parser) parsePrintArguments(line int) (Node, error) {
	current := p.tokenizer.Current()
	switch current.Kind {
	case KindKeyword:
		if current.Line == line {
			return nil, fmt.Errorf("Expected variable, string, number, or new line. Line: %d column: %d", current.Line, current.Column)
		}
		return nil, nil
	case KindInteger:
		return IntegerNumber{Value: current.Value, tokens: []Token{current}}, nil
	case KindReal:
		return RealNumber{Value: current.Value, tokens: []Token{current}}, nil
	case KindString:
		return String{Value: current.Value, tokens: []Token{current}}, nil
	case KindIdentifier:
		return p.parseVariable()
	}

	return nil, nil
}
