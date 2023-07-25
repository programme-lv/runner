package gatherers


type Gatherer interface {
	// compilation
	SetCompilationOutput(stdout string, stderr string)
	FinishCompilationMetrics(cpuTimeSec float64, wallTimeSec float64,
		memoryKb int64, exitCode int64)

	// execution
	AppendExecutionOutput(stdout string, stderr string)
	FinishExecutionMetrics(cpuTimeSec float64, wallTimeSec float64,
		memoryKb int64, exitCode int64)

	// error
	FinishWithError(err string)
}
