package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type PlayOpts struct {
	Title      string
	Fullscreen bool
	HWDec      string
	Player     string
}

func Play(url string, opts PlayOpts) error {
	switch opts.Player {
	case "vlc":
		return playWithVLC(url, opts)
	case "kmplayer":
		return playWithKMPlayer(url, opts)
	default:
		return playWithMPV(url, opts)
	}
}

// ─── MPV ──────────────────────────────────────────────────────────────────────

func playWithMPV(url string, opts PlayOpts) error {
	if _, err := exec.LookPath("mpv"); err != nil {
		return fmt.Errorf("mpv nao encontrado. Instale: https://mpv.io/installation/")
	}

	hwdec := opts.HWDec
	if hwdec == "" {
		hwdec = "no"
	}

	args := []string{
		url,
		"--no-terminal",
		"--really-quiet",
		"--hwdec=" + hwdec,
		"--vd-lavc-threads=0",
		"--vd-lavc-o=flags2=+ignoreDecodeErrors",
		"--demuxer-lavf-o=fflags=+genpts+discardcorrupt",
		"--cache=yes",
		"--cache-secs=15",
		"--demuxer-max-bytes=100MiB",
		"--demuxer-readahead-secs=5",
		"--stream-lavf-o=timeout=5000000",
	}
	if opts.Title != "" {
		args = append(args, "--title="+opts.Title, "--force-media-title="+opts.Title)
	}
	if opts.Fullscreen {
		args = append(args, "--fs")
	}
	return exec.Command("mpv", args...).Start()
}

// ─── VLC ──────────────────────────────────────────────────────────────────────

func findVLC() (string, error) {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidates := []string{
			filepath.Join(exeDir, "players", "vlc", "vlc-portable.exe"),
			filepath.Join(exeDir, "players", "vlc", "vlc.exe"),
			filepath.Join(exeDir, "players", "vlc", "vlc"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}

	systemNames := []string{"vlc", "cvlc"}
	if runtime.GOOS == "darwin" {
		systemNames = append(systemNames, "/Applications/VLC.app/Contents/MacOS/VLC")
	}
	if runtime.GOOS == "windows" {
		systemNames = append(systemNames,
			`C:\Program Files\VideoLAN\VLC\vlc.exe`,
			`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`,
		)
	}
	for _, name := range systemNames {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
		if filepath.IsAbs(name) {
			if _, err := os.Stat(name); err == nil {
				return name, nil
			}
		}
	}

	return "", fmt.Errorf(
		"vlc nao encontrado.\n" +
			"  Opcao 1 — coloque em: players/vlc/vlc-portable.exe\n" +
			"  Opcao 2 — instale: https://www.videolan.org/vlc/",
	)
}

func playWithVLC(url string, opts PlayOpts) error {
	bin, err := findVLC()
	if err != nil {
		return err
	}
	args := []string{
		url,
		"--no-qt-error-dialogs",
		"--quiet",
		"--network-caching=5000",
		"--live-caching=5000",
		"--file-caching=5000",
		"--avcodec-skiploopfilter=0",
		"--avcodec-skip-frame=0",
		"--avcodec-skip-idct=0",
		"--avcodec-error-resilience=1",
		"--sout-avcodec-strict=-2",
	}
	if opts.Title != "" {
		args = append(args, "--meta-title="+opts.Title)
	}
	if opts.Fullscreen {
		args = append(args, "--fullscreen")
	}
	return exec.Command(bin, args...).Start()
}

// ─── KMPlayer ─────────────────────────────────────────────────────────────────

func findKMPlayer() (string, error) {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidates := []string{
			filepath.Join(exeDir, "players", "kmplayer", "KMPlayerPortable.exe"),
			filepath.Join(exeDir, "players", "kmplayer", "KMPlayer.exe"),
			filepath.Join(exeDir, "players", "kmplayer", "kmplayer.exe"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}

	systemNames := []string{"KMPlayer", "kmplayer"}
	if runtime.GOOS == "windows" {
		systemNames = append(systemNames,
			`C:\Program Files\KMPlayer\KMPlayer.exe`,
			`C:\Program Files (x86)\KMPlayer\KMPlayer.exe`,
		)
	}
	for _, name := range systemNames {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
		if filepath.IsAbs(name) {
			if _, err := os.Stat(name); err == nil {
				return name, nil
			}
		}
	}

	return "", fmt.Errorf(
		"kmplayer nao encontrado.\n" +
			"  Opcao 1 — coloque em: players/kmplayer/KMPlayerPortable.exe\n" +
			"  Opcao 2 — instale: https://www.kmplayer.com/",
	)
}

func playWithKMPlayer(url string, opts PlayOpts) error {
	bin, err := findKMPlayer()
	if err != nil {
		return err
	}
	// KMPlayer aceita a URL diretamente como argumento
	args := []string{url}
	if opts.Fullscreen {
		args = append(args, "/fullscreen")
	}
	return exec.Command(bin, args...).Start()
}