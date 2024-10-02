# make-sense-of-MapReduce

Implement a simple MapReduce in Golang.

Use Docker environment to simulate the word count task.

MapReduce is More Than Just a Programming Framework:

It’s a programming model that simplifies processing large data sets across distributed clusters by abstracting the complexity of parallelization, fault tolerance, data distribution, and load balancing.

## Architecture

### Master Node

- Coordinates the MapReduce job.
- Splits the input data into chunks.
- Assigns Map and Reduce tasks to Worker nodes.
- Collects intermediate results and produces the final output.

### Worker Node

- Perform Map tasks: process data chunks and produce intermediate key-value pairs.
- Perform Reduce tasks: aggregate intermediate data into produce final results.
- Communicate with the Master node to receive tasks and send results.

## Usage

To spin up and run the task:

```bash
make build
make run
```

To combine and view the results:

```bash
cat mr-out-* > final-output.txt
cat final-output.txt
```

Use `Ctrl C` to shut the process and then:

```bash
make stop
```
