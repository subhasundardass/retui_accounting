package context

import (
	"context"

	"github.com/subhasundardass/retui/internal/config"
	"github.com/subhasundardass/retui/internal/database"
	"github.com/subhasundardass/retui/retui"
)

var DefaultContext = retui.CreateContext[*AppContext](nil)

// Use returns the current context
func Use() *AppContext {
	return retui.UseContext(DefaultContext)
}

// AppContext is the value threaded through the render tree via retui's
// context mechanism. It previously only carried three loose fields
// (appName, darkMode, userName) derived by hand from Config in
// Bootstrap.setContext — anything added to Config later (Version,
// Theme, DateFormat, Timezone, ...) was invisible to screens even
// though it had been loaded. Screens now get the whole *config.Config
// via Config()/cfg, so nothing loaded at startup is stranded again.
type AppContext struct {
	Context    context.Context
	appName    string
	darkMode   bool
	userName   string
	toggleDark func()
	cfg        *config.Config

	DB *database.DB
}

type AppContextValues struct {
	Context    context.Context
	AppName    string
	DarkMode   bool
	UserName   string
	ToggleDark func()
	DB         *database.DB

	// Config is the fully loaded application configuration. When set,
	// it becomes the source of truth for AppName()/Theme()/etc. — pass
	// AppName/DarkMode above only if you need to override a specific
	// field without touching Config itself (e.g. a runtime dark-mode
	// toggle layered on top of the on-disk theme default).
	Config *config.Config
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
	c.cfg = v.Config

	// Backfill from Config for anything the caller left unset, so a
	// partially-populated AppContextValues (as older call sites still
	// send) doesn't lose data that Config actually has.
	if v.Config != nil {
		if c.appName == "" {
			c.appName = v.Config.AppName
		}
		if v.UserName == "" && c.userName == "" {
			c.userName = "Guest"
		}
	}
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

// Config returns the fully loaded application configuration, or nil if
// none was set (e.g. AppContext used outside of Bootstrap). Screens
// should prefer this over the individual convenience getters below when
// they need more than one field, to avoid a long chain of nil-checked
// one-liners.
func (c *AppContext) Config() *config.Config {
	if c == nil {
		return nil
	}
	return c.cfg
}

// Version returns the configured app version, or "" if unset.
func (c *AppContext) Version() string {
	if c == nil || c.cfg == nil {
		return ""
	}
	return c.cfg.Version
}

// Theme returns the configured theme name (e.g. "dark", "light").
func (c *AppContext) Theme() string {
	if c == nil || c.cfg == nil {
		return "dark"
	}
	return c.cfg.Theme
}

// DefaultPage returns the screen ID to show when none is selected.
func (c *AppContext) DefaultPage() string {
	if c == nil || c.cfg == nil {
		return "dashboard"
	}
	return c.cfg.DefaultPage
}

// DateFormat returns the configured display date format.
func (c *AppContext) DateFormat() string {
	if c == nil || c.cfg == nil {
		return "DD/MM/YYYY"
	}
	return c.cfg.DateFormat
}

// Timezone returns the configured timezone label.
func (c *AppContext) Timezone() string {
	if c == nil || c.cfg == nil {
		return "UTC"
	}
	return c.cfg.Timezone
}

// Debug reports whether debug mode is enabled.
func (c *AppContext) Debug() bool {
	if c == nil || c.cfg == nil {
		return false
	}
	return c.cfg.Debug
}
