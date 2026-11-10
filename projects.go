package main

import (
	"github.com/MrBlindMouse/Portico/components"
	"github.com/MrBlindMouse/Portico/installation"
)

const projectsMenuKey = "projects_menu"

type portfolioProject struct {
	label string
	scene components.SceneID
}

var portfolioProjects = []portfolioProject{
	{label: "Portico", scene: "project-portico"},
	{label: "Off-Leach.club", scene: "project-off-leach"},
	{label: "Flit-PKM", scene: "project-flit-pkm"},
	{label: "Oshkelosh", scene: "project-oshkelosh"},
	{label: "para_scope", scene: "project-para-scope"},
	{label: "Roomz", scene: "project-roomz"},
}

// Projects lists portfolio work and opens a project detail scene.
func Projects() *installation.Builder {
	b := installation.NewBuilder()

	b.Action("projects-up", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, projectsMenuKey, -1, len(projectsMenu()))
		},
	})
	b.Action("projects-down", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, projectsMenuKey, 1, len(projectsMenu()))
		},
	})
	b.Action("projects-go", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			menu := projectsMenu()
			n := len(menu)
			if n == 0 {
				return nil
			}
			idx := installation.ListIndex(ctx, projectsMenuKey)
			if idx >= n {
				idx = n - 1
			}
			if idx == n-1 {
				return []components.Effect{components.TransitionEffect("landing")}
			}
			if idx >= len(portfolioProjects) {
				return nil
			}
			return []components.Effect{
				components.SetLocalEffect("flash", ""),
				components.TransitionEffect(portfolioProjects[idx].scene),
			}
		},
	})
	b.Action("projects-back", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("landing")}
		},
	})

	b.Scene("projects", components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			return renderKeyPage(ctx, "Projects", projectsListIntro(), projectsMenuKey, projectsMenu(), pageMenuHint)
		},
		Inputs: pageNavInputs("projects-up", "projects-down", "projects-go", "projects-back"),
	})

	for _, p := range portfolioProjects {
		registerProjectDetail(b, p)
	}

	return b
}

func projectsMenu() []string {
	out := make([]string, len(portfolioProjects))
	for i, p := range portfolioProjects {
		out[i] = p.label
	}
	return append(out, "Back to landing")
}

func projectsListIntro() []string {
	return []string{
		"Pick a project to read more.",
		"Details are placeholders for now.",
	}
}

func registerProjectDetail(b *installation.Builder, p portfolioProject) {
	prefix := string(p.scene)
	menu := []string{"Back to projects"}
	menuKey := prefix + "_menu"
	n := len(menu)

	b.Action(components.ActionID(prefix+"-up"), components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, menuKey, -1, n)
		},
	})
	b.Action(components.ActionID(prefix+"-down"), components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, menuKey, 1, n)
		},
	})
	b.Action(components.ActionID(prefix+"-go"), components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("projects")}
		},
	})
	b.Action(components.ActionID(prefix+"-back"), components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("projects")}
		},
	})

	label := p.label
	b.Scene(p.scene, components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			body := []string{
				"Write-up coming soon.",
			}
			return renderKeyPage(ctx, label, body, menuKey, menu, pageMenuHint)
		},
		Inputs: pageNavInputs(
			components.ActionID(prefix+"-up"),
			components.ActionID(prefix+"-down"),
			components.ActionID(prefix+"-go"),
			components.ActionID(prefix+"-back"),
		),
	})
}
