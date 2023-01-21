# Hints

_contains spoilers, proceed with caution!_

## Tests
Below is the order in which you should pass the tests:
	
1. `Test_encodeHeader`
2. `Test_encodeKV`
3. `TestDiskStore_Set`
4. `TestDiskStore_Get` and `TestDiskStore_GetInvalid`
5. `TestDiskStore_SetWithPersistence`

## Tasks

I have relisted the tasks here again, but with more details. If you have difficulty understanding the steps, the following should help you.

### Read the paper
The word 'paper' might scare you, but [Bitcask's paper](https://riak.com/assets/bitcask-intro.pdf) is very approachable. It is only six pages long, half of them being diagrams. 

### Header

| Test | Test_encodeHeader |
|------|-------------------|

The next step is to implement a fixed-sized header similar to Bitcask. Every record in our DB contains a header holding metadata and helps our DB read values from the disk. The DB will read a bunch of bytes from the disk, so we need information on how many bytes to read and from which byte offset. 
	 
**Some more details:**	

The header contains three fields timestamp, key size, and value size. 

| Field      | Type | Size |
|------------|------|------|
| timestamp  | int  | 4B   |
| key_size   | int  | 4B   |
| value_size | int  | 4B   |
				
We need to implement a function which takes all these three fields and serialises them to bytes. The function signature looks like this:

```go
func encodeHeader(timestamp uint32, keySize uint32, valueSize uint32)
```

Then we also need to write the opposite of the above:

```go
func decodeHeader(header []byte) (uint32, uint32, uint32)
```

**More Hints:**
- Read this [comment](https://github.com/avinassh/go-caskdb/blob/0ae4fab/format.go#L3,L37) to understand why do we need serialiser methods 
- Not sure how to come up with a file format? Read the comment in the [format module](https://github.com/avinassh/go-caskdb/blob/0ae4fab/format.go#L41,L63)

### Key Value Serialisers

| Test | Test_encodeKV |
|------|---------------|

Now we will write encode and decode methods for key and value. 

The method signatures: 
```go
func encodeKV(timestamp uint32, key string, value string) (int, []byte)
func decodeKV(data []byte) (uint32, string, string)
```

Note that `encodeKV` method returns the bytes and the bytes' size.

### Storing to Disk

| Test | TestDiskStore_Set |
|------|-------------------|

This step involves figuring out the persistence layer, saving the data to the disk, and keeping the pointer to the inserted record in the memory. 

So, implement the `DiskStore.Set` class in `disk_store.go`

**Hints:**
- Some meta info on the DiskStore and inner workings of the DiskStore are [here](https://github.com/avinassh/go-caskdb/blob/0ae4fab/disk_store.go#L24,L63).

### Start up tasks

| Test | TestDiskStore_SetWithPersistence |
|------|----------------------------------|

DiskStore is a persistent key-value store, so we need to load the existing keys into the `keyDir` at the start of the database. 
