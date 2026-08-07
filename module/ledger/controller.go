package ledger

import (
	"fmt"

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

func (c *LedgerController) GetGroup(id int) (*LedgerGroupState, error) {
	group, err := c.repo.GetGroup(c.ctx.Ctx(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to load group %d: %w", id, err)
	}

	return &LedgerGroupState{
		Code:        group.Code,
		Name:        group.Name,
		Nature:      string(group.Nature),
		Description: group.Description,
	}, nil
}

// -Create or Update
func (c *LedgerController) CreateOrUpdate(id int, in LedgerGroupState) (*ent.Ledger_Group, error) {

	if in.Mode == ModeUpdate {
		return c.repo.GroupUpdate(c.ctx.Ctx(), id, in)
	}
	return c.repo.GroupCreate(c.ctx.Ctx(), in)
}

// func (c *LedgerController) EditGroup(id int) {
// 	group, err := c.repo.GetGroup(c.ctx.Ctx(), id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load Group %d: %w", id, err)
// 	}
// 	// return comp, nil

// 	state:= LedgerGroupState{
// 		Code: ,
// 	}
// }
