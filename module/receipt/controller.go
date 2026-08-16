package receipt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
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

func (c *Controller) List(offset, limit int) ([]*ent.Journal, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	journals, err := c.journalService.List(c.ctx.Context, journal.ListFilter{
		Limit:  limit,
		Offset: offset,
		Type:   journal.VoucherRV,
	})
	if err != nil {
		retui.Error(err)
		return nil, err
	}

	retui.Infof("Loaded %d journals", len(journals))
	return journals, nil
}

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
	vIn, err := buildVoucherInput(in)
	if err != nil {
		return nil, err
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

// buildVoucherInput maps the receipt form state into a journal.VoucherInput.
// The receipt account (cash/bank) is debited — money coming in — and each
// party line is credited, representing the source of the money.
func buildVoucherInput(in FormState) (journal.VoucherInput, error) {
	if in.RcptAccount == 0 {
		return journal.VoucherInput{}, fmt.Errorf("receipt account is required")
	}

	date, err := time.Parse("02/01/2006", strings.TrimSpace(in.Date))
	if err != nil {
		return journal.VoucherInput{}, fmt.Errorf("invalid date format, expected DD/MM/YYYY: %w", err)
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
			continue // skip blank trailing rows
		}
		lines = append(lines, journal.LineInput{
			LedgerID:    l.Ledger,
			Credit:      float64(l.Amount),
			Description: l.Remarks,
		})
	}

	return journal.VoucherInput{
		Type:        journal.VoucherRV, // Receipt Voucher
		VoucherNo:   in.VcNo,
		Date:        date,
		VoucherDate: date,
		ReferenceNo: in.Reference,
		Narration:   in.Narration,
		Status:      journal.StatusDraft,
		Lines:       lines,
	}, nil
}
