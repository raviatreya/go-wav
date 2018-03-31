package wav

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/davecgh/go-spew/spew"
)

const headerSize = 44

// WaveHeader is wave header
type WaveHeader struct {
	ChunkSize     uint32
	SubChunkSize  uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	SubChunk2Size uint32
}

func (parser *Wav) readRiffChunk(buffer []byte) ([]byte, error) {
	if "RIFF" != string(buffer[:4]) {
		return buffer, errors.New("This is not wav file")
	}
	fmt.Println("chunkID")
	spew.Dump(buffer[:4])

	buffer = buffer[4:]

	parser.header.ChunkSize = binary.LittleEndian.Uint32(buffer[:4])

	fmt.Println("\nchunkSize")
	spew.Dump(buffer[:4])
	spew.Dump(parser.header.ChunkSize)

	buffer = buffer[4:]

	if "WAVE" != string(buffer[:4]) {
		log.Fatal("This is not WAVE file!\n")
	}

	fmt.Println("\nformat")
	spew.Dump(buffer[:4])

	buffer = buffer[4:]
	return buffer, nil
}

func (parser *Wav) readFmtSubChunk(buffer []byte) ([]byte, error) {
	subchunkID := buffer[:4]
	if "fmt " != string(buffer[:4]) {
		return buffer, errors.New("This is not wav file")
	}

	fmt.Println("\nsubchunkID")
	spew.Dump(subchunkID)

	buffer = buffer[4:]

	parser.header.SubChunkSize = binary.LittleEndian.Uint32(buffer[:4])

	fmt.Println("\nsubChunkSize")
	spew.Dump(buffer[:4])
	spew.Dump(parser.header.SubChunkSize)

	buffer = buffer[4:]

	parser.header.AudioFormat = binary.LittleEndian.Uint16(buffer[:2])

	fmt.Println("\naudioFormat")
	spew.Dump(buffer[:2])
	spew.Dump(parser.header.AudioFormat)

	buffer = buffer[2:]

	parser.header.NumChannels = binary.LittleEndian.Uint16(buffer[:2])

	fmt.Println("\numChannels")
	spew.Dump(buffer[:2])
	spew.Dump(parser.header.NumChannels)

	buffer = buffer[2:]

	parser.header.SampleRate = binary.LittleEndian.Uint32(buffer[:4])

	fmt.Println("\nsampleRate")
	spew.Dump(buffer[:4])
	spew.Dump(parser.header.SampleRate)

	buffer = buffer[4:]

	parser.header.ByteRate = binary.LittleEndian.Uint32(buffer[:4])

	fmt.Println("byte rate")
	spew.Dump(buffer[:4])
	spew.Dump(parser.header.ByteRate)

	buffer = buffer[4:]

	parser.header.BlockAlign = binary.LittleEndian.Uint16(buffer[:2])

	fmt.Println("\nblockAlign")
	spew.Dump(buffer[:2])
	spew.Dump(parser.header.BlockAlign)

	buffer = buffer[2:]

	parser.header.BitsPerSample = binary.LittleEndian.Uint16(buffer[:2])

	fmt.Println("\n bitsPerSample")
	spew.Dump(buffer[:2])
	spew.Dump(parser.header.BitsPerSample)

	buffer = buffer[2:]

	return buffer, nil
}

func (parser *Wav) readDataSubChunk(buffer []byte) []byte {
	fmt.Println("\n subchunk2ID")
	spew.Dump(buffer[:4])

	buffer = buffer[4:]

	parser.header.SubChunk2Size = binary.LittleEndian.Uint32(buffer[:4])
	fmt.Println("\n subChunk2 size")
	spew.Dump(buffer[:4])
	spew.Dump(parser.header.SubChunk2Size)
	fmt.Println("")

	buffer = buffer[4:]

	return buffer
}

func (parser *Wav) Parse() error {
	buffer := make([]byte, headerSize)
	_, err := io.ReadAtLeast(parser.reader, buffer, headerSize)
	if err != nil {
		return err
	}

	buffer, err = parser.readRiffChunk(buffer)
	if err != nil {
		return err
	}
	buffer, err = parser.readFmtSubChunk(buffer)
	if err != nil {
		return err
	}
	buffer = parser.readDataSubChunk(buffer)

	/* Reset Read Position */
	parser.reader.Seek(0, 0)

	return nil
}

func (parser *Wav) GetHeader() *WaveHeader {
	return &parser.header
}
