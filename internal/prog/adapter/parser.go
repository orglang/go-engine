package adapter

import (
	"fmt"
	"slices"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"orglang/go-engine/internal/adt/uniqsym"
	pooltypedef "orglang/go-engine/internal/pool/typedef/core"
	proctypedef "orglang/go-engine/internal/proc/typedef/core"
	coreprog "orglang/go-engine/internal/prog/core"
)

type programExpression interface {
	programExpression()
}

type programAST struct {
	Expressions []programExpression `parser:"@@*"`
}

type poolDefAST struct {
	TypeQN  string     `parser:"'pool' 'type' @Ident"`
	TypeExpr poolExprAST `parser:"'{' @@? '}'"`
}

func (poolDefAST) programExpression() {}

type procDefAST struct {
	TypeQN  string     `parser:"'proc' 'type' @Ident"`
	TypeExpr procExprAST `parser:"'{' @@? '}'"`
}

func (procDefAST) programExpression() {}

type exprKind string

func (k *exprKind) Parse(lex *lexer.PeekingLexer) error {
	value := exprKind(lex.Peek().Value)
	if slices.Contains([]exprKind{"one", "link", "plus", "with", "tensor", "lolli", "up", "down"}, value) {
		*k = value
	}
	return nil
}

type poolExprAST struct {
	Kind   exprKind       `parser:"@@"`
	Link   *linkAST      `parser:"( 'link' @@"`
	Plus   *poolLaborAST `parser:"| 'plus' @@"`
	With   *poolLaborAST `parser:"| 'with' @@"`
	Tensor *poolResAST   `parser:"| 'tensor' @@"`
	Lolli  *poolResAST   `parser:"| 'lolli' @@"`
	Up     *shiftAST     `parser:"| 'up' @@"`
	Down   *shiftAST     `parser:"| 'down' @@"`
	One    bool          `parser:"| 'one' )"`
}

type procExprAST struct {
	Kind   exprKind       `parser:"@@"`
	Link   *linkAST      `parser:"( 'link' @@"`
	Tensor *procProdAST  `parser:"| 'tensor' @@"`
	Lolli  *procProdAST  `parser:"| 'lolli' @@"`
	Plus   *procSumAST   `parser:"| 'plus' @@"`
	With   *procSumAST   `parser:"| 'with' @@"`
	Up     *procShiftAST `parser:"| 'up' @@"`
	Down   *procShiftAST `parser:"| 'down' @@"`
	One    bool          `parser:"| 'one' )"`
}

type linkAST struct {
	TypeQN string `parser:"@Ident"`
}

type poolLaborAST struct {
	ProcQNs []string   `parser:"'(' @Ident* ')'"`
	ContExp poolExprAST `parser:"'{' @@? '}'"`
}

type poolResAST struct {
	ValExp  poolExprAST `parser:""`
	ContExp poolExprAST `parser:"'{' @@? '}'"`
}

type shiftAST struct {
	Cont poolExprAST `parser:"'{' @@? '}'"`
}

type procShiftAST struct {
	Cont procExprAST `parser:"'{' @@? '}'"`
}

type procProdAST struct {
	Val  procExprAST `parser:"'{' @@? '}'"`
	Cont procExprAST `parser:"'{' @@? '}'"`
}

type procSumAST struct {
	Choices []procChoiceAST `parser:"'{' @@* '}'"`
}

type procChoiceAST struct {
	Label string     `parser:"'@' @Ident"`
	Cont  procExprAST `parser:"'{' @@? '}'"`
}

func parseProgram(text string) (programAST, error) {
	parser, err := participle.Build[programAST](participle.Union[programExpression](poolDefAST{}, procDefAST{}))
	if err != nil {
		return programAST{}, err
	}
	ast, err := parser.ParseString("", text)
	if err != nil {
		return programAST{}, err
	}
	return *ast, nil
}

