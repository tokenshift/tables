package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"gitlab.com/tokenshift/tables/dice"
)

var Args struct {
	RollSpecs []string `kong:"arg,required,name='rollspec',help='One or more roll specs, like 1d6+2.'"`
	Verbose   bool     `kong:"short='v',help='Whether to output per-die roll information, or just the end result.'"`
}

func main() {
	ctx := kong.Parse(&Args, kong.UsageOnError())

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, arg := range Args.RollSpecs {
		parsed, err := dice.Parse(arg)
		ctx.FatalIfErrorf(err)
		results := parsed.Roll(r)

		if len(os.Args) > 2 {
			if Args.Verbose {
				fmt.Println(strings.TrimSpace(arg), "=", results, "=", results.Sum())
			} else {
				fmt.Println(results.Sum())
			}
		} else {
			if Args.Verbose {
				fmt.Println(results, "=", results.Sum())
			} else {
				fmt.Println(results.Sum())
			}
		}
	}
}
