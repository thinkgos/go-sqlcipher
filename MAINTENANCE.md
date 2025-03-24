# MAINTENANCE

To maintain this code properly, the following three repositories have to be
tracked for changes (maintenance details below):

- https://github.com/mattn/go-sqlite
- https://github.com/sqlcipher/sqlcipher
- https://github.com/libtom/libtomcrypt

## Update code from https://github.com/mattn/go-sqlite3

Current release: `v1.14.24`

Use `./track_go-sqlite3.sh`

## Update code from https://github.com/sqlcipher/sqlcipher

Current release: `v4.6.1`

Execute:

```bash
./configure
make
```

Track files:

- `sqlite3.h`
- `sqlite3.c`

## Update code from https://github.com/libtom/libtomcrypt

Current HEAD: `a6b9aff7aab857fe1b491710a5c5b9e2be49cb08`
(from `develop` branch, `2025-03-02`)

Use `./track_libtomcrypt.sh`
