package mapreduce

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
)

// doMap manages one map task: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
/*
	在这里我将文件一次读取，一次传给mapF()，这可能不是一种好的方式，
	如果考虑文件的大小及内存占用，也许将文件分块并分几次调用mapF()可能是更好的方式
 */
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// print details of this task
	fmt.Printf("---Map: job name = %s, input file = %s, map task id = %d, nReduce = %d---\n",
		jobName, inFile, mapTaskNumber, nReduce)
	// open file
	if fileBytes,err:=ioutil.ReadFile(inFile);err==nil{
		// create intermed files
		intermedFiles:=make([]*os.File,nReduce)
		for i:=0;i<nReduce;i++{
			intermedFileName:=reduceName(jobName,mapTaskNumber,i)
			intermedFiles[i],err=os.Create(intermedFileName)
			if err!=nil{
				log.Fatal("doMap err when create intermed files",err)
			}
		}
		// main part of this func
		kvPairs:=mapF(inFile,string(fileBytes))
		for _,kv:=range kvPairs{
			index:=ihash(kv.Key)%nReduce
			enc:=json.NewEncoder(intermedFiles[index])
			enc.Encode(kv)
		}
		for _,intermedFile:=range intermedFiles{
			err:=intermedFile.Close()
			if err!=nil{
				log.Fatal("doMap err when close the output file:",err)
			}
		}
	}else{
		// cant open this file
		log.Fatal("doMap err when open the input file:",err)
	}
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
