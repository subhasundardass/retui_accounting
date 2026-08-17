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
	ctx     *appctx.AppContext
	service Service
}

func NewController(ctx *appctx.AppContext) *Controller {
	repo := NewRepository(ctx.DB.Client)
	return &Controller{
		ctx:     ctx,
		service: NewService(repo),
	}
}

// ---- Screen navigation ----

func (*Controller) ShowJournal(id int) {
	retui.SetFocus("journal_view")
	retui.PushScreen("journal_view", retui.ScreenParams{"journalID": id})
}

// ---- Reads ----

func (c *Controller) Get(id int) (*ent.Journal, error) {
	j, err := c.service.Get(c.ctx.Context, id)
	if err != nil {
		retui.Error(err)
		return nil, err
	}
	return j, nil
}

func (c *Controller) List(offset, limit int) ([]*ent.Journal, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	journals, err := c.service.List(c.ctx.Context, ListFilter{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		retui.Error(err)
		return nil, err
	}

	retui.Infof("Loaded %d journals", len(journals))
	return journals, nil
}

// ---- Write ----

// Save creates or updates a journal voucher from form input. Validation
// (balance check, required fields, etc.) happens inside c.service.Save —
// the controller's only job is mapping FormState → VoucherInput.
func (c *Controller) Save(mode Mode, id int, input FormState) (*ent.Journal, error) {
	in, err := buildVoucherInput(input)
	if err != nil {
		return nil, err
	}

	j, err := c.service.Save(c.ctx.Context, mode, id, in)
	if err != nil {
		retui.Error(err)
		return nil, err
	}

	retui.Infof("Journal %s saved successfully.", j.VoucherNo)
	return j, nil
}

// buildVoucherInput maps UI form state into the generic VoucherInput
// shape the journal Service understands. This is where DD/MM/YYYY
// strings become time.Time and free-text fields get trimmed.
func buildVoucherInput(input FormState) (VoucherInput, error) {

	if err := ValidateFormShape(input); err != nil {
		return VoucherInput{}, err
	}

	vcNo := strings.TrimSpace(input.VcNo)
	if vcNo == "" {
		return VoucherInput{}, fmt.Errorf("voucher number is required")
	}

	date, err := time.Parse("02/01/2006", strings.TrimSpace(input.VcDate))
	if err != nil {
		return VoucherInput{}, fmt.Errorf("invalid date format, expected DD/MM/YYYY: %w", err)
	}

	lines := make([]LineInput, 0, len(input.Lines))
	for _, l := range input.Lines {
		if l.LedgerID == 0 && l.Debit == 0 && l.Credit == 0 {
			continue // skip blank trailing rows from the UI
		}
		lines = append(lines, LineInput{
			LedgerID:    l.LedgerID,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Description: strings.TrimSpace(l.Remarks),
		})
	}

	return VoucherInput{
		Type:        VoucherJV,
		VoucherNo:   vcNo,
		Date:        date,
		VoucherDate: date,
		ReferenceNo: strings.TrimSpace(input.VcReference),
		Narration:   strings.TrimSpace(input.VcNarration),
		Status:      StatusDraft,
		Lines:       lines,
	}, nil
}
