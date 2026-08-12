package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Ledger holds the schema definition for the Ledger entity.
type Ledger struct {
	ent.Schema
}

func (Ledger) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Ledger) Fields() []ent.Field {
	return []ent.Field{
		// Explicit ID field (optional)
		field.Int("id"),

		field.Int("group_id"),

		field.String("code").
			NotEmpty().
			Unique(),

		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.String("alias").
			Optional().
			Default("").
			MaxLen(255),

		field.String("description").
			Optional().
			Default(""),

		field.Float("balance").
			Default(0.00).
			Comment("Current balance"),

		field.Enum("party_type").
			Values("CUSTOMER", "SUPPLIER", "BOTH", "INTERNAL").
			Default("INTERNAL"),

		field.String("address_line1").Optional().Default("").MaxLen(255),
		field.String("address_line2").Optional().Default("").MaxLen(255),
		field.String("city").Optional().Default("").MaxLen(100),

		// Foreign keys as int to match State and Country IDs
		field.Int("state_id").Optional(),
		field.Int("country_id").Optional(),

		field.String("pincode").Optional().Default("").MaxLen(10),

		field.String("phone").Optional().Default("").MaxLen(20),
		field.String("mobile").Optional().Default("").MaxLen(20),
		field.String("email").Optional().Default("").MaxLen(255),
		field.String("contact_person").Optional().Default("").MaxLen(255),

		field.Enum("gst_registration_type").
			Values("REGULAR", "COMPOSITION", "UNREGISTERED", "CONSUMER", "SEZ", "OVERSEAS").
			Default("UNREGISTERED"),
		field.String("gstin").Optional().Default("").MaxLen(15),
		field.String("pan").Optional().Default("").MaxLen(10),

		field.String("bank_name").Optional().Default("").MaxLen(255),
		field.String("bank_account_no").Optional().Default("").MaxLen(30),
		field.String("bank_ifsc").Optional().Default("").MaxLen(11),
		field.String("bank_branch").Optional().Default("").MaxLen(255),

		field.Bool("is_system").
			Default(false).
			Comment("Built-in system ledger"),

		field.Bool("is_party").
			Default(false),

		field.Bool("is_bank").
			Default(false),

		field.Bool("is_cash").
			Default(false),

		field.Bool("is_active").
			Default(true),
	}
}

func (Ledger) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
		index.Fields("name"),
		index.Fields("group_id"),
		index.Fields("gstin"),
		index.Fields("state_id"),
		index.Fields("party_type"),
		index.Fields("group_id", "name"),
		index.Fields("country_id"),
		index.Fields("state_id", "country_id"),
	}
}

func (Ledger) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Ledger_Group.Type).
			Field("group_id").
			Required().
			Unique(),

		edge.To("state", State.Type).
			Field("state_id").
			Unique(),

		edge.To("country", Country.Type).
			Field("country_id").
			Unique(),

		edge.To("journal_lines", Journal_Line.Type),
	}
}
