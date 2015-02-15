# Mirror

A small, machine-independent utility to synchronise folders across machines and file-systems.

It is similar to rsync, but it can transfer between varying file system types (e.g. S3) and across operating systems.

[![Build Status](https://travis-ci.org/mefellows/mirror.svg?branch=prototype)](https://travis-ci.org/mefellows/mirror)
[![Coverage Status](https://coveralls.io/repos/mefellows/mirror/badge.svg?branch=prototype)](https://coveralls.io/r/mefellows/mirror?branch=prototype)

## Status

Highly experimental and not Production ready.

## Features

* Copy from a local source directory to a remote file-system-like structure, including S3

## Usage

```
mirror sync -exclude="\.bak$" -exclude="git" -src=. -dest=s3://<mybucket>/myfolder/
```
