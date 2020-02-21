[![GoDoc](https://godoc.org/github.com/crspybits/action-times?status.svg)](https://godoc.org/github.com/crspybits/action-times)
[![Go Report Card](https://goreportcard.com/badge/crspybits/action-times)](https://goreportcard.com/report/crspybits/action-times) [test coverage](https://gocover.io/github.com/crspybits/action-times)

# action-times
Coding test for JumpCloud

## go version
Developed and tested with `go version go1.12.5 darwin/amd64`


## Example usage in a `main.go`:

```
package main

import "fmt"
import ats "github.com/crspybits/action-times"

func main() {
  ats.AddAction(`{"action":"jump", "time":100}`)
  ats.AddAction(`{"action":"jump", "time":200}`)
  ats.AddAction(`{"action":"run", "time":75}`)
  ats.AddAction(`{"action":"bling", "time":800}`)

  s := ats.GetStats()
  fmt.Print("result: " + s + "\n")
}
```
Put this in a file named main.go, and do:
1) go get -d github.com/crspybits/action-times
2) go run main.go

## Testing
With this package installed in a directory named `action-times` under `src` in your $GOPATH:
1) Change to that directory
2) Run `go test`
