# goprompt

**goprompt** is a minimalistic library that wraps `github.com/chzyer/readline` 
in order to provide simple prompts and selection mechanisms to CLI applications.

## Usage

This library provides two main structures: `Prompt` and `Select`:

### Prompt

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/heyvito/goprompt"
)

func main() {
	prompt := goprompt.Prompt{
		Label:        "Name",
		DefaultValue: "Paul Appleseed",
		Description:  "Something explaining why we are prompting for the user's name",
		Validation: func(s string) bool {
			// Return true to allow value to be accepted.
			// Let's suppose we want to check whether the name contains at least
			// one space.
			return strings.Contains(s, " ")
		},
	}

	result, err := prompt.Run()
	if err != nil {
		// Error showing prompt
		panic(err)
	}
	if result.Cancelled {
		fmt.Println("Aborted.")
		os.Exit(1)
	}
}
```

### Select

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/heyvito/goprompt"
)

func main() {
	sel := goprompt.Select{
		Label:        "Select an item",
		Description:  "Something explaining why we are prompting the user to select an item",
	}

	result, err := sel.Run()
	if err != nil {
		// Error showing prompt
		panic(err)
	}
	if result.Cancelled {
		fmt.Println("Aborted.")
		os.Exit(1)
	}
	fmt.Printf("You selected option %s (index %d)\n", result.SelectedValue, result.SelectedIndex)
}
```

## License

```
MIT License

Copyright (c) 2021 Victor Gama

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
