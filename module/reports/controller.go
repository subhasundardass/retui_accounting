package reports

import (
	"github.com/subhasundardass/retui/ent"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/module/ledger"
	"github.com/subhasundardass/retui/retui"
)

type Controller struct {
	ctx            *appctx.AppContext
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
		ledgerService:  ledgerService,
		journalService: journalService,
	}

	return controller
}

func (c *Controller) GetCashBookEntries() ([]*ent.Journal, error) {

	cashLedgerID, err := c.ledgerService.GetLedger(c.ctx.Ctx(), 10)
	if err != nil {
		retui.Debug("Failed to load journals:", err)
		return nil, err
	}

	return c.journalService.List(c.ctx.Ctx(), journal.ListFilter{
		LedgerID: cashLedgerID.ID,
	})

}

func (c *Controller) BankBook() {

}
