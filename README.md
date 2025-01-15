# go-pretty-diff
I couldn't find a good library for generating pretty HTML diffs in Go, so I decided to write my own.

## Installation

```bash
go get github.com/bebiksior/go-pretty-diff
```

## Usage

```go
package main

import (
    "fmt"
    prettydiff "github.com/bebiksior/go-pretty-diff"
)

func main() {
    diff, err := prettydiff.ParseUnifiedDiff(`--- a.txt
+++ b.txt
@@ -1,3 +1,3 @@
-old line
+new line
 unchanged`)
    if err != nil {
        panic(err)
    }

    html, err := prettydiff.GenerateHTML(diff)
    if err != nil {
        panic(err)
    }

    fmt.Println(html)
}
```
