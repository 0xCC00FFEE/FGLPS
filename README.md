# Fast TCP Port Scanner

A highly concurrent TCP port scanner.

## Run Tests with Code Coverage (Linux)
`go test -cover`

## Compile (Linux)
`go build -v -o fglps`

## Run (Linux)

### Scan a single host
`./fglps -host localhost`

### See the built-in help
`./fglps -help`

### Usage Information

```
Usage of ./fglps:
  -firstPort int
        First port of port range to scan (1-65535) (default 1)
  -host string
        Host to scan
  -lastPort int
        Last port of port range to scan (1-65535) (default 65535)
  -portTimeout int
        Port timeout in seconds. (default 5)
  -threads int
        Thread count. (maximum simultaneous port scans) (default 65535)
```