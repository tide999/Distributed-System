package main

import (
	"fmt"
	"mapreduce"
	"os"
	"strings"
	"unicode"
	"strconv"
)

// The mapping function is called once for each piece of the input.
// In this framework, the key is the name of the file that is being processed,
// and the value is the file's contents. The return value should be a slice of
// key/value pairs, each represented by a mapreduce.KeyValue.
func mapF(document string, value string) (res []mapreduce.KeyValue) {
	// TODO: you have to write this function
	//arr := strings.FieldsFunc(value, unicode.IsSpace)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	arr := strings.FieldsFunc(value, f)
	res = make([]mapreduce.KeyValue, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = mapreduce.KeyValue{arr[i], "1"}
	}
	return res
}

// The reduce function is called once for each key generated by Map, with a
// list of that key's string value (merged across all inputs). The return value
// should be a single output value for that key.
func reduceF(key string, values []string) string {
	// TODO: you also have to write this function
	sum := 0
	tmp := 0
	var err error
	for _, v := range values {
		tmp, err = strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reduceF string to int Error: %s\n", err)
			panic(err.Error())
		}
		sum += tmp
	}
	return strconv.Itoa(sum)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *mapreduce.Master
		if os.Args[2] == "sequential" {
			mr = mapreduce.Sequential("wcseq", os.Args[3:], 3, mapF, reduceF)
		} else {
			mr = mapreduce.Distributed("wcseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		mapreduce.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100)
	}
}
