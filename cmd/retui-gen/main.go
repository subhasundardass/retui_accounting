package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	makeCmd := flag.String("make", "", "Generate: screen:name")
	flag.Parse()

	if *makeCmd == "" {
		fmt.Println("Usage: retui-gen -make screen:name")
		fmt.Println("Example: retui-gen -make screen:general")
		os.Exit(1)
	}

	parts := strings.SplitN(*makeCmd, ":", 2)
	if len(parts) != 2 || parts[0] != "screen" {
		fmt.Println("Error: Use format 'screen:name'")
		os.Exit(1)
	}

	name := parts[1]
	generateScreen(name)
}

func generateScreen(name string) {
	// Create directory
	dir := filepath.Join("ui/screens", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Generate files
	files := map[string]string{
		"screen.go":     screenTemplate,
		"controller.go": controllerTemplate,
		"components.go": componentsTemplate,
	}

	for filename, tmpl := range files {
		path := filepath.Join(dir, filename)
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("Error creating %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer f.Close()

		t := template.Must(template.New(filename).Parse(tmpl))
		t.Execute(f, map[string]string{
			"Name":  name,
			"Title": strings.Title(name),
		})
	}

	fmt.Printf("✅ Screen '%s' created at internal/ui/screens/%s\n", name, name)
	fmt.Println("\n📝 Register in internal/ui/registry.go:")
	fmt.Printf(`"%s": {
    ID:     "%s",
    Title:  "%s",
    Render: %s.Screen,
},`, name, name, strings.Title(name), name)
}

// ─── Templates ──────────────────────────────────────────────────────────────

const screenTemplate = `package {{.Name}}

import (
	"github.com/subhasundardass/retui/internal/context"
	"github.com/subhasundardass/retui/retui"
)

func Screen(ctx *context.AppContext, props retui.Props) retui.Element {

	if retui.CurrentKey.Code == retui.KeyEscape && !retui.IsAnyOverlayOpen() {
		retui.PopScreen()
		retui.FocusPrev()
		return retui.Box(retui.Props{}, retui.NewStyle())
	}

	// repo := repository.NewRepository(ctx.DB)
	controller := NewController(ctx, repo)
	components := NewComponents(controller)
	
	return components.RenderScreen()
}`

const controllerTemplate = `package {{.Name}}


type Controller struct{
	window *window.Window
	ctx    *context.AppContext
	repo   *repository.YourRepository
}

func NewController(ctx *context.AppContext, repo *repository.YourRepository) *Controller {
	return &Controller{
		ctx:  ctx,
		repo: repo,
		
	}
}

`

const componentsTemplate = `package {{.Name}}

import (
	"github.com/subhasundardass/retui/retui"
)

type Components struct {
	controller *Controller
}

func NewComponents(controller *Controller) *Components {
	return &Components{controller: controller}
}

func (c *Components) bindKeys() {
	if retui.IsFocused("screen_id") {
		switch retui.CurrentKey.Code {
		}
	}
}

func (c *Components) RenderScreen() retui.Element {
	//--Key bindings
	c.bindKeys()

	return retui.Box(
		retui.Props{Direction: retui.Column, Gap: 1},
		retui.NewStyle(),

		retui.Text("{{.Name}}", retui.NewStyle().Bold(true).Foreground(retui.Cyan)),
	)
}`
