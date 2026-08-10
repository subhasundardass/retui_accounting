package views

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/ui"
)

func Register(ctx *context.AppContext) {

	receiptBook := NewComponent(ctx)

	//receipt_book
	ui.Register("receipt_book", ui.Screen{
		ID:     "receipt_book",
		Title:  "Receipt Book",
		Render: receiptBook.ReceiptBook,
	})

}
