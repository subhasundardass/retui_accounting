package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// PartyMaster holds the schema definition for the PartyMaster entity.
type PartyMaster struct {
	ent.Schema
}

func (PartyMaster) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Annotations of the User.
func (PartyMaster) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "party_master"},
	}
}

// Fields of the PartyMaster.
func (PartyMaster) Fields() []ent.Field {
	return []ent.Field{
		field.Int("ledger_id").Unique(),

		field.Enum("type").
			Values("CUSTOMER", "SUPPLIER", "BOTH"),

		field.String("display_name").
			NotEmpty(),

		field.String("legal_name").
			Optional().
			Nillable(),

		// Tax registration (Tally: "Tax Registration Details")
		field.Enum("gst_registration_type").
			Values("REGULAR", "COMPOSITION", "UNREGISTERED", "CONSUMER", "SEZ", "OVERSEAS").
			Default("UNREGISTERED"),
		field.String("gstin").Optional().Default("").MaxLen(15),
		field.String("pan").Optional().Default("").MaxLen(10),

		// needed for the 45-day MSME payment rule — worth adding proactively)
		field.Bool("is_msme").Default(false),
		field.String("msme_number").Optional().Default("").MaxLen(30),

		// Bank details for payments/NEFT (optional but common ask)
		field.String("bank_name").Optional().Default("").MaxLen(255),
		field.String("bank_account_no").Optional().Default("").MaxLen(30),
		field.String("bank_ifsc").Optional().Default("").MaxLen(11),
		field.String("bank_branch").Optional().Default("").MaxLen(255),

		field.String("notes").Optional().Default(""),

		// Communication
		field.String("contact_person").
			Optional().
			Nillable(),

		field.String("mobile").
			Optional().
			Nillable(),

		field.String("phone").
			Optional().
			Nillable(),

		field.String("email").
			Optional().
			Nillable(),

		field.String("website").
			Optional().
			Nillable(),

		field.Float("credit_limit").
			Default(0),

		field.Int("credit_days").
			Default(0),

		field.Float("opening_balance").
			Default(0),

		field.String("address").
			Optional().
			Nillable(),

		field.String("city").
			Optional().
			Nillable(),

		field.String("state").
			Optional().
			Nillable(),

		field.String("country").
			Default("India"),

		field.String("pincode").
			Optional().
			Nillable(),
	}
}

func (PartyMaster) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ledger_id").
			Unique(),

		index.Fields("gstin"),
		index.Fields("state_code"),
		index.Fields("party_type"),
	}
}

// Edges of the PartyMaster.
func (PartyMaster) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ledger", Ledger.Type).
			Ref("party").
			Unique(),
	}
}
