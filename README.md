# AWS Crypt

A simple utility CLI that leverages AWS KMS Keys to encrypt and decrypt file contents.

## Usage

### Encrypt

Writes to an encrypted file

```bash
aws-crypt -action encrypt -alias alias/my-key -file /tmp/secret
Wrote out encrypted file to: /tmp/secret.enc
```

### Decrypt

Outputs contents to STDOUT

```bash
aws-crypt -action decrypt -alias alias/my-key -file /tmp/secret.enc
My secret contents
```
