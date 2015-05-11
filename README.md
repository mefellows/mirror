# Mirror

A small, machine-independent utility to synchronise folders across machines and file-systems.

It is similar to rsync, but is bi-directional and can transfer between varying file system types (e.g. S3) and operating systems.

[![Build Status](https://travis-ci.org/mefellows/mirror.svg?branch=master)](https://travis-ci.org/mefellows/mirror)
[![Coverage Status](https://coveralls.io/repos/mefellows/mirror/badge.svg?branch=master)](https://coveralls.io/r/mefellows/mirror?branch=master)

## Status

Highly experimental and not Production ready.

## Features

* Sync a local source directory with a remote file-system-like structure, including S3 - can sync in either direction, initiated from either side.
* Copy from a local source directory to a remote file-system-like structure, including S3

### Sync/Copy To/From S3

Ensure your AWS Credentials are loaded in the [appropriate](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) environment variables or files:

```
bin/mirror sync --src /tmp/dat1 --dest s3://mybucket.s3.amazonaws.com/dat2
```

### Remote FS sync (ala rsync):

On remote server, start the mirror daemon. By default, it will listen on port 8123:

```
mirror daemon --insecure
```

On the client, specify the hostname of the remote server:

```
mirror remote --hostname myhost.com --src /tmp/foo --dest /tmp/bar
```

### Remote FS sync with SSL enabled

The use of SSL is recommended when transferring files between remote file systems, let's
see how we go about this with mirror.

To begin with, start the daemon with `mirror daemon`.

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
mirror pki --importCA --certFile ca.crt
mirror pki --importClientCert --certFile client.crt --importClientKey --keyFile client.pem
```

We can now communicate securely with the remote host over a trusted connection:

```
bin/mirror sync --host myserver.com --src /tmp/dat1 --dest /var/backups/dat1
```

