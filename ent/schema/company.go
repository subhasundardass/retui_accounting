package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Company struct {
	ent.Schema
}

func (Company) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(150),

		field.String("code").
			Unique().
			NotEmpty().
			MaxLen(20),

		field.String("legal_name").
			Optional().
			MaxLen(200),

		field.String("email").
			Optional(),

		field.String("phone").
			Optional().
			MaxLen(30),

		field.String("website").
			Optional(),

		field.String("tax_id").
			Optional().
			MaxLen(50),

		field.String("gstin").
			Optional().
			MaxLen(20),

		field.String("pan").
			Optional().
			MaxLen(20),

		field.String("currency").
			Default("INR"),

		field.String("timezone").
			Default("Asia/Kolkata"),

		field.String("logo").
			Optional(),

		field.Text("address").
			Optional(),

		field.String("city").
			Optional(),

		field.String("state").
			Optional(),

		field.String("country").
			Default("India"),

		field.String("postal_code").
			Optional(),

		field.Bool("active").
			Default(true),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Company) Edges() []ent.Edge {
	return nil
}
