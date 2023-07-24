package runner

import (
	"io"
	"strings"

	"github.com/programme-lv/runner/internal/languages"
	"github.com/programme-lv/runner/pkg/isolate"
	"golang.org/x/exp/slog"
)

type Language languages.ProgrammingLanguage

type Gatherer interface {
	// compilation
	SetCompilationOutput(stdout string, stderr string)
	FinishCompilationMetrics(cpuTimeSec float64, wallTimeSec float64,
		memoryKb int64, exitCode int)

	// execution
	AppendExecutionOutput(stdout string, stderr string)
	FinishExecutionMetrics(cpuTimeSec float64, wallTimeSec float64,
		memoryKb int64, exitCode int)

	// error
	FinishWithError(err string)
}

type Runner struct {
	logger   *slog.Logger
	gatherer Gatherer
}

func NewRunner(gatherer Gatherer) *Runner {
	return &Runner{
		logger: slog.Default(),
	}
}

func (r *Runner) Run(code string, language Language, stdin string) {
    logger := r.logger

	isolate, err := isolate.NewIsolate()
	if err != nil {
        logger = logger.With(slog.String("error", err.Error()))
        errMsg := "failed to create isolate class"
		r.logger.Error(errMsg)
		r.gatherer.FinishWithError(errMsg)
		return
	}

	box, err := isolate.NewBox()
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
        process.LogOutput() // TODO: replace by logging to gatherer
		err = process.Wait()
		if err != nil {
            logger = logger.With(slog.String("error", err.Error()))
            errMsg := "failed to compile code"
            logger.Error(errMsg)
            r.gatherer.FinishWithError(errMsg)
			return
		}
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

    process.LogOutput() // TODO: replace by logging to gatherer
	err = process.Wait()
	if err != nil {
        logger = logger.With(slog.String("error", err.Error()))
        errMsg := "failed to run code"
        logger.Error(errMsg)
        r.gatherer.FinishWithError(errMsg)
		return
	}
}
