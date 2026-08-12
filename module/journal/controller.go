package journal

import (
	"fmt"
	"strings"
	"time"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

type Controller struct {
	ctx  *appctx.AppContext
	repo *JournalRepository
}

func NewController(ctx *appctx.AppContext) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: NewRepository(ctx.DB.Client),
	}
}

func (c *Controller) ListWithPagination(offset, limit int) ([]*ent.Journal, error) {
	if limit <= 0 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	journals, err := c.repo.ListWithPagination(c.ctx.Context, offset, limit)
	retui.Infof("Loaded %d journals", len(journals))
	if err != nil {
		retui.Error(err)
		return nil, err
	}

	return journals, nil
}

// ShowJournal
func (*Controller) ShowJournal(id int) {
	// retui.Debugf("ID=========%d", id)
	retui.SetFocus("journal_view")
	retui.PushScreen("journal_view", retui.ScreenParams{"journalID": id})
}

// Save Journal
func (c *Controller) SaveJournal(input FormState) (*ent.Journal, error) {

	// Validation first — before any DB work
	if err := ValidateJournal(input); err != nil {
		return nil, err
	}

	voucherDate, err := time.Parse("02/01/2006", input.VcDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected DD/MM/YYYY: %w", err)
	}

	journal := &ent.Journal{
		Date:          voucherDate,
		VoucherType:   "JV",
		VoucherNo:     input.VcNo,
		VoucherDate:   voucherDate,
		ReferenceNo:   &input.VcReference,
		Narration:     &input.VcNarration,
		JournalStatus: "DRAFT",
	}

	var lines []JournalLine
	for _, l := range input.Lines {
		lines = append(lines, JournalLine{
			LedgerCode: l.LedgerCode,
			Debit:      l.Debit,
			Credit:     l.Credit,
			Remarks:    l.Remarks,
		})
	}

	jrnl, err := c.repo.CreateNew(c.ctx.Context, journal, lines)
	if err != nil {
		return nil, err
	}

	retui.Infof("Journal %s saved successfully.", jrnl.VoucherNo)
	return jrnl, nil
}

// ValidateJournal validates header + lines
func ValidateJournal(input FormState) error {
	// -- Header validation
	if strings.TrimSpace(input.VcNo) == "" {
		return fmt.Errorf("voucher number is required")
	}
	if strings.TrimSpace(input.VcDate) == "" {
		return fmt.Errorf("voucher date is required")
	}
	// Validate date format
	if _, err := time.Parse("02/01/2006", input.VcDate); err != nil {
		return fmt.Errorf("invalid date format, expected DD/MM/YYYY")
	}

	// -- Line validation
	if len(input.Lines) < 2 {
		return fmt.Errorf("journal must have at least 2 lines")
	}

	var totalDebit, totalCredit float64
	for i, line := range input.Lines {
		lineNo := i + 1

		if strings.TrimSpace(line.LedgerCode) == "" {
			return fmt.Errorf("line %d: ledger is required", lineNo)
		}
		if line.Debit < 0 {
			return fmt.Errorf("line %d: debit cannot be negative", lineNo)
		}
		if line.Credit < 0 {
			return fmt.Errorf("line %d: credit cannot be negative", lineNo)
		}
		if line.Debit == 0 && line.Credit == 0 {
			return fmt.Errorf("line %d: debit or credit must be entered", lineNo)
		}
		if line.Debit > 0 && line.Credit > 0 {
			return fmt.Errorf("line %d: cannot have both debit and credit", lineNo)
		}

		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	// -- Balance check
	if totalDebit != totalCredit {
		return fmt.Errorf(
			"journal is not balanced — debit %.2f, credit %.2f (difference: %.2f)",
			totalDebit, totalCredit, totalDebit-totalCredit,
		)
	}

	return nil
}
