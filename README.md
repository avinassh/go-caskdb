![logo](https://github.com/avinassh/py-caskdb/raw/master/assets/logo.svg)
# CaskDB - Disk based Log Structured Hash Table Store

![made-with-go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)
[![build](https://github.com/avinassh/py-caskdb/actions/workflows/build.yml/badge.svg)](https://github.com/avinassh/go-caskdb/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/avinassh/py-caskdb/branch/master/graph/badge.svg?token=9SA8Q4L7AZ)](https://codecov.io/gh/avinassh/py-caskdb)
![GitHub License](https://img.shields.io/github/license/avinassh/go-caskdb)
[![twitter@iavins](https://img.shields.io/twitter/follow/iavins?style=social)](https://twitter.com/iavins)

![architecture](https://user-images.githubusercontent.com/640792/167299554-0fc44510-d500-4347-b680-258e224646fa.png)

CaskDB is a  disk-based, embedded, persistent, key-value store based on the [Riak's bitcask paper](https://riak.com/assets/bitcask-intro.pdf), written in Go. It is more focused on the educational capabilities than using it in production. The file format is platform, machine, and programming language independent. Say, the database file created from Go on macOS should be compatible with Rust on Windows.

This project aims to help anyone, even a beginner in databases, build a persistent database in a few hours. There are no external dependencies; only the Go standard library is enough.

If you are interested in writing the database yourself, head to the workshop section.

## Features
- Low latency for reads and writes
- High throughput
- Easy to back up / restore
- Simple and easy to understand
- Store data much larger than the RAM

## Limitations
Most of the following limitations are of CaskDB. However, there are some due to design constraints by the Bitcask paper.

- Single file stores all data, and deleted keys still take up the space
- CaskDB does not offer range scans
- CaskDB requires keeping all the keys in the internal memory. With a lot of keys, RAM usage will be high
- Slow startup time since it needs to load all the keys in memory

## Community

[![CaskDB Discord](https://img.shields.io/discord/851000331721900053)](https://discord.gg/HzthUYkrPp)

Consider joining the Discord community to build and learn KV Store with peers.


## Dependencies
CaskDB does not require any external libraries to run. Go standard library is enough.

## Installation
```shell
go get github.com/avinassh/go-caskdb
```

## Usage

```go
store, _ := NewDiskStore("books.db")
store.Set("othello", "shakespeare")
author := store.Get("othello")
```

## Cask DB (Python)
This project is a Go version of the [same project in Python](https://github.com/avinassh/py-caskdb). 

## Prerequisites
The workshop is for intermediate-advanced programmers. Knowing basics of Go helps, and you can build the database in any language you wish.

Not sure where you stand? You are ready if you have done the following in any language:
- If you have used a dictionary or hash table data structure
- Converting an object (class, struct, or dict) to JSON and converting JSON back to the things
- Open a file to write or read anything. A common task is dumping a dictionary contents to disk and reading back

## Workshop
**NOTE:** I don't have any [workshops](workshop.md) scheduled shortly. [Follow me on Twitter](https://twitter.com/iavins/) for updates. [Drop me an email](http://scr.im/avii) if you wish to arrange a workshop for your team/company.

CaskDB comes with a full test suite and a wide range of tools to help you write a database quickly. [A Github action](https://github.com/avinassh/go-caskdb/blob/master/.github/workflows/build.yml) is present with an automated tests runner. Fork the repo, push the code, and pass the tests!

Throughout the workshop, you will implement the following:
- Serialiser methods take a bunch of objects and serialise them into bytes. Also, the procedures take a bunch of bytes and deserialise them back to the things.
- Come up with a data format with a header and data to store the bytes on the disk. The header would contain metadata like timestamp, key size, and value.
- Store and retrieve data from the disk
- Read an existing CaskDB file to load all keys

### Tasks
1. Read [the paper](https://riak.com/assets/bitcask-intro.pdf). Fork this repo and checkout the `start-here` branch
2. Implement the fixed-sized header, which can encode timestamp (uint, 4 bytes), key size (uint, 4 bytes), value size (uint, 4 bytes) together
3. Implement the key, value serialisers, and pass the tests from `format_test.go`
4. Figure out how to store the data on disk and the row pointer in the memory. Implement the get/set operations. Tests for the same are in `disk_store_test.go`
5. Code from the task #2 and #3 should be enough to read an existing CaskDB file and load the keys into memory

Run `make test` to run the tests locally. Push the code to Github, and tests will run on different OS: ubuntu, mac, and windows.

Not sure how to proceed? Then check the [hints](hints.md) file which contains more details on the tasks and hints.

### Hints
- Not sure how to come up with a file format? Read the comment in the [format file](format.go)

## What next?
I often get questions about what is next after the basic implementation. Here are some challenges (with different levels of difficulties)

### Level 1:
- Crash safety: the bitcask paper stores CRC in the row, and while fetching the row back, it verifies the data
- Key deletion: CaskDB does not have a delete API. Read the paper and implement it
- Instead of using a hash table, use a data structure like the red-black tree to support range scans
- CaskDB accepts only strings as keys and values. Make it generic and take other data structures like int or bytes.

### Level 2:
- Hint file to improve the startup time. The paper has more details on it
- Implement an internal cache which stores some of the key-value pairs. You may explore and experiment with different cache eviction strategies like LRU, LFU, FIFO etc.
- Split the data into multiple files when the files hit a specific capacity

### Level 3:
- Support for multiple processes
- Garbage collector: keys which got updated and deleted remain in the file and take up space. Write a garbage collector to remove such stale data
- Add SQL query engine layer
- Store JSON in values and explore making CaskDB as a document database like Mongo
- Make CaskDB distributed by exploring algorithms like raft, paxos, or consistent hashing

## Line Count

```shell
$ tokei -f format.go disk_store.go
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 Go                      2          320          133          168           19
-------------------------------------------------------------------------------
 format.go                          111           35           67            9
 disk_store.go                      209           98          101           10
===============================================================================
 Total                   2          320          133          168           19
===============================================================================
```

## Contributing
All contributions are welcome. Please check [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## Community Contributions

| Author                                          | Feature   | PR                                                 |
|-------------------------------------------------|-----------|----------------------------------------------------|
| [PaulisMatrix](https://github.com/PaulisMatrix) | Delete Op | [#6](https://github.com/avinassh/go-caskdb/pull/6) |
| [PaulisMatrix](https://github.com/PaulisMatrix) | Checksum  | [#7](https://github.com/avinassh/go-caskdb/pull/7) |

## License
The MIT license. Please check `LICENSE` for more details.
