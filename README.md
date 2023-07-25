# Programme.lv code runner

The `programme.lv/runner` go module is meant to execute arbitrary code with
provided time and memory constraints as well as standart input, streaming
the results ( compilation stdout, stderr, execution stdout, etc.) back to the user.

The runner can be either executed through command line with a few arguments
or retrieve jobs from a RabbitMQ queue and stream back the results.

Note that this is not the module that evaluates user submissions.
See [tester](https://github.com/programme-lv/tester).

## Requirements

To sandbox the code execution the runner uses `isolate` tool.
It is required during runtime.
See [isolate](https://github.com/ioi/isolate).

After installing `isolate` the following commands should be available:
- `isolate --cg --init`;
- `isolate --cg --run /usr/bin/env pwd`;

Isolate also provides the `isolate-check-environment` utility.

Compilation of the project requires `go` to be installed.

## Command line usage

When in root of the repository, run:
```bash
go run ./cmd/runner [options] file
```

For example:
```bash
go run ./cmd/runner \
    --time 1 \
    --mem 256 \
    --lang go \
    --stdin ./test/testdata/runner/hello.in \
    ./test/testdata/runner/hello.go
```

The following options are available:
- `--time` - time limit in seconds;
- `--mem` - memory limit in megabytes;
- `--lang` - language of the code file;
- `--stdin` - path to the file containing standart input.

## Programming languages

Programming languages and other tools can be configured through
`./configs/languages.json` file. The file is read on each run.
Here's an example of the file:
```json
[
    {
        "id": "python3.10",
        "full_name": "Python 3.10",
        "code_filename": "main.py",
        "compile_cmd": null,
        "execute_cmd": "python3.10 main.py",
        "env_version_cmd": "python3.10 --version",
        "hello_world_code": "print(\"Hello, World!\")",
        "monaco_id": "python"
    },
    {
        "id": "cpp17",
        "full_name": "C++17 (GNU G++)",
        "code_filename": "main.cpp",
        "compile_cmd": "g++ -std=c++17 -o main main.cpp",
        "execute_cmd": "./main",
        "env_version_cmd": "g++ --version",
        "hello_world_code": "#include <iostream>\nint main() { std::cout << \"Hello, World!\"; }",
        "monaco_id": "cpp"
    }
]
```

When in production and receiving jobs from RabbitMQ the
runner will fetch programming language information from the database
before each run. Database connection string is configured through
`./configs/general.toml` file.


## OOP architecture

`Runner` itself is a class that takes in:
- a `Gatherer` interface;
- an `Isolate` instance;
- `code` string;
- programming language;
- `stdin` string.

To compile and execute the code in question `Runner` creates
an `IsolateBox` using the `Isolate` instance.

After compilation and during execution `Runner` reports
output, metrics and error to gatherer.

### `IsolateBox` and `IsolateProcess`

`Run` method params:
- cpu time limit;
- memory limit;
- stdin stream.

THe `Run` method returns an `IsolateProcess` pointer.

The pointer can be used to call a method that awaits the finish
of the execution and returns runtime metrics as well as providing
stdout and stderr streams.

Metrics provided by the `IsolateProcess`:
- total memory use in kilobytes by the control group;
- whether the program was killed by the out-of-memory killer;
- number of context switches forced by the kernel;
- number of context switches caused by the process giving up the CPU; 
- exitcode returned by the program;
- whether the program exited after receiving fatal signal;
- whether the program was terminated by the sandbox;
- maximum resident set size of the process in kilobytes;
- status message not intended for machine processing;
- status code:
  - RE - run-time error, i.e., exited with non-zero exit code;
  - SG - program died on signal;
  - TO - timed out;
  - XX - internal error of the sandbox.
- cpu run time of the program in fractional seconds;
- wall clock time of the program in fractional seconds.


### `Gatherer` interface

`Gatherer` collects feedback and streams it back to the user be it through
the command line or through websockets or anything else.

Currently `Gatherer` has the following methods:
- SetCompilationOutput(stdout string, stderr string)
- FinishCompilationMetrics(cpuTimeSec float64, wallTimeSec float64, memoryKb int64, exitCode int)
- AppendExecutionOutput(stdout string, stderr string)
- FinishExecutionMetrics(cpuTimeSec float64, wallTimeSec float64, memoryKb int64, exitCode int)
- FinishWithError(err string)

