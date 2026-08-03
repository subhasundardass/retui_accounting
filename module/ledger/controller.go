package ledger

import (
	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
)

type LedgerController struct {
	ctx  *appctx.AppContext
	repo *Repository
}

func NewController(ctx *appctx.AppContext) *LedgerController {
	return &LedgerController{
		ctx:  ctx,
		repo: NewRepository(ctx.DB.Client),
	}
}

func (c *LedgerController) List(groupID int) ([]*ent.Ledger, error) {
	var (
		ledgers []*ent.Ledger
		err     error
	)

	if groupID > 0 {
		ledgers, err = c.repo.ListByGroup(c.ctx.Context, groupID)
	} else {
		ledgers, err = c.repo.List(c.ctx.Context)
	}

	return ledgers, err
}

//==== GROUPS

func (c *LedgerController) Groups() ([]*ent.Ledger_Group, error) {
	var (
		ledgers []*ent.Ledger_Group
		err     error
	)

	ledgers, err = c.repo.Groups(c.ctx.Context)

	return ledgers, err
}
