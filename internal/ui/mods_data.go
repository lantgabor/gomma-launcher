package ui

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Mod is one row of the mods table. The backend will provide these later.
type Mod struct {
	LoadOrder   int     `json:"-"` // 1-based, from the row's position in the list
	Name        string  `json:"name"`
	Downloading float64 `json:"downloading"` // 0..1; 0 hides the progress bar
	Extracting  float64 `json:"extracting"`
	Status      string  `json:"status"` // shown as-is; "Mismatch" and "Error" are coloured
}

// placeholder_mods.json holds example rows, in load order, until the backend
// provides real data.
//
//go:embed placeholder_mods.json
var placeholderModsJSON []byte

// PlaceholderMods returns the example rows from placeholder_mods.json.
func PlaceholderMods() []Mod {
	var mods []Mod
	if err := json.Unmarshal(placeholderModsJSON, &mods); err != nil {
		panic(fmt.Sprintf("placeholder_mods.json: %v", err))
	}
	for i := range mods {
		mods[i].LoadOrder = i + 1
	}
	return mods
}
