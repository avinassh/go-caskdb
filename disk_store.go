package caskdb

import (
	"os"
	"time"
)

const defaultWhence = 0

type DiskStore struct {
	fileName      string
	file          *os.File
	writePosition int
	keyDir        map[string]KeyEntry
}

func NewDiskStore(fileName string) *DiskStore {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return &DiskStore{fileName, file, 0, make(map[string]KeyEntry)}
}

func (d *DiskStore) Get(key string) string {
	kv, ok := d.keyDir[key]
	if !ok {
		return ""
	}
	// TODO: handle errors
	d.file.Seek(int64(kv.position), defaultWhence)
	data := make([]byte, kv.totalSize)
	// TODO: handle errors
	n, err := d.file.Read(data)
	if n != int(kv.totalSize) || err != nil {
		panic("read error")
	}
	_, _, value := decodeKV(data)
	return value
}

func (d *DiskStore) Set(key string, value string) {
	timestamp := uint32(time.Now().Unix())
	size, data := encodeKV(timestamp, key, value)
	d.write(data)
	d.keyDir[key] = NewKeyEntry(timestamp, uint32(d.writePosition), uint32(size))
	d.writePosition += size
}

func (d *DiskStore) Close() bool {
	// TODO: handle errors
	d.file.Sync()
	if err := d.file.Close(); err != nil {
		// TODO: log the error
		return false
	}
	return true
}

func (d *DiskStore) write(data []byte) {
	// TODO: handle errors
	_, err := d.file.Write(data)
	if err != nil {
		panic(err)
	}
	// TODO: handle errors
	d.file.Sync()
}

func (d *DiskStore) name() {

}
