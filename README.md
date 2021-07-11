# Fast TCP Port Scanner

A highly concurrent TCP port scanner.

## Compile (Linux)
`go build -v -o fglps`

## Run (Linux)

### Scan a single host
`./fglps -host localhost`

### See the built-in help
`./fglps -help`

## TODO
- Add support for IP range scanning, instead of single host scanning.
- Add support for customizing the packets TCP headers.
- Add support for Packet Header Fragmentation scanning method.
- Check host validity before starting the actual scan.
