package sqlbuilder

import (
	"fmt"
	"strings"
)

// Clause represents a SQL where Clause
type Clause interface {
	// interpret interprets Clause into string
	interpret(sb *SelectBuilder) string
	// Not negatives Clause
	Not() *notClause
	// And connects several Clause into an andClause
	And(clause ...Clause) *andClause
	// Or connects several Clause into an orClause
	Or(clause ...Clause) *orClause
}

// newAndClause creates an *andClause
func newAndClause(augend Clause, addend ...Clause) *andClause {
	return &andClause{
		Augend: augend,
		Addend: addend,
	}
}

// andClause represents a SQL AND Clause
type andClause struct {
	Augend Clause
	Addend []Clause
}

func (a *andClause) interpret(sb *SelectBuilder) string {
	andExpr := make([]string, 0, len(a.Addend)+1)
	andExpr = append(andExpr, a.Augend.interpret(sb))
	for _, c := range a.Addend {
		andExpr = append(andExpr, c.interpret(sb))
	}
	return fmt.Sprintf("(%v)", strings.Join(andExpr, " AND "))
}

func (a *andClause) Not() *notClause {
	return newNotClause(a)
}

func (a *andClause) And(clause ...Clause) *andClause {
	return newAndClause(a, clause...)
}

func (a *andClause) Or(clause ...Clause) *orClause {
	return newOrClause(a, clause...)
}

// newOrClause creates an *orClause
func newOrClause(augend Clause, addend ...Clause) *orClause {
	return &orClause{
		Augend: augend,
		Addend: addend,
	}
}

// orClause represents a SQL OR Clause
type orClause struct {
	Augend Clause
	Addend []Clause
}

func (o *orClause) interpret(sb *SelectBuilder) string {
	orExpr := make([]string, 0, len(o.Addend)+1)
	orExpr = append(orExpr, o.Augend.interpret(sb))
	for _, c := range o.Addend {
		orExpr = append(orExpr, c.interpret(sb))
	}
	return fmt.Sprintf("(%v)", strings.Join(orExpr, " OR "))
}


func (o *orClause) Not() *notClause {
	return newNotClause(o)
}

func (o *orClause) And(clause ...Clause) *andClause {
	return newAndClause(o, clause...)
}

func (o *orClause) Or(clause ...Clause) *orClause {
	return newOrClause(o, clause...)
}

// newNotClause creates a notClause
func newNotClause(clause Clause) *notClause {
	return &notClause{
		clause,
	}
}

// notClause represents a SQL NOT Clause
type notClause struct {
	negend Clause
}

func (n *notClause) interpret(sb *SelectBuilder) string {
	return fmt.Sprintf("(NOT %v)", n.negend.interpret(sb))
}

func (n *notClause) Not() *notClause {
	return newNotClause(n)
}

func (n *notClause) And(clause ...Clause) *andClause {
	return newAndClause(n, clause...)
}

func (n *notClause) Or(clause ...Clause) *orClause {
	return newOrClause(n, clause...)
}

// operate interprets basicClause into string
type operate func(b *basicClause) string

// newBasicClause creates a *basicClause
func newBasicClause(field string, operate operate) *basicClause {
	return &basicClause{
		field:   field,
		operate: operate,
	}
}

// basicClause represents a specific basic SQL where Clause
type basicClause struct {
	field   string
	operate operate
	operand []interface{}
	sb      *SelectBuilder
}

func (b *basicClause) interpret(sb *SelectBuilder) string {
	b.sb = sb
	return b.operate(b)
}

func (b *basicClause) Not() *notClause {
	return newNotClause(b)
}

func (b *basicClause) And(clause ...Clause) *andClause {
	return newAndClause(b, clause...)
}

func (b *basicClause) Or(clause ...Clause) *orClause {
	return newOrClause(b, clause...)
}

// SetOperand set basicClause operand
func (b *basicClause) SetOperand(value ...interface{}) *basicClause {
	b.operand = value
	return b
}

// zeroOperandClause is a basicClause which has zero operand
type zeroOperandClause struct {
	*basicClause
}

// SetOperand accepts zero operand
func (z *zeroOperandClause) SetOperand() *zeroOperandClause {
	return z
}

// oneOperandClause is a basicClause which has one operand
type oneOperandClause struct {
	*basicClause
}

