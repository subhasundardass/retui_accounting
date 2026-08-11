package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	receipt := NewComponent(ctx)
	form := NewFormComponent(ctx)

	//receipt_book
	ui.Register("receipt_book", ui.Screen{
		ID:     "receipt_book",
		Title:  "Receipt Book",
		Render: receipt.ReceiptBook,
	})

	//receipt_new
	ui.Register("receipt_entry", ui.Screen{
		ID:     "receipt_entry",
		Title:  "Receipt New",
		Render: form.Receipt,
	})

}
