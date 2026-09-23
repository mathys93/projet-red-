package audio

import (
	"os"
	"path/filepath"
)

var Volume = 250

var BattleTrack = ""

var BossTracks = map[string]string{
	"DIO":            "musique/KwikFlip.mp3",
	"Diavolo":        "",
	"Yoshikage Kira": "",
	"Enrico Pucci":   "",
	"Iggy":           "",
}

func TrackFor(name string) string {
	if t := BossTracks[name]; t != "" {
		return t
	}
	return BattleTrack
}

func resolve(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.FromSlash(path)
	if filepath.IsAbs(path) {
		return path
	}
	if _, err := os.Stat(path); err == nil {
		return path
	}
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), path)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return path
}
