package hack

import (
	"io"
	"net"
	"testing"
)

func BenchmarkWithHack(b *testing.B) {
	client, srv := net.Pipe()
	done := make(chan struct{})
	req := []byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\n")
	read := make([]byte, 4096)
	b.SetBytes(int64(len(req) * 30))

	l := MalformedHostHeaderOverrideConn{client, true}
	go func() {
		for {
			if _, err := srv.Write(req); err != nil {
				srv.Close()
				break
			}
			l.first = true // make sure each subsequent run uses the hack parsing
		}
		close(done)
	}()

	for i := 0; i < b.N; i++ {
		for i := 0; i < 30; i++ {
			if n, err := l.Read(read); err != nil && err != io.EOF || n != len(req) {
				b.Fatalf("read: %d - %d, err: %v\n%s", n, len(req), err, string(read[:n]))
			}
		}
	}
	l.Close()
	<-done
}

func BenchmarkNoHack(b *testing.B) {
	client, srv := net.Pipe()
	done := make(chan struct{})
	req := []byte("GET /foo\nHost: /var/run/docker.sock\nUser-Agent: Docker\n")
	read := make([]byte, 4096)
	b.SetBytes(int64(len(req) * 30))

	go func() {
		for {
			if _, err := srv.Write(req); err != nil {
				srv.Close()
				break
			}
		}
		close(done)
	}()

	for i := 0; i < b.N; i++ {
		for i := 0; i < 30; i++ {
			if _, err := client.Read(read); err != nil && err != io.EOF {
				b.Fatal(err)
			}
		}
	}
	client.Close()
	<-done
}
