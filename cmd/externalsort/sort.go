package main

import (
	"bufio"
	"fmt"
	"imooc.com/sb/gointro/pipeline"
	"os"
	"strconv"
)

func main() {
	// 这个做的是单机版的，用来模拟数据量为80G以上的外部排序的情况
	// 如果数据量小，用这种方法速度肯定是慢的，因为涉及到goroutine之间的通信、等待
	//p := createPipeline("large.in", 800000000, 4)

	p := createNetWorkPipeline("large.in", 800000000, 4)

	writeToFile(p, "large.out")
	printFile("large.out")

}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}


func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)

}

// createPipeline filename文件名，fileSize文件大小，chunkCount把文件分成多少块
func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	pipeline.Init()

	chunkSize := fileSize / chunkCount

	sortResults := []<-chan int{}

	for i:=0;i<chunkCount;i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemSort(source))
	}

	return pipeline.MergeN(sortResults...)
}

// createPipeline filename文件名，fileSize文件大小，chunkCount把文件分成多少块
func createNetWorkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	pipeline.Init()

	chunkSize := fileSize / chunkCount

	sortAddr := []string{}

	for i:=0;i<chunkCount;i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(7000 + i)
		pipeline.NetWorkSink(addr, pipeline.InMemSort(source)) // 开了一个server，等连接
		sortAddr = append(sortAddr, addr)
	}

	sortResults := []<-chan int{}
	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetWorkSource(addr)) // 开了一个client，连接Server
	}

	return pipeline.MergeN(sortResults...)
}