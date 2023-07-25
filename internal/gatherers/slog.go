package gatherers

import "golang.org/x/exp/slog"

type SlogGatherer struct {
}

func NewSlogGatherer() *SlogGatherer {
	return &SlogGatherer{}
}

func (g *SlogGatherer) SetCompilationOutput(stdout string, stderr string) {
	slog.Info("compilation output",
		slog.String("stdout", stdout),
		slog.String("stderr", stderr))
}

func (g *SlogGatherer) FinishCompilationMetrics(
	cpuTimeSec float64, wallTimeSec float64,
	memoryKb int64, exitCode int) {
	slog.Info("compilation metrics",
		slog.Float64("cpu_time_sec", cpuTimeSec),
		slog.Float64("wall_time_sec", wallTimeSec),
		slog.Int64("memory_kb", memoryKb),
		slog.Int("exit_code", exitCode))
}

func (g *SlogGatherer) AppendExecutionOutput(stdout string, stderr string) {
	slog.Info("execution output",
		slog.String("stdout", stdout),
		slog.String("stderr", stderr))
}

func (g *SlogGatherer) FinishExecutionMetrics(
    cpuTimeSec float64, wallTimeSec float64,
    memoryKb int64, exitCode int) {
    slog.Info("execution metrics",
        slog.Float64("cpu_time_sec", cpuTimeSec),
        slog.Float64("wall_time_sec", wallTimeSec),
        slog.Int64("memory_kb", memoryKb),
        slog.Int("exit_code", exitCode))
}

func (g *SlogGatherer) FinishWithError(err string) {
    slog.Error("finished with error", slog.String("error", err))
}

var _ Gatherer = (*SlogGatherer)(nil)
