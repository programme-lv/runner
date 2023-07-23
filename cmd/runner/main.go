package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

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
	Stdin  string
	Code    string
}

func parseArguments () Args {
    flag.Parse()

    // read code file
    var code string
    if *codePathArg == "" {
        slog.Error("No code file provided")
        os.Exit(1)
    } else {
        content, err := os.ReadFile(*codePathArg)
        if err != nil {
            slog.Error("Failed to read code file: %s", err)
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
            slog.Error("Failed to read stdin file: %s", err)
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
    
    slog.Info("Received arguments",
        slog.Float64("time limit", args.TimeLim),
        slog.Int("memory limit", args.MemLim),
        slog.String("language", args.Lang),
        slog.String("stdin", args.Stdin),
        slog.String("code", args.Code))
}
