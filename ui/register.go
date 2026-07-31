// ui/register.go

package ui

import (
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

type RegisterScreen struct {
	ID     string
	Title  string
	Render func(ctx *appctx.AppContext) retui.Element
}

var routes = map[string]RegisterScreen{}

func Register(id string, screen RegisterScreen) {
	if _, exists := routes[id]; exists {
		panic("screen already registered: " + id)
	}
	routes[id] = screen
}

func GetRegisterScreen(id string) (RegisterScreen, bool) {
	s, ok := routes[id]
	return s, ok
}
