package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
	"github.com/programme-lv/runner/internal/gatherers"
	"github.com/programme-lv/runner/internal/languages"
	"github.com/programme-lv/runner/internal/runner"
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
	TimeLim  float64
	MemLim   int
	Lang     string
	Stdin    string
	Code     string
	Filename string
}

func parseArguments() Args {
	flag.Parse()

	if *codePathArg == "" {
		slog.Error("no code file provided")
		os.Exit(1)
	}
	code := string(readFile(*codePathArg))
	filename := filepath.Base(*codePathArg)

	var stdin string
	if *stdinPathArg != "" {
		stdin = string(readFile(*stdinPathArg))
	}

	return Args{
		TimeLim:  float64(*timeLimitArg),
		MemLim:   *memLimitArg,
		Lang:     *langArg,
		Stdin:    stdin,
		Code:     code,
		Filename: filename,
	}
}

func main() {
	args := parseArguments()

	// colorful logging
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	slog.Info("using arguments",
		slog.Float64("time limit", args.TimeLim),
		slog.Int("memory limit", args.MemLim),
		slog.String("language", args.Lang),
		slog.String("stdin", args.Stdin),
		slog.String("code", args.Code))

	var languageProvider languages.LanguageProvider
	languageProvider, err := languages.NewJsonLanguageProvider("./configs/languages.json")
	if err != nil {
		slog.Error("failed to create language provider", slog.String("error", err.Error()))
        return
	}

	var language languages.ProgrammingLanguage
	if args.Lang != "" {
		var err error
		language, err = languageProvider.GetLanguage(args.Lang)
		if err != nil {
			slog.Error("failed to get programming language", slog.String("error", err.Error()))
            return
		}
	} else if args.Filename != "" {
		var err error
		extension := filepath.Ext(args.Filename)
		language, err = languageProvider.FindByFileExtension(extension)
		if err != nil {
			slog.Error("failed to get programming language", slog.String("error", err.Error()))
            return
		}
	} else {
		slog.Error("no language provided")
        return
	}

	slog.Info("found language", slog.String("language", fmt.Sprintf("%+v", language)))

    gatherer := gatherers.NewSlogGatherer()
    isolate, err := isolate.NewIsolate()
    if err != nil {
        slog.Error("failed to create isolate", slog.String("error", err.Error()))
        return
    }

    runner := runner.NewRunner(gatherer, isolate)
    runner.Run(args.Code, language, args.Stdin)

    slog.Info("finished running")
}

func readFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		slog.Error("failed to read file", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return content
}
