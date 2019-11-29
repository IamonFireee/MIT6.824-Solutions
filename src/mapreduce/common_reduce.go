package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// open the intermedFiles
	intermedFiles:=make([]*os.File,nMap)
	for i:=0;i<nMap;i++{
		fileName:=reduceName(jobName,i,reduceTaskNumber)
		var err error
		intermedFiles[i],err=os.Open(fileName)
		if err!=nil{
			log.Fatal("doReduce cant open intermedFile",err)
		}
	}
	// put k/v into kvsPairs
	kvsPairs:=make(map[string][]string)
	for _,intermedFile:=range intermedFiles{
		defer intermedFile.Close()
		dec:=json.NewDecoder(intermedFile)
		for{
			var kvs KeyValue
			err:=dec.Decode(&kvs)
			if err!=nil{
				break
			}
			kvsPairs[kvs.Key]=append(kvsPairs[kvs.Key],kvs.Value)
		}
	}
	keys:=make([]string,0,len(kvsPairs))
	for key,_:=range kvsPairs{
		keys=append(keys,key)
	}
	sort.Strings(keys)
	outPutFile,err:=os.Create(outFile)
	if err!=nil{
		log.Fatal("doReduce cant create outfile",err)
	}
	defer outPutFile.Close()
	enc:=json.NewEncoder(outPutFile)
	for _,k:=range keys{
		v:=reduceF(k,kvsPairs[k])
		err:=enc.Encode(KeyValue{k,v})
		if err!=nil{
			log.Fatal("doReduce cant encode k/v",err)
		}
	}
}
