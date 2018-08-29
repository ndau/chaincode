package main

type SymType int

const (
	SymNumber    SymType = iota
	SymBytes     SymType = iota
	SimTImestamp SymType = iota
)

type Symbol struct {
	N string
	T SymType
}

type SymbolTable map[string]Symbol
