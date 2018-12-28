## Overview

Vault-export is a Hashicorp Vault secret exporter. It recursively scans secrets under a given path and gives output as to easily get the write commands:

```
vault write secret/path/to/write key1=value1 key2=value2
```

## Installation

```
go get -u github.com/serhatcetinkaya/vault-export
```

In order to use vault-export you should provide a configuration file with vault addr and token to your vault cluster:

```
$ cat .auth.yaml
token: ACCESS-TOKEN
vault_addr: ADDRESS-OF-YOUR-VAULT-SERVER
```

If `VAULT_EXPORTER_CONFIG_FILE` environmental variable is specified it reads the config from given file otherwise it tries to read `$(pwd)/.auth.yaml`

## Usage

Read specific path:

```
$ vault-export -k secret/path/to/key
vault write secret/path/to/key key=value
```

Read all secrets:

```
$ vault-export -k secret/
vault write secret/path/to/key key=value
vault write secret/path/to/anotherkey anotherkey=anothervalue
.
.
.
```

