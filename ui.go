package main

import (
	"strings"

	"github.com/MrBlindMouse/Portico/components"
	"github.com/MrBlindMouse/Portico/installation"
	"github.com/MrBlindMouse/Portico/render"
)

const (
	pageFramePct    = 90
	landingMenuHint = "↑↓ jk ws  ·  Enter → l Space"
	pageMenuHint    = "↑↓ jk ws  ·  Enter → l Space  ·  ← h b back"
)

func pageNavInputs(up, down, goAct, back components.ActionID) map[string]components.ActionID {
	return map[string]components.ActionID{
		"up":    up,
		"down":  down,
		"k":     up,
		"j":     down,
		"w":     up,
		"s":     down,
		"enter": goAct,
		"right": goAct,
		"l":     goAct,
		" ":     goAct,
		"left":  back,
		"h":     back,
		"b":     back,
	}
}

func pageFrameRect(w, h int) (ox, oy, frameW, frameH int) {
	frameW = w * pageFramePct / 100
	frameH = h * pageFramePct / 100
	if frameW < 40 {
		frameW = 40
	}
	if frameH < 8 {
		frameH = 8
	}
	if w < frameW {
		frameW = w
	}
	if h < frameH {
		frameH = h
	}
	ox = (w - frameW) / 2
	oy = (h - frameH) / 2
	return ox, oy, frameW, frameH
}

// renderKeyPage draws a framed page with body lines and a key-driven menu list.
func renderKeyPage(ctx components.RenderContext, title string, body []string, menuKey string, menu []string, hint string) string {
	w, h := ctx.Size()
	if h < 6 || w < 24 {
		return render.Box(title, strings.Join(body, "\n")+"\n\n"+installation.RenderList(ctx, menuKey, menu, 0)+"\n\n"+hint, w)
	}

	ox, oy, frameW, frameH := pageFrameRect(w, h)
	inner := frameW - 2
	col := ox + 1
	bottom := oy + frameH - 2

	s := render.NewScreen(w, h)
	drawPaneBorder(s, ox, oy, frameW, frameH)

	row := oy + 1
	drawCenteredRole(s, row, col, inner, title, render.RoleBold)
	row += 2

	if flash, ok := ctx.Local("flash"); ok {
		if msg, ok := flash.(string); ok && msg != "" {
			s.DrawWithin(row, col, inner, 1, render.Truncate(msg, inner))
			row += 2
		}
	}

	for _, line := range body {
		if row >= bottom-2 {
			break
		}
		s.DrawWithin(row, col, inner, 1, render.Truncate(line, inner))
		row++
	}

	if len(menu) > 0 && row < bottom {
		row++
		listH := bottom - row
		if listH < 1 {
			listH = 1
		}
		s.DrawList(row, col, inner, listH, menu, installation.ListIndex(ctx, menuKey), 0)
	}

	if hint != "" {
		drawMenuHint(s, bottom, col, inner, hint)
	}
	return s.Show(ctx)
}

// renderLinePage draws a framed page with body text and a placed line editor (line mode).
func renderLinePage(ctx components.RenderContext, title string, body []string, prompt string) string {
	w, h := ctx.Size()
	if h < 6 || w < 24 {
		var b strings.Builder
		b.WriteString(strings.Join(body, "\n"))
		b.WriteString("\n\n")
		b.WriteString(prompt)
		return render.Box(title, strings.TrimRight(b.String(), "\n"), w)
	}

	ox, oy, frameW, frameH := pageFrameRect(w, h)
	inner := frameW - 2
	col := ox + 1
	bottom := oy + frameH - 2
	if prompt == "" {
		prompt = "> "
	}

	s := render.NewScreen(w, h)
	drawPaneBorder(s, ox, oy, frameW, frameH)

	row := oy + 1
	drawCenteredRole(s, row, col, inner, title, render.RoleBold)
	row += 2

	if flash, ok := ctx.Local("flash"); ok {
		if msg, ok := flash.(string); ok && msg != "" {
			drawCenteredRole(s, row, col, inner, msg, render.RolePlain)
			row += 2
		}
	}

	for _, line := range body {
		if row >= bottom-2 {
			break
		}
		s.DrawWithin(row, col, inner, 1, render.Truncate(line, inner))
		row++
	}

	s.Input(bottom, col, inner, prompt)
	return s.Show(ctx)
}

func drawCenteredRole(s *render.Screen, row, col, inner int, plain string, role render.Role) {
	text := render.Truncate(plain, inner)
	vis := len([]rune(text))
	s.DrawWithin(row, col, inner, 1, render.AlignCenter(text, inner))
	if role == render.RolePlain || vis == 0 {
		return
	}
	startCol := col + (inner-vis)/2
	s.MarkRole(row, startCol, vis, role)
}

func drawMenuHint(s *render.Screen, row, col, inner int, hint string) {
	if row < 0 || inner < 1 || hint == "" {
		return
	}
	s.DrawWithin(row, col, inner, 1, render.AlignCenter(render.Truncate(hint, inner), inner))
}

func drawPaneBorder(s *render.Screen, col, row, pw, ph int) {
	inner := pw - 2
	if inner < 1 || ph < 2 {
		return
	}
	bar := strings.Repeat("─", inner)
	s.DrawWithin(row, col, pw, 1, "┌"+bar+"┐")
	s.DrawWithin(row+ph-1, col, pw, 1, "└"+bar+"┘")
	side := "│" + strings.Repeat(" ", inner) + "│"
	for r := 1; r < ph-1; r++ {
		s.DrawWithin(row+r, col, pw, 1, side)
	}
}