func programASTToCore(ast programAST) (coreprog.Spec, error) {
	var result coreprog.Spec
	for _, expression := range ast.Expressions {
		switch expression := expression.(type) {
		case poolDefAST:
			qn, err := uniqsym.ConvertFromString(expression.TypeQN)
			if err != nil {
				return coreprog.Spec{}, err
			}
			expr, err := poolExprToCore(expression.TypeExpr)
			if err != nil {
				return coreprog.Spec{}, err
			}
			result.Pools = append(result.Pools, pooltypedef.DefSpec{TypeQN: qn, TypeExpr: exp})
		case procDefAST:
			qn, err := uniqsym.ConvertFromString(expression.TypeQN)
			if err != nil {
				return coreprog.Spec{}, err
			}
			expr, err := procExprToCore(expression.TypeExpr)
			if err != nil {
				return coreprog.Spec{}, err
			}
			result.Procs = append(result.Procs, proctypedef.DefSpec{TypeQN: qn, TypeExpr: exp})
		default:
			return coreprog.Spec{}, fmt.Errorf("unsupported program expression: %T", expression)
		}
	}
	return result, nil
}

func poolExprToCore(ast poolExprAST) (pooltypedef.ExpSpec, error) {
	switch ast.Kind {
	case "one":
		return pooltypedef.OneSpec{}, nil
	case "link":
		qn, err := uniqsym.ConvertFromString(ast.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return pooltypedef.LinkSpec{TypeQN: qn}, nil
	case "plus", "with":
		labor := ast.PlusOrWith()
		qns := make([]uniqsym.ADT, len(labor.ProcQNs))
		for i, value := range labor.ProcQNs {
			qn, err := uniqsym.ConvertFromString(value)
			if err != nil {
				return nil, err
			}
			qns[i] = qn
		}
		cont, err := poolExprToCore(labor.ContExp)
		if err != nil {
			return nil, err
		}
		if ast.Kind == "plus" {
			return pooltypedef.PlusSpec{ProcQNs: qns, ContExp: cont}, nil
		}
		return pooltypedef.WithSpec{ProcQNs: qns, ContExp: cont}, nil
	case "up", "down":
		cont, err := poolExprToCore(ast.Shift().Cont)
		if err != nil {
			return nil, err
		}
		if ast.Kind == "up" {
			return pooltypedef.UpSpec{ContExp: cont}, nil
		}
		return pooltypedef.DownSpec{ContExp: cont}, nil
	default:
		return nil, fmt.Errorf("unsupported pool type expression kind %q", ast.Kind)
	}
}

func (a poolExprAST) PlusOrWith() *poolLaborAST {
	if a.Kind == "plus" {
		return a.Plus
	}
	return a.With
}

func (a poolExprAST) Shift() *shiftAST {
	if a.Kind == "up" {
		return a.Up
	}
	return a.Down
}

func procExprToCore(ast procExprAST) (proctypedef.ExpSpec, error) {
	switch ast.Kind {
	case "one":
		return proctypedef.OneSpec{}, nil
	case "link":
		qn, err := uniqsym.ConvertFromString(ast.Link.TypeQN)
		if err != nil {
			return nil, err
		}
		return proctypedef.LinkSpec{TypeQN: qn}, nil
	case "tensor", "lolli":
		prod := ast.TensorOrLolli()
		val, err := procExprToCore(prod.Val)
		if err != nil {
			return nil, err
		}
		cont, err := procExprToCore(prod.Cont)
		if err != nil {
			return nil, err
		}
		if ast.Kind == "tensor" {
			return proctypedef.TensorSpec{Val: val, Cont: cont}, nil
		}
		return proctypedef.LolliSpec{Val: val, Cont: cont}, nil
	case "plus", "with":
		sum := ast.Sum()
		choices := make(map[uniqsym.ADT]proctypedef.ExpSpec, len(sum.Choices))
		for _, choice := range sum.Choices {
			label, err := uniqsym.ConvertFromString(choice.Label)
			if err != nil {
				return nil, err
			}
			cont, err := procExprToCore(choice.Cont)
			if err != nil {
				return nil, err
			}
			choices[label] = cont
		}
		if ast.Kind == "plus" {
			return proctypedef.PlusSpec{Choices: choices}, nil
		}
		return proctypedef.WithSpec{Choices: choices}, nil
	default:
		return nil, fmt.Errorf("unsupported process type expression kind %q", ast.Kind)
	}
}

func (a procExprAST) TensorOrLolli() *procProdAST {
	if a.Kind == "tensor" {
		return a.Tensor
	}
	return a.Lolli
}

func (a procExprAST) Sum() *procSumAST {
	if a.Kind == "plus" {
		return a.Plus
	}
	return a.With
}
