// ui/register.go

package ui

import (
	"fmt"

	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

type Screen struct {
	ID     string
	Title  string
	Render func(ctx *appctx.AppContext) retui.Element
}

var routes = map[string]Screen{}

func Register(key string, screen Screen) {
	if _, exists := routes[key]; exists {
		panic("screen already registered: " + key)
	}
	if screen.ID != key {
		panic(fmt.Sprintf("screen ID %q does not match registration key %q", screen.ID, key))
	}
	routes[key] = screen
}

func GetScreen(id string) (Screen, bool) {
	s, ok := routes[id]
	return s, ok
}
