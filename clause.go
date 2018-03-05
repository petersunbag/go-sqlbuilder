package sqlbuilder

// Clause represents a specific basic SQL where Clause
type Clause struct {
	Builder
}

func (c *Clause) Not() *Clause {
	c.Builder = Build("NOT $?", c.Builder)
	return c
}

func (c *Clause) And(others ...*Clause) *Clause {
	builder := c.Builder
	for _, o := range others {
		builder = Build("$? AND $?", builder, o.Builder)
	}

	c.Builder = Build("($?)", builder)
	return c
}

func (c *Clause) Or(others ...*Clause) *Clause {
	builder := c.Builder
	for _, o := range others {
		builder = Build("$? OR $?", builder, o.Builder)
	}

	c.Builder = Build("($?)", builder)
	return c
}

// Interpret interprets Clause into string
func Interpret(clause *Clause, sb *SelectBuilder) string {
	sql, _ := clause.Build()
	sb.Where(sb.Var(clause.Builder))
	return sql
}

// operate interprets Clause into string
type operate func(field string, value ...interface{}) string

// newOperation creates an *operation
func newOperation(field string, operate operate) *operation {
	return &operation{
		field,
		operate,
	}
}

// operation stores field and operate of clause
type operation struct {
	field   string
	operate operate
}

// NewClause creates *Clause with operand value
func (o *operation) NewClause(value ...interface{}) *Clause {
	builder := Build(o.operate(o.field, value...), value...)
	return &Clause{builder}
}

// newZeroOperation creates a *zeroOperandOperation
func newZeroOperation(field string, operate operate) *zeroOperandOperation {
	return &zeroOperandOperation{
		newOperation(field, operate),
	}
}

// zeroOperandOperation can create *Clause with zero operand
type zeroOperandOperation struct {
	*operation
}

// NewClause creates *Clause with zero operand
func (z *zeroOperandOperation) NewClause() *Clause {
	return z.operation.NewClause()
}

// newOneOperandOperation creates a *oneOperandOperation
func newOneOperandOperation(field string, operate operate) *oneOperandOperation {
	return &oneOperandOperation{
		newOperation(field, operate),
	}
}

// oneOperandOperation can create *Clause with one operand
type oneOperandOperation struct {
	*operation
}

// NewClause creates *Clause with one operand v
func (o *oneOperandOperation) NewClause(v interface{}) *Clause {
	return o.operation.NewClause(v)
}

// newTwoOperandOperation creates a *twoOperandOperation
func newTwoOperandOperation(field string, operate operate) *twoOperandOperation {
	return &twoOperandOperation{
		newOperation(field, operate),
	}
}

// twoOperandOperation can create *Clause with two operand
type twoOperandOperation struct {
	*operation
}

// NewClause creates *Clause with operand v1, v2
func (t *twoOperandOperation) NewClause(v1, v2 interface{}) *Clause {
	return t.operation.NewClause(v1, v2)
}

var (
	isNull operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).IsNull(field)
	}

	notNull operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).IsNotNull(field)
	}

	e operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).E(field, value[0])
	}

	ne operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).NE(field, value[0])
	}

	g operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).G(field, value[0])
	}

	ge operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).GE(field, value[0])
	}

	l operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).L(field, value[0])
	}

	le operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).LE(field, value[0])
	}

	like operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).Like(field, value[0])
	}

	notLike operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).NotLike(field, value[0])
	}

	between operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).Between(field, value[0], value[1])
	}

	notBetween operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).NotBetween(field, value[0], value[1])
	}

	in operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).In(field, value...)
	}

	notIn operate = func(field string, value ...interface{}) string {
		return (&Cond{&Args{}}).NotIn(field, value...)
	}
)

// NewIsNullOperation creates a operation which can create Clause that represents "field IS NULL"
func NewIsNullOperation(field string) *zeroOperandOperation {
	return newZeroOperation(field, isNull)
}

// NewNotNullOperation creates operation which can create Clause that represents "field IS NOT NULL"
func NewNotNullOperation(field string) *zeroOperandOperation {
	return newZeroOperation(field, notNull)
}

// NewEqualOperation creates operation which can create Clause that represents "field = value"
func NewEqualOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, e)
}

// NewNotEqualOperation creates operation which can create Clause that represents "field != value"
func NewNotEqualOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, ne)
}

// NewGreaterThanOperation creates operation which can create Clause that represents "field > value"
func NewGreaterThanOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, g)
}

// NewGreaterEqualThanOperation creates operation which can create Clause that represents "field >= value"
func NewGreaterEqualThanOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, ge)
}

// NewLessThanOperation creates operation which can create Clause that represents "field < value"
func NewLessThanOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, l)
}

// NewLessEqualThanOperation creates operation which can create Clause that represents "field <= value"
func NewLessEqualThanOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, le)
}

// NewLikeOperation creates operation which can create Clause that represents "field LIKE value"
func NewLikeOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, like)
}

// NewNotLikeOperation creates operation which can create Clause that represents "field NOT LIKE value"
func NewNotLikeOperation(field string) *oneOperandOperation {
	return newOneOperandOperation(field, notLike)
}

// NewBetweenOperation creates operation which can create Clause that represents "field BETWEEN lower AND upper"
func NewBetweenOperation(field string) *twoOperandOperation {
	return newTwoOperandOperation(field, between)
}

// NewNotBetweenOperation creates operation which can create Clause that represents "field NOT BETWEEN lower AND upper"
func NewNotBetweenOperation(field string) *twoOperandOperation {
	return newTwoOperandOperation(field, notBetween)
}

// NewInOperation creates operation which can create Clause that represents "field IN (value...)"
func NewInOperation(field string) *operation {
	return newOperation(field, in)
}

// NewNotInOperation creates operation which can create Clause that represents "field NOT IN (value...)"
func NewNotInOperation(field string) *operation {
	return newOperation(field, notIn)
}
