# Programme.lv code runner

The `programme.lv/runner` go module is meant to execute arbitrary code with
provided time and memory constraints as well as standart input, streaming
the results ( compilation stdout, stderr, execution stdout, etc.) back to the user.

The runner can be either executed through command line with a few arguments
or retrieve jobs from a RabbitMQ queue and stream back the results.

Note that this is not the module that evaluates user submissions.
See [tester](https://github.com/programme-lv/tester).

## Command line usage

When in root of the repository, run:
```bash
./cmd/runner/runner [options] file
```

For example:
```bash
./cmd/runner/runner \
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

`Runner` itself is a class that takes in a `Gatherer` and an `Executable`.

### `Executable`

`Executable` is an interface that has the methods `run` and `stop`.

`run` method params:
- cpu time limit;
- memory limit;
- stdin stream.

`run` method returns information about the execution, i.e.,:
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


### `Gatherer`

`Gatherer` collects feedback and streams it back to the user be it through
the command line or through websockets.

## Roadmap

- implement ioi `isolate` interface;
- placing code file in the compiler `IsolateEnvironment`;
- create `IsolatedExecutable` class;
- implement executing arbitrary code file.

