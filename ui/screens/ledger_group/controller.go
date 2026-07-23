package ledger_group

import (
	"github.com/subhasundardass/retui/ent"
	context "github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/internal/repository"
	"github.com/subhasundardass/retui/retui"
	"github.com/subhasundardass/retui/retui/window"
)

type Controller struct {
	window *window.Window
	ctx    *context.AppContext

	repo *repository.LedgerGroupRepository
}

// EditGroupState holds the state for the edit form
type EditGroupState struct {
	id          int
	name        string
	code        string
	description string
	focusIndex  int
	totalFields int
}

func NewController(ctx *context.AppContext, repo *repository.LedgerGroupRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
	}
}

func (c *Controller) LoadGroups() []*ent.Ledger_Group {

	// Load from repository
	data, err := c.repo.List(c.ctx.Context)
	if err != nil {
		retui.Errorf("Failed to load groups: %v", err)
		data = []*ent.Ledger_Group{}
	}

	return data
}

func (c *Controller) ShowLedgers(groupID int) {
	retui.SetFocus("ledger_list")
	retui.PushScreen("ledger_list", retui.ScreenParams{"groupID": groupID})
}

func (c *Controller) HandleSave(state *EditGroupState) {

	// Validation
	if state.name == "" {
		window.AlertError("Error!", retui.Text("Group name is required", retui.NewStyle()))
		return
	}
	if state.code == "" {
		window.AlertError("Error!", retui.Text("Code is required", retui.NewStyle()))
	}

	// TODO: Update to database
	data := ent.Ledger_Group{
		Code:        state.code,
		Name:        state.name,
		Description: state.description,
	}

	_, err := c.repo.Update(c.ctx.Context, state.id, data)
	if err != nil {
		window.AlertError("Error!", retui.Text("Data update Error!", retui.NewStyle()))
		return
	}

	window.CloseFocusedWindow()
}
