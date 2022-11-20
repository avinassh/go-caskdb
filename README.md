![logo](https://github.com/avinassh/py-caskdb/raw/master/assets/logo.svg)
# CaskDB - Disk based Log Structured Hash Table Store

![made-with-python](https://img.shields.io/badge/Made%20with-Python-1f425f.svg)
[![build](https://github.com/avinassh/py-caskdb/actions/workflows/build.yml/badge.svg)](https://github.com/avinassh/py-caskdb/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/avinassh/py-caskdb/branch/master/graph/badge.svg?token=9SA8Q4L7AZ)](https://codecov.io/gh/avinassh/py-caskdb)
[![MIT license](https://camo.githubusercontent.com/f7358a0a5a91ec17974d36c9d426073a0ac67a958b22319be1ba5aa32542c28d/68747470733a2f2f62616467656e2e6e65742f6769746875622f6c6963656e73652f4e61657265656e2f5374726170646f776e2e6a73)](https://github.com/avinassh/py-caskdb/blob/master/LICENSE)
[![twitter@iavins](https://img.shields.io/twitter/follow/iavins?style=social)](https://twitter.com/iavins)

![architecture](https://user-images.githubusercontent.com/640792/167299554-0fc44510-d500-4347-b680-258e224646fa.png)

CaskDB is a  disk-based, embedded, persistent, key-value store based on the [Riak's bitcask paper](https://riak.com/assets/bitcask-intro.pdf), written in Python. It is more focused on the educational capabilities than using it in production. The file format is platform, machine, and programming language independent. Say, the database file created from Python on macOS should be compatible with Rust on Windows.

This project aims to help anyone, even a beginner in databases, build a persistent database in a few hours. There are no external dependencies; only the Go standard library is enough.