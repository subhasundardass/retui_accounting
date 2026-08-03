package journal

import (
	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

type JournalController struct {
	ctx  *appctx.AppContext
	repo *JournalRepository
}

func NewController(ctx *appctx.AppContext) *JournalController {
	return &JournalController{
		ctx:  ctx,
		repo: NewRepository(ctx.DB.Client),
	}
}

func (c *JournalController) ListWithPagination(offset, limit int) ([]*ent.Journal, error) {
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
func (*JournalController) ShowJournal(id int) {
	// retui.Debugf("ID=========%d", id)
	retui.SetFocus("journal_view")
	retui.PushScreen("journal_view", retui.ScreenParams{"journalID": id})
}
