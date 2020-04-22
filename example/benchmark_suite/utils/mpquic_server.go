package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"strconv"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "0.0.0.0:4242"

var fakeBuff = make([]byte, 1024*1024*100) // 100 MB

func main() {
	for i := 0; i < cap(fakeBuff); i++ {
		fakeBuff[i] = 'A'
	}

	quicConfig := &quic.Config{
		CreatePaths:      true,
		HandshakeTimeout: 30 * time.Second,
		IdleTimeout:      60 * time.Second,
	}
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), quicConfig)

	if err != nil {
		print("Unable to make socket")
		panic(err)
	}

	sess, err := listener.Accept()
	if err != nil {
		print("Couldnt make Session")
	}
	fmt.Print("Connection made with ")
	fmt.Println(sess.RemoteAddr())

	finished := make(chan bool)

	stream1, err := sess.AcceptStream()
	if err != nil {
		fmt.Println("Couldnt make Stream 1")
	}
	go func() { handleStream(stream1, finished) }()

	<-finished

}

func handleStream(stream quic.Stream, finished chan bool) error {
	buff := make([]byte, 1000)

	for {
		var chunkSizeString = ""
		bytesRead, err := stream.Read(buff)
		if err != nil {
			break
		}

		for buff[bytesRead-1] != '\n' {
			chunkSizeString += string(buff[:bytesRead])
		}
		chunkSizeString += string(buff[:bytesRead-1]) // -1 to ignore newline

		chunkSize, err := strconv.Atoi(chunkSizeString)

		if err != nil {
			fmt.Println("Invalid chunk size sent")
		}
		fmt.Printf("Chunk size : %d\n", chunkSize)

		if chunkSize == 0 {
			break
		}

		sendChunckOnStream(stream, chunkSize)

	}

	finished <- true
	return nil
}

func sendChunckOnStream(stream quic.Stream, size int) {
	for bytesToSend := size; bytesToSend > 0; {
		bytesSent, _ := stream.Write(fakeBuff[:bytesToSend])
		bytesToSend = bytesToSend - bytesSent
	}

}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
