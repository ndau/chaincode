package main

import (
	"fmt"
	"strings"
)

// This file provides helpers to make errors in chasm code more palatable.

// ErrorPosition defines the raw error position data.
type ErrorPosition struct {
	name   string
	line   int
	col    int
	offset int
}

// ErrorPositioner is an interface that can be used to tell if an error provides
// position data in the source file.
type ErrorPositioner interface {
	ErrorPos() ErrorPosition
}

func (p *parserError) ErrorPos() ErrorPosition {
	return ErrorPosition{
		name:   p.prefix,
		line:   p.pos.line,
		col:    p.pos.col,
		offset: p.pos.offset,
	}
}

func describeError(err error, source string) string {
	if e, ok := err.(ErrorPositioner); ok {
		lines := strings.Split(source, "\n")
		ep := e.ErrorPos()
		ntabs := strings.Count(lines[ep.line-1], "\t")
		prefix := "     " + strings.Repeat("\t", ntabs) + strings.Repeat(" ", ep.col-ntabs)
		return fmt.Sprintf("%s\n%4d: %s\n%s\n", err.Error(), ep.line, lines[ep.line-1], prefix+"^")
	}
	fmt.Printf("NOT ErrorPositioner: %#v\n", err)
	return err.Error()
}

func describeErrors(err error, source string) string {
	if el, ok := err.(errList); ok {
		s := ""
		for _, e := range el {
			s += describeError(e, source)
		}
		return s
	}
	fmt.Printf("NOT errList: %#v\n", err)
	return describeError(err, source)
}
