package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	payment := NewComponent(ctx)
	form := NewFormComponent(ctx)

	//payment_book
	ui.Register("payment_book", ui.Screen{
		ID:     "payment_book",
		Title:  "Payment Book",
		Render: payment.PaymentBook,
	})

	//payment_new
	ui.Register("payment_entry", ui.Screen{
		ID:     "payment_entry",
		Title:  "Payment New",
		Render: form.Payment,
	})

}
