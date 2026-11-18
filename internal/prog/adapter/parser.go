package adapter

import (
	"fmt"
	"slices"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"orglang/go-engine/internal/adt/uniqsym"
	coreprog "orglang/go-engine/internal/prog/core"
	pooltypedef "orglang/go-engine/internal/pool/typedef/core"
	pooltypeexp "orglang/go-engine/internal/pool/typeexp/core"
	proctypedef "orglang/go-engine/internal/proc/typedef/core"
	proctypeexp "orglang/go-engine/internal/proc/typeexp/core"
)

type programExpression interface { programExpression() }

type programAST struct {
	Expressions []programExpression `parser:"@@*"`
}

type poolDefAST struct {
	TypeQN string       `parser:"'pool' 'type' @Ident"`
	TypeExp poolExpAST  `parser:"'{' @@? '}'"`
}
func (poolDefAST) programExpression() {}

type procDefAST struct {
	TypeQN string      `parser:"'proc' 'type' @Ident"`
	TypeExp procExpAST `parser:"'{' @@? '}'"`
}
func (procDefAST) programExpression() {}

type expKind string
func (k *expKind) Parse(lex *lexer.PeekingLexer) error {
	value := expKind(lex.Peek().Value)
	if slices.Contains([]expKind{"one", "link", "plus", "with", "tensor", "lolli", "up", "down"}, value) {
		*k = value
	}
	return nil
}

type poolExpAST struct {
	Kind expKind       `parser:"@@"`
	Link *linkAST      `parser:"( 'link' @@"`
	Plus *poolLaborAST `parser:"| 'plus' @@"`
	With *poolLaborAST `parser:"| 'with' @@"`
	Tensor *poolResAST `parser:"| 'tensor' @@"`
	Lolli *poolResAST  `parser:"| 'lolli' @@"`
	Up *shiftAST       `parser:"| 'up' @@"`
	Down *shiftAST     `parser:"| 'down' @@ )"`
}
type procExpAST struct {
	Kind expKind       `parser:"@@"`
	Link *linkAST      `parser:"( 'link' @@"`
	Tensor *procProdAST `parser:"| 'tensor' @@"`
	Lolli *procProdAST  `parser:"| 'lolli' @@"`
	Plus *procSumAST    `parser:"| 'plus' @@"`
	With *procSumAST    `parser:"| 'with' @@"`
	Up *shiftAST        `parser:"| 'up' @@"`
	Down *shiftAST      `parser:"| 'down' @@ )"`
}
type linkAST struct { TypeQN string `parser:"@Ident"` }
type poolLaborAST struct {
	ProcQNs []string  `parser:"'(' @Ident* ')'"`
	ContExp poolExpAST `parser:"'{' @@? '}'"`
}
type poolResAST struct {
	ValExp poolExpAST `parser:""`
	ContExp poolExpAST `parser:"'{' @@? '}'"`
}
type shiftAST struct { Cont poolExpAST `parser:"'{' @@? '}'"` }
type procProdAST struct {
	Val procExpAST `parser:"'{' @@? '}'"`
	Cont procExpAST `parser:"'{' @@? '}'"`
}
type procSumAST struct { Choices []procChoiceAST `parser:"'{' @@* '}'"` }
type procChoiceAST struct {
	Label string `parser:"'@' @Ident"`
	Cont procExpAST `parser:"'{' @@? '}'"`
}

func parseProgram(text string) (programAST, error) {
	parser, err := participle.Build[programAST](participle.Union[programExpression](poolDefAST{}, procDefAST{}))
	if err != nil { return programAST{}, err }
	ast, err := parser.ParseString("", text)
	if err != nil { return programAST{}, err }
	return *ast, nil
}

