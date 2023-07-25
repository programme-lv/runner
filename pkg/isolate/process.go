package isolate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
    "log"
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
    timeSec float64
    timeWallSec float64
    maxRssKb int
    cswVoluntary int
    cswForced int
    cgMemKb int
    exitCode int
    status string
    message string
}

type IsolateProcess struct {
	cmd          *exec.Cmd
	stdout       io.ReadCloser
	stderr       io.ReadCloser
	metaFilePath string
}

func (process *IsolateProcess) Wait() error {
	err := process.cmd.Wait()
	if err != nil {
		return err
	}
    // read metaFilePaht
    content, err := os.ReadFile(process.metaFilePath)
    if err != nil {
        return err
    }
    
    log.Println(string(content))
	return nil
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
