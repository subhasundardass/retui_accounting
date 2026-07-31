package app

import "github.com/subhasundardass/retui/module/company"

func (b *Bootstrap) registerModules() {

	company.Register(b.AppCtx)
}
