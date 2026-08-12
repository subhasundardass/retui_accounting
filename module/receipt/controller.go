package receipt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui/components"
)

type Controller struct {
	ctx            *appctx.AppContext
	repo           *Repository
	ledgerService  ledger.Service
	journalService journal.Service
}

func NewController(ctx *appctx.AppContext) *Controller {
	ledgerRepo := ledger.NewRepository(ctx.DB.Client)
	ledgerService := ledger.NewService(ledgerRepo)

	journalRepo := journal.NewRepository(ctx.DB.Client)
	journalService := journal.NewService(journalRepo)

	controller := &Controller{
		ctx:            ctx,
		repo:           NewRepository(ctx.DB.Client),
		ledgerService:  ledgerService,
		journalService: journalService,
	}

	return controller
}

// func (c *Controller) loadCashBankDropdown() error {
// 	ledgers, err := c.ledgerService.GetCashBankLedgers(c.ctx.Context)
// 	if err != nil {
// 		return fmt.Errorf("failed to load cash/bank ledgers: %w", err)
// 	}

// 	options := make([]components.SelectOption, 0, len(ledgers))

// 	for _, item := range ledgers {
// 		options = append(options, components.SelectOption{
// 			Label: item.Name,
// 			Value: strconv.Itoa(item.ID),
// 		})
// 	}

// 	c.cashBankOptions = options

// 	return options
// }

func (c *Controller) CashBankOptions() []components.SelectOption {
	ledgers, err := c.ledgerService.GetCashBankLedgers(c.ctx.Context)
	if err != nil {
		return nil
	}

	options := make([]components.SelectOption, 0, len(ledgers))
	for _, item := range ledgers {
		options = append(options, components.SelectOption{
			Label: item.Name,
			Value: strconv.Itoa(item.ID),
		})
	}

	// c.cashBankOptions = options
	return options
}

// --Save
func (c *Controller) Save(mode FormMode, id int, in FormState) (*ent.Journal, error) {
	vType := journal.VoucherJV

	date, err := time.Parse("02/01/2006", in.Date) // matches your DD/MM/YYYY format
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	lines := make([]journal.LineInput, 0, len(in.Lines)+1)

	// Receipt account (cash/bank) is debited — money coming in.
	lines = append(lines, journal.LineInput{
		LedgerID: in.RcptAccount,
		Debit:    in.Amount,
	})

	// Each party line is credited — the source of the money.
	for _, l := range in.Lines {
		if l.Ledger == 0 {
			continue
		}
		lines = append(lines, journal.LineInput{
			LedgerID:    l.Ledger,
			Credit:      float64(l.Amount),
			Description: l.Remarks,
		})
	}

	vIn := journal.VoucherInput{
		Type:        vType,
		VoucherNo:   in.VcNo,
		ReferenceNo: in.Reference,
		Date:        date,
		Narration:   in.Narration,
		Lines:       lines,
	}

	jMode := journal.ModeCreate
	if mode == ModeUpdate {
		jMode = journal.ModeUpdate
	}

	return c.journalService.Save(c.ctx.Ctx(), jMode, id, vIn)
}

// ValidateForm checks the receipt form for basic consistency before
// it's persisted. It does not touch the database.
func ValidateForm(in FormState) error {
	if in.RcptAccount == 0 {
		return fmt.Errorf("receipt account is required")
	}
	if in.Date == "" {
		return fmt.Errorf("date is required")
	}
	if in.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	if len(in.Lines) == 0 {
		return fmt.Errorf("at least one party line is required")
	}

	var lineTotal float64
	activeLines := 0
	for i, l := range in.Lines {
		if l.Ledger == 0 {
			continue // skip blank trailing rows
		}
		if l.Amount <= 0 {
			return fmt.Errorf("line %d: amount must be greater than zero", i+1)
		}
		lineTotal += float64(l.Amount)
		activeLines++
	}

	if activeLines == 0 {
		return fmt.Errorf("at least one party line must have a ledger selected")
	}

	if !amountsEqual(in.Amount, lineTotal) {
		return fmt.Errorf(
			"header amount (%.2f) does not match total of party lines (%.2f)",
			in.Amount, lineTotal,
		)
	}

	return nil
}

// amountsEqual compares two money values with a small epsilon to avoid
// float rounding false-negatives.
func amountsEqual(a, b float64) bool {
	const epsilon = 0.005
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

// buildJournalInput maps the UI form state into whatever shape the
// repository layer expects for persistence. Adjust field names to
// match your actual ent.JournalCreateInput / repo signature.
func buildJournalInput(in FormState) ent.Journal {
	lines := make([]ent.Journal, 0, len(in.Lines))
	for _, l := range in.Lines {
		if l.Ledger == 0 {
			continue
		}
		lines = append(lines, ent.Journal{
			// LedgerID: l.Ledger,
			// Amount:   float64(l.Amount),
			// Remarks:  l.Remarks,
		})
	}

	return ent.Journal{
		// VcNo:       in.VcNo,
		// Reference:  in.Reference,
		// Date:       in.Date,
		// Amount:     in.Amount,
		// RcptLedger: in.RcptAccount,
		// Narration:  in.Narration,
		// Lines:      lines,
	}
}
