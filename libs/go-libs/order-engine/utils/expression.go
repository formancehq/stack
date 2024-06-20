package utils

type LogicalOperator string




const (
    And LogicalOperator = "AND"
    Or  LogicalOperator = "OR"
    Not LogicalOperator = "NOT"
    Eq  LogicalOperator = "="
    Ne  LogicalOperator = "!="
    Gt  LogicalOperator = ">"
    Lt  LogicalOperator = "<"
    Ge  LogicalOperator = ">="
    Le  LogicalOperator = "<="
)

type NumericalOperator string

const (
    Add NumericalOperator = "+"
    Sub NumericalOperator = "-"
    Mul NumericalOperator = "*"
    Div NumericalOperator = "/"
    Min NumericalOperator = "MIN"
    Max NumericalOperator = "MAX"
    Per NumericalOperator = "PER"
	
)

const (
	Balance NumericalOperator = "BALANCE"
	AggBalance NumericalOperator = "AGGBALANCE"
	MDAggBalance NumericalOperator = "MDAGGBALANCE"
)

var Operators = struct {
    And LogicalOperator
    Or LogicalOperator
    Not LogicalOperator
    Eq LogicalOperator
    Ne LogicalOperator
    Gt  LogicalOperator 
    Lt  LogicalOperator 
    Ge  LogicalOperator 
    Le  LogicalOperator 

    Add NumericalOperator
    Sub NumericalOperator 
    Mul NumericalOperator 
    Div NumericalOperator 
    Min NumericalOperator 
    Max NumericalOperator 
    Per NumericalOperator

    Balance NumericalOperator 
	AggBalance NumericalOperator 
	MDAggBalance NumericalOperator 
    
} {
    And : And,
    Or:Or,
    Not: Not,
    Ne:Ne,
    Gt:Gt,
    Lt:Lt,
    Ge : Ge,
    Le : Le,

    Add : Add,
    Sub : Sub, 
    Mul : Mul,
    Div : Div, 
    Min : Min, 
    Max : Max,
    Per : Per,
    
    Balance:Balance,
	AggBalance: AggBalance,
	MDAggBalance: MDAggBalance, 
    
}





type Operand interface{}


type LogicalExpression struct {
    Operator LogicalOperator
    Operands []Operand
}

type ValueExpression struct {
    Operator NumericalOperator
    Operands []Operand
}

func (v ValueExpression) New(operator NumericalOperator, operands ...Operand) *ValueExpression{
    return &ValueExpression{Operator: operator, Operands: operands}
}

func (v LogicalExpression) New(operator LogicalOperator, operands ...Operand) *LogicalExpression{
    return &LogicalExpression{Operator: operator, Operands: operands}
}

func NewLogicalExpression(operator LogicalOperator, operands ...Operand) LogicalExpression {
    return LogicalExpression{Operator: operator, Operands: operands}
}


func NewValueExpression(operator NumericalOperator, operands ...Operand) ValueExpression {
    return ValueExpression{Operator: operator, Operands: operands}
}
