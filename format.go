package caskdb

import (
	"bytes"
	"encoding/binary"
)

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

// headerSize specifies the total header size. Our key value pair, when stored on disk
// looks like this:
//
//	┌────────┬───────────┬───────────┬─────────────┬───────┬───────┐
//	│ meta   | timestamp │ key_size  | value_size  │  key  │ value │
//	└────────┴───────────┴───────────┴─────────────┴───────┘───────┘
//
// This is analogous to a typical database's row (or a record). The total length of
// the row is variable, depending on the contents of the key and value.
//
// The first four fields form the header:
//
//	┌───────────┬────────────────┬───────────────┬─────────────────┐
//	│ meta(1B)  │ timestamp(4B)  | key_size(4B)  │ value_size(4B)  │
//	└───────────┴────────────────┴───────────────┴─────────────────┘
//
// The first byte stores the metadata about the kv record.
// The rest three fields store unsigned integers of size 4 bytes giving our header a fixed length of 14 bytes.
// Timestamp field stores the time the record we inserted in unix epoch seconds.
// Key size and value size fields store the length of bytes occupied by the key and value.
// We can use it for marking a record as tombstone by setting its MSB to 1.
const headerSize = 13

// The maximum integer stored by 4 bytes is 4,294,967,295 (2 ** 32 - 1), roughly ~4.2GB.
// So, the size of each key or value cannot exceed this. Theoretically, a single row can be as large as ~8.4GB.
const (
	MaxKeySize   = 1<<32 - 1
	MaxValueSize = 1<<32 - 1
)

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
	Meta      uint8
	TimeStamp uint32
	KeySize   uint32
	ValueSize uint32
}

type Record struct {
	Header     Header
	Key        string
	Value      string
	RecordSize uint32
}

func NewKeyEntry(timestamp uint32, position uint32, totalSize uint32) KeyEntry {
	return KeyEntry{timestamp, position, totalSize}
}

func (h *Header) EncodeHeader(buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.LittleEndian, &h.Meta)
	binary.Write(buf, binary.LittleEndian, &h.TimeStamp)
	binary.Write(buf, binary.LittleEndian, &h.KeySize)
	binary.Write(buf, binary.LittleEndian, &h.ValueSize)
	return err
}

func (h *Header) DecodeHeader(buf []byte) error {
	err := binary.Read(bytes.NewReader(buf[0:1]), binary.LittleEndian, &h.Meta)
	binary.Read(bytes.NewReader(buf[1:5]), binary.LittleEndian, &h.TimeStamp)
	binary.Read(bytes.NewReader(buf[5:9]), binary.LittleEndian, &h.KeySize)
	binary.Read(bytes.NewReader(buf[9:13]), binary.LittleEndian, &h.ValueSize)
	return err
}

func (h *Header) MarkTombStone() {
	// setting the MSB to 1
	h.Meta = h.Meta | (1 << 7)
}

func (h *Header) IsTombStone() bool {
	// checking if MSB is set to 1
	return (1 << 7) == (h.Meta & (1 << 7))
}

func NewHeader(buf []byte) (*Header, error) {
	h := &Header{}
	err := h.DecodeHeader(buf)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (r *Record) EncodeKV(buf *bytes.Buffer) error {
	r.Header.EncodeHeader(buf)
	buf.WriteString(r.Key)
	_, err := buf.Write([]byte(r.Value))
	return err
}

func (r *Record) DecodeKV(buf []byte) error {
	err := r.Header.DecodeHeader(buf[:headerSize])
	r.Key = string(buf[headerSize : headerSize+r.Header.KeySize])
	r.Value = string(buf[headerSize+r.Header.KeySize : headerSize+r.Header.KeySize+r.Header.ValueSize])
	r.RecordSize = headerSize + r.Header.KeySize + r.Header.ValueSize
	return err
}

func (r *Record) Size() uint32 {
	return r.RecordSize
}
