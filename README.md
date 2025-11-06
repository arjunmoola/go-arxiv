
## go-arxiv (garx)

A lightweight terminal UI (TUI) to search arXiv from your terminal, built in Go using the Charm toolkit.

### Features
- Search arXiv via the official Atom API
- Simple TUI: enter a query and press Enter to fetch results
- Displays titles and summaries in a scrollable list
- Quit anytime with Ctrl+C

### Requirements
- Go 1.25.1 or later

### Install

Using Go:

```bash
go install github.com/arjunmoola/go-arxiv/cmd/garx@latest
```

Using Makefile:

```bash
make install
```

Build locally:

```bash
make build
# binary at ./bin/garx
```

Clean build artifacts:

```bash
make clean
```

### Run

If installed to your `GOBIN`:

```bash
garx
```

From the repo (without install):

```bash
go run ./cmd/garx
```

Usage inside the TUI:
- Type your search query
- Press Enter to search
- Press Ctrl+C to exit

### Project Layout

```text
bin/
  garx
client/
  client.go
cmd/
  garx/
    main.go
internal/
  app/
    app.go
go.mod
go.sum
Makefile
```

### Library usage (optional)

You can use the client directly to query arXiv in your own Go code.

```go
package main

import (
	"context"
	"fmt"

	arx "github.com/arjunmoola/go-arxiv/client"
)

func main() {
	c := arx.New()

	// Search all fields for "graph neural networks"
	feed, err := c.Search(context.Background(), arx.AllOfTheAbove, "graph neural networks")
	if err != nil {
		panic(err)
	}

	fmt.Println("Results:", feed.TotalResults.Value)
	for _, e := range feed.Entries {
		fmt.Println(e.Title)
	}
}
```

Key types and helpers:
- `Client.Search(ctx, prefix, query, ops...)` performs a search.
- Field prefixes: `Title`, `AuthorPrefix`, `Abstract`, `AllOfTheAbove`, etc.
- Optional operators: `WithMaxResults(n)`, `WithSortby(arx.RELEVANCE|LASTUPDATED|SUBMITTED)`, `WithSortOrder(arx.ASCENDING|DESCENDING)`, and query combinators `WithAnd`, `WithOr`, `WithAndNot`.

### Module

- Module path: `github.com/arjunmoola/go-arxiv`
- Binary: `github.com/arjunmoola/go-arxiv/cmd/garx`
