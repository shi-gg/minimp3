package minimp3

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
)

func TestRaceStress(t *testing.T) {
	t.Parallel()

	mp3, err := os.ReadFile("./test.mp3")
	if err != nil {
		t.Fatal(err)
	}

	const concurrency = 32
	const chunkSize = 128

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			reader := &chunkedReader{data: mp3, chunkSize: chunkSize}
			dec, err := NewDecoder(reader)
			if err != nil {
				t.Error(err)
				return
			}
			defer dec.Close()

			<-dec.Started()

			buf := make([]byte, 4096)
			for {
				_, err := dec.Read(buf)
				if err == io.EOF {
					return
				}
				if err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()
}

type chunkedReader struct {
	data      []byte
	chunkSize int
	pos       int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	end := r.pos + r.chunkSize
	if end > len(r.data) {
		end = len(r.data)
	}
	n := copy(p, r.data[r.pos:end])
	r.pos += n
	return n, nil
}

func BenchmarkDecoder(b *testing.B) {
	mp3, err := os.ReadFile("./test.mp3")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(mp3)
		dec, err := NewDecoder(reader)
		if err != nil {
			b.Fatal(err)
		}
		<-dec.Started()

		buf := make([]byte, 4096)
		for {
			_, err := dec.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
		dec.Close()
	}
}
