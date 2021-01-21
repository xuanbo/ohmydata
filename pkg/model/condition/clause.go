package condition

type (
	// Op 操作
	Op uint8
	// Combine 组合
	Combine uint8
)

const (
	// OpEq =
	OpEq Op = iota
	// OpNotEq !=
	OpNotEq
	// OpGt >
	OpGt
	// OpGte >=
	OpGte
	// OpLt <
	OpLt
	// OpLte <=
	OpLte
	// OpLike like
	OpLike
	// OpNotLike not like
	OpNotLike
	// OpIn in
	OpIn
	// OpNotIn not in
	OpNotIn
	// OpIsNull is null
	OpIsNull
	// OpIsNotNull is not null
	OpIsNotNull
)

const (
	// CombineAnd and
	CombineAnd Combine = iota
	// CombineOr or
	CombineOr
)

var (
	emtpy = new(Clause)
)

// Clause 短语
type Clause struct {
	*SingleClause
	*CombineClause
	err error
}

// IsEmpty 是否为空
func (c *Clause) IsEmpty() bool {
	if c == emtpy {
		return true
	}
	return c.SingleClause == nil && c.CombineClause == nil
}

// WrapSingleClause 包装
func WrapSingleClause(clause *SingleClause) *Clause {
	if clause == nil {
		return emtpy
	}
	return &Clause{SingleClause: clause}
}

// WrapCombineClause 包装
func WrapCombineClause(clause *CombineClause) *Clause {
	if clause == nil || len(clause.Clauses) == 0 {
		return emtpy
	}
	return &Clause{CombineClause: clause}
}

// SingleClause 单一短语
type SingleClause struct {
	Name  string      `json:"name"`
	Op    Op          `json:"op"`
	Value interface{} `json:"value"`
}

// Eq =
func Eq(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpEq, Value: value}}
}

// NotEq !=
func NotEq(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpNotEq, Value: value}}
}

// Gt >
func Gt(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpGt, Value: value}}
}

// Gte >=
func Gte(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpGte, Value: value}}
}

// Lt <
func Lt(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpLt, Value: value}}
}

// Lte <=
func Lte(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpGte, Value: value}}
}

// Like like
func Like(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpLike, Value: value}}
}

// NotLike not like
func NotLike(name string, value interface{}) *Clause {
	if name == "" || value == nil {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpNotLike, Value: value}}
}

// In in
func In(name string, value []interface{}) *Clause {
	if name == "" || len(value) == 0 {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpIn, Value: value}}
}

// NotIn not in
func NotIn(name string, value []interface{}) *Clause {
	if name == "" || len(value) == 0 {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpNotIn, Value: value}}
}

// IsNull is null
func IsNull(name string) *Clause {
	if name == "" {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpIsNull}}
}

// IsNotNull is not in
func IsNotNull(name string) *Clause {
	if name == "" {
		return emtpy
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: OpIsNotNull}}
}

// CombineClause 多短语
type CombineClause struct {
	Combine Combine   `json:"combine"`
	Clauses []*Clause `json:"clauses"`
}

// NewCombineClause 创建
func NewCombineClause(combine Combine) *CombineClause {
	return &CombineClause{Combine: combine, Clauses: make([]*Clause, 0, 8)}
}

// Add add clause
func (c *CombineClause) Add(clause *Clause, others ...*Clause) {
	var clauses *Clause
	if c.Combine == CombineAnd {
		clauses = And(clause, others...)
	} else if c.Combine == CombineOr {
		clauses = Or(clause, others...)
	} else {
		clauses = emtpy
	}
	if clauses.IsEmpty() {
		return
	}
	c.Clauses = append(c.Clauses, clauses.Clauses...)
}

// And and
func And(left *Clause, rights ...*Clause) *Clause {
	clauses := make([]*Clause, 0, 8)
	if !(left == nil || left.IsEmpty()) {
		clauses = append(clauses, left)
	}
	for _, clause := range rights {
		if clause == nil || clause.IsEmpty() {
			continue
		}
		clauses = append(clauses, clause)
	}
	if len(clauses) == 0 {
		return emtpy
	}
	return &Clause{CombineClause: &CombineClause{Combine: CombineAnd, Clauses: clauses}}
}

// Or or
func Or(left *Clause, rights ...*Clause) *Clause {
	clauses := make([]*Clause, 0, 8)
	if !(left == nil || left.IsEmpty()) {
		clauses = append(clauses, left)
	}
	for _, clause := range rights {
		if clause == nil || clause.IsEmpty() {
			continue
		}
		clauses = append(clauses, clause)
	}
	if len(clauses) == 0 {
		return emtpy
	}
	return &Clause{CombineClause: &CombineClause{Combine: CombineOr, Clauses: clauses}}
}
