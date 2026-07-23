package ledger

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

// State - keep it simple, just a struct
type State struct {
	// Ledgers        []*ent.Ledger     `json:"ledgers"`
	// SelectedLedger *ent.Ledger       `json:"current"`
	Errors   map[string]string `json:"errors"`
	IsDirty  bool              `json:"isDirty"`
	IsLoaded bool              `json:"isLoaded"`
}

// Controller - just a thin wrapper
type Controller struct {
	window      *window.Window
	ctx         *context.AppContext
	repo        *repository.LedgerRepository
	journalRepo *repository.JournalRepository
	state       State // embed state directly
}

func NewController(
	ctx *context.AppContext,
	repo *repository.LedgerRepository,
	journalRepo *repository.JournalRepository,
) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
		state: State{
			Errors:   make(map[string]string),
			IsDirty:  false,
			IsLoaded: false,
		},
	}
}

func (c *Controller) GetState() State {
	return c.state
}

// ----------Business Logic-------------------
func (c *Controller) GetLedgers(groupID int) []*ent.Ledger {

	var (
		ledgers []*ent.Ledger
		err     error
	)

	if groupID > 0 {
		ledgers, err = c.repo.ListByGroup(c.ctx.Context, groupID)
	} else {
		ledgers, err = c.repo.List(c.ctx.Context)
	}

	if err != nil {
		c.state.Errors["load"] = err.Error()
		return nil
	}

	// c.state.Ledgers = ledgers
	c.state.IsLoaded = true
	delete(c.state.Errors, "load")

	return ledgers
}

// Show Journals

func (c *Controller) GetJournalsByLedger(ledgerID int) {

	retui.PushScreen("journal_line", retui.ScreenParams{"ledgerID": ledgerID})
	retui.SetFocus("journal_line")

}
