package isolate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slog"
)

/*
time:0.112
time-wall:0.103
max-rss:18984
csw-voluntary:1513
csw-forced:23
cg-mem:38248
exitcode:0
*/

/*
time:0.002
time-wall:0.045
max-rss:2624
csw-voluntary:6
csw-forced:2
cg-mem:38248
exitcode:2
status:RE
message:Exited with error status 2
*/

/*
time:0.115
time-wall:0.125
max-rss:18444
csw-voluntary:1597
csw-forced:28
cg-mem:38248
status:TO
message:Time limit exceeded
*/

type IsolateMetrics struct {
    TimeSec float64
    TimeWallSec float64
    MaxRssKb int64
    CswVoluntary int64
    CswForced int64
    CgMemKb int64
    ExitCode int64
    Status string
    Message string
}

type IsolateProcess struct {
	cmd          *exec.Cmd
	stdout       io.ReadCloser
	stderr       io.ReadCloser
	metaFilePath string
}

func (process *IsolateProcess) Wait() (*IsolateMetrics, error) {
	err := process.cmd.Wait()
	if err != nil {
		return nil, err
	}
    // read metaFilePaht
    content, err := os.ReadFile(process.metaFilePath)
    if err != nil {
        return nil, err
    }
   
    slog.Info("meta file content", slog.String("content", string(content)))

    // parse metrics
    lines := strings.Split(string(content), "\n")
    metrics := &IsolateMetrics{}

    for _, line := range lines {
        if line == "" {
            continue
        }

        parts := strings.Split(line, ":")
        if len(parts) != 2 {
            slog.Info("invalid meta file line", slog.String("line", line))
            return nil, fmt.Errorf("invalid meta file line: %s", line)
        }
        
        key := parts[0]
        value := parts[1]

        switch key {
        case "time":
            fmt.Sscanf(value, "%f", &metrics.TimeSec)
        case "time-wall":
            fmt.Sscanf(value, "%f", &metrics.TimeWallSec)
        case "max-rss":
            fmt.Sscanf(value, "%d", &metrics.MaxRssKb)
        case "csw-voluntary":
            fmt.Sscanf(value, "%d", &metrics.CswVoluntary)
        case "csw-forced":
            fmt.Sscanf(value, "%d", &metrics.CswForced)
        case "cg-mem":
            fmt.Sscanf(value, "%d", &metrics.CgMemKb)
        case "exitcode":
            fmt.Sscanf(value, "%d", &metrics.ExitCode)
        case "status":
            metrics.Status = value
        case "message":
            metrics.Message = value
        case "":
            // ignore
        default:
            slog.Info("unknown meta file line", slog.String("line", line))
        }
    }

    slog.Info("metrics",
        slog.Float64("time", metrics.TimeSec),
        slog.Float64("time-wall", metrics.TimeWallSec),
        slog.Int64("max-rss", metrics.MaxRssKb),
        slog.Int64("csw-voluntary", metrics.CswVoluntary),
        slog.Int64("csw-forced", metrics.CswForced),
        slog.Int64("cg-mem", metrics.CgMemKb),
        slog.Int64("exitcode", metrics.ExitCode),
        slog.String("status", metrics.Status),
        slog.String("message", metrics.Message))
    
	return metrics, nil
}

func (process *IsolateProcess) LogOutput() {
	stdoutScanner := bufio.NewScanner(process.Stdout())
	stderrScanner := bufio.NewScanner(process.Stderr())
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf(stdoutScanner.Text())
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			fmt.Printf(stderrScanner.Text())
		}
	}()
}

func (process *IsolateProcess) Stdout() io.ReadCloser {
	return process.stdout
}

func (process *IsolateProcess) Stderr() io.ReadCloser {
	return process.stderr
}
