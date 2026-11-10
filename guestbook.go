// Copied from Portico plan guestbook. Edit freely.
package main

import (
	"fmt"
	"strings"

	"github.com/MrBlindMouse/Portico/components"
	"github.com/MrBlindMouse/Portico/installation"
	"github.com/MrBlindMouse/Portico/render"
)

const guestbookSpace = components.SpaceID("guestbook")
const guestbookListKey = "entries"
const guestbookMenuKey = "guestbook_menu"

// Guestbook is one shared board. Root scene is "guestbook".
func Guestbook() *installation.Builder {
	menu := []string{"Sign guestbook", "Back to landing"}
	n := len(menu)

	b := installation.NewBuilder()
	b.Space(guestbookSpace, components.Space{
		OnLeave: []components.Effect{components.BroadcastEffect(guestbookSpace)},
	})
	b.OnSessionStart(components.JoinSpaceEffect(guestbookSpace))
	b.OnSessionEnd(components.LeaveSpaceEffect(guestbookSpace))
	b.Root("guestbook")

	b.Action("guestbook-up", components.Action{
		Run: guestbookWithHydrate(func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, guestbookMenuKey, -1, n)
		}),
	})
	b.Action("guestbook-down", components.Action{
		Run: guestbookWithHydrate(func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, guestbookMenuKey, 1, n)
		}),
	})
	b.Action("guestbook-go", components.Action{
		Run: guestbookWithHydrate(func(ctx components.ActionContext) []components.Effect {
			idx := installation.ListIndex(ctx, guestbookMenuKey)
			if idx >= n {
				idx = n - 1
			}
			if idx == 0 {
				return []components.Effect{
					components.SetLocalEffect("flash", ""),
					components.TransitionEffect("guestbook-sign"),
				}
			}
			return []components.Effect{components.TransitionEffect("landing")}
		}),
	})
	b.Action("guestbook-back", components.Action{
		Run: guestbookWithHydrate(func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("landing")}
		}),
	})

	b.Action("guestbook-sign-submit", components.Action{
		Run: guestbookWithHydrate(func(ctx components.ActionContext) []components.Effect {
			text := strings.TrimSpace(ctx.InputLine())
			if text == "" {
				return []components.Effect{
					components.SetLocalEffect("flash", "message required"),
				}
			}
			effects := guestbookSignEffects(ctx, text)
			effects = append(effects, components.TransitionEffect("guestbook"))
			return effects
		}),
	})
	b.Action("guestbook-sign-cancel", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{
				components.SetLocalEffect("flash", ""),
				components.TransitionEffect("guestbook"),
			}
		},
	})

	b.Scene("guestbook", components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			return renderGuestbookPage(ctx, guestbookMenuKey, menu, pageMenuHint)
		},
		Inputs: pageNavInputs("guestbook-up", "guestbook-down", "guestbook-go", "guestbook-back"),
	})
	b.Scene("guestbook-sign", components.Scene{
		OnLine:  "guestbook-sign-submit",
		OnEmpty: "guestbook-sign-cancel",
		Render: func(ctx components.RenderContext) string {
			name := ctx.Identity().String()
			body := []string{
				"Signing as " + name + ".",
				"",
				"Type your message below, then press Enter.",
				"",
				"Leave the line empty and press Enter to go back without signing.",
			}
			return renderLinePage(ctx, "Sign guestbook", body, "> ")
		},
	})
	return b
}

func guestbookWithHydrate(run func(components.ActionContext) []components.Effect) func(components.ActionContext) []components.Effect {
	return func(ctx components.ActionContext) []components.Effect {
		var effects []components.Effect
		effects = append(effects, guestbookHydrateEffects(ctx)...)
		effects = append(effects, run(ctx)...)
		return effects
	}
}

func guestbookHydrateEffects(ctx components.ActionContext) []components.Effect {
	if _, ok := ctx.Space().Local(guestbookSpace, "entries_loaded"); ok {
		return nil
	}
	entries, err := loadGuestbookEntries()
	if err != nil {
		return []components.Effect{
			components.SpaceSetEffect(guestbookSpace, "entries_loaded", true),
			components.SetLocalEffect("flash", "could not load saved signatures"),
		}
	}
	effects := []components.Effect{
		components.SpaceSetEffect(guestbookSpace, "entries_loaded", true),
	}
	for _, e := range entries {
		effects = append(effects, components.AppendSpaceListEffect(guestbookSpace, guestbookListKey, e))
	}
	if len(entries) > 0 {
		effects = append(effects, components.BroadcastEffect(guestbookSpace))
	}
	return effects
}

func guestbookEntries(ctx components.RenderContext) []string {
	entries := ctx.SpaceList(guestbookSpace, guestbookListKey)
	if len(entries) > 0 {
		out := make([]string, len(entries))
		for i, e := range entries {
			out[i] = fmt.Sprint(e)
		}
		return out
	}
	saved, err := loadGuestbookEntries()
	if err != nil {
		return []string{"(could not load saved signatures)"}
	}
	return saved
}

