package app

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	company "github.com/subhasundardass/retui/module/company/views"
	"github.com/subhasundardass/retui/module/dashboard"
	journal "github.com/subhasundardass/retui/module/journal/views"
	ledger "github.com/subhasundardass/retui/module/ledger/views"
	payment "github.com/subhasundardass/retui/module/payment/views"
	receipt "github.com/subhasundardass/retui/module/receipt/views"
	reports "github.com/subhasundardass/retui/module/reports/views"
	statement "github.com/subhasundardass/retui/module/statement/views"
)

func (b *Bootstrap) registerModules() {
	registerFns := []func(*appctx.AppContext){
		dashboard.Register,
		company.Register,
		ledger.Register,
		journal.Register,
		receipt.Register,
		payment.Register,
		reports.Register,
		statement.Register,
	}

	for _, register := range registerFns {
		register(b.AppCtx)
	}
}
