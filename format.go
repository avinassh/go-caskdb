package caskdb

// format file provides encode/decode functions for serialisation and deserialisation
// operations
//
// format methods are generic and does not have any disk or memory specific code.
//
// The disk storage deals with bytes; you cannot just store a string or object without
// converting it to bytes. The programming languages provide abstractions where you
// don't have to think about all this when storing things in memory (i.e. RAM).
// Consider the following example where you are storing stuff in a hash table:
//
//    books = {}
//    books["hamlet"] = "shakespeare"
//    books["anna karenina"] = "tolstoy"
//
// In the above, the language deals with all the complexities:
//
//    - allocating space on the RAM so that it can store data of `books`
//    - whenever you add data to `books`, convert that to bytes and keep it in the memory
//    - whenever the size of `books` increases, move that to somewhere in the RAM so that
//      we can add new items
//
// Unfortunately, when it comes to disks, we have to do all this by ourselves, write
// code which can allocate space, convert objects to/from bytes and many other operations.
//
// This file has two functions which help us with serialisation of data.
//
//    encodeKV - takes the key value pair and encodes them into bytes
//    decodeKV - takes a bunch of bytes and decodes them into key value pairs
//
//**workshop note**
//
//For the workshop, the functions will have the following signature:
//
//    func encodeKV(timestamp uint32, key string, value string) (int, []byte)
//    func decodeKV(data []byte) (uint32, string, string)

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"time"
)

// headerSize specifies the total header size. Our key value pair, when stored on disk
// looks like this:
//
//	┌────────┬─────────────┬──────────┬───────────┬────────────┬───────┬─────────┐
//	|   crc  │ timestamp   │ expiry   | key_size  │ value_size │ key   │ value   │
//	└────────┴─────────────┴──────────┴───────────┴────────────┴───────┴─────────┴
//
// This is analogous to a typical database's row (or a record). The total length of
// the row is variable, depending on the contents of the key and value.
//
// The first five fields form the header:
//
//	┌────────────┬──────────────┬─────────────┬────────────────┐────────────────┐
//	|   crc(4B)  │ timestamp(4B)│ expiry(4B)  | key_size(4B)   │ value_size(4B) │
//	└────────────┴──────────────┴─────────────┘────────────────┘────────────────┘
//
// These four fields store unsigned integers of size 4 bytes, giving our header a
// fixed length of 16 bytes.
// crc(CheckSum) field stores the checksum to verify if the stored value is valid or not.
// Timestamp field stores the time the record we inserted in unix epoch seconds.
// Expiry field stores the time after which the record will expiry.
// Key size and Value size fields store the length of bytes occupied by the key and value. The maximum integer
// stored by 4 bytes is 4,294,967,295 (2 ** 32 - 1), roughly ~4.2GB. So, the size of
// each key or value cannot exceed this. Theoretically, a single row can be as large
// as ~8.4GB.
const headerSize = 20

// For deletion we will write a special "tombstone" value instead of actually deleting the key or storing this in the header.
const TombStoneVal = "tombstone"

// KeyEntry keeps the metadata about the KV, specially the position of
// the byte offset in the file. Whenever we insert/update a key, we create a new
// KeyEntry object and insert that into keyDir.
type KeyEntry struct {
	// Timestamp at which we wrote the KV pair to the disk. The value
	// is current time in seconds since the epoch.
	timestamp uint32
	// The position is the byte offset in the file where the data
	// exists
	position uint32
	// Total size of bytes of the value. We use this value to know
	// how many bytes we need to read from the file
	totalSize uint32
}

type Header struct {
	CheckSum   uint32
	ExpiryTime uint32
	TimeStamp  uint32
	KeySize    uint32
	ValueSize  uint32
}

type Record struct {
	Header Header
	Key    string
	Value  string
}

func NewKeyEntry(timestamp uint32, position uint32, totalSize uint32) KeyEntry {
	return KeyEntry{timestamp, position, totalSize}
}

func (h *Header) EncodeHeader(buf *bytes.Buffer) error {
	return binary.Write(buf, binary.LittleEndian, h)
}

func (h *Header) DecodeHeader(buf []byte) error {
	return binary.Read(bytes.NewReader(buf), binary.LittleEndian, h)
}

func (h *Header) CalculateCheckSum(value string) uint32 {
	return crc32.ChecksumIEEE([]byte(value))
}

func (r *Record) EncodeKV(buf *bytes.Buffer) error {
	err := r.Header.EncodeHeader(buf)
	buf.WriteString(r.Key)
	buf.WriteString(r.Value)
	return err
}

func (r *Record) DecodeKV(buf []byte) error {
	err := r.Header.DecodeHeader(buf[:headerSize])
	r.Key = string(buf[headerSize : headerSize+r.Header.KeySize])
	r.Value = string(buf[headerSize+r.Header.KeySize : headerSize+r.Header.KeySize+r.Header.ValueSize])
	return err
}

func (r *Record) VerifyCheckSum() bool {
	return crc32.ChecksumIEEE([]byte(r.Value)) == r.Header.CheckSum
}

func (r *Record) IsExpired() bool {
	if r.Header.ExpiryTime == 0 {
		//no expiry set
		return false
	}
	return uint32(time.Now().Unix()) > r.Header.ExpiryTime
}
