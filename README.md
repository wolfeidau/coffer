# coffer

This command line tool is will retrieve an encrypted bundle of data from an s3 bucket and create files on the local host.

A typical use case for coffer is you have a docker container which needs to retrieve on startup some file based secrets and apply them prior to starting a service. This is quite common requirement with continuous integration agents running in docker containers.



# bundle format

The bundle file is a YAML file containing a list of files, there mode, owner and optionally some content. 

```yaml
files:
  "/home/user/myfile2" :
    mode: "000755"
    owner: root
    group: root
    content: |
      # this is my file
      # with content
```

# environment

The command reads the following environment variables.

* `AWS_ACCESS_KEY`
* `AWS_SECRET_KEY`
* `COFFER_SECRET`
* `S3_BUCKET`

# usage

Sub commands for this tool are:

* encypt
* decrypt
* validate
* upload
* sync --dry-run

# status

I am still building this at the moment.

# License

This code is released under the MIT license see the LICENSE.md file for more details. 
