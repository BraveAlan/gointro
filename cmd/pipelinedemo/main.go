package main

import (
	"bufio"
	"imooc.com/sb/gointro/pipline"
	"os"
	"fmt"
)

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArraySource(3, 2, 5, 1, 0)),
		pipeline.InMemSort(pipeline.ArraySource(9, 10,50,25,32)),
	)
	for v := range p {
		fmt.Println(v)
	}
}

func main() {
	const filename = "small.in"
	const n = 64
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file) // bufio会读得快一点
	pipeline.WriterSink(writer, p)
	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v:= range p {
		fmt.Println(v)
		if count == 100 {
			break
		}
		count++
	}
}
