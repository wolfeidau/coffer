# coffer

This command line tool is designed to simplify storage and retrieval of secrets in [Amazon Web Services](https://aws.amazon.com).

It uses the following services:

* [Simple Storage Service](https://aws.amazon.com/s3/) (S3) to store secrets encrypted in files
* [Key Management Service](https://aws.amazon.com/kms/) (KMS) to manage encryption keys which encrypt/decrypt your secrets

A typical use case for coffer is you have a docker container which needs to retrieve on startup some file based secrets and apply them prior to starting a service. This is quite common requirement with continuous integration agents running in docker containers.

# coffer bundle format

coffer uses a a YAML file file to package a bunch of files together. The format of this file is illustrated below.

coffer has the ability to synchronise the files described in this bundle with the filesystem, creating/updating and changing the mode of the files.

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

* `AWS_REGION` the AWS region
* `AWS_PROFILE` the AWS profile to use
* `COFFER_ALIAS` the alias name of the file in KMS
* `S3_BUCKET` the S3 bucket which the file will be uploaded

# usage

Sub commands for this tool are:

* encrypt, this encrypts the coffer file.
* decrypt, this decrypts the coffer file, required at the moment if you want to edit it.
* upload, uploads the coffer to s3, ensuring that only encrypted data gets uploaded.
* download, pull down a coffer and validates it, file is only saved if it is decrypts and is valid.
* sync, sync a coffer with the file system, this creates/modifies/chmods files based on the information in the yaml.

# example

Before you start.

* Create a bucket in S3, I suggest something like `XXXX-coffers` in the same region as your KMS key.
* Create a KMS key see [Creating Keys](http://docs.aws.amazon.com/kms/latest/developerguide/create-keys.html) with the alias `coffer`, note this needs to be in the same region as your S3 bucket.
* Make an IAM role in AWS for your servers permitting access to the S3 bucket and KMS key (see the IAM policy below).

Create a coffer file with some SSH keys in it.

```
cat > buildkite.coffer <<EOF
files:
  "/var/lib/buildkite-agent/.ssh/id_rsa":
    mode: 0600
    content: |
        -----BEGIN RSA PRIVATE KEY-----
        XXXX
        -----END RSA PRIVATE KEY-----
EOF
```

Encrypt and Upload the coffer file to S3.

```
AWS_PROFILE=XXXX AWS_REGION=us-west-2 coffer --coffer-file buildkite.coffer upload --bucket="XXXX-coffers"
```

# IAM Role

If you want to give systems permission to access your coffer key in KMS use the following role. Note you will need to grab the ARN of your key from KMS.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:Get*",
        "s3:List*"
      ],
      "Resource": [
        "arn:aws:s3:::XXXX-coffers/*"
      ]
    },
    {
      "Sid": "Allow use of the key",
      "Effect": "Allow",
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "arn:aws:kms:us-west-2:XXXX:key/XXXX-XXXX-XXXX-XXXX-XXXX"
    }
  ]
}
```

# KMS

You can list your key aliases using the AWS CLI.

```
aws --profile XXXX kms list-aliases
```

# encryption

This now uses `golang.org/x/crypto/nacl/secretbox` which is a great little library designed to help people do message encryption correctly.

# change log

# 2.0

* Changed file format, now uses YAML as a container for meta data and encrypted payload
* Added a version and name field
* Added support for KMS to remove the need for a secret

# License

This code is released under the MIT license see the LICENSE.md file for more details.
