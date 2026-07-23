package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Ledger_Group holds the schema definition for the Ledger_Group entity.
type Ledger_Group struct {
	ent.Schema
}

// Fields of the Ledger_Group.
func (Ledger_Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			Unique(),

		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.Enum("nature").
			Values(
				"ASSET",
				"LIABILITY",
				"EQUITY",
				"INCOME",
				"EXPENSE",
			),

		field.Bool("is_system").
			Default(false),

		field.String("description").
			Optional().
			Default(""),
	}
}

func (Ledger_Group) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Unique(),

		index.Fields("name"),

		index.Fields("nature"),
	}
}

// Edges of the Ledger_Group.
func (Ledger_Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("parent", Ledger_Group.Type).
			Unique().
			From("children"),

		edge.To("ledgers", Ledger.Type),
	}
}
