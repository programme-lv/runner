package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/programme-lv/runner/pkg/isolate"
	"golang.org/x/exp/slog"
)

var (
	timeLimitArg = flag.Int("time", 1, "time limit in seconds")
	memLimitArg  = flag.Int("mem", 256, "memory limit in megabytes")
	langArg      = flag.String("lang", "", "language of the code file")
	stdinPathArg = flag.String("stdin", "", "path to the file containing standard input")
	codePathArg  = flag.String("code", "", "path to the code file")
)

type Args struct {
	TimeLim float64
	MemLim  int
	Lang    string
	Stdin   string
	Code    string
}

func parseArguments() Args {
	flag.Parse()

	// read code file
	var code string
	if *codePathArg == "" {
		slog.Error("no code file provided")
		os.Exit(1)
	} else {
		content, err := os.ReadFile(*codePathArg)
		if err != nil {
			slog.Error("failed to read code file: %s", err)
			os.Exit(1)
		}
		code = string(content)
	}

	// guess programming language
	var lang string
	if *langArg == "" {
		extension := filepath.Ext(*codePathArg)
		lang = strings.TrimPrefix(extension, ".")
	}

	var stdin string
	if *stdinPathArg != "" {
		content, err := os.ReadFile(*stdinPathArg)
		if err != nil {
			slog.Error("failed to read stdin file: %s", err)
			os.Exit(1)
		}
		stdin = string(content)
	}

	return Args{
		TimeLim: float64(*timeLimitArg),
		MemLim:  *memLimitArg,
		Lang:    lang,
		Stdin:   stdin,
		Code:    code,
	}
}

func main() {
	args := parseArguments()

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	slog.Info("received arguments",
		slog.Float64("time limit", args.TimeLim),
		slog.Int("memory limit", args.MemLim),
		slog.String("language", args.Lang),
		slog.String("stdin", args.Stdin),
		slog.String("code", args.Code))

	isolate, err := isolate.NewIsolate()
	if err != nil {
		slog.Error("failed to create isolate: %s", err)
		os.Exit(1)
	}

    box, err := isolate.NewBox()
    box.Id()
    slog.Info("created box", slog.Int("id", box.Id()))
}
