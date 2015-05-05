# coffer

This command line tool is will retrieve an encrypted coffer of data from an s3 bucket and create files on the local host.

A typical use case for coffer is you have a docker container which needs to retrieve on startup some file based secrets and apply them prior to starting a service. This is quite common requirement with continuous integration agents running in docker containers.

# coffer format

The coffer file is a YAML file containing a list of files, there mode and some content. At the moment all files are created with the normal default of the running user.

```yaml
files:
  "/home/user/myfile2" :
    mode: 0755
    content: |
      # this is my file
      # with content
```

# environment

The command reads the following environment variables.

* `AWS_ACCESS_KEY`
* `AWS_SECRET_KEY`
* `AWS_REGION`
* `COFFER_SECRET`
* `S3_BUCKET`

# usage

Sub commands for this tool are:

* encypt, this encrypts the coffer file.
* decrypt, this decrypts the coffer file, required at the moment if you want to edit it.
* upload, uploads the coffer to s3, ensuring that only encrypted data gets uploaded.
* download, pull down a coffer and validates it, file is only saved if it is decrypts and is valid.
* sync, sync a coffer with the file system, this creates/modifies/chmods files based on the information in the yaml.

# encryption

This now uses `golang.org/x/crypto/nacl/secretbox` which is a great little library designed to help people do message encryption correctly.

# License

This code is released under the MIT license see the LICENSE.md file for more details. 
