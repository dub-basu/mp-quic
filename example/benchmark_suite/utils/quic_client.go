package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"strconv"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "100.0.0.1:4242"

var message string

// var blockSize = 1024 * 100
var blockSizePointer = flag.Int("blocksize", 102400, "size of block in bytes")
var delayMilliseconds = flag.Int("delayMilli", 20, "delay in milliseconds")
var numBlocks = flag.Int("numBlocks", 10, "number of blocks to be requested")

func main() {
	flag.Parse()
	quicConfig := &quic.Config{
		CreatePaths:      false,
		HandshakeTimeout: 30 * time.Second,
		IdleTimeout:      60 * time.Second,
	}

	sess, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, quicConfig)
	if err != nil {
		fmt.Println("Error Connecting")
		panic(err)
	}

	stream1, err := sess.OpenStream()
	finished := make(chan bool)

	requestChunkedData(stream1, finished)
}

func requestChunkedData(stream quic.Stream, finished chan bool) error {
	start := time.Now()
	for i := 0; i < *numBlocks; i++ {
		bytesToRead := *blockSizePointer
		sizeBytesString := strconv.Itoa(bytesToRead) + "\n"
		stream.Write([]byte(sizeBytesString))

		buff := make([]byte, *blockSizePointer)
		for bytesToRead > 0 {
			bytesRead, err := stream.Read(buff)
			if err != nil {
				fmt.Println(err)
				fmt.Println("ERROR requestrequestChunkedData")
				return err
			}
			bytesToRead = bytesToRead - bytesRead
		}
		time.Sleep(time.Duration(*delayMilliseconds) * time.Millisecond)
	}
	print(time.Since(start))
	// print("Closing stream")
	stream.Close()
	// finished <- true
	return nil
}
