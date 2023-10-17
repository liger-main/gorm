package clause

type Values struct {
	Columns []Column
	Values  [][]interface{}
}

// Name from clause name
func (Values) Name() string {
	return "VALUES"
}

// Build build from clause
func (values Values) Build(builder Builder) {
	if len(values.Columns) == 0 && len(values.Values) == 0 {
		builder.WriteString("DEFAULT VALUES")
		return
	}
	wroteColumns := false
	if len(values.Columns) > 0 {
		builder.WriteByte('(')
		for idx, column := range values.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteQuoted(column)
		}
		builder.WriteByte(')')
		wroteColumns = true
	}
	if len(values.Values) > 0 {
		if wroteColumns {
			builder.WriteString(" ")
		}
		builder.WriteString("VALUES ")

		for idx, value := range values.Values {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteByte('(')
			builder.AddVar(builder, value...)
			builder.WriteByte(')')
		}
	}
}

// MergeClause merge values clauses
func (values Values) MergeClause(clause *Clause) {
	clause.Name = ""
	clause.Expression = values
}
