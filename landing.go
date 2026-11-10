// Copied from Portico plan landing. Edit freely.
package main

import (
	"strings"
	"time"

	"github.com/MrBlindMouse/Portico/components"
	"github.com/MrBlindMouse/Portico/installation"
	"github.com/MrBlindMouse/Portico/render"
)

// Landing is the portfolio entrance. Root scene is "landing".
func Landing(destinations map[string]components.SceneID) *installation.Builder {
	items := landingItems(destinations)
	n := len(items)

	b := installation.NewBuilder()
	b.Action("landing-up", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, "", -1, n)
		},
	})
	b.Action("landing-down", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, "", 1, n)
		},
	})
	b.Action("landing-go", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			if n == 0 {
				return nil
			}
			idx := installation.ListIndex(ctx, "")
			if idx >= n {
				idx = n - 1
			}
			return []components.Effect{components.TransitionEffect(items[idx].scene)}
		},
	})

	b.Root("landing")
	b.Scene("landing", components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			labels := make([]string, len(items))
			for i, it := range items {
				labels[i] = it.label
			}
			return renderPortfolioLanding(ctx, labels)
		},
		Inputs: map[string]components.ActionID{
			"up":    "landing-up",
			"down":  "landing-down",
			"k":     "landing-up",
			"j":     "landing-down",
			"w":     "landing-up",
			"s":     "landing-down",
			"enter": "landing-go",
			"right": "landing-go",
			"l":     "landing-go",
			" ":     "landing-go",
		},
	})
	return b
}

type landingEntry struct {
	command string
	label   string
}

type landingItem struct {
	label string
	scene components.SceneID
}

var portfolioLandingMenu = []landingEntry{
	{command: "browse projects", label: "Browse Projects"},
	{command: "sign guestbook", label: "Sign guestbook"},
	{command: "read blog", label: "Read Blog"},
}

func landingItems(destinations map[string]components.SceneID) []landingItem {
	out := make([]landingItem, 0, len(portfolioLandingMenu))
	for _, item := range portfolioLandingMenu {
		scene, ok := destinations[item.command]
		if !ok {
			continue
		}
		out = append(out, landingItem{label: item.label, scene: scene})
	}
	return out
}

const landingSplitMinW = 64 // below this: one full-terminal frame, welcome over menu

func renderPortfolioLanding(ctx components.RenderContext, menuLabels []string) string {
	w, h := ctx.Size()
	if h < 6 || w < 24 {
		return render.Box("MrBlindMouse", portfolioLandingPlain(ctx, menuLabels), w)
	}

	s := render.NewScreen(w, h)
	if w < landingSplitMinW {
		drawLandingStacked(s, 0, 0, w, h, ctx, menuLabels)
		return s.Show(ctx)
	}

	ox, oy, frameW, frameH := pageFrameRect(w, h)
	if frameW < landingSplitMinW {
		frameW = landingSplitMinW
		ox = (w - frameW) / 2
	}
	drawLandingSplit(s, ox, oy, frameW, frameH, ctx, menuLabels)
	return s.Show(ctx)
}

func drawLandingSplit(s *render.Screen, ox, oy, frameW, frameH int, ctx components.RenderContext, menuLabels []string) {
	leftW := frameW / 2
	rightW := frameW - leftW
	drawPaneBorder(s, ox, oy, leftW, frameH)
	drawPaneBorder(s, ox+leftW, oy, rightW, frameH)

	leftInner := leftW - 2
	rightInner := rightW - 2
	innerH := frameH - 2
	leftCol := ox + 1
	rightCol := ox + leftW + 1
	bottomInner := oy + frameH - 2

	welcomeStart := oy + 1
	if blockH := welcomeBlockHeight(ctx); innerH > blockH {
		welcomeStart = oy + 1 + (innerH-blockH)/2
	}
	drawWelcomeBlock(s, welcomeStart, leftCol, leftInner, ctx)

	if len(menuLabels) > 0 && rightInner > 0 {
		menuStart := oy + 2
		if innerH > len(menuLabels)+3 {
			menuStart = oy + 1 + (innerH-len(menuLabels)-2)/2
		}
		listH := bottomInner - menuStart
		if listH < 1 {
			listH = 1
		}
		drawLandingMenu(s, menuStart, rightCol, rightInner, listH, menuLabels, ctx)
	}
	drawMenuHint(s, bottomInner, rightCol, rightInner, landingMenuHint)
}

func drawLandingStacked(s *render.Screen, ox, oy, frameW, frameH int, ctx components.RenderContext, menuLabels []string) {
	drawPaneBorder(s, ox, oy, frameW, frameH)
	inner := frameW - 2
	col := ox + 1
	bottomInner := oy + frameH - 2

	menuStart := drawWelcomeBlock(s, oy+1, col, inner, ctx)
	menuStart++ // blank line between welcome and menu

	if len(menuLabels) > 0 && inner > 0 && menuStart < bottomInner {
		listH := bottomInner - menuStart
		if listH < 1 {
			listH = 1
		}
		drawLandingMenu(s, menuStart, col, inner, listH, menuLabels, ctx)
	}
	drawMenuHint(s, bottomInner, col, inner, landingMenuHint)
}

func drawLandingMenu(s *render.Screen, row, col, width, height int, labels []string, ctx components.RenderContext) {
	s.DrawList(row, col, width, height, labels, installation.ListIndex(ctx, ""), 0)
}

func welcomeBlockHeight(ctx components.RenderContext) int {
	h := 9
	if flash, ok := ctx.Local("flash"); ok {
		if msg, ok := flash.(string); ok && msg != "" {
			h += 2
		}
	}
	return h
}

func drawWelcomeBlock(s *render.Screen, startRow, col, inner int, ctx components.RenderContext) int {
	name := ctx.Identity().String()
	row := startRow
	if flash, ok := ctx.Local("flash"); ok {
		if msg, ok := flash.(string); ok && msg != "" {
			drawCenteredRole(s, row, col, inner, msg, render.RolePlain)
			row += 2
		}
	}
	drawCenteredRole(s, row, col, inner, "Welcome, "+name+"!", render.RoleBold)
	row += 2
	drawCenteredRole(s, row, col, inner, "You've made it to MrBlindMouse's portfolio.", render.RolePlain)
	row += 2
	when := time.Now().In(landingSASTLoc())
	drawCenteredRole(s, row, col, inner, when.Format("Mon, 2 Jan 2006"), render.RoleMuted)
	row += 2
	drawCenteredRole(s, row, col, inner, when.Format("15:04:05 SAST"), render.RoleMuted)
	return row + 1
}

func portfolioLandingPlain(ctx components.RenderContext, menuLabels []string) string {
	var b strings.Builder
	b.WriteString("Welcome, ")
	b.WriteString(ctx.Identity().String())
	b.WriteString("!\n\n")
	b.WriteString("You've made it to MrBlindMouse's portfolio.\n\n")
	b.WriteString(landingSASTDateTime())
	b.WriteString("\n\n")
	b.WriteString(installation.RenderList(ctx, "", menuLabels, 0))
	b.WriteString("\n\n")
	b.WriteString(landingMenuHint)
	return strings.TrimRight(b.String(), "\n")
}

func landingSASTLoc() *time.Location {
	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return time.FixedZone("SAST", 2*3600)
	}
	return loc
}

func landingSASTDateTime() string {
	when := time.Now().In(landingSASTLoc())
	return when.Format("Mon, 2 Jan 2006") + "\n" + when.Format("15:04:05 SAST")
}

