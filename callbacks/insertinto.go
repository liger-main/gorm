package callbacks

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertInto(config *Config) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Error != nil {
			return
		}

		stmt := db.Statement

		clauseSelect := clause.Select{
			Distinct: db.Statement.Distinct,
			Columns:  make([]clause.Column, len(db.Statement.Selects)),
		}
		for idx, name := range db.Statement.Selects {
			if db.Statement.Schema == nil {
				clauseSelect.Columns[idx] = clause.Column{Name: name, Raw: true}
			} else if f := db.Statement.Schema.LookUpField(name); f != nil {
				clauseSelect.Columns[idx] = clause.Column{Name: f.DBName}
			} else {
				clauseSelect.Columns[idx] = clause.Column{Name: name, Raw: true}
			}
		}
		stmt.AddClauseIfNotExists(clauseSelect)

		stmt.AddClauseIfNotExists(clause.From{
			TableExpr: stmt.TableExpr,
			Tables:    []clause.Table{{Name: stmt.Table}},
		})

		stmt.SQL.Grow(100)
		stmt.Build(stmt.BuildClauses...)

		if !db.DryRun && db.Error == nil {
			result, err := stmt.ConnPool.ExecContext(stmt.Context, stmt.SQL.String(), stmt.Vars...)
			if db.AddError(err) == nil {
				db.RowsAffected, _ = result.RowsAffected()
			}

			return
		}
	}
}
