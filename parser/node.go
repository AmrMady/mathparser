package parser

import "math/big"

type NodeType int

const (
	NodeTypeConstant NodeType = iota
	NodeTypeBinaryOp
	NodeTypeSumSeries
	NodeTypeFunctionCall
	NodeTypeUnaryOp
	NodeTypeConditional
	NodeTypeLoop
	NodeTypeFunctionDefinition
	NodeTypeSummation
	NodeTypeVariable
	NodeInfinity
	NodeFactorial
)

type Node interface {
	Type() NodeType
}

type ConstantNode struct {
	Value *big.Float
}

func (n *ConstantNode) Type() NodeType { return NodeTypeConstant }

type BinaryOpNode struct {
	Left  Node
	Op    string
	Right Node
}

func (n *BinaryOpNode) Type() NodeType { return NodeTypeBinaryOp }

type SumSeriesNode struct {
	Term    Node
	VarName string
	Start   Node
	End     Node
}

func (n *SumSeriesNode) Type() NodeType { return NodeTypeSumSeries }

type FunctionCallNode struct {
	Name      string
	Arguments []Node
}

func (n *FunctionCallNode) Type() NodeType { return NodeTypeFunctionCall }

type UnaryOpNode struct {
	Operand Node
	Op      string // For negation, this could be "-"
}

func (n *UnaryOpNode) Type() NodeType { return NodeTypeUnaryOp }

type ConditionalNode struct {
	Condition   Node
	Consequence Node
	Alternative Node
}

func (n *ConditionalNode) Type() NodeType { return NodeTypeConditional }

type LoopNode struct {
	Variable  string
	Start     Node
	End       Node
	Increment Node
	Body      Node
}

func (n *LoopNode) Type() NodeType { return NodeTypeLoop }

type FunctionDefinitionNode struct {
	Name       string
	Parameters []string
	Body       Node
}

func (n *FunctionDefinitionNode) Type() NodeType { return NodeTypeFunctionDefinition }

type SummationNode struct {
	Variable   string
	Start      Node
	End        Node
	Expression Node
}

func (n *SummationNode) Type() NodeType { return NodeTypeSummation }

type VariableNode struct {
	Name string
}

func (n *VariableNode) Type() NodeType { return NodeTypeVariable }

type InfinityNode struct{}

func (n *InfinityNode) Type() NodeType { return NodeInfinity }

type FactorialNode struct {
	Operand Node
}

func (n *FactorialNode) Type() NodeType { return NodeFactorial }
