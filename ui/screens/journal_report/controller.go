package journal_report

import (
	"github.com/subhasundardass/retui/ent"
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

type State struct {
	IsDirty  bool
	Errors   map[string]string
	IsLoaded bool
}

type Controller struct {
	window *window.Window
	ctx    *context.AppContext
	repo   *repository.JournalRepository
	state  State // embed state directly
}

func NewController(ctx *context.AppContext, repo *repository.JournalRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
		state: State{
			Errors:   make(map[string]string),
			IsDirty:  false,
			IsLoaded: false,
		},
	}
}

func (c *Controller) GetState() State {
	return c.state
}

func (c *Controller) GetJournalsPaginated(offset, limit int) ([]*ent.Journal, error) {
	if limit <= 0 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	journals, err := c.repo.ListWithPagination(c.ctx.Context, offset, limit)
	retui.Infof("Loaded %d journals", len(journals))
	if err != nil {
		retui.Error(err)
		return nil, err
	}

	return journals, nil
}

// ShowJournal
func (*Controller) ShowJournal(id int) {
	// retui.Debugf("ID=========%d", id)
	retui.SetFocus("journal_view")
	retui.PushScreen("journal_view", retui.ScreenParams{"journalID": id})
}
