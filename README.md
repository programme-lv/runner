# Programme.lv code runner

The programme.lv/runner go module is meant to execute arbitrary code with
provided time and memory constraints as well as standart input, streaming
the results ( compilation stdout, stderr, execution stdout, etc.) back to the user.

The runner can be either executed through command line with a few arguments
or retrieve jobs from a RabbitMQ queue and stream back the results.

Note that this is not the module that evaluates user submissions

## OOP architecture

`Runner` itself is a class that takes in a `Gatherer` and an `Executable`.

### `Executable`

`Executable` is an interface that has the methods `run` and `stop`.

`run` method params:
    - cpu time limit;
    - memory limit;
    - stdin stream
The `run` method has params for retrieving constraints and stdin stream.

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

- implement ioi `isolate` interface
- placing code file in the compiler `IsolateEnvironment`
- create `IsolatedExecutable` class
- implement executing arbitrary code file

