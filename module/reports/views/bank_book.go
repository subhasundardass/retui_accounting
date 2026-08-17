package views

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/module/receipt"
	"github.com/subhasundardass/retui/retui"
)

type BankBookComponent struct {
	controller *receipt.Controller
	ctx        *appctx.AppContext
}

func NewBankBookComponent(ctx *appctx.AppContext) *BankBookComponent {
	return &BankBookComponent{
		controller: receipt.NewController(ctx),
		ctx:        ctx,
	}
}

func (c *BankBookComponent) bindKeys() {
	key := retui.CurrentKey
	if key == (retui.Key{}) || key.Consumed {
		return
	}
	if retui.CapturedFocus() != "" {
		return
	}

	switch retui.CurrentKey.Code {
	case retui.KeyEscape:
		retui.PopScreen()

	}
}

func (c *BankBookComponent) Book(ctx *appctx.AppContext) retui.Element {
	c.bindKeys()
	return retui.Element{}
}
