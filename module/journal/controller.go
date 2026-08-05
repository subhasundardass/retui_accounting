package journal

import (
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
	voucherDate, err := time.Parse("02/01/2006", input.VcDate)
	if err != nil {
		return nil, err
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
