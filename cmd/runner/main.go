package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/programme-lv/runner/internal/languages"
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
	}

	var language languages.ProgrammingLanguage
	if args.Lang != "" {
		var err error
		language, err = languageProvider.GetLanguage(args.Lang)
		if err != nil {
			slog.Error("failed to get programming language", slog.String("error", err.Error()))
			os.Exit(1)
		}
	} else if args.Filename != "" {
		var err error
		extension := filepath.Ext(args.Filename)
		language, err = languageProvider.FindByFileExtension(extension)
		if err != nil {
			slog.Error("failed to get programming language", slog.String("error", err.Error()))
			os.Exit(1)
		}
	} else {
		slog.Error("no language provided")
	}

	slog.Info("found language", slog.String("language", fmt.Sprintf("%+v", language)))

	isolate, err := isolate.NewIsolate()
	if err != nil {
		slog.Error("failed to create isolate class", slog.String("error", err.Error()))
		os.Exit(1)
	}

	box, err := isolate.NewBox()
	box.Id()
	slog.Info("created box", slog.Int("id", box.Id()))

	// place the code file in the box
	err = box.AddFile(language.CodeFilename, []byte(args.Code))
	if err != nil {
		slog.Error("failed to add code file to box", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if language.CompileCmd != nil {
        slog.Info("compiling code")
		stdin := io.NopCloser(strings.NewReader(args.Stdin))
		process, err := box.Run(*language.CompileCmd, stdin, nil)
		if err != nil {
			slog.Error("failed to compile code", slog.String("error", err.Error()))
			os.Exit(1)
		}
        process.LogOutput()
        err = process.Wait()
        if err != nil {
            slog.Error("failed to compile code", slog.String("error", err.Error()))
            os.Exit(1)
        }
	}

    slog.Info("running code")

    stdin := io.NopCloser(strings.NewReader(args.Stdin))
    process, err := box.Run(language.ExecuteCmd, stdin, nil)
    if err != nil {
        slog.Error("failed to run code", slog.String("error", err.Error()))
        os.Exit(1)
    }
    process.LogOutput()
    err = process.Wait()
    if err != nil {
        slog.Error("failed to run code", slog.String("error", err.Error()))
        os.Exit(1)
    }
}

func readFile(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		slog.Error("failed to read file", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return content
}