const guestbookRecordsTitle = "Records of Legends"

func renderGuestbookPage(ctx components.RenderContext, menuKey string, menu []string, hint string) string {
	w, h := ctx.Size()
	entries := guestbookEntries(ctx)
	if h < 8 || w < 48 {
		body := append([]string{"Legend, want to leave a mark?", ""}, entries...)
		if len(entries) == 0 {
			body = []string{"Legend, want to leave a mark?", "", "(no signatures yet)"}
		}
		return renderKeyPage(ctx, guestbookRecordsTitle, body, menuKey, menu, pageMenuHint)
	}

	ox, oy, frameW, frameH := pageFrameRect(w, h)
	inner := frameW - 2
	col := ox + 1
	bottom := oy + frameH - 2

	s := render.NewScreen(w, h)
	drawPaneBorder(s, ox, oy, frameW, frameH)

	row := oy + 1
	if flash, ok := ctx.Local("flash"); ok {
		if msg, ok := flash.(string); ok && msg != "" {
			s.DrawWithin(row, col, inner, 1, render.Truncate(msg, inner))
			row += 2
		}
	}

	splitRow := row
	contentBottom := bottom - 1 // hint row
	if contentBottom < splitRow+4 {
		contentBottom = bottom
	}

	leftInner := inner * 42 / 100
	if leftInner < 22 {
		leftInner = 22
	}
	if leftInner > inner-24 {
		leftInner = inner / 2
	}
	gutter := 1
	rightCol := col + leftInner + gutter
	rightInner := inner - leftInner - gutter
	if rightInner < 12 {
		rightInner = 12
		rightCol = col + inner - rightInner
		leftInner = rightCol - col - gutter
	}

	panelTop := splitRow
	panelH := contentBottom - panelTop + 1
	if panelH < 5 {
		panelH = 5
	}
	panelW := rightInner + 2
	drawPaneBorder(s, rightCol, panelTop, panelW, panelH)

	drawCenteredRole(s, panelTop+1, rightCol+1, rightInner, guestbookRecordsTitle, render.RoleBold)
	entryCol := rightCol + 1
	innerPad := rightInner
	entryStart := panelTop + 3
	entryEnd := panelTop + panelH - 2
	maxLines := entryEnd - entryStart + 1
	if maxLines < 1 {
		maxLines = 1
	}

	visible := entries
	if len(visible) > maxLines {
		visible = visible[len(visible)-maxLines:]
	}
	entryRow := entryStart
	if len(entries) == 0 {
		drawCenteredRole(s, entryStart, entryCol, innerPad, "(no signatures yet)", render.RoleMuted)
	} else {
		for _, line := range visible {
			if entryRow > entryEnd {
				break
			}
			s.DrawWithin(entryRow, entryCol, innerPad, 1, render.Truncate(line, innerPad))
			entryRow++
		}
	}
	if len(entries) > len(visible) {
		drawCenteredRole(s, entryEnd, entryCol, innerPad,
			fmt.Sprintf("… %d earlier signature(s)", len(entries)-len(visible)), render.RoleMuted)
	}

	leftRow := splitRow
	s.DrawWithin(leftRow, col, leftInner, 1, render.Truncate("Legend, want to leave a mark?", leftInner))
	leftRow += 2
	if last, ok := ctx.Persisted("last_signature"); ok {
		s.DrawWithin(leftRow, col, leftInner, 1, render.Truncate("Your last: "+fmt.Sprint(last), leftInner))
		leftRow += 2
	}
	menuStart := leftRow
	listH := contentBottom - menuStart
	if listH < 1 {
		listH = 1
	}
	if len(menu) > 0 {
		s.DrawList(menuStart, col, leftInner, listH, menu, installation.ListIndex(ctx, menuKey), 0)
	}
	if hint != "" {
		drawMenuHint(s, bottom, col, inner, hint)
	}
	return s.Show(ctx)
}

func guestbookSignEffects(ctx components.ActionContext, text string) []components.Effect {
	text = strings.TrimSpace(text)
	if text == "" {
		return []components.Effect{components.SetLocalEffect("flash", "message required")}
	}
	line := fmt.Sprintf("%s: %s", ctx.Identity().String(), text)

	var saved []string
	for _, e := range ctx.Space().List(guestbookSpace, guestbookListKey) {
		saved = append(saved, fmt.Sprint(e))
	}
	saved = append(saved, line)
	_ = saveGuestbookEntries(saved) // ponytail: disk write failure does not block the board; entries still live in the space until restart

	effects := []components.Effect{
		components.SetLocalEffect("flash", ""),
		components.JoinSpaceEffect(guestbookSpace),
		components.AppendSpaceListEffect(guestbookSpace, guestbookListKey, line),
		components.BroadcastEffect(guestbookSpace),
	}
	if ctx.Identity().Fingerprint() != "" {
		effects = append(effects, components.PersistEffect("last_signature", line))
	}
	return effects
}
