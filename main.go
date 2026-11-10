package main

import (
	"flag"
	"log"
	"os"

	portico "github.com/MrBlindMouse/Portico"
	"github.com/MrBlindMouse/Portico/components"
)

func main() {
	addr := flag.String("addr", ":2222", "SSH listen address")
	flag.Parse()

	inst, err := portico.New().
		Action("go-landing", components.Action{
			Run: func(ctx components.ActionContext) []components.Effect {
				_ = ctx
				return []components.Effect{components.TransitionEffect("landing")}
			},
		}).
		Merge(Projects()).
		Merge(Guestbook()).
		Merge(Blog()).
		Merge(Landing(map[string]components.SceneID{
			"browse projects": "projects",
			"sign guestbook":  "guestbook",
			"read blog":       "blog-list",
		})).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	dataDir := ".portico-data"
	if v := os.Getenv("PORTICO_DATA"); v != "" {
		dataDir = v
	}
	setGuestbookDataDir(dataDir)
	inspectAddr := "127.0.0.1:9922"
	opts := portico.Options{
		DataDir:     dataDir,
		AssetDir:    "assets",
		InspectAddr: inspectAddr,
		StatusAddr:  "127.0.0.1:9923",
		OnReady: func(h *portico.Host) {
			startLandingClockRepaint(h, inspectAddr)
		},
	}
	if err := portico.Serve(*addr, inst, opts); err != nil {
		log.Fatal(err)
	}
}
