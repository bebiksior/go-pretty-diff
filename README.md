# go-pretty-diff
Library for generating good looking HTML diffs in Go from unified diffs.

![CleanShot 2025-01-15 at 19 51 56](https://github.com/user-attachments/assets/27e1b856-5d14-4fbc-8ff4-dffc2971921e)


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
