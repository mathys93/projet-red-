//go:build !windows

package audio

func Play(path string) {}

func Stop() {}

func SetVolume(v int) { Volume = max(0, min(1000, v)) }

func SetEnabled(on bool) { Enabled = on }
