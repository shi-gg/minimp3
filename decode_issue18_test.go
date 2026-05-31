package minimp3

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

const delayDuration = 100 * time.Millisecond

func TestIssue18DelayPutData(t *testing.T) {
	t.Parallel()

	reader, writer := io.Pipe()
	dec, err := NewDecoder(reader)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()

	go func() {
		time.Sleep(delayDuration)
		file, err := os.ReadFile("test.mp3")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = io.Copy(writer, bytes.NewReader(file))
		if err != nil {
			t.Error(err)
			return
		}
		writer.Close()
	}()

	data, err := io.ReadAll(dec)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 44928 {
		t.Errorf("unexpected pcm length: got %d, want 44928", len(data))
	}
}

func TestIssue18GracefulExit(t *testing.T) {
	t.Parallel()

	reader, writer := io.Pipe()
	dec, err := NewDecoder(reader)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()

	go func() {
		time.Sleep(delayDuration)
		writer.Close()
	}()

	data, err := io.ReadAll(dec)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Errorf("expected no data on graceful exit, got %d bytes", len(data))
	}
}
