package main

import "net/rpc"

const (
	MasterAddress = "master:1234"
)

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	NoTask
)

type Task struct {
	TaskType TaskType
	TaskID   int
	Filename string
	NReduce  int
}

type TaskResponse struct {
	Message string
}

type KeyValue struct {
	Key   string
	Value string
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	client, err := rpc.Dial("tcp", MasterAddress)
	if err != nil {
		return false
	}
	defer client.Close()
	err = client.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	return false
}
