package context

import (
	"context"

	"github.com/subhasundardass/retui/internal/database"
	"github.com/subhasundardass/retui/retui"
)

var DefaultContext = retui.CreateContext[*AppContext](nil)

// Use returns the current context
func Use() *AppContext {
	return retui.UseContext(DefaultContext)
}

type AppContext struct {
	Context    context.Context
	appName    string
	darkMode   bool
	userName   string
	toggleDark func()

	DB *database.DB
}

type AppContextValues struct {
	Context    context.Context
	AppName    string
	DarkMode   bool
	UserName   string
	ToggleDark func()
	DB         *database.DB
}

func (c *AppContext) Set(v AppContextValues) {
	// Add nil check
	if c == nil {
		return
	}
	c.Context = v.Context
	c.appName = v.AppName
	c.darkMode = v.DarkMode
	c.userName = v.UserName
	c.toggleDark = v.ToggleDark
	c.DB = v.DB
}

func (c *AppContext) Ctx() context.Context {
	if c == nil || c.Context == nil {
		return context.Background()
	}
	return c.Context
}

func (c *AppContext) AppName() string {
	if c == nil {
		return "App"
	}
	return c.appName
}

func (c *AppContext) IsDarkMode() bool {
	if c == nil {
		return false
	}
	return c.darkMode
}

func (c *AppContext) UserName() string {
	if c == nil {
		return "Guest"
	}
	return c.userName
}

func (c *AppContext) ToggleDark() {
	if c == nil || c.toggleDark == nil {
		return
	}
	c.toggleDark()
}
