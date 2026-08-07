package companies

// import (
// 	"github.com/subhasundardass/retui/internal/context"
// 	"github.com/subhasundardass/retui/internal/repository"
// 	"github.com/subhasundardass/retui/retui"
// )

// func Screen(ctx *context.AppContext, props retui.Props) retui.Element {

// 	if retui.CurrentKey.Code == retui.KeyEscape && !retui.IsAnyOverlayOpen() {
// 		retui.PopScreen()
// 		retui.FocusPrev()
// 		return retui.Box(retui.Props{}, retui.NewStyle())
// 	}

// 	repo := repository.NewCompanyRepository(ctx.DB)
// 	controller := NewController(ctx, repo)
// 	components := NewComponents(controller)

// 	return components.RenderScreen()
// }
