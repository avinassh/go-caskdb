package caskdb

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"time"
)

// defaultWhence helps us with `file.Seek` method to move our cursor to certain byte offset for read
// or write operations. The method takes two parameters file.Seek(offset, whence).
// The offset says the byte offset and whence says the direction:
//
// whence 0 - beginning of the file
// whence 1 - current cursor position
// whence 2 - end of the file
//
// read more about it here:
// https://pkg.go.dev/os#File.Seek
const defaultWhence = 0

// DiskStore is a Log-Structured Hash Table as described in the BitCask paper. We
// keep appending the data to a file, like a log. DiskStorage maintains an in-memory
// hash table called KeyDir, which keeps the row's location on the disk.
//
// The idea is simple yet brilliant:
//   - Write the record to the disk
//   - Update the internal hash table to point to that byte offset
//   - Whenever we get a read request, check the internal hash table for the address,
//     fetch that and return
//
// KeyDir does not store values, only their locations.
//
// The above approach solves a lot of problems:
//   - Writes are insanely fast since you are just appending to the file
//   - Reads are insanely fast since you do only one disk seek. In B-Tree backed
//     storage, there could be 2-3 disk seeks
//
// However, there are drawbacks too:
//   - We need to maintain an in-memory hash table KeyDir. A database with a large
//     number of keys would require more RAM
//   - Since we need to build the KeyDir at initialisation, it will affect the startup
//     time too
//   - Deleted keys need to be purged from the file to reduce the file size
//
// Read the paper for more details: https://riak.com/assets/bitcask-intro.pdf
//
// DiskStore provides two simple operations to get and set key value pairs. Both key
// and value need to be of string type, and all the data is persisted to disk.
// During startup, DiskStorage loads all the existing KV pair metadata, and it will
// throw an error if the file is invalid or corrupt.
//
// Note that if the database file is large, the initialisation will take time
// accordingly. The initialisation is also a blocking operation; till it is completed,
// we cannot use the database.
//
// Typical usage example:
//
//		store, _ := NewDiskStore("books.db")
//	   	store.Set("othello", "shakespeare")
//	   	author := store.Get("othello")
type DiskStore struct {
	// file object pointing the file_name
	file *os.File
	// current cursor position in the file where the data can be written
	writePosition int
	// keyDir is a map of key and KeyEntry being the value. KeyEntry contains the position
	// of the byte offset in the file where the value exists. key_dir map acts as in-memory
	// index to fetch the values quickly from the disk
	keyDir map[string]KeyEntry
}

func isFileExists(fileName string) bool {
	// https://stackoverflow.com/a/12518877
	if _, err := os.Stat(fileName); err == nil || errors.Is(err, fs.ErrExist) {
		return true
	}
	return false
}

func NewDiskStore(fileName string) (*DiskStore, error) {
	ds := &DiskStore{keyDir: make(map[string]KeyEntry)}
	// if the file exists already, then we will load the key_dir
	if isFileExists(fileName) {
		ds.initKeyDir(fileName)
	}
	// we open the file in following modes:
	//	os.O_APPEND - says that the writes are append only.
	// 	os.O_RDWR - says we can read and write to the file
	// 	os.O_CREATE - creates the file if it does not exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	ds.file = file
	return ds, nil
}

func (d *DiskStore) Get(key string) (interface{}, error) {
	// Get retrieves the value from the disk and returns. If the key does not
	// exist then it returns an empty string
	//
	// How get works?
	//	1. Check if there is any KeyEntry record for the key in keyDir
	//	2. Return an empty string if key doesn't exist or if the key has been deleted
	//	3. If it exists, then read KeyEntry.totalSize bytes starting from the
	//     KeyEntry.position from the disk
	//	4. Decode the bytes into valid KV pair and return the value
	//
	kEntry, ok := d.keyDir[key]
	if !ok {
		return "", ErrKeyNotExist
	}
	// move the current pointer to the right offset
	_, err := d.file.Seek(int64(kEntry.position), defaultWhence)
	if err != nil {
		return "", ErrSeekFailed
	}

	data := make([]byte, kEntry.totalSize)
	_, err = io.ReadFull(d.file, data)
	if err != nil {
		return "", ErrReadFailed
	}

	result := &Record{}
	result.DecodeKV(data)

	//decode value
	value, err := result.DecodeValue()
	if err != nil {
		log.Fatalf("error in decoding the value: %v", err)
		return "", ErrDecodingFailed
	}

	//check for expiry
	if result.IsExpired() {
		return "", ErrKeyNotFound
	}

	//check if checkSum matches and we dont have any corrupt value
	if !result.VerifyCheckSum() {
		return "", ErrKeyValueCorrupted
	}

	//check if its tombestone value
	if value == TombStoneVal {
		return TombStoneVal, ErrKeyNotFound
	}

	return value, nil
}

func (d *DiskStore) Set(key string, value interface{}) error {
	// Set stores the key and value on the disk
	//
	// The steps to save a KV to disk is simple:
	// 1. Encode the KV into bytes
	// 2. Write the bytes to disk by appending to the file
	// 3. Update KeyDir with the KeyEntry of this key

	//prepare kv record
	h := Header{
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
	}
	r := &Record{
		Header: h,
		Key:    key,
	}
	err := r.EncodeValue(value)
	if err != nil {
		log.Fatalf("error in encoding value: %v", err)
		return ErrInvalidValue
	}
	r.Header.ValueSize = uint32(len(r.Value))
	r.Header.CheckSum = r.Header.CalculateCheckSum(r.Value)

	buf := bytes.NewBuffer(make([]byte, headerSize))
	err = r.EncodeKV(buf)
	if err != nil {
		log.Fatalf("error in encoding kv record: %v", err)
		return ErrEncodingFailed
	}

	d.write(buf.Bytes())
	size := headerSize + r.Header.KeySize + r.Header.ValueSize
	d.keyDir[key] = NewKeyEntry(r.Header.TimeStamp, uint32(d.writePosition), size)
	// update last write position, so that next record can be written from this point
	d.writePosition += int(size)
	return nil
}

