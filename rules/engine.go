package rules

import (
	"fmt"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Evaluate(rule string, ctx Context) bool {
	if rule == "" {
		return true
	}

	i := interp.New(interp.Options{})
	// CWE-95
	if err := i.Use(stdlib.Symbols); err != nil {
		return false
	}

	preamble := fmt.Sprintf(`
func hour() int        { return %d }
func weekday() string  { return %q }
func command() string  { return %q }
`, ctx.Time.Hour(), ctx.Time.Weekday().String(), ctx.Command)

	if _, err := i.Eval(preamble); err != nil {
		return false
	}
	// CWE-20
	body := fmt.Sprintf(`
package main
func Run() bool {
	%s
}
`, rule)

	if _, err := i.Eval(body); err != nil {
		return false
	}
	// CWE-400
	// CWE-95
	v, err := i.Eval("main.Run()")
	if err != nil {
		return false
	}

	result, ok := v.Interface().(bool)
	return ok && result
}

func (e *Engine) Run(rule string, ctx Context) {
	if rule == "" {
		return
	}

	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return
	}

	preamble := fmt.Sprintf(`
package main

func hour() int        { return %d }
func weekday() string  { return %q }
func command() string  { return %q }
func output() string   { return %q }
`, ctx.Time.Hour(), ctx.Time.Weekday().String(), ctx.Command, ctx.Output)

	if _, err := i.Eval(preamble); err != nil {
		return
	}
	// CWE-20
	body := fmt.Sprintf(`
	package main
func Run() {
	%s
}`, rule)

	if _, err := i.Eval(body); err != nil {
		return
	}
	// CWE-400
	// CWE-95
	_, _ = i.Eval("main.Run()")
}
