package mapreduce

import (
	"encoding/json"
    "log"
	"os"
	"sort"
)

// doReduce manages one reduce task: it should read the intermediate
// files for the task, sort the intermediate key/value pairs by key,
// call the user-defined reduce function (reduceF) for each key, and
// write reduceF's output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outputFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
    // Make a decoder slice with the number of map tasks.
	var decoders = make([]*json.Decoder, nMap)

    // Read the intermediate file for each map task
	for i := 0; i < nMap; i++ {
	    // reduceName(jobName, m, reduceTask) yields the file name from map task m.
		fileName := reduceName(jobName, i, reduceTask)

        // Open read-only.
		fd, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
		defer fd.Close()
		if err != nil {
            log.Fatal(err)
			return
		}

        // Call a decoder
		decoders[i] = json.NewDecoder(fd)
	}

	// Unmarshal all intermediate files and collate key-values
	kvs := make(map[string][]string)
	for i := 0; i < nMap; i++ {
		var kv *KeyValue
		for {
			err := decoders[i].Decode(&kv)
			if err != nil {
				break
			}
			kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
		}
	}

	// Sort the intermediate key/value pairs by key,
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create an output file
	fd, err := os.OpenFile(outputFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer fd.Close()
	if err != nil {
        log.Fatal(err)
		return
	}

    // Call the reduce function (reduceF) for each key and write its output to disk
	encoder := json.NewEncoder(fd)
	for _, key := range keys {
		encoder.Encode(KeyValue{key, reduceF(key, kvs[key])})
	}
}
