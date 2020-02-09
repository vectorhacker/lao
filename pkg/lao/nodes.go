package lao

// NodeType type of the node
type NodeType int

// NodeType
const (
	_ NodeType = iota
	TypeIfThen
	TypePrint
	TypeRem
	TypeeAssignment
	TypeRead
	TypeEnd
)

type Node interface {
	// Tokens returns the list of tokens that make up that node
	Tokens() []Token
}

type AssignmentStatement struct {
	Variable             Variable
	ArithmeticExpression Node
	tokens               []Token
}

func (a AssignmentStatement) Tokens() []Token {
	return a.tokens
}

// ArithmeticOperator operations for math
type ArithmeticOperator int

// ArithmeticOperator
const (
	_ ArithmeticOperator = iota
	ArithmeticAdd
	ArithmeticSubtract
	ArithmeticDivision
	ArithmeticMultiplication
)

// ArithmeticExpression expression
type ArithmeticExpression struct {
	Left     Node
	Right    Node
	Operator ArithmeticOperator
	tokens   []Token
}

func (a ArithmeticExpression) Tokens() []Token {
	return a.tokens
}

// BinaryOperator opeators for realtion
type BinaryOperator int

// RelationalOperator
const (
	_ BinaryOperator = iota
	LessThan
	LessThanEqual
	GreaterThan
	GreaterThanEqual
	Equal
	NotEqual
	And
	Not
	Or
)

var precedence = map[BinaryOperator]int{
	Or:               1,
	And:              2,
	Not:              3,
	Equal:            4,
	NotEqual:         4,
	LessThan:         5,
	LessThanEqual:    5,
	GreaterThan:      5,
	GreaterThanEqual: 5,
}

// PrintStatement node
type PrintStatement struct {
	tokens     []Token
	Argumenent Node
}

func (a PrintStatement) Tokens() []Token {
	return a.tokens
}

// VariableType types of variable.
type VariableType int

// VariableType
const (
	_ VariableType = iota
	VariableReal
	VariableInteger
	VariableString
)

// Variable node
type Variable struct {
	Name   string
	Type   VariableType
	tokens []Token
}

func (a Variable) Tokens() []Token {
	return a.tokens
}

// ReadStatement node
type ReadStatement struct {
	Token    Token
	Variable Variable
	tokens   []Token
}

func (r ReadStatement) Tokens() []Token {
	return r.tokens
}

// ConditionalExpression node
type ConditionalExpression struct {
	Left     Node
	Right    Node
	Operator BinaryOperator
	tokens   []Token
}

func (r ConditionalExpression) Tokens() []Token {
	return r.tokens
}

// IfStatement node
type IfStatement struct {
	Condition     ConditionalExpression
	ThenStatement Node
	tokens        []Token
}

func (r IfStatement) Tokens() []Token {
	return r.tokens
}

type RealNumber struct {
	Value  string
	tokens []Token
}

func (r RealNumber) Tokens() []Token {
	return r.tokens
}

type IntegerNumber struct {
	Value  string
	tokens []Token
}

func (r IntegerNumber) Tokens() []Token {
	return r.tokens
}

type String struct {
	Value  string
	tokens []Token
}

func (r String) Tokens() []Token {
	return r.tokens
}

type RemStatement struct {
	tokens []Token
}

func (r RemStatement) Tokens() []Token {
	return r.tokens
}

type EndStatement struct {
	tokens []Token
}

func (r EndStatement) Tokens() []Token {
	return r.tokens
}
