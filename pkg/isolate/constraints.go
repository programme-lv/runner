package isolate

import (
	"fmt"
)

type RuntimeConstraints struct {
    CpuTimeLimInSec float64
    ExtraCpuTimeLimInSec float64
    WallTimeLimInSec float64
    MemoryLimitInKB int
    MaxProcesses int
}

func DefaultRuntimeConstraints() RuntimeConstraints {
    return RuntimeConstraints{
        CpuTimeLimInSec: 1.0,
        ExtraCpuTimeLimInSec: 0.5,
        WallTimeLimInSec: 10.0,
        MemoryLimitInKB: 256000,
    }
}

func (constraints *RuntimeConstraints) AllIsolateArgs() []string {
    return []string{
        constraints.MemLimArg(),
        constraints.CpuTimeLimArg(),
        constraints.ExtraCpuTimeLimArg(),
        constraints.WallTimeLimArg(),
        constraints.MaxProcessesArg(),
    }
}

func (constraints *RuntimeConstraints) MemLimArg() string {
    return fmt.Sprintf("--mem=%d", constraints.MemoryLimitInKB)
}

func (constraints *RuntimeConstraints) CpuTimeLimArg() string {
    return fmt.Sprintf("--time=%f", constraints.CpuTimeLimInSec)
}

func (constraints *RuntimeConstraints) ExtraCpuTimeLimArg() string {
    return fmt.Sprintf("--extra-time=%f", constraints.ExtraCpuTimeLimInSec)
}

func (constraints *RuntimeConstraints) WallTimeLimArg() string {
    return fmt.Sprintf("--wall-time=%f", constraints.WallTimeLimInSec)
}

func (constraints *RuntimeConstraints) MaxProcessesArg() string {
    return fmt.Sprintf("--processes=%d", constraints.MaxProcesses)
}

