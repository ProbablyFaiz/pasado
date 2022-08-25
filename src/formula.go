package src

import "fmt"

type FormulaNode interface {
	Eval(ctx *EvalContext) (interface{}, error)
	GetRefs() []*ReferenceNode
	ToFormula() string
}

type LiteralNode struct {
	Value interface{}
}

type ReferenceNode struct {
	Row          int
	Col          int
	Sheet        *Sheet // If nil, then the cell is in the current sheet
	ResolvedUuid cellId // If nil, then the cell has not been resolved yet
}

type FunctionNode struct {
	Name string
	Args []FormulaNode
}

type NilNode struct {
}

type EvalContext struct {
	Cell *Cell
	// TODO: Add more stuff here
}

// GetRefs implementations

func (ln *LiteralNode) GetRefs() []*ReferenceNode {
	return make([]*ReferenceNode, 0)
}

func (rn *ReferenceNode) GetRefs() []*ReferenceNode {
	return []*ReferenceNode{rn}
}

func (fn *FunctionNode) GetRefs() []*ReferenceNode {
	refs := make([]*ReferenceNode, 0)
	for _, arg := range fn.Args {
		refs = append(refs, arg.GetRefs()...)
	}
	return refs
}

func (_ *NilNode) GetRefs() []*ReferenceNode {
	return make([]*ReferenceNode, 0)
}

// ToFormula implementations

func (ln *LiteralNode) ToFormula() string {
	return fmt.Sprintf("%v", ln.Value)
}

func (rn *ReferenceNode) ToFormula() string {
	return getA1Notation(rn.Row, rn.Col)
}

func (fn *FunctionNode) ToFormula() string {
	formula := fn.Name + "("
	for i, arg := range fn.Args {
		formula += arg.ToFormula()
		if i != len(fn.Args)-1 {
			formula += ", "
		}
	}
	formula += ")"
	return formula
}

func (_ *NilNode) ToFormula() string {
	return ""
}

// Eval implementations in src/eval.go