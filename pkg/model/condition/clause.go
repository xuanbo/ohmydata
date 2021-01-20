package condition

import "errors"

type (
	// Op 操作
	Op uint8
	// Combine 组合
	Combine uint8
)

const (
	// Eq =
	Eq Op = iota
	// NotEq !=
	NotEq
	// Gt >
	Gt
	// Gte >=
	Gte
	// Lt <
	Lt
	// Lte <=
	Lte
	// Like like
	Like
	// NotLike not like
	NotLike
	// In in
	In
	// NotIn not in
	NotIn
	// IsNull is null
	IsNull
	// IsNotNull is not null
	IsNotNull
)

const (
	// And and
	And Combine = iota
	// Or or
	Or
)

var (
	// ErrNameNil is nil
	ErrNameNil = errors.New("name is nil")
	// ErrClauseNil is nil
	ErrClauseNil = errors.New("clause is nil")
)

// Clause 短语
type Clause struct {
	*SingleClause
	*CombineClause
	err error
}

// IsError 是否错误
func (s Clause) IsError() bool {
	return s.err != nil
}

// Error 错误
func (s Clause) Error() error {
	return s.err
}

// SingleClause 单一短语
type SingleClause struct {
	Name  string      `json:"name"`
	Op    Op          `json:"op"`
	Value interface{} `json:"value"`
}

// Eq =
func (Clause) Eq(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Eq, Value: value}}
}

// NotEq !=
func (Clause) NotEq(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: NotEq, Value: value}}
}

// Gt >
func (Clause) Gt(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Gt, Value: value}}
}

// Gte >=
func (Clause) Gte(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Gte, Value: value}}
}

// Lt <
func (Clause) Lt(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Lt, Value: value}}
}

// Lte <=
func (Clause) Lte(name string, value interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Gte, Value: value}}
}

// Like like
func (Clause) Like(name string, value string) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: Like, Value: value}}
}

// NotLike not like
func (Clause) NotLike(name string, value string) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: NotLike, Value: value}}
}

// In in
func (Clause) In(name string, value []interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: In, Value: value}}
}

// NotIn not in
func (Clause) NotIn(name string, value []interface{}) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: NotIn, Value: value}}
}

// IsNull is null
func (Clause) IsNull(name string) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: IsNull}}
}

// IsNotNull is not in
func (Clause) IsNotNull(name string) *Clause {
	if name == "" {
		return &Clause{err: ErrNameNil}
	}
	return &Clause{SingleClause: &SingleClause{Name: name, Op: IsNotNull}}
}

// CombineClause 多短语
type CombineClause struct {
	Combine Combine   `json:"combine"`
	Clauses []*Clause `json:"clauses"`
}

// And and
func (Clause) And(left *Clause, right *Clause, others ...*Clause) *Clause {
	if left == nil || right == nil {
		return &Clause{err: ErrClauseNil}
	}
	clauses := make([]*Clause, 0, 8)
	clauses = append(clauses, left)
	clauses = append(clauses, right)
	clauses = append(clauses, others...)
	return &Clause{CombineClause: &CombineClause{Combine: And, Clauses: clauses}}
}

// Or or
func (CombineClause) Or(left *Clause, right *Clause, others ...*Clause) *Clause {
	if left == nil || right == nil {
		return &Clause{err: ErrClauseNil}
	}
	clauses := make([]*Clause, 0, 8)
	clauses = append(clauses, left)
	clauses = append(clauses, right)
	clauses = append(clauses, others...)
	return &Clause{CombineClause: &CombineClause{Combine: Or, Clauses: clauses}}
}
