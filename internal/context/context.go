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
	Context       context.Context
	currentPage   string
	appName       string
	darkMode      bool
	userName      string
	navigateTo    func(string)
	toggleDark    func()
	pushScreen    func(string)
	replaceScreen func(string)
	popScreen     func() string
	getStack      func() []string
	getCurrent    func() string
	DB            *database.DB
}

type AppContextValues struct {
	Context       context.Context
	CurrentPage   string
	AppName       string
	DarkMode      bool
	UserName      string
	NavigateTo    func(string)
	ToggleDark    func()
	PushScreen    func(string)
	ReplaceScreen func(string)
	PopScreen     func() string
	GetStack      func() []string
	GetCurrent    func() string
	DB            *database.DB
}

func (c *AppContext) Set(v AppContextValues) {
	// Add nil check
	if c == nil {
		return
	}
	c.Context = v.Context
	c.currentPage = v.CurrentPage
	c.appName = v.AppName
	c.darkMode = v.DarkMode
	c.userName = v.UserName
	c.navigateTo = v.NavigateTo
	c.toggleDark = v.ToggleDark
	c.pushScreen = v.PushScreen
	c.replaceScreen = v.ReplaceScreen
	c.popScreen = v.PopScreen
	c.getStack = v.GetStack
	c.getCurrent = v.GetCurrent
	c.DB = v.DB
}

func (c *AppContext) Ctx() context.Context {
	if c == nil || c.Context == nil {
		return context.Background()
	}
	return c.Context
}

// func (c *AppContext) Page() string {
// 	if c == nil {
// 		return "home"
// 	}
// 	return c.currentPage
// }

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

// func (c *AppContext) NavigateTo(page string) {
// 	if c == nil {
// 		return
// 	}
// 	c.navigateTo(page)
// }

func (c *AppContext) ToggleDark() {
	if c == nil || c.toggleDark == nil {
		return
	}
	c.toggleDark()
}
