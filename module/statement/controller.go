package statement

import (
	"context"
	"fmt"
	"time"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/journal"
	"github.com/subhasundardass/retui/module/ledger"
)

type Controller struct {
	ctx            *appctx.AppContext
	repo           *Repository
	ledgerService  ledger.Service
	journalService journal.Service
}

// Result is the fully computed statement handed back to the view layer.
type Result struct {
	Rows           []Entry
	OpeningBalance float64
	ClosingBalance float64
}

// dateLayout matches the DD/MM/YYYY format used by components.DateInput.
const dateLayout = "02/01/2006"

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

// Load validates the raw filter inputs, fetches data, and returns a fully
// computed Result (running balances included). It never mutates FormState
// directly — the view decides how to apply the result.
func (c *Controller) Load(ledgerAc int, fromRaw, toRaw string) (ledger.Result, error) {
	if ledgerAc <= 0 {
		return ledger.Result{}, fmt.Errorf("please select a ledger account")
	}

	from, err := time.Parse(dateLayout, fromRaw)
	if err != nil {
		return ledger.Result{}, fmt.Errorf("invalid from date, expected DD/MM/YYYY")
	}
	to, err := time.Parse(dateLayout, toRaw)
	if err != nil {
		return ledger.Result{}, fmt.Errorf("invalid to date, expected DD/MM/YYYY")
	}
	if to.Before(from) {
		return ledger.Result{}, fmt.Errorf("to date cannot be before from date")
	}
	// Include the whole "to" day.
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	queryCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	statement, err := c.ledgerService.GetStatement(queryCtx, ledgerAc, from, to)
	if err != nil {
		return ledger.Result{}, err
	}

	// running := entries.OpeningBalance
	// for i := range entries {
	// 	running += entries[i].Debit - entries[i].Credit
	// 	entries[i].Balance = running
	// }

	// return Result{
	// 	Rows:           entries,
	// 	OpeningBalance: opening,
	// 	ClosingBalance: running,
	// }, nil

	return statement, err
}
