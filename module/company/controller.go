package company

import (
	"fmt"

	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/config"
	appctx "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

type Controller struct {
	ctx  *appctx.AppContext
	repo *Repository
}

func NewController(ctx *appctx.AppContext) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: NewRepository(ctx.DB.Client),
	}
}

// --Initialize company
func (c *Controller) Initialize() (*StartupResult, error) {
	settings, err := config.LoadSettings()
	if err != nil {
		return nil, fmt.Errorf("load company settings: %w", err)
	}

	// Try the company used in the previous session.
	if settings.LastCompanyID > 0 {
		company, err := c.repo.Get(
			c.ctx.Ctx(),
			settings.LastCompanyID,
		)

		if err == nil {
			c.ctx.SetCompanyID(company.ID)

			return &StartupResult{
				State:   StartupReady,
				Company: company,
			}, nil
		}

		// Previous company no longer exists.
		// Ignore stale setting and continue startup.
		retui.Debugf(
			"Last company %d no longer exists; selecting another company",
			settings.LastCompanyID,
		)
	}

	// No valid previous company. Check how many companies exist.
	count, err := c.repo.Count(c.ctx.Ctx())
	if err != nil {
		return nil, fmt.Errorf("count companies: %w", err)
	}

	switch count {
	case 0:
		return &StartupResult{
			State: StartupCreateCompany,
		}, nil

	case 1:
		companies, err := c.repo.List(c.ctx.Ctx())
		if err != nil {
			return nil, fmt.Errorf("load company: %w", err)
		}

		company := companies[0]

		if err := config.SaveLastCompany(company.ID); err != nil {
			return nil, fmt.Errorf("save last company: %w", err)
		}

		c.ctx.SetCompanyID(company.ID)

		return &StartupResult{
			State:   StartupReady,
			Company: company,
		}, nil

	default:
		return &StartupResult{
			State: StartupSelectCompany,
		}, nil
	}
}

func (c *Controller) List() ([]*ent.Company, error) {
	return c.repo.List(c.ctx.Ctx())
}

func (c *Controller) Edit(id int) (*ent.Company, error) {
	comp, err := c.repo.Get(c.ctx.Ctx(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to load company %d: %w", id, err)
	}
	return comp, nil
}

func (c *Controller) Save(mode FormMode, id int, data FormState) (*ent.Company, error) {
	switch mode {
	case ModeCreate:
		if err := ValidateCreate(data); err != nil {
			return nil, err
		}
		return c.repo.Create(c.ctx.Ctx(), data)
	case ModeUpdate:
		if err := ValidateUpdate(data); err != nil {
			return nil, err
		}
		return c.repo.Update(c.ctx.Ctx(), id, data)
	default:
		return nil, fmt.Errorf("unknown save mode: %v", mode)
	}
}
