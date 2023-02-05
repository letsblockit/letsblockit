package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type filterLintCmd struct{}

const (
	successMessage         = "Thanks for your contribution! The automated tests passed, we will review your PR shortly!"
	failureMessageTemplate = `Thanks for your contribution!

Unfortunately, the automated tests failed. Please check the output below and [the syntax documentation](https://github.com/letsblockit/letsblockit/blob/main/data/filters/README.md) to fix these.
Don't hesitate to comment if you need help addressing these failures!

### Test output:
` + "```\n%s\n```"
)

func (c *filterLintCmd) Run() error {
	cmd := exec.Command("go", "test", "./src/filters")
	out, err := cmd.Output()
	if err == nil {
		fmt.Print(successMessage)
	} else {
		fmt.Printf(failureMessageTemplate, strings.TrimSpace(string(out)))
	}
	fmt.Println()
	return nil
}
