package example

import "github.com/subhasundardass/retui/retui"

// Screen describes an example screen registered in the example application.
type Screen struct {
	ID     string
	Title  string
	Render func(props retui.Props) retui.Element
}

func GetScreen(id string) (Screen, bool) {
	screen, ok := Registry[id]
	return screen, ok
}

// Registry contains all example screens.
var Registry = map[string]Screen{
	"basic-inputs": {ID: "basic-inputs", Title: "Basic Inputs", Render: BasicInputExample},
	"list":         {ID: "list", Title: "List", Render: ListExample},
	"windows":      {ID: "windows", Title: "Windows", Render: WindowsExample},
	"other":        {ID: "other", Title: "Other", Render: OtherExample},
	"state":        {ID: "state", Title: "State", Render: CounterExample},
}

// ReturnToExampleRoot demonstrates the current PopTo navigation API.
// It removes all routes above the root example screen when that route exists.
func ReturnToExampleRoot() {
	retui.PopTo("home")
}
