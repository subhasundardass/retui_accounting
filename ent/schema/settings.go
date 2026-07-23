package schema

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Settings holds the schema definition for the Settings entity.
type Settings struct {
	ent.Schema
}

// Mixin of the Settings.
func (Settings) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Settings.
func (Settings) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable(),

		// Key field with constraints
		field.String("key").
			Unique().
			MaxLen(100).
			MinLen(1).
			Match(regexp.MustCompile(`^[a-zA-Z0-9_\.\-]+$`)).
			Comment("Settings key (unique identifier)"),

		// Value field - flexible for different types
		field.String("value").
			MaxLen(5000).
			Comment("Settings value"),
	}
}

// Indexes of the Settings.
func (Settings) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint on key
		index.Fields("key").
			Unique(),
	}
}

// Edges of the Settings.
func (Settings) Edges() []ent.Edge {
	return nil
}
