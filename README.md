<h1 align="center">Claude-Go</h1>

<p align="center">
        <a href="https://github.com/freeeverett/anthropic-sdk-go/master/LICENSE"><img src="https://img.shields.io/github/license/freeeverett/claude-go?style=flat-square" alt="MIT"></a>
        <a href="#"><img src="https://img.shields.io/github/go-mod/go-version/freeeverett/claude-go?label=Go%20Version&style=flat-square" alt="Go Version"></a>
    </p>
<p align="center">Golang SDK for Anthropic Claude AI</p>

<br>

## Unofficial Claude SDK, Keep Update ...

**This is under testing and improvement, please do not use it in important situations**

**Only support Claude 3**

### Usage

```shell
 go get github.com/freeeverett/claude-go@v0.0.2
```
### Simple Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/freeeverett/claude-go"
)

func main() {
	// A simple example
	apiKey := ""
	cli := claude.New(apiKey)
	text := "Tell a joke"
	r, err := cli.CreateSimpleMessage(context.TODO(),text)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}
```

For other usage, please refer to [server_test.go](./server_test.go)

### todo

- [ ] Unified Error Format