// SetOperand accepts one operand
func (o *oneOperandClause) SetOperand(v interface{}) *oneOperandClause {
	o.basicClause.SetOperand(v)
	return o
}

// twoOperandSpec is a basicClause which has two operands
type twoOperandSpec struct {
	*basicClause
}

// SetOperand accepts two operands
func (t *twoOperandSpec) SetOperand(v1, v2 interface{}) *twoOperandSpec {
	t.basicClause.SetOperand(v1, v2)
	return t
}

var (
	isNull operate = func(b *basicClause) string { return b.sb.IsNull(b.field) }

	notNull operate = func(b *basicClause) string { return b.sb.IsNotNull(b.field) }

	e operate = func(b *basicClause) string { return b.sb.E(b.field, b.operand[0]) }

	ne operate = func(b *basicClause) string { return b.sb.NE(b.field, b.operand[0]) }

	g operate = func(b *basicClause) string { return b.sb.G(b.field, b.operand[0]) }

	ge operate = func(b *basicClause) string { return b.sb.GE(b.field, b.operand[0]) }

	l operate = func(b *basicClause) string { return b.sb.L(b.field, b.operand[0]) }

	le operate = func(b *basicClause) string { return b.sb.LE(b.field, b.operand[0]) }

	like operate = func(b *basicClause) string { return b.sb.Like(b.field, b.operand[0]) }

	notLike operate = func(b *basicClause) string { return b.sb.NotLike(b.field, b.operand[0]) }

	between operate = func(b *basicClause) string { return b.sb.Between(b.field, b.operand[0], b.operand[1]) }

	notBetween operate = func(b *basicClause) string { return b.sb.NotBetween(b.field, b.operand[0], b.operand[1]) }

	in operate = func(b *basicClause) string { return b.sb.In(b.field, b.operand...) }

	notIn operate = func(b *basicClause) string { return b.sb.NotIn(b.field, b.operand...) }
)

// NewIsNullClause creates a Clause which represents "field IS NULL"
func NewIsNullClause(field string) *zeroOperandClause {
	return &zeroOperandClause{
		newBasicClause(field, isNull),
	}
}

// NewNotNullClause creates Clause which represents "field IS NOT NULL"
func NewNotNullClause(field string) *zeroOperandClause {
	return &zeroOperandClause{
		newBasicClause(field, notNull),
	}
}

// NewEqualClause creates Clause which represents "field = value"
func NewEqualClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, e),
	}
}

// NewNotEqualClause creates Clause which represents "field != value"
func NewNotEqualClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, ne),
	}
}

// NewGreaterThanClause creates Clause which represents "field > value"
func NewGreaterThanClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, g),
	}
}

// NewGreaterEqualThanClause creates Clause which represents "field >= value"
func NewGreaterEqualThanClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, ge),
	}
}

// NewLessThanClause creates Clause which represents "field < value"
func NewLessThanClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, l),
	}
}

// NewLessEqualThanClause creates Clause which represents "field <= value"
func NewLessEqualThanClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, le),
	}
}

// NewLikeClause creates Clause which represents "field LIKE value"
func NewLikeClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, like),
	}
}

// NewNotLikeClause creates Clause which represents "field NOT LIKE value"
func NewNotLikeClause(field string) *oneOperandClause {
	return &oneOperandClause{
		newBasicClause(field, notLike),
	}
}

// NewBetweenClause creates Clause which represents "field BETWEEN lower AND upper"
func NewBetweenClause(field string) *twoOperandSpec {
	return &twoOperandSpec{
		newBasicClause(field, between),
	}
}

// NewNotBetweenClause creates Clause which represents "field NOT BETWEEN lower AND upper"
func NewNotBetweenClause(field string) *twoOperandSpec {
	return &twoOperandSpec{
		newBasicClause(field, notBetween),
	}
}

// NewInClause creates Clause which represents "field IN (value...)"
func NewInClause(field string) *basicClause {
	return newBasicClause(field, in)
}

// NewNotInClause creates Clause which represents "field NOT IN (value...)"
func NewNotInClause(field string) *basicClause {
	return newBasicClause(field, notIn)
}

// Interpret interprets Clause into string
func Interpret(clause Clause, sb *SelectBuilder) string {
	return clause.interpret(sb)
}
