package clause

type CTE struct {
	Alias       string
	Expressions []Expression
	IsRecursive bool
}

type With struct {
	CommonTableExpressions []CTE
}

func (w With) Name() string {
	return "WITH"
}

func (w With) Build(builder Builder) {
	for i, t := range w.CommonTableExpressions {
		if i > 0 {
			builder.WriteString(", ")
		}
		if t.IsRecursive {
			builder.WriteString("RECURSIVE ")
		}
		builder.WriteQuoted(t.Alias)
		builder.WriteString(" AS (")
		for i, expr := range t.Expressions {
			if i > 0 {
				builder.WriteString(" UNION ")
			}
			builder.WriteString("( ")
			expr.Build(builder)
			builder.WriteString(" )")
		}
		builder.WriteString(")")
	}
}

func (w With) MergeClause(clause *Clause) {
	clause.Expression = w
}
