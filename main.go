package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ColeVanOphem/golem/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOptions{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       "local/" + listenAddr + "_network",
		PathTransformFunc: DefaultPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":8000", "")
	s2 := makeServer(":8080", "")
	s3 := makeServer(":8081", ":8000", ":8080")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(1 * time.Second)

	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(1 * time.Second)

	go s3.Start()
	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("data_%d.txt", i)
		data := bytes.NewReader([]byte("my very important information"))

		s3.Store(key, data)

		if err := s3.store.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("[%s] finished using file (%s) with content \"%s\"\n", s3.Transport.Addr(), key, string(b))
		time.Sleep(1 * time.Second)
	}

	for i := 10; i < 20; i++ {
		key := fmt.Sprintf("data_%d.txt", i)
		data := bytes.NewReader([]byte("my very important information"))

		s2.Store(key, data)

		if err := s2.store.Delete(s2.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s2.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("[%s] finished using file (%s) with content \"%s\"\n", s2.Transport.Addr(), key, string(b))
		time.Sleep(1 * time.Second)
	}

	for i := 20; i < 30; i++ {
		key := fmt.Sprintf("data_%d.txt", i)
		data := bytes.NewReader([]byte("my very important information"))

		s1.Store(key, data)

		if err := s1.store.Delete(s1.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s1.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("[%s] finished using file (%s) with content \"%s\"\n", s1.Transport.Addr(), key, string(b))
		time.Sleep(1 * time.Second)
	}
}
