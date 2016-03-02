# Mirror

A small, machine-independent utility to synchronise folders across machines and file-systems.

It is similar to rsync, but is bi-directional and can transfer between varying file system types (e.g. S3) and operating systems.

[![Build Status](https://travis-ci.org/mefellows/mirror.svg?branch=master)](https://travis-ci.org/mefellows/mirror)
[![Coverage Status](https://coveralls.io/repos/mefellows/mirror/badge.svg?branch=master)](https://coveralls.io/r/mefellows/mirror?branch=master)

## Getting Mirror

Mirror is go gettable with `go get github.com/mefellows/mirror`. You can also download the [binary releases](https://github.com/mefellows/mirror/releases).

## Status

Beta: features are working as expected, but it's fairly rough around the edges. In use in non-mission critical systems.

## Features

* Sync a local source directory with a remote file-system-like structure, including S3 - can sync in either direction, initiated from either side.
* Watch a directory and _continuously_ synchronise changes to a another file system (local, remote or S3)

### Remote FS sync (ala rsync):

On remote server, start the mirror daemon. By default, it will listen on port 8123:

```
mirror daemon --insecure
```

On the client, specify the hostname of the remote server:

```
mirror sync --src /tmp/foo --dest mirror://mydomain.com/tmp/bar
```

#### Watch for changes:

Simply add the `--watch` flag, and mirror will watch for changes, and continuously synchronise them to the target:

```
mirror sync --src /tmp/foo --dest mirror://mydomain.com/tmp/bar --watch
```

#### Exclude files

The `--exclude` flag accepts a POSIX regular expression that can be used to filter files to be synced:

```
mirror sync --src /tmp/foo --dest /tmp/bar --exclude ".git" --exclude "^ignore" --exclude "tmp$"
```

The `--exclude` flag may be specified multiple times.

### Remote FS sync with SSL enabled

The use of SSL is recommended when transferring files between remote file systems, let's
see how we go about this with mirror.

To begin with, setup the PKI for the host with `mirror pki --caHost mydomain.com`, being careful to replace `mydomain.com` with your own domain/hostname.

You may now start the daemon with `mirror daemon`.

This will create a Certificate Authority (CA) in `~/.mirror.d/ca/`,
and client certificates in `~/.mirror.d/certs`. We will need
this to communicate securely across the network, let's extract them:

```
mirror pki --outputCA > ca.crt
mirror pki --outputClientCert > client.crt
mirror pki --outputClientKey > client.pem
```

Download these files to your local system. From the client, we import them back in:

```
mirror pki --importCA  ca.crt
mirror pki --importClientCert client.crt --importClientKey client.pem
```

We can now communicate securely with the remote host over a trusted connection:

```
bin/mirror sync --src /tmp/dat1 --dest mirror://myserver/var/backups/dat1
```

### Sync/Copy To/From S3

Ensure your AWS Credentials are loaded in the [appropriate](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) environment variables or files:

```
bin/mirror sync --src /tmp/dat1 --dest s3://mybucket.s3.amazonaws.com/dat2
```