func programASTToCore(ast programAST) (coreprog.Spec, error) {
	var result coreprog.Spec
	for _, expression := range ast.Expressions {
		switch expression := expression.(type) {
		case poolDefAST:
			qn, err := uniqsym.ConvertFromString(expression.TypeQN); if err != nil { return coreprog.Spec{}, err }
			exp, err := poolExpToCore(expression.TypeExp); if err != nil { return coreprog.Spec{}, err }
			result.Pools = append(result.Pools, pooltypedef.DefSpec{TypeQN: qn, TypeExp: exp})
		case procDefAST:
			qn, err := uniqsym.ConvertFromString(expression.TypeQN); if err != nil { return coreprog.Spec{}, err }
			exp, err := procExpToCore(expression.TypeExp); if err != nil { return coreSpec{}, err }
			result.Procs = append(result.Procs, proctypedef.DefSpec{TypeQN: qn, TypeExp: exp})
		default:
			return coreprog.Spec{}, fmt.Errorf("unsupported program expression: %T", expression)
		}
	}
	return result, nil
}

func poolExpToCore(ast poolExpAST) (pooltypeexp.ExpSpec, error) {
	switch ast.Kind {
	case "one": return pooltypeexp.OneSpec{}, nil
	case "link":
		qn, err := uniqsym.ConvertFromString(ast.Link.TypeQN); if err != nil { return nil, err }
		return pooltypeexp.LinkSpec{TypeQN: qn}, nil
	case "plus", "with":
		qns := make([]uniqsym.ADT, len(ast.PlusOrWith().ProcQNs))
		for i, value := range ast.PlusOrWith().ProcQNs { qn, err := uniqsym.ConvertFromString(value); if err != nil { return nil, err }; qns[i] = qn }
		cont, err := poolExpToCore(ast.PlusOrWith().ContExp); if err != nil { return nil, err }
		if ast.Kind == "plus" { return pooltypeexp.PlusSpec{ProcQNs:qns, ContExp:cont}, nil }
		return pooltypeexp.WithSpec{ProcQNs:qns, ContExp:cont}, nil
	case "up", "down":
		cont, err := poolExpToCore(ast.Shift().Cont); if err != nil { return nil, err }
		if ast.Kind == "up" { return pooltypeexp.UpSpec{ContExp:cont}, nil }
		return pooltypeexp.DownSpec{ContExp:cont}, nil
	default: return nil, fmt.Errorf("unsupported pool type expression kind %q", ast.Kind)
	}
}
func (a poolExpAST) PlusOrWith() *poolLaborAST { if a.Kind == "plus" { return a.Plus }; return a.With }
func (a poolExpAST) Shift() *shiftAST { if a.Kind == "up" { return a.Up }; return a.Down }

func procExpToCore(ast procExpAST) (proctypeexp.ExpSpec, error) {
	switch ast.Kind {
	case "one": return proctypeexp.OneSpec{}, nil
	case "link":
		qn, err := uniqsym.ConvertFromString(ast.Link.TypeQN); if err != nil { return nil, err }
		return proctypeexp.LinkSpec{TypeQN: qn}, nil
	case "tensor", "lolli":
		val, err := procExpToCore(ast.TensorOrLolli().Val); if err != nil { return nil, err }
		cont, err := procExpToCore(ast.TensorOrLolli().Cont); if err != nil { return nil, err }
		if ast.Kind == "tensor" { return proctypeexp.TensorSpec{Val:val, Cont:cont}, nil }
		return proctypeexp.LolliSpec{Val:val, Cont:cont}, nil
	case "plus", "with":
		choices := make(map[uniqsym.ADT]proctypeexp.ExpSpec, len(ast.Sum().Choices))
		for _, choice := range ast.Sum().Choices {
			label, err := uniqsym.ConvertFromString(choice.Label); if err != nil { return nil, err }
			cont, err := procExpToCore(choice.Cont); if err != nil { return nil, err }
			if _, exists := choices[label]; exists { return nil, fmt.Errorf("duplicate process type label %q", choice.Label) }
			choices[label] = cont
		}
		if ast.Kind == "plus" { return proctypeexp.PlusSpec{Choices:choices}, nil }
		return proctypeexp.WithSpec{Choices:choices}, nil
	default: return nil, fmt.Errorf("unsupported process type expression kind %q", ast.Kind)
	}
}
func (a procExpAST) TensorOrLolli() *procProdAST { if a.Kind == "tensor" { return a.Tensor }; return a.Lolli }
func (a procExpAST) Sum() *procSumAST { if a.Kind == "plus" { return a.Plus }; return a.With }
