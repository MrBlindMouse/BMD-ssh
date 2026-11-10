package main

import (
	"encoding/json"
	"net/http"
	"time"

	portico "github.com/MrBlindMouse/Portico"
)

const landingSceneID = "landing"

// startLandingClockRepaint repaints only sessions on the landing scene so the
// welcome clock can tick without redrawing guestbook, blog, and projects.
func startLandingClockRepaint(h *portico.Host, inspectAddr string) {
	base := "http://" + inspectAddr
	client := &http.Client{Timeout: 2 * time.Second}
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for range tick.C {
			if err := repaintSessionsOnScene(h, client, base, landingSceneID); err != nil {
				return
			}
		}
	}()
}

type inspectSessionSnap struct {
	ID    string `json:"id"`
	Scene string `json:"scene"`
}

func repaintSessionsOnScene(h *portico.Host, client *http.Client, base, scene string) error {
	resp, err := client.Get(base + "/inspect/sessions")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	var sessions []inspectSessionSnap
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil
	}
	for _, s := range sessions {
		if s.Scene != scene {
			continue
		}
		if err := h.Submit(s.ID, nil); err != nil {
			return err
		}
	}
	return nil
}
