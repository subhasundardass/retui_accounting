package ledger

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/components"
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

func (c *LedgerController) LedgerFilterOptions(query string) []components.SelectOption {
	ledgers, err := c.repo.Search(c.ctx.Context, query, 10)
	if err != nil {
		retui.Debug("LedgerFilterOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, len(ledgers))
	for i, l := range ledgers {
		opts[i] = components.SelectOption{Label: l.Name, Value: l.Code}
	}

	return opts
}

// ==== GROUPS
func (c *LedgerController) LedgerGroupFilterOptions(query string) []components.SelectOption {
	groups, err := c.repo.SearchGroup(c.ctx.Context, query, 10)
	if err != nil {
		retui.Debug("GroupFilterOptions error: " + err.Error())
		return nil
	}

	opts := make([]components.SelectOption, len(groups))
	for i, l := range groups {
		opts[i] = components.SelectOption{Label: l.Name, Value: l.Code}
	}

	return opts
}

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
func (c *LedgerController) CreateOrUpdate(mode FormMode, id int, in LedgerGroupState) (*ent.Ledger_Group, error) {

	// retui.Debugf("value: %v", in)

	if mode == ModeUpdate {
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
