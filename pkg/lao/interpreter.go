package lao

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Interpreter executes the AST
type Interpreter interface {
	Execute([]Node) error
}

func NewInterpreter(out io.Writer) Interpreter {
	return &interpreter{
		out:     out,
		symbols: map[string]interface{}{},
	}
}

type interpreter struct {
	symbols map[string]interface{}
	out     io.Writer
}

func (i *interpreter) evalauteArithmeticExpression(
	vType VariableType,
	exp interface{},
) (interface{}, error) {

	switch e := exp.(type) {
	case ArithmeticExpression:
		left, err := i.evalauteArithmeticExpression(vType, e.Left)
		if err != nil {
			return nil, err
		}
		right, err := i.evalauteArithmeticExpression(vType, e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator {
		case ArithmeticAdd:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l + r, nil
				case float64:
					return float64(l) + r, nil
				case string:
					return fmt.Sprintf("%d%s", l, r), nil
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l + float64(r), nil
				case float64:
					return l + r, nil
				case string:
					return fmt.Sprintf("%.6f%s", l, r), nil
				}
			case string:
				switch r := right.(type) {
				case int:
					return fmt.Sprintf("%s%d", l, r), nil
				case float64:
					return fmt.Sprintf("%s%.6f", l, r), nil
				case string:
					return l + r, nil
				}
			}
		case ArithmeticSubtract:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l - r, nil
				case float64:
					return float64(l) - r, nil
				case string:
					return nil, fmt.Errorf("Cannot substract string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l - float64(r), nil
				case float64:
					return l - r, nil
				case string:
					return nil, fmt.Errorf("Cannot substract string")
				}
			case string:
				return nil, fmt.Errorf("Cannot substract from string")
			}
		case ArithmeticDivision:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l / r, nil
				case float64:
					return float64(l) / r, nil
				case string:
					return nil, fmt.Errorf("Cannot divide string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l / float64(r), nil
				case float64:
					return l / r, nil
				case string:
					return nil, fmt.Errorf("Cannot divide string")
				}
			case string:
				return nil, fmt.Errorf("Cannot divide string")
			}
		case ArithmeticMultiplication:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l * r, nil
				case float64:
					return float64(l) * r, nil
				case string:
					return nil, fmt.Errorf("Cannot multiply string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l * float64(r), nil
				case float64:
					return l * r, nil
				case string:
					return nil, fmt.Errorf("Cannot multiply string")
				}
			case string:
				return nil, fmt.Errorf("Cannot multiply string")
			}
		}
	case Variable:
		// TODO implement read variable
		if value, ok := i.symbols[e.Name]; ok {
			return value, nil
		}
		return nil, fmt.Errorf("No variable named %s", e.Name)
	case IntegerNumber:
		return strconv.Atoi(e.Value)
	case RealNumber:
		return strconv.ParseFloat(e.Value, 64)
	case String:
		return strings.ReplaceAll(e.Value, "\"", ""), nil
	}

	return nil, nil
}

func (i *interpreter) interpretAssignment(a AssignmentStatement) error {
	value, err := i.evalauteArithmeticExpression(a.Variable.Type, a.ArithmeticExpression)
	if err != nil {
		return err
	}

	switch value.(type) {
	case int:
		if a.Variable.Type != VariableInteger {
			return fmt.Errorf("invalid assignment to variable type integer")
		}
	case float64:
		if a.Variable.Type != VariableReal {
			return fmt.Errorf("invalid assignment to variable type real")
		}
	case string:
		if a.Variable.Type != VariableString {
			return fmt.Errorf("invalid assignment to variable type string")
		}
	}

	i.symbols[a.Variable.Name] = value
	return nil
}

