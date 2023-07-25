package runner

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/programme-lv/runner/internal/gatherers"
	"github.com/programme-lv/runner/internal/languages"
	"github.com/programme-lv/runner/pkg/isolate"
	"golang.org/x/exp/slog"
)

type Language = languages.ProgrammingLanguage
type Gatherer = gatherers.Gatherer

type Runner struct {
	logger   *slog.Logger
	gatherer Gatherer
    isolate *isolate.Isolate
}

func NewRunner(gatherer Gatherer, isolate *isolate.Isolate) *Runner {
	return &Runner{
		logger: slog.Default(),
        isolate: isolate,
        gatherer: gatherer,
	}
}

func (r *Runner) Run(code string, language Language, stdin string) {
    logger := r.logger

	box, err := r.isolate.NewBox()
    logger = logger.With(slog.Int("box", box.Id()))
	logger.Info("created box")

	err = box.AddFile(language.CodeFilename, []byte(code))
	if err != nil {
        logger = logger.With(slog.String("error", err.Error()))
        errMsg := "failed to add code file to box"
		logger.Error(errMsg)
		r.gatherer.FinishWithError(errMsg)
		return
	}
	logger.Info("added code file to box")

	if language.CompileCmd != nil {
		logger.Info("compiling code")
		stdinReader := io.NopCloser(strings.NewReader(stdin))
		process, err := box.Run(*language.CompileCmd, stdinReader, nil)
		if err != nil {
            logger = logger.With(slog.String("error", err.Error()))
            errMsg := "failed to compile code"
            logger.Error(errMsg)
            r.gatherer.FinishWithError(errMsg)
			return
		}
        stdoutStream := process.Stdout()
        stderrStream := process.Stderr()

        stdout, err := io.ReadAll(stdoutStream)
        if err != nil {
            logger = logger.With(slog.String("error", err.Error()))
            errMsg := "failed to read compilation stdout"
            logger.Error(errMsg)
            r.gatherer.FinishWithError(errMsg)
            return
        }

        stderr, err := io.ReadAll(stderrStream)
        if err != nil {
            logger = logger.With(slog.String("error", err.Error()))
            errMsg := "failed to read compilation stderr"
            logger.Error(errMsg)
            r.gatherer.FinishWithError(errMsg)
            return
        }

        r.gatherer.SetCompilationOutput(string(stdout), string(stderr))

        metrics, err := process.Wait()

		if err != nil {
            logger = logger.With(slog.String("error", err.Error()))
            errMsg := "failed to compile code"
            logger.Error(errMsg)
            r.gatherer.FinishWithError(errMsg)
			return
		}

        r.gatherer.FinishCompilationMetrics(
            metrics.TimeSec,
            metrics.TimeWallSec,
            metrics.CgMemKb,
            metrics.ExitCode,
        )
	}

	logger.Info("running code")

	stdinReader := io.NopCloser(strings.NewReader(stdin))
	process, err := box.Run(language.ExecuteCmd, stdinReader, nil)
	if err != nil {
        logger = logger.With(slog.String("error", err.Error()))
        errMsg := "failed to run code"
        logger.Error(errMsg)
        r.gatherer.FinishWithError(errMsg)
		return
	}

    stdoutStream := process.Stdout()
    stderrStream := process.Stderr()

    stdoutScanner := bufio.NewScanner(stdoutStream)
    stderrScanner := bufio.NewScanner(stderrStream)

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        for stdoutScanner.Scan() {
            r.gatherer.AppendExecutionOutput(stdoutScanner.Text(), "")
        }
    }()

    go func() {
        defer wg.Done()
        for stderrScanner.Scan() {
            r.gatherer.AppendExecutionOutput("", stderrScanner.Text())
        }
    }()

    metrics, err := process.Wait()
	if err != nil {
        logger = logger.With(slog.String("error", err.Error()))
        errMsg := "failed to run code"
        logger.Error(errMsg)
        r.gatherer.FinishWithError(errMsg)
		return
	}

    r.gatherer.FinishExecutionMetrics(
        metrics.TimeSec,
        metrics.TimeWallSec,
        metrics.CgMemKb,
        metrics.ExitCode,
    )
}
