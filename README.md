# action-times
Coding test for JumpCloud


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