func (i *interpreter) evaluateExpression(expr interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case ConditionalExpression:
		left, err := i.evaluateExpression(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := i.evaluateExpression(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator {
		case Not:
			r, ok := right.(bool)
			if !ok {
				return nil, fmt.Errorf("unable to convert expression to boolean")
			}
			return !r, nil
		case And:
			l, lok := left.(bool)
			r, rok := right.(bool)
			if !(lok && rok) {
				return nil, fmt.Errorf("unable to convert expression to boolean")
			}
			return l && r, nil
		case Or:
			l, lok := left.(bool)
			r, rok := right.(bool)
			if !(lok && rok) {
				return nil, fmt.Errorf("unable to convert expression to boolean")
			}
			return l || r, nil
		case LessThan:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l < r, nil
				case float64:
					return float64(l) < r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l < float64(r), nil
				case float64:
					return l < r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case string:
				switch r := right.(type) {
				case int:
					return nil, fmt.Errorf("cannot compare integer and string")
				case float64:
					return nil, fmt.Errorf("cannot compare integer and string")
				case string:
					return l < r, nil
				}
			}
		case LessThanEqual:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l <= r, nil
				case float64:
					return float64(l) <= r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l <= float64(r), nil
				case float64:
					return l <= r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case string:
				switch r := right.(type) {
				case int:
					return nil, fmt.Errorf("cannot compare integer and string")
				case float64:
					return nil, fmt.Errorf("cannot compare integer and string")
				case string:
					return l <= r, nil
				}
			}
		case Equal:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l == r, nil
				case float64:
					return float64(l) == r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l == float64(r), nil
				case float64:
					return l == r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case string:
				switch r := right.(type) {
				case int:
					return nil, fmt.Errorf("cannot compare integer and string")
				case float64:
					return nil, fmt.Errorf("cannot compare integer and string")
				case string:
					return l == r, nil
				}
			}
		case GreaterThan:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l > r, nil
				case float64:
					return float64(l) > r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l > float64(r), nil
				case float64:
					return l > r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case string:
				switch r := right.(type) {
				case int:
					return nil, fmt.Errorf("cannot compare integer and string")
				case float64:
					return nil, fmt.Errorf("cannot compare integer and string")
				case string:
					return l > r, nil

				}
			}
		case GreaterThanEqual:
			switch l := left.(type) {
			case int:
				switch r := right.(type) {
				case int:
					return l >= r, nil
				case float64:
					return float64(l) >= r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case float64:
				switch r := right.(type) {
				case int:
					return l >= float64(r), nil
				case float64:
					return l >= r, nil
				case string:
					return nil, fmt.Errorf("cannot compare integer and string")
				}
			case string:
				switch r := right.(type) {
				case int:
					return nil, fmt.Errorf("cannot compare integer and string")
				case float64:
					return nil, fmt.Errorf("cannot compare integer and string")
				case string:
					return l >= r, nil

				}
			}
		}
	case Variable:
		return i.symbols[e.Name], nil
	case IntegerNumber:
		return strconv.Atoi(e.Value)
	case RealNumber:
		return strconv.ParseFloat(e.Value, 64)
	case String:
		return strings.ReplaceAll(e.Value, "\"", ""), nil
	}

	return false, nil
}

func (i *interpreter) interpretIf(ifStatement IfStatement) error {

	condition, err := i.evaluateExpression(ifStatement.Condition)
	if err != nil {
		return err
	}

	cond, ok := condition.(bool)
	if !ok {
		return fmt.Errorf("Invalid condition")
	}
	if cond {
		return i.evaluateStatement(ifStatement.ThenStatement)
	}

	return nil
}

func (i *interpreter) interpretPrint(print PrintStatement) error {

	switch a := print.Argumenent.(type) {
	case Variable:
		v, ok := i.symbols[a.Name]
		if !ok {
			return fmt.Errorf("Unable to find variable  %s", a.Name)
		}

		switch a.Type {
		case VariableInteger:
			fmt.Fprintf(i.out, "%d\n", v)
		case VariableString:
			fmt.Fprintf(i.out, "%s\n", v)
		case VariableReal:
			fmt.Fprintf(i.out, "%.6f\n", v)
		}
	case String:
		fmt.Fprintln(i.out, strings.ReplaceAll(a.Value, "\"", ""))
	case IntegerNumber:
		fmt.Fprintln(i.out, a.Value)
	case RealNumber:
		fmt.Fprintln(i.out, a.Value)
	default:
		fmt.Fprintln(i.out)
	}

	return nil
}

func (i *interpreter) interpretRead(read ReadStatement) error {

	switch read.Variable.Type {
	case VariableInteger:
		var temp int
		if _, err := fmt.Scanf("%d\n", &temp); err != nil {
			return err
		}
		i.symbols[read.Variable.Name] = temp
	case VariableReal:
		var temp float64
		if _, err := fmt.Scanf("%.6f\n", &temp); err != nil {
			return err
		}
		i.symbols[read.Variable.Name] = temp
	case VariableString:
		var temp string
		if _, err := fmt.Scanf("%s\n", &temp); err != nil {
			return err
		}
		i.symbols[read.Variable.Name] = temp
	}

	return nil
}

func (i *interpreter) evaluateStatement(statement Node) error {

	switch s := statement.(type) {
	case RemStatement:
	case AssignmentStatement:
		return i.interpretAssignment(s)
	case IfStatement:
		return i.interpretIf(s)
	case PrintStatement:
		return i.interpretPrint(s)
	case ReadStatement:
		return i.interpretRead(s)
	case EndStatement:
		return io.EOF
	}
	return nil
}

func (i *interpreter) Execute(statements []Node) error {
	for _, statement := range statements {
		err := i.evaluateStatement(statement)
		if err != nil {
			return err
		}
	}

	return nil
}