func (d *DiskStore) SetX(key string, value interface{}, expiry time.Duration) error {
	// Set but with expiry

	//prepare  kv record
	h := Header{
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
	}
	r := &Record{
		Header: h,
		Key:    key,
	}
	err := r.EncodeValue(value)
	if err != nil {
		log.Fatalf("error in encoding the value: %v", err)
		return ErrInvalidValue
	}
	r.Header.ValueSize = uint32(len(r.Value))
	r.Header.CheckSum = r.Header.CalculateCheckSum(r.Value)
	r.Header.ExpiryTime = uint32(time.Now().Add(expiry).Unix())

	buf := bytes.NewBuffer(make([]byte, headerSize))
	err = r.EncodeKV(buf)
	if err != nil {
		log.Fatalf("error in encoding kv record: %v", err)
		return ErrEncodingFailed
	}

	d.write(buf.Bytes())
	size := headerSize + r.Header.KeySize + r.Header.ValueSize
	d.keyDir[key] = NewKeyEntry(r.Header.TimeStamp, uint32(d.writePosition), size)
	// update last write position, so that next record can be written from this point
	d.writePosition += int(size)
	return nil
}

func (d *DiskStore) Delete(key string) error {
	// for delete operation, simply write a special tombstone value
	h := Header{
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
	}
	r := &Record{
		Header: h,
		Key:    key,
	}
	err := r.EncodeValue(TombStoneVal)
	if err != nil {
		log.Fatalf("error in encoding the value: %v", err)
		return ErrInvalidValue
	}
	r.Header.ValueSize = uint32(len(r.Value))
	r.Header.CheckSum = r.Header.CalculateCheckSum(r.Value)

	buf := bytes.NewBuffer(make([]byte, headerSize))
	err = r.EncodeKV(buf)
	if err != nil {
		log.Fatalf("error in encoding the kv record: %v", err)
		return ErrEncodingFailed
	}

	d.write(buf.Bytes())
	size := headerSize + r.Header.KeySize + r.Header.ValueSize
	// key is already present, it will update with our new value
	d.keyDir[key] = NewKeyEntry(r.Header.TimeStamp, uint32(d.writePosition), size)
	d.writePosition += int(size)

	return nil
}

func (d *DiskStore) Close() bool {
	// before we close the file, we need to safely write the contents in the buffers
	// to the disk. Check documentation of DiskStore.write() to understand
	// following the operations
	// TODO: handle errors
	err := d.file.Sync()
	if err != nil {
		log.Fatalf("error while doing a fsync: %v", err)
	}
	if err := d.file.Close(); err != nil {
		log.Printf("error in closing the file: %v", err)
		return false
	}
	return true
}

func (d *DiskStore) write(data []byte) {
	// saving stuff to a file reliably is hard!
	// if you would like to explore and learn more, then
	// start from here: https://danluu.com/file-consistency/
	// and read this too: https://lwn.net/Articles/457667/
	if _, err := d.file.Write(data); err != nil {
		panic(err)
	}
	// calling fsync after every write is important, this assures that our writes
	// are actually persisted to the disk
	if err := d.file.Sync(); err != nil {
		panic(err)
	}
}

func (d *DiskStore) initKeyDir(existingFile string) {
	// we will initialise the keyDir by reading the contents of the file, record by
	// record. As we read each record, we will also update our keyDir with the
	// corresponding KeyEntry
	//
	// NOTE: this method is a blocking one, if the DB size is yuge then it will take
	// a lot of time to startup
	file, _ := os.Open(existingFile)
	defer file.Close()
	for {
		header := make([]byte, headerSize)
		_, err := io.ReadFull(file, header)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading from the file: %v", err)
		}
		h := &Header{}
		h.DecodeHeader(header)
		key := make([]byte, h.KeySize)
		value := make([]byte, h.ValueSize)

		_, err = io.ReadFull(file, key)
		if err != nil {
			log.Fatalf("error while reading the key: %v", err)
			break
		}

		_, err = io.ReadFull(file, value)
		if err != nil {
			log.Fatalf("error while reading the value: %v", err)
			break
		}

		totalSize := headerSize + h.KeySize + h.ValueSize
		d.keyDir[string(key)] = NewKeyEntry(h.TimeStamp, uint32(d.writePosition), totalSize)
		d.writePosition += int(totalSize)
	}
}

func (d *DiskStore) ListKeys(existingFile string) <-chan Record {
	// read k,v one by one in the chan
	file, _ := os.Open(existingFile)
	result := make(chan Record)

	go func(file *os.File) {
		for {
			header := make([]byte, headerSize)
			_, err := io.ReadFull(file, header)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading from the file: %v", err)
			}
			h := &Header{}
			h.DecodeHeader(header)
			key := make([]byte, h.KeySize)
			value := make([]byte, h.ValueSize)

			_, err = io.ReadFull(file, key)
			if err != nil {
				log.Fatalf("error while reading the key: %v", err)
				break
			}

			_, err = io.ReadFull(file, value)
			if err != nil {
				log.Fatalf("error while reading the value: %v", err)
				break
			}

			r := Record{
				Header: *h,
				Key:    string(key),
				Value:  value,
			}
			result <- r
		}

		close(result)
		file.Close()
	}(file)

	return result
}
