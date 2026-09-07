package runner

import (
	"context"
	"time"
)

// warmSource is the smallest program that pulls in the parts of the standard
// library every submission touches — reading stdin, formatting, printing.
const warmSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	_ = r
	_ = sort.Ints
	_ = strconv.Itoa
	_ = strings.TrimSpace
	fmt.Fprint(os.Stdout, "")
}
`

// WarmGo compiles a throwaway program so the shared build cache is populated
// before a learner's submission depends on it.
//
// Without this the first Go submission after a deploy pays for compiling the
// standard library — around twenty seconds, against roughly one when warm — and
// is killed by its own time limit. The learner sees a timeout on correct code,
// once, for reasons that have nothing to do with them.
//
// Safe to call repeatedly, and cheap once the cache is warm.
func WarmGo(ctx context.Context, e Executor) error {
	if e == nil || !e.Enabled() {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	_, err := e.Run(ctx, Request{
		Lang: LangGo,
		Code: warmSource,
		// Generous on purpose: this is the compile everybody else is spared.
		TimeoutMs: 150000,
	})
	return err
}
