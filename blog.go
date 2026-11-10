// Copied from Portico plan blog. Edit freely.
package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MrBlindMouse/Portico/components"
	"github.com/MrBlindMouse/Portico/installation"
	"github.com/MrBlindMouse/Portico/render"
)

const blogMenuKey = "blog_menu"
const blogPostMenuKey = "blog_post_menu"

// Blog lists post assets and opens one.
func Blog() *installation.Builder {
	b := installation.NewBuilder()

	b.Action("blog-list-up", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			menu := blogPostMenu(ctx)
			return installation.MoveList(ctx, blogMenuKey, -1, len(menu))
		},
	})
	b.Action("blog-list-down", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			menu := blogPostMenu(ctx)
			return installation.MoveList(ctx, blogMenuKey, 1, len(menu))
		},
	})
	b.Action("blog-list-go", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			ids := blogPostIDs(ctx)
			menu := blogPostMenu(ctx)
			n := len(menu)
			if n == 0 {
				return nil
			}
			idx := installation.ListIndex(ctx, blogMenuKey)
			if idx >= n {
				idx = n - 1
			}
			if idx == n-1 {
				return []components.Effect{components.TransitionEffect("landing")}
			}
			if len(ids) == 0 || idx >= len(ids) {
				return []components.Effect{components.SetLocalEffect("flash", "no posts yet")}
			}
			id := ids[idx]
			effects := []components.Effect{
				components.SetLocalEffect("flash", ""),
				components.SetLocalEffect("post_id", id),
				components.TransitionEffect("blog-post"),
			}
			if ctx.Identity().Fingerprint() != "" {
				effects = append(effects, components.PersistEffect("blog_last", id))
			}
			return effects
		},
	})
	b.Action("blog-list-back", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("landing")}
		},
	})

	b.Action("blog-post-up", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, blogPostMenuKey, -1, 1)
		},
	})
	b.Action("blog-post-down", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			return installation.MoveList(ctx, blogPostMenuKey, 1, 1)
		},
	})
	b.Action("blog-post-go", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("blog-list")}
		},
	})
	b.Action("blog-post-back", components.Action{
		Run: func(ctx components.ActionContext) []components.Effect {
			_ = ctx
			return []components.Effect{components.TransitionEffect("blog-list")}
		},
	})

	b.Scene("blog-list", components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			menu := blogPostMenu(ctx)
			return renderKeyPage(ctx, "Blog", blogListIntro(ctx), blogMenuKey, menu, pageMenuHint)
		},
		Inputs: pageNavInputs("blog-list-up", "blog-list-down", "blog-list-go", "blog-list-back"),
	})
	b.Scene("blog-post", components.Scene{
		Keys: true,
		Render: func(ctx components.RenderContext) string {
			menu := []string{"Back to posts"}
			title, body := blogPostBody(ctx)
			return renderKeyPage(ctx, title, body, blogPostMenuKey, menu, pageMenuHint)
		},
		Inputs: pageNavInputs("blog-post-up", "blog-post-down", "blog-post-go", "blog-post-back"),
	})
	b.Root("blog-list")
	return b
}

func blogPostIDs(ctx components.RenderContext) []string {
	var ids []string
	for name := range ctx.Assets() {
		rest, ok := strings.CutPrefix(name, "post-")
		if ok && rest != "" {
			ids = append(ids, rest)
		}
	}
	sort.Strings(ids)
	return ids
}

func blogPostMenu(ctx components.RenderContext) []string {
	ids := blogPostIDs(ctx)
	menu := make([]string, 0, len(ids)+1)
	for _, id := range ids {
		menu = append(menu, id)
	}
	if len(menu) == 0 {
		menu = append(menu, "(no posts)")
	}
	menu = append(menu, "Back to landing")
	return menu
}

func blogListIntro(ctx components.RenderContext) []string {
	var lines []string
	if last, ok := ctx.Persisted("blog_last"); ok {
		lines = append(lines, "Last opened: "+fmt.Sprint(last), "")
	}
	lines = append(lines, "Choose a post below.")
	return lines
}

func blogPostBody(ctx components.RenderContext) (title string, lines []string) {
	id, _ := ctx.Local("post_id")
	s, _ := id.(string)
	if s == "" {
		return "Post", []string{"no post selected"}
	}
	title = "Post: " + s
	text, ok := ctx.Asset("post-" + s)
	if !ok {
		return title, []string{fmt.Sprintf("missing asset post-%s", s)}
	}
	w, _ := ctx.Size()
	width := w * pageFramePct / 100
	if width < 40 {
		width = 40
	}
	width -= 4
	if width < 20 {
		width = 20
	}
	return title, strings.Split(render.Wrap(text, width), "\n")
}
