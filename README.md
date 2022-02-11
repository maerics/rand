# rand

Print random bytes from a secure source to stdout.

## Usage
```
Print random bytes from a secure source to stdout.

Usage:
  rand LENGTH_BYTES [flags]

Flags:
  -a, --base64         print random bytes encoded as Base64
  -b, --binary         print random bytes directly without formatting or trailing newline
  -h, --help           help for rand
  -n, --omit-newline   do not print the trailing newline character
  -s, --seed string    seed value as a decimal 64bit integer using an insecure random source
```

## Examples
```sh
$ rand 16
6548caf97a6c07132cf2eeeb2264270b
$ rand 32 --base64
kPRmfq/2UPzUk7YAXADZzSpN8K6Z8JwweXkewYxj5lw=
$ rand 8 --binary | xxd
00000000: 0da2 41b4 7c3f a06d                      ..A.|?.m
$ rand 8 --seed=123
f1405ced8b9968ba
$ rand 16 --seed=123
f1405ced8b9968baf9109259515bf702
$ rand 32 --seed=123
f1405ced8b9968baf9109259515bf7025a291b00ff7bfd6a4cdb51d40f4b367c
```
