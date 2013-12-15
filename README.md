# argf

This is a small Go library that reads lines of text from either a file or files specified in `os.Args`, or, if
none are given, from stdin. It's similar to Ruby's `ARGF` or Perl's diamond operator.

This is mainly intended for tiny script-like CLI tools. Note that multiple goroutines should not call its
functions concurrently.

Here's how to write a simple version of `cat` using argf:

``` go
package main

import (
  "fmt"
  "github.com/cespare/argf"
  "os"
)

func main() {
  for argf.Scan() {
    fmt.Println(argf.String())
  }
  if err := argf.Error(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
```

See full package documentation [at godoc.org](http://godoc.org/github.com/cespare/argf).
