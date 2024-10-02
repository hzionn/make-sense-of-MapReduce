package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
)

type Master struct {
	mu           sync.Mutex
	files        []string
	nReduce      int
	mapTasks     map[int]bool
	reduceTasks  map[int]bool
	nMapComplete int
}

func (m *Master) AssignTask(args *Task, reply *Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, file := range m.files {
		if !m.mapTasks[i] {
			m.mapTasks[i] = true
			reply.TaskType = MapTask
			reply.TaskID = i
			reply.Filename = file
			reply.NReduce = m.nReduce
			return nil
		}
	}

	if m.nMapComplete < len(m.files) {
		reply.TaskType = NoTask
		return nil
	}

	for i := 0; i < m.nReduce; i++ {
		if !m.reduceTasks[i] {
			m.reduceTasks[i] = true
			reply.TaskType = ReduceTask
			reply.TaskID = i
			reply.NReduce = m.nReduce
			return nil
		}
	}

	reply.TaskType = NoTask
	return nil
}

func (m *Master) TaskCompleted(args *Task, reply *TaskResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if args.TaskType == MapTask {
		m.nMapComplete++
	}

	reply.Message = "Task recorded as complete."
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run master.go input_data/input1.txt input_data/input2.txt ...")
		os.Exit(1)
	}

	m := Master{
		files:       os.Args[1:],
		nReduce:     3,
		mapTasks:    make(map[int]bool),
		reduceTasks: make(map[int]bool),
	}

	rpc.Register(&m)
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer l.Close()

	fmt.Println("Master is running...")
	rpc.Accept(l)
}
