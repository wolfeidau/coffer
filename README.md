# coffer

This command line tool is will retrieve an encrypted bundle of data from an s3 bucket and create files on the local host.

A typical use case for coffer is you have a docker container which needs to retrieve on startup some file based secrets and apply them prior to starting a service. This is quite common requirement with continuous integration agents running in docker containers.

# bundle format

The bundle file is a YAML file containing a list of files, there mode and some content. At the moment all files are created with the normal default of the running user.

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

* encypt
* decrypt
* upload
* download
* sync

# encryption

This utility currently uses GCM and AES256 see https://github.com/wolfeidau/coffer/blob/master/crypto.go.

# disclaimer

This code has not been reviewed for security so use at your own risk. 

# contributions

Suggestions are welcome especially around the method I have used to encrypt and decrypt the yaml bundle file.

# License

This code is released under the MIT license see the LICENSE.md file for more details. 
