package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// State holds the schema definition for the State entity.
type State struct {
	ent.Schema
}

// Fields of the State.
func (State) Fields() []ent.Field {
	return []ent.Field{
		// Explicitly define ID as int (optional since it's the default)
		field.Int("id"),

		field.String("name").
			NotEmpty().
			MaxLen(100),

		field.String("code").
			MaxLen(10).
			Optional(),

		field.String("gst_code").
			MaxLen(5).
			Optional(),

		field.Int("country_id"),
	}
}

// Edges of the State.
func (State) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("country", Country.Type).
			Ref("states").
			Field("country_id").
			Unique().
			Required(),
	}
}

func (State) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("country_id"),
		index.Fields("country_id", "name").Unique(),
		index.Fields("gst_code"),
	}
}
