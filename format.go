package caskdb

import "encoding/binary"

const headerSize = 12

type KeyEntry struct {
	timestamp uint32
	position  uint32
	totalSize uint32
}

func NewKeyEntry(timestamp uint32, position uint32, totalSize uint32) KeyEntry {
	return KeyEntry{timestamp, position, totalSize}
}

func encodeHeader(timestamp uint32, keySize uint32, valueSize uint32) []byte {
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[0:4], timestamp)
	binary.LittleEndian.PutUint32(header[4:8], keySize)
	binary.LittleEndian.PutUint32(header[8:12], valueSize)
	return header
}

func decodeHeader(header []byte) (uint32, uint32, uint32) {
	timestamp := binary.LittleEndian.Uint32(header[0:4])
	keySize := binary.LittleEndian.Uint32(header[4:8])
	valueSize := binary.LittleEndian.Uint32(header[8:12])
	return timestamp, keySize, valueSize
}

func encodeKV(timestamp uint32, key string, value string) (int, []byte) {
	header := encodeHeader(timestamp, uint32(len(key)), uint32(len(value)))
	data := append([]byte(key), []byte(value)...)
	return headerSize + len(data), append(header, data...)
}

func decodeKV(data []byte) (uint32, string, string) {
	timestamp, keySize, valueSize := decodeHeader(data[0:headerSize])
	key := string(data[headerSize : headerSize+keySize])
	value := string(data[headerSize+keySize : headerSize+keySize+valueSize])
	return timestamp, key, value
}
