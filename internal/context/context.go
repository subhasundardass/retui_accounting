package context

import (
	"context"

	"github.com/subhasundardass/retui/internal/config"
	"github.com/subhasundardass/retui/internal/database"
	"github.com/subhasundardass/retui/retui"
)

var DefaultContext = retui.CreateContext[*AppContext](nil)

// Use returns the current application context.
func Use() *AppContext {
	return retui.UseContext(DefaultContext)
}

// AppContext contains application-wide runtime state threaded through
// the RetUI render tree.
type AppContext struct {
	Context context.Context

	appName    string
	darkMode   bool
	userName   string
	toggleDark func()

	cfg *config.Config

	DB *database.DB

	// Active company for the current application session.
	companyID int
}

// AppContextValues contains the values used to initialize AppContext.
type AppContextValues struct {
	Context    context.Context
	AppName    string
	DarkMode   bool
	UserName   string
	ToggleDark func()

	DB *database.DB

	// CompanyID is the company currently active in the application.
	CompanyID int

	// Config is the fully loaded application configuration.
	Config *config.Config
}

// Set initializes the application context.
func (c *AppContext) Set(v AppContextValues) {
	if c == nil {
		return
	}

	c.Context = v.Context
	c.appName = v.AppName
	c.darkMode = v.DarkMode
	c.userName = v.UserName
	c.toggleDark = v.ToggleDark
	c.DB = v.DB
	c.companyID = v.CompanyID
	c.cfg = v.Config

	// Use Config as the fallback source for application values.
	if v.Config != nil {
		if c.appName == "" {
			c.appName = v.Config.AppName
		}
	}

	if c.userName == "" {
		c.userName = "Guest"
	}
}

// Ctx returns the application context.
func (c *AppContext) Ctx() context.Context {
	if c == nil || c.Context == nil {
		return context.Background()
	}

	return c.Context
}

// CompanyID returns the currently active company ID.
//
// A value of 0 means no company is currently selected.
func (c *AppContext) CompanyID() int {
	if c == nil {
		return 0
	}

	return c.companyID
}

// SetCompanyID changes the currently active company.
//
// Persistence of the selected company should be handled by config.Settings;
// this method only changes the runtime application state.
func (c *AppContext) SetCompanyID(id int) {
	if c == nil {
		return
	}

	c.companyID = id
}

// AppName returns the configured application name.
func (c *AppContext) AppName() string {
	if c == nil || c.appName == "" {
		return "App"
	}

	return c.appName
}

// IsDarkMode reports whether dark mode is currently enabled.
func (c *AppContext) IsDarkMode() bool {
	if c == nil {
		return false
	}

	return c.darkMode
}

// UserName returns the current user name.
func (c *AppContext) UserName() string {
	if c == nil || c.userName == "" {
		return "Guest"
	}

	return c.userName
}

// ToggleDark toggles the current dark-mode state.
func (c *AppContext) ToggleDark() {
	if c == nil || c.toggleDark == nil {
		return
	}

	c.toggleDark()
}

// Config returns the fully loaded application configuration.
func (c *AppContext) Config() *config.Config {
	if c == nil {
		return nil
	}

	return c.cfg
}

// Version returns the configured application version.
func (c *AppContext) Version() string {
	if c == nil || c.cfg == nil {
		return ""
	}

	return c.cfg.Version
}

// DateFormat returns the configured display date format.
func (c *AppContext) DateFormat() string {
	if c == nil || c.cfg == nil {
		return "DD/MM/YYYY"
	}

	return c.cfg.DateFormat
}

// Timezone returns the configured timezone.
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
