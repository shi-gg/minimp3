package minimp3

import (
	"log"
	"os"
	"testing"
)

func TestIssue21DataRace(t *testing.T) {
	file, err := os.Open("test.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dec, err := NewDecoder(file)
	if err != nil {
		t.Fatal(err)
	}
	defer dec.Close()

	<-dec.Started()

	dec.decoderLocker.Lock()
	sampleRate := dec.SampleRate
	channels := dec.Channels
	dec.decoderLocker.Unlock()

	log.Printf("Convert audio sample rate: %d, channels: %d\n", sampleRate, channels)
}
