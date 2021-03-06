package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
)

// doMap manages one map task: it should read one of the input files
// (inFile), call the user-defined map function (mapF) for that file's
// contents, and partition mapF's output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTask int, // which map task this is
	inputFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)

	// mapF() is the map function provided by the application.
	// The first argument should be the input file name,
	// though the map function typically ignores it.
	// The second argument should be the entire input file contents.
	// mapF() returns a slice containing the key/value pairs for reduce
	mapF func(filename string, contents string) []KeyValue,
) {
	// Read the input file.
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Call mapF function for that file's content.
	// Need to convert data to string because returned data is []uint8
	kvs := mapF(inputFile, string(data))

	// Do the partition mapF's output into nReduce intermediate files.
	// Create a slice with nReduce number of Encoders .
	var encoders = make([]*json.Encoder, nReduce)
	var fd *os.File

	// Open an intermediate file for each reduce task and give it an Encoder.
	for i := 0; i < nReduce; i++ {
		// The file name includes both the map task number and the reduce task number.
		// Use the filename generated by reduceName(jobName, mapTask, r)
		// as the intermediate file for reduce task r.
		fileName := reduceName(jobName, mapTask, i)

		// Create a file if none exists, open it write-only, append while writing.
		fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		defer fd.Close()
		if err != nil {
			log.Fatal(err)
			return
		}

		encoders[i] = json.NewEncoder(fd)
	}

	// For every key-value pair of mapF output,
	// call ihash() on each key and mod nReduce to pick the immediate file r.
	// Use that file to encode the key value content
	for _, kv := range kvs {
		r := ihash(kv.Key) % nReduce
		err = encoders[r].Encode(kv)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
