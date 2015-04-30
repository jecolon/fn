# fn
fn is a package, command line, and web app to fix filenames as per http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html

## Usage
```go
package main

import(
  "fmt"
  
  "github.com/jecolon/fn"
)

func main() {
  fmt.Println(fn.FixForShell("- \a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" ' áéíóúñ .txt- "))
  fmt.Println(fn.FixForURL("- \a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" ' áéíóúñ .txt- "))
}
```
