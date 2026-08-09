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

type BillingType string

const (
	BillingNone      BillingType = "NONE"       // no bill tracking (e.g. expense/capital ledgers)
	BillingBillWise  BillingType = "BILL_WISE"  // Tally's "Maintain balances bill-by-bill"
	BillingOnAccount BillingType = "ON_ACCOUNT" // lump-sum, no bill references
)

func (Ledger) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Ledger) Fields() []ent.Field {
	return []ent.Field{

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

		// Accounting
		// field.Float("opening_balance").
		// 	Optional().
		// 	Default(0.00),

		// field.Time("opening_balance_date").
		// 	Optional(),

		// field.Enum("opening_balance_type").
		// 	Values("DR", "CR").
		// 	Optional().
		// 	Default("DR"),

		field.Float("balance").
			Default(0.00).
			Comment("Current balance"),

		// Party classification
		field.Enum("party_type").
			Values("CUSTOMER", "SUPPLIER", "BOTH", "INTERNAL").
			Default("INTERNAL"),

		// Address (Tally: multi-line mailing address)
		field.String("address_line1").Optional().Default("").MaxLen(255),
		field.String("address_line2").Optional().Default("").MaxLen(255),
		field.String("city").Optional().Default("").MaxLen(100),
		field.String("state").Optional().Default("").MaxLen(12),
		// ^ needed to auto-decide CGST+SGST vs IGST (same state vs inter-state)
		field.String("country").Optional().Default("India").MaxLen(100),
		field.String("pincode").Optional().Default("").MaxLen(10),

		// Contact
		field.String("phone").Optional().Default("").MaxLen(20),
		field.String("mobile").Optional().Default("").MaxLen(20),
		field.String("email").Optional().Default("").MaxLen(255),
		field.String("contact_person").Optional().Default("").MaxLen(255),

		// Tax registration (Tally: "Tax Registration Details")
		field.Enum("gst_registration_type").
			Values("REGULAR", "COMPOSITION", "UNREGISTERED", "CONSUMER", "SEZ", "OVERSEAS").
			Default("UNREGISTERED"),
		field.String("gstin").Optional().Default("").MaxLen(15),
		field.String("pan").Optional().Default("").MaxLen(10),

		// Bank details for payments/NEFT (optional but common ask)
		field.String("bank_name").Optional().Default("").MaxLen(255),
		field.String("bank_account_no").Optional().Default("").MaxLen(30),
		field.String("bank_ifsc").Optional().Default("").MaxLen(11),
		field.String("bank_branch").Optional().Default("").MaxLen(255),

		// --- Billing / credit control (Tally: "Maintain balances bill-by-bill") ---
		// field.Enum("billing_type").
		// 	GoType(BillingType("")).
		// 	Default(string(BillingNone)),
		// field.Int("credit_period_days").Optional().Default(0),
		// field.Float("credit_limit").Optional().Default(0),

		// Status
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

		index.Fields("code").
			Unique(),

		index.Fields("name"),

		index.Fields("group_id"),

		index.Fields("gstin"),
		index.Fields("state"),
		index.Fields("party_type"),

		index.Fields("group_id", "name"),
	}
}

func (Ledger) Edges() []ent.Edge {
	return []ent.Edge{

		edge.To("group", Ledger_Group.Type).
			Field("group_id").
			Required().
			Unique(),

		edge.To("journal_lines", Journal_Line.Type),
	}
}
