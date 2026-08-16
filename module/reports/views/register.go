package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	cashbook := NewCashBookComponent(ctx)
	bankbook := NewBankBookComponent(ctx)

	//reports
	ui.Register("bank_book", ui.Screen{
		ID:     "bank_book",
		Title:  "Bank Book",
		Render: bankbook.Book,
	})

	ui.Register("cash_book", ui.Screen{
		ID:     "cash_book",
		Title:  "Cash Book",
		Render: cashbook.Book,
	})

}
