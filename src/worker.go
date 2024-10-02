// worker.go
package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

func main() {
	for {
		task := RequestTask()
		switch task.TaskType {
		case MapTask:
			fmt.Println("Received Map task:", task.TaskID)
			DoMapTask(task)
		case ReduceTask:
			fmt.Println("Received Reduce task:", task.TaskID)
			DoReduceTask(task)
		case NoTask:
			fmt.Println("No task assigned. Sleeping...")
			time.Sleep(time.Second * 2)
		}
	}
}

func RequestTask() Task {
	args := Task{}
	reply := Task{}
	ok := call("Master.AssignTask", &args, &reply)
	if !ok {
		log.Fatal("Failed to call Master.AssignTask")
	}
	return reply
}

func DoMapTask(task Task) {
	// Read file
	content, err := ioutil.ReadFile(task.Filename)
	if err != nil {
		log.Fatalf("Cannot read %v", task.Filename)
	}

	// Map function
	kva := Map(task.Filename, string(content))

	// Partition and write intermediate files
	for i := 0; i < task.NReduce; i++ {
		intermediate := []KeyValue{}
		for _, kv := range kva {
			if ihash(kv.Key)%task.NReduce == i {
				intermediate = append(intermediate, kv)
			}
		}
		intermediateFile := fmt.Sprintf("mr-%d-%d", task.TaskID, i)
		file, _ := os.Create(intermediateFile)
		enc := json.NewEncoder(file)
		for _, kv := range intermediate {
			enc.Encode(&kv)
		}
		file.Close()
	}

	// Notify Master
	args := task
	reply := TaskResponse{}
	call("Master.TaskCompleted", &args, &reply)
	fmt.Println(reply.Message)
}

func DoReduceTask(task Task) {
	intermediate := []KeyValue{}

	// Read intermediate files
	for i := 0; ; i++ {
		filename := fmt.Sprintf("mr-%d-%d", i, task.TaskID)
		file, err := os.Open(filename)
		if err != nil {
			break
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
		file.Close()
	}

	// Sort by key
	sort.Slice(intermediate, func(i, j int) bool {
		return intermediate[i].Key < intermediate[j].Key
	})

	// Reduce function
	oname := fmt.Sprintf("output_data/mr-out-%d", task.TaskID)
	ofile, err := os.Create(oname)
	if err != nil {
		log.Fatalf("Cannot create output file %v", oname)
	}

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := Reduce(intermediate[i].Key, values)
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)
		i = j
	}
	ofile.Close()

	// Notify Master
	args := task
	reply := TaskResponse{}
	call("Master.TaskCompleted", &args, &reply)
	fmt.Println(reply.Message)
}

func Map(filename string, contents string) []KeyValue {
	var kva []KeyValue
	words := splitWords(contents)
	for _, word := range words {
		kva = append(kva, KeyValue{Key: word, Value: "1"})
	}
	return kva
}

func Reduce(key string, values []string) string {
	return fmt.Sprintf("%d", len(values))
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}

func splitWords(text string) []string {
	// Simple split by non-letter characters
	words := []string{}
	word := ""
	for _, r := range text {
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			word += string(r)
		} else if word != "" {
			words = append(words, word)
			word = ""
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}
