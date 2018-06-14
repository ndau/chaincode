package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Script",
			pos:  position{line: 6, col: 1, offset: 23},
			expr: &actionExpr{
				pos: position{line: 6, col: 11, offset: 33},
				run: (*parser).callonScript1,
				expr: &seqExpr{
					pos: position{line: 6, col: 11, offset: 33},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 6, col: 11, offset: 33},
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 11, offset: 33},
								name: "EOL",
							},
						},
						&labeledExpr{
							pos:   position{line: 6, col: 16, offset: 38},
							label: "p",
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 18, offset: 40},
								name: "Preamble",
							},
						},
						&labeledExpr{
							pos:   position{line: 6, col: 27, offset: 49},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 32, offset: 54},
								name: "Code",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 6, col: 37, offset: 59},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Preamble",
			pos:  position{line: 8, col: 1, offset: 94},
			expr: &actionExpr{
				pos: position{line: 8, col: 14, offset: 107},
				run: (*parser).callonPreamble1,
				expr: &seqExpr{
					pos: position{line: 8, col: 14, offset: 107},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 8, col: 14, offset: 107},
							expr: &ruleRefExpr{
								pos:  position{line: 8, col: 14, offset: 107},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 8, col: 17, offset: 110},
							val:        "context",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 8, col: 27, offset: 120},
							expr: &ruleRefExpr{
								pos:  position{line: 8, col: 27, offset: 120},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 8, col: 30, offset: 123},
							val:        ":",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 8, col: 34, offset: 127},
							expr: &ruleRefExpr{
								pos:  position{line: 8, col: 34, offset: 127},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 8, col: 37, offset: 130},
							label: "cc",
							expr: &ruleRefExpr{
								pos:  position{line: 8, col: 40, offset: 133},
								name: "ContextConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 8, col: 56, offset: 149},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 10, col: 1, offset: 202},
			expr: &actionExpr{
				pos: position{line: 10, col: 10, offset: 211},
				run: (*parser).callonCode1,
				expr: &seqExpr{
					pos: position{line: 10, col: 10, offset: 211},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 10, col: 10, offset: 211},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 10, offset: 211},
								name: "EOL",
							},
						},
						&labeledExpr{
							pos:   position{line: 10, col: 15, offset: 216},
							label: "f",
							expr: &oneOrMoreExpr{
								pos: position{line: 10, col: 17, offset: 218},
								expr: &ruleRefExpr{
									pos:  position{line: 10, col: 17, offset: 218},
									name: "FunctionDef",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionDef",
			pos:  position{line: 12, col: 1, offset: 250},
			expr: &actionExpr{
				pos: position{line: 12, col: 16, offset: 265},
				run: (*parser).callonFunctionDef1,
				expr: &seqExpr{
					pos: position{line: 12, col: 16, offset: 265},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 12, col: 16, offset: 265},
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 16, offset: 265},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 12, col: 19, offset: 268},
							val:        "func",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 12, col: 26, offset: 275},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 12, col: 28, offset: 277},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 30, offset: 279},
								name: "FunctionName",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 12, col: 43, offset: 292},
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 43, offset: 292},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 12, col: 46, offset: 295},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 12, col: 50, offset: 299},
							label: "s",
							expr: &oneOrMoreExpr{
								pos: position{line: 12, col: 52, offset: 301},
								expr: &ruleRefExpr{
									pos:  position{line: 12, col: 52, offset: 301},
									name: "Line",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 12, col: 58, offset: 307},
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 58, offset: 307},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 12, col: 61, offset: 310},
							val:        "}",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 12, col: 65, offset: 314},
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 65, offset: 314},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 22, col: 1, offset: 591},
			expr: &choiceExpr{
				pos: position{line: 23, col: 7, offset: 605},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 23, col: 7, offset: 605},
						run: (*parser).callonLine2,
						expr: &seqExpr{
							pos: position{line: 23, col: 7, offset: 605},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 23, col: 7, offset: 605},
									expr: &ruleRefExpr{
										pos:  position{line: 23, col: 7, offset: 605},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 23, col: 10, offset: 608},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 23, col: 13, offset: 611},
										name: "Operation",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 23, col: 23, offset: 621},
									name: "EOL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 24, col: 7, offset: 650},
						run: (*parser).callonLine9,
						expr: &ruleRefExpr{
							pos:  position{line: 24, col: 7, offset: 650},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Operation",
			pos:  position{line: 27, col: 1, offset: 681},
			expr: &choiceExpr{
				pos: position{line: 28, col: 7, offset: 700},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 28, col: 7, offset: 700},
						name: "ConstDef",
					},
					&ruleRefExpr{
						pos:  position{line: 29, col: 7, offset: 715},
						name: "Opcode",
					},
				},
			},
		},
		{
			name: "ConstDef",
			pos:  position{line: 32, col: 1, offset: 729},
			expr: &actionExpr{
				pos: position{line: 33, col: 7, offset: 747},
				run: (*parser).callonConstDef1,
				expr: &seqExpr{
					pos: position{line: 33, col: 7, offset: 747},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 33, col: 7, offset: 747},
							label: "k",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 9, offset: 749},
								name: "Constant",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 33, col: 18, offset: 758},
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 18, offset: 758},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 33, col: 21, offset: 761},
							val:        "=",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 33, col: 25, offset: 765},
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 25, offset: 765},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 28, offset: 768},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 30, offset: 770},
								name: "Value",
							},
						},
					},
				},
			},
		},
		{
			name: "Opcode",
			pos:  position{line: 46, col: 1, offset: 1203},
			expr: &choiceExpr{
				pos: position{line: 47, col: 7, offset: 1218},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 47, col: 7, offset: 1218},
						run: (*parser).callonOpcode2,
						expr: &litMatcher{
							pos:        position{line: 47, col: 7, offset: 1218},
							val:        "nop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 48, col: 7, offset: 1268},
						run: (*parser).callonOpcode4,
						expr: &litMatcher{
							pos:        position{line: 48, col: 7, offset: 1268},
							val:        "drop2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 7, offset: 1322},
						run: (*parser).callonOpcode6,
						expr: &litMatcher{
							pos:        position{line: 49, col: 7, offset: 1322},
							val:        "drop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 50, col: 7, offset: 1374},
						run: (*parser).callonOpcode8,
						expr: &litMatcher{
							pos:        position{line: 50, col: 7, offset: 1374},
							val:        "dup2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 51, col: 7, offset: 1426},
						run: (*parser).callonOpcode10,
						expr: &litMatcher{
							pos:        position{line: 51, col: 7, offset: 1426},
							val:        "dup",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 52, col: 7, offset: 1476},
						run: (*parser).callonOpcode12,
						expr: &litMatcher{
							pos:        position{line: 52, col: 7, offset: 1476},
							val:        "swap",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 53, col: 7, offset: 1528},
						run: (*parser).callonOpcode14,
						expr: &litMatcher{
							pos:        position{line: 53, col: 7, offset: 1528},
							val:        "over",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 54, col: 7, offset: 1580},
						run: (*parser).callonOpcode16,
						expr: &seqExpr{
							pos: position{line: 54, col: 7, offset: 1580},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 54, col: 7, offset: 1580},
									val:        "pick",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 54, col: 14, offset: 1587},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 54, col: 16, offset: 1589},
									label: "b",
									expr: &ruleRefExpr{
										pos:  position{line: 54, col: 18, offset: 1591},
										name: "Byte",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 55, col: 7, offset: 1652},
						run: (*parser).callonOpcode22,
						expr: &seqExpr{
							pos: position{line: 55, col: 7, offset: 1652},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 55, col: 7, offset: 1652},
									val:        "roll",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 55, col: 14, offset: 1659},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 55, col: 16, offset: 1661},
									label: "b",
									expr: &ruleRefExpr{
										pos:  position{line: 55, col: 18, offset: 1663},
										name: "Byte",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 56, col: 7, offset: 1724},
						run: (*parser).callonOpcode28,
						expr: &litMatcher{
							pos:        position{line: 56, col: 7, offset: 1724},
							val:        "ret",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 57, col: 7, offset: 1774},
						run: (*parser).callonOpcode30,
						expr: &litMatcher{
							pos:        position{line: 57, col: 7, offset: 1774},
							val:        "fail",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 58, col: 7, offset: 1826},
						run: (*parser).callonOpcode32,
						expr: &litMatcher{
							pos:        position{line: 58, col: 7, offset: 1826},
							val:        "zero",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 59, col: 7, offset: 1878},
						run: (*parser).callonOpcode34,
						expr: &litMatcher{
							pos:        position{line: 59, col: 7, offset: 1878},
							val:        "false",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 60, col: 7, offset: 1932},
						run: (*parser).callonOpcode36,
						expr: &seqExpr{
							pos: position{line: 60, col: 7, offset: 1932},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 60, col: 7, offset: 1932},
									val:        "pushb",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 60, col: 15, offset: 1940},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 60, col: 17, offset: 1942},
									label: "ba",
									expr: &ruleRefExpr{
										pos:  position{line: 60, col: 20, offset: 1945},
										name: "Bytes",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 61, col: 7, offset: 1982},
						run: (*parser).callonOpcode42,
						expr: &seqExpr{
							pos: position{line: 61, col: 7, offset: 1982},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 61, col: 7, offset: 1982},
									val:        "push",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 61, col: 14, offset: 1989},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 16, offset: 1991},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 18, offset: 1993},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 62, col: 7, offset: 2043},
						run: (*parser).callonOpcode48,
						expr: &litMatcher{
							pos:        position{line: 62, col: 7, offset: 2043},
							val:        "one",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 7, offset: 2093},
						run: (*parser).callonOpcode50,
						expr: &litMatcher{
							pos:        position{line: 63, col: 7, offset: 2093},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 64, col: 7, offset: 2145},
						run: (*parser).callonOpcode52,
						expr: &litMatcher{
							pos:        position{line: 64, col: 7, offset: 2145},
							val:        "neg1",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 7, offset: 2197},
						run: (*parser).callonOpcode54,
						expr: &seqExpr{
							pos: position{line: 65, col: 7, offset: 2197},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 65, col: 7, offset: 2197},
									val:        "pusht",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 65, col: 15, offset: 2205},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 65, col: 17, offset: 2207},
									label: "ts",
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 20, offset: 2210},
										name: "Timestamp",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 66, col: 7, offset: 2268},
						run: (*parser).callonOpcode60,
						expr: &litMatcher{
							pos:        position{line: 66, col: 7, offset: 2268},
							val:        "pushl",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 67, col: 7, offset: 2322},
						run: (*parser).callonOpcode62,
						expr: &litMatcher{
							pos:        position{line: 67, col: 7, offset: 2322},
							val:        "now",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 68, col: 7, offset: 2372},
						run: (*parser).callonOpcode64,
						expr: &seqExpr{
							pos: position{line: 68, col: 7, offset: 2372},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 68, col: 7, offset: 2372},
									val:        "pusha",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 68, col: 15, offset: 2380},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 68, col: 17, offset: 2382},
									label: "a",
									expr: &ruleRefExpr{
										pos:  position{line: 68, col: 19, offset: 2384},
										name: "Address",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 69, col: 7, offset: 2433},
						run: (*parser).callonOpcode70,
						expr: &litMatcher{
							pos:        position{line: 69, col: 7, offset: 2433},
							val:        "rand",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 70, col: 7, offset: 2485},
						run: (*parser).callonOpcode72,
						expr: &litMatcher{
							pos:        position{line: 70, col: 7, offset: 2485},
							val:        "add",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 71, col: 7, offset: 2535},
						run: (*parser).callonOpcode74,
						expr: &litMatcher{
							pos:        position{line: 71, col: 7, offset: 2535},
							val:        "sub",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 72, col: 7, offset: 2585},
						run: (*parser).callonOpcode76,
						expr: &litMatcher{
							pos:        position{line: 72, col: 7, offset: 2585},
							val:        "divmod",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 73, col: 7, offset: 2641},
						run: (*parser).callonOpcode78,
						expr: &litMatcher{
							pos:        position{line: 73, col: 7, offset: 2641},
							val:        "muldiv",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 74, col: 7, offset: 2697},
						run: (*parser).callonOpcode80,
						expr: &litMatcher{
							pos:        position{line: 74, col: 7, offset: 2697},
							val:        "mul",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 75, col: 7, offset: 2747},
						run: (*parser).callonOpcode82,
						expr: &litMatcher{
							pos:        position{line: 75, col: 7, offset: 2747},
							val:        "div",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 76, col: 7, offset: 2797},
						run: (*parser).callonOpcode84,
						expr: &litMatcher{
							pos:        position{line: 76, col: 7, offset: 2797},
							val:        "mod",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 77, col: 7, offset: 2847},
						run: (*parser).callonOpcode86,
						expr: &litMatcher{
							pos:        position{line: 77, col: 7, offset: 2847},
							val:        "not",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 78, col: 7, offset: 2897},
						run: (*parser).callonOpcode88,
						expr: &litMatcher{
							pos:        position{line: 78, col: 7, offset: 2897},
							val:        "neg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 79, col: 7, offset: 2947},
						run: (*parser).callonOpcode90,
						expr: &litMatcher{
							pos:        position{line: 79, col: 7, offset: 2947},
							val:        "inc",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 80, col: 7, offset: 2997},
						run: (*parser).callonOpcode92,
						expr: &litMatcher{
							pos:        position{line: 80, col: 7, offset: 2997},
							val:        "dec",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 81, col: 7, offset: 3047},
						run: (*parser).callonOpcode94,
						expr: &litMatcher{
							pos:        position{line: 81, col: 7, offset: 3047},
							val:        "eq",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 82, col: 7, offset: 3095},
						run: (*parser).callonOpcode96,
						expr: &litMatcher{
							pos:        position{line: 82, col: 7, offset: 3095},
							val:        "gt",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 83, col: 7, offset: 3143},
						run: (*parser).callonOpcode98,
						expr: &litMatcher{
							pos:        position{line: 83, col: 7, offset: 3143},
							val:        "lt",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 7, offset: 3191},
						run: (*parser).callonOpcode100,
						expr: &litMatcher{
							pos:        position{line: 84, col: 7, offset: 3191},
							val:        "index",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 85, col: 7, offset: 3245},
						run: (*parser).callonOpcode102,
						expr: &litMatcher{
							pos:        position{line: 85, col: 7, offset: 3245},
							val:        "len",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 86, col: 7, offset: 3295},
						run: (*parser).callonOpcode104,
						expr: &litMatcher{
							pos:        position{line: 86, col: 7, offset: 3295},
							val:        "append",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 87, col: 7, offset: 3351},
						run: (*parser).callonOpcode106,
						expr: &litMatcher{
							pos:        position{line: 87, col: 7, offset: 3351},
							val:        "extend",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 88, col: 7, offset: 3407},
						run: (*parser).callonOpcode108,
						expr: &litMatcher{
							pos:        position{line: 88, col: 7, offset: 3407},
							val:        "slice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 89, col: 7, offset: 3461},
						run: (*parser).callonOpcode110,
						expr: &litMatcher{
							pos:        position{line: 89, col: 7, offset: 3461},
							val:        "fieldl",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 90, col: 7, offset: 3517},
						run: (*parser).callonOpcode112,
						expr: &litMatcher{
							pos:        position{line: 90, col: 7, offset: 3517},
							val:        "field",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 91, col: 7, offset: 3571},
						run: (*parser).callonOpcode114,
						expr: &litMatcher{
							pos:        position{line: 91, col: 7, offset: 3571},
							val:        "choice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 92, col: 7, offset: 3627},
						run: (*parser).callonOpcode116,
						expr: &litMatcher{
							pos:        position{line: 92, col: 7, offset: 3627},
							val:        "wchoice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 93, col: 7, offset: 3685},
						run: (*parser).callonOpcode118,
						expr: &litMatcher{
							pos:        position{line: 93, col: 7, offset: 3685},
							val:        "sort",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 94, col: 7, offset: 3737},
						run: (*parser).callonOpcode120,
						expr: &litMatcher{
							pos:        position{line: 94, col: 7, offset: 3737},
							val:        "lookup",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 95, col: 7, offset: 3793},
						run: (*parser).callonOpcode122,
						expr: &seqExpr{
							pos: position{line: 95, col: 7, offset: 3793},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 95, col: 7, offset: 3793},
									val:        "call",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 95, col: 14, offset: 3800},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 95, col: 16, offset: 3802},
									label: "n",
									expr: &ruleRefExpr{
										pos:  position{line: 95, col: 18, offset: 3804},
										name: "FunctionName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 95, col: 31, offset: 3817},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 95, col: 33, offset: 3819},
									label: "b",
									expr: &ruleRefExpr{
										pos:  position{line: 95, col: 35, offset: 3821},
										name: "Byte",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 96, col: 7, offset: 3892},
						run: (*parser).callonOpcode131,
						expr: &litMatcher{
							pos:        position{line: 96, col: 7, offset: 3892},
							val:        "enddef",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 97, col: 7, offset: 3948},
						run: (*parser).callonOpcode133,
						expr: &litMatcher{
							pos:        position{line: 97, col: 7, offset: 3948},
							val:        "ifz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 98, col: 7, offset: 3998},
						run: (*parser).callonOpcode135,
						expr: &litMatcher{
							pos:        position{line: 98, col: 7, offset: 3998},
							val:        "ifnz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 99, col: 7, offset: 4050},
						run: (*parser).callonOpcode137,
						expr: &litMatcher{
							pos:        position{line: 99, col: 7, offset: 4050},
							val:        "else",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 100, col: 7, offset: 4102},
						run: (*parser).callonOpcode139,
						expr: &litMatcher{
							pos:        position{line: 100, col: 7, offset: 4102},
							val:        "endif",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 101, col: 7, offset: 4156},
						run: (*parser).callonOpcode141,
						expr: &litMatcher{
							pos:        position{line: 101, col: 7, offset: 4156},
							val:        "sum",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 102, col: 7, offset: 4206},
						run: (*parser).callonOpcode143,
						expr: &litMatcher{
							pos:        position{line: 102, col: 7, offset: 4206},
							val:        "avg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 103, col: 7, offset: 4256},
						run: (*parser).callonOpcode145,
						expr: &litMatcher{
							pos:        position{line: 103, col: 7, offset: 4256},
							val:        "max",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 104, col: 7, offset: 4306},
						run: (*parser).callonOpcode147,
						expr: &litMatcher{
							pos:        position{line: 104, col: 7, offset: 4306},
							val:        "min",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Timestamp",
			pos:  position{line: 107, col: 1, offset: 4357},
			expr: &actionExpr{
				pos: position{line: 107, col: 14, offset: 4370},
				run: (*parser).callonTimestamp1,
				expr: &labeledExpr{
					pos:   position{line: 107, col: 14, offset: 4370},
					label: "ts",
					expr: &seqExpr{
						pos: position{line: 107, col: 18, offset: 4374},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 107, col: 18, offset: 4374},
								name: "Date",
							},
							&litMatcher{
								pos:        position{line: 107, col: 23, offset: 4379},
								val:        "T",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 107, col: 27, offset: 4383},
								name: "Time",
							},
							&litMatcher{
								pos:        position{line: 107, col: 32, offset: 4388},
								val:        "Z",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "Date",
			pos:  position{line: 108, col: 1, offset: 4421},
			expr: &seqExpr{
				pos: position{line: 108, col: 9, offset: 4429},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 108, col: 9, offset: 4429},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 15, offset: 4435},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 21, offset: 4441},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 27, offset: 4447},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 108, col: 33, offset: 4453},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 37, offset: 4457},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 43, offset: 4463},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 108, col: 49, offset: 4469},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 53, offset: 4473},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 108, col: 59, offset: 4479},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Time",
			pos:  position{line: 109, col: 1, offset: 4485},
			expr: &seqExpr{
				pos: position{line: 109, col: 10, offset: 4494},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 109, col: 10, offset: 4494},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 109, col: 16, offset: 4500},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 109, col: 22, offset: 4506},
						val:        ":",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 109, col: 26, offset: 4510},
						val:        "[0:9]",
						chars:      []rune{'0', ':', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 109, col: 32, offset: 4516},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 109, col: 38, offset: 4522},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 109, col: 42, offset: 4526},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 109, col: 48, offset: 4532},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "ContextConstant",
			pos:  position{line: 111, col: 1, offset: 4539},
			expr: &choiceExpr{
				pos: position{line: 112, col: 7, offset: 4564},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 112, col: 7, offset: 4564},
						run: (*parser).callonContextConstant2,
						expr: &litMatcher{
							pos:        position{line: 112, col: 7, offset: 4564},
							val:        "TEST",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 113, col: 7, offset: 4604},
						run: (*parser).callonContextConstant4,
						expr: &litMatcher{
							pos:        position{line: 113, col: 7, offset: 4604},
							val:        "NODE_PAYOUT",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 114, col: 7, offset: 4657},
						run: (*parser).callonContextConstant6,
						expr: &litMatcher{
							pos:        position{line: 114, col: 7, offset: 4657},
							val:        "EAI_TIMING",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 115, col: 7, offset: 4708},
						run: (*parser).callonContextConstant8,
						expr: &litMatcher{
							pos:        position{line: 115, col: 7, offset: 4708},
							val:        "NODE_QUALITY",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 116, col: 7, offset: 4763},
						run: (*parser).callonContextConstant10,
						expr: &litMatcher{
							pos:        position{line: 116, col: 7, offset: 4763},
							val:        "MARKET_PRICE",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 119, col: 1, offset: 4819},
			expr: &choiceExpr{
				pos: position{line: 120, col: 7, offset: 4833},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 120, col: 7, offset: 4833},
						name: "Integer",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 7, offset: 4847},
						name: "ConstantRef",
					},
				},
			},
		},
		{
			name: "ConstantRef",
			pos:  position{line: 124, col: 1, offset: 4866},
			expr: &actionExpr{
				pos: position{line: 124, col: 16, offset: 4881},
				run: (*parser).callonConstantRef1,
				expr: &labeledExpr{
					pos:   position{line: 124, col: 16, offset: 4881},
					label: "k",
					expr: &ruleRefExpr{
						pos:  position{line: 124, col: 18, offset: 4883},
						name: "Constant",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 125, col: 1, offset: 5001},
			expr: &actionExpr{
				pos: position{line: 125, col: 12, offset: 5012},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 125, col: 12, offset: 5012},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 125, col: 12, offset: 5012},
							expr: &litMatcher{
								pos:        position{line: 125, col: 12, offset: 5012},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 125, col: 17, offset: 5017},
							expr: &charClassMatcher{
								pos:        position{line: 125, col: 17, offset: 5017},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Bytes",
			pos:  position{line: 126, col: 1, offset: 5082},
			expr: &choiceExpr{
				pos: position{line: 127, col: 7, offset: 5097},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 127, col: 7, offset: 5097},
						run: (*parser).callonBytes2,
						expr: &labeledExpr{
							pos:   position{line: 127, col: 7, offset: 5097},
							label: "b",
							expr: &oneOrMoreExpr{
								pos: position{line: 127, col: 9, offset: 5099},
								expr: &ruleRefExpr{
									pos:  position{line: 127, col: 9, offset: 5099},
									name: "Byte",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 128, col: 7, offset: 5165},
						run: (*parser).callonBytes6,
						expr: &seqExpr{
							pos: position{line: 128, col: 7, offset: 5165},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 128, col: 7, offset: 5165},
									expr: &ruleRefExpr{
										pos:  position{line: 128, col: 7, offset: 5165},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 128, col: 10, offset: 5168},
									val:        "\"",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 128, col: 14, offset: 5172},
									label: "s",
									expr: &oneOrMoreExpr{
										pos: position{line: 128, col: 16, offset: 5174},
										expr: &charClassMatcher{
											pos:        position{line: 128, col: 16, offset: 5174},
											val:        "[^\"]",
											chars:      []rune{'"'},
											ignoreCase: false,
											inverted:   true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 128, col: 22, offset: 5180},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Address",
			pos:  position{line: 131, col: 1, offset: 5234},
			expr: &actionExpr{
				pos: position{line: 131, col: 12, offset: 5245},
				run: (*parser).callonAddress1,
				expr: &seqExpr{
					pos: position{line: 131, col: 12, offset: 5245},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 131, col: 12, offset: 5245},
							val:        "nd",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 131, col: 17, offset: 5250},
							expr: &charClassMatcher{
								pos:        position{line: 131, col: 17, offset: 5250},
								val:        "[2-9a-km-np-zA-KM-NP-Z]",
								ranges:     []rune{'2', '9', 'a', 'k', 'm', 'n', 'p', 'z', 'A', 'K', 'M', 'N', 'P', 'Z'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Byte",
			pos:  position{line: 133, col: 1, offset: 5316},
			expr: &choiceExpr{
				pos: position{line: 134, col: 7, offset: 5330},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 134, col: 7, offset: 5330},
						run: (*parser).callonByte2,
						expr: &seqExpr{
							pos: position{line: 134, col: 7, offset: 5330},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 134, col: 7, offset: 5330},
									expr: &ruleRefExpr{
										pos:  position{line: 134, col: 7, offset: 5330},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 134, col: 10, offset: 5333},
									val:        "0x",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 134, col: 15, offset: 5338},
									val:        "[0-9A-Fa-f]",
									ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrOneExpr{
									pos: position{line: 134, col: 27, offset: 5350},
									expr: &charClassMatcher{
										pos:        position{line: 134, col: 27, offset: 5350},
										val:        "[0-9A-Fa-f]",
										ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 135, col: 7, offset: 5430},
						run: (*parser).callonByte10,
						expr: &seqExpr{
							pos: position{line: 135, col: 7, offset: 5430},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 135, col: 7, offset: 5430},
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 7, offset: 5430},
										name: "_",
									},
								},
								&oneOrMoreExpr{
									pos: position{line: 135, col: 10, offset: 5433},
									expr: &charClassMatcher{
										pos:        position{line: 135, col: 10, offset: 5433},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Constant",
			pos:  position{line: 137, col: 1, offset: 5530},
			expr: &actionExpr{
				pos: position{line: 137, col: 13, offset: 5542},
				run: (*parser).callonConstant1,
				expr: &seqExpr{
					pos: position{line: 137, col: 13, offset: 5542},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 137, col: 13, offset: 5542},
							val:        "[A-Z]",
							ranges:     []rune{'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 137, col: 19, offset: 5548},
							expr: &charClassMatcher{
								pos:        position{line: 137, col: 19, offset: 5548},
								val:        "[A-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionName",
			pos:  position{line: 138, col: 1, offset: 5611},
			expr: &actionExpr{
				pos: position{line: 138, col: 17, offset: 5627},
				run: (*parser).callonFunctionName1,
				expr: &seqExpr{
					pos: position{line: 138, col: 17, offset: 5627},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 138, col: 17, offset: 5627},
							val:        "[A-Za-z]",
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 138, col: 26, offset: 5636},
							expr: &charClassMatcher{
								pos:        position{line: 138, col: 26, offset: 5636},
								val:        "[A-Za-z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 141, col: 1, offset: 5694},
			expr: &oneOrMoreExpr{
				pos: position{line: 141, col: 6, offset: 5699},
				expr: &charClassMatcher{
					pos:        position{line: 141, col: 6, offset: 5699},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 143, col: 1, offset: 5707},
			expr: &seqExpr{
				pos: position{line: 143, col: 8, offset: 5714},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 143, col: 8, offset: 5714},
						expr: &ruleRefExpr{
							pos:  position{line: 143, col: 8, offset: 5714},
							name: "_",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 143, col: 11, offset: 5717},
						expr: &ruleRefExpr{
							pos:  position{line: 143, col: 11, offset: 5717},
							name: "Comment",
						},
					},
					&choiceExpr{
						pos: position{line: 143, col: 21, offset: 5727},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 143, col: 21, offset: 5727},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 143, col: 30, offset: 5736},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 143, col: 39, offset: 5745},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 143, col: 46, offset: 5752},
								val:        "\n",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 145, col: 1, offset: 5760},
			expr: &seqExpr{
				pos: position{line: 145, col: 12, offset: 5771},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 145, col: 12, offset: 5771},
						val:        ";",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 145, col: 16, offset: 5775},
						expr: &charClassMatcher{
							pos:        position{line: 145, col: 16, offset: 5775},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 147, col: 1, offset: 5785},
			expr: &notExpr{
				pos: position{line: 147, col: 8, offset: 5792},
				expr: &anyMatcher{
					line: 147, col: 9, offset: 5793,
				},
			},
		},
	},
}

func (c *current) onScript1(p, code interface{}) (interface{}, error) {
	return newScript(p, code)
}

func (p *parser) callonScript1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onScript1(stack["p"], stack["code"])
}

func (c *current) onPreamble1(cc interface{}) (interface{}, error) {
	return newPreambleNode(cc.(vm.ContextByte))
}

func (p *parser) callonPreamble1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreamble1(stack["cc"])
}

func (c *current) onCode1(f interface{}) (interface{}, error) {
	return f, nil
}

func (p *parser) callonCode1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCode1(stack["f"])
}

func (c *current) onFunctionDef1(n, s interface{}) (interface{}, error) {
	fm := c.globalStore["functions"].(map[string]int)
	name := n.(string)
	ctr := c.globalStore["functionCounter"].(int)
	fm[name] = ctr
	ctr++
	c.globalStore["functionCounter"] = ctr
	return newFunctionDef(name, s)

}

func (p *parser) callonFunctionDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionDef1(stack["n"], stack["s"])
}

func (c *current) onLine2(op interface{}) (interface{}, error) {
	return op, nil
}

func (p *parser) callonLine2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine2(stack["op"])
}

func (c *current) onLine9() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonLine9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine9()
}

func (c *current) onConstDef1(k, v interface{}) (interface{}, error) {
	if key, ok := k.(string); ok {
		cm := c.globalStore["constants"].(map[string]string)
		cm[key] = v.(string)
		return v, nil
	}
	return nil, errors.New("Bad const def")

}

func (p *parser) callonConstDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstDef1(stack["k"], stack["v"])
}

func (c *current) onOpcode2() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNop)
}

func (p *parser) callonOpcode2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode2()
}

func (c *current) onOpcode4() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDrop2)
}

func (p *parser) callonOpcode4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode4()
}

func (c *current) onOpcode6() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDrop)
}

func (p *parser) callonOpcode6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode6()
}

func (c *current) onOpcode8() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDup2)
}

func (p *parser) callonOpcode8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode8()
}

func (c *current) onOpcode10() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDup)
}

func (p *parser) callonOpcode10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode10()
}

func (c *current) onOpcode12() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSwap)
}

func (p *parser) callonOpcode12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode12()
}

func (c *current) onOpcode14() (interface{}, error) {
	return newUnitaryOpcode(vm.OpOver)
}

func (p *parser) callonOpcode14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode14()
}

func (c *current) onOpcode16(b interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpPick, b.(string))
}

func (p *parser) callonOpcode16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode16(stack["b"])
}

func (c *current) onOpcode22(b interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpRoll, b.(string))
}

func (p *parser) callonOpcode22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode22(stack["b"])
}

func (c *current) onOpcode28() (interface{}, error) {
	return newUnitaryOpcode(vm.OpRet)
}

func (p *parser) callonOpcode28() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode28()
}

func (c *current) onOpcode30() (interface{}, error) {
	return newUnitaryOpcode(vm.OpFail)
}

func (p *parser) callonOpcode30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode30()
}

func (c *current) onOpcode32() (interface{}, error) {
	return newUnitaryOpcode(vm.OpZero)
}

func (p *parser) callonOpcode32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode32()
}

func (c *current) onOpcode34() (interface{}, error) {
	return newUnitaryOpcode(vm.OpFalse)
}

func (p *parser) callonOpcode34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode34()
}

func (c *current) onOpcode36(ba interface{}) (interface{}, error) {
	return newPushB(ba)
}

func (p *parser) callonOpcode36() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode36(stack["ba"])
}

func (c *current) onOpcode42(v interface{}) (interface{}, error) {
	return newPushOpcode(v.(string))
}

func (p *parser) callonOpcode42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode42(stack["v"])
}

func (c *current) onOpcode48() (interface{}, error) {
	return newUnitaryOpcode(vm.OpOne)
}

func (p *parser) callonOpcode48() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode48()
}

func (c *current) onOpcode50() (interface{}, error) {
	return newUnitaryOpcode(vm.OpTrue)
}

func (p *parser) callonOpcode50() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode50()
}

func (c *current) onOpcode52() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNeg1)
}

func (p *parser) callonOpcode52() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode52()
}

func (c *current) onOpcode54(ts interface{}) (interface{}, error) {
	return newPushTimestamp(ts.(string))
}

func (p *parser) callonOpcode54() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode54(stack["ts"])
}

func (c *current) onOpcode60() (interface{}, error) {
	return newUnitaryOpcode(vm.OpPushL)
}

func (p *parser) callonOpcode60() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode60()
}

func (c *current) onOpcode62() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNow)
}

func (p *parser) callonOpcode62() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode62()
}

func (c *current) onOpcode64(a interface{}) (interface{}, error) {
	return newPushAddr(a.(string))
}

func (p *parser) callonOpcode64() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode64(stack["a"])
}

func (c *current) onOpcode70() (interface{}, error) {
	return newUnitaryOpcode(vm.OpRand)
}

func (p *parser) callonOpcode70() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode70()
}

func (c *current) onOpcode72() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAdd)
}

func (p *parser) callonOpcode72() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode72()
}

func (c *current) onOpcode74() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSub)
}

func (p *parser) callonOpcode74() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode74()
}

func (c *current) onOpcode76() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDivMod)
}

func (p *parser) callonOpcode76() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode76()
}

func (c *current) onOpcode78() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMulDiv)
}

func (p *parser) callonOpcode78() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode78()
}

func (c *current) onOpcode80() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMul)
}

func (p *parser) callonOpcode80() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode80()
}

func (c *current) onOpcode82() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDiv)
}

func (p *parser) callonOpcode82() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode82()
}

func (c *current) onOpcode84() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMod)
}

func (p *parser) callonOpcode84() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode84()
}

func (c *current) onOpcode86() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNot)
}

func (p *parser) callonOpcode86() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode86()
}

func (c *current) onOpcode88() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNeg)
}

func (p *parser) callonOpcode88() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode88()
}

func (c *current) onOpcode90() (interface{}, error) {
	return newUnitaryOpcode(vm.OpInc)
}

func (p *parser) callonOpcode90() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode90()
}

func (c *current) onOpcode92() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDec)
}

func (p *parser) callonOpcode92() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode92()
}

func (c *current) onOpcode94() (interface{}, error) {
	return newUnitaryOpcode(vm.OpEq)
}

func (p *parser) callonOpcode94() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode94()
}

func (c *current) onOpcode96() (interface{}, error) {
	return newUnitaryOpcode(vm.OpGt)
}

func (p *parser) callonOpcode96() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode96()
}

func (c *current) onOpcode98() (interface{}, error) {
	return newUnitaryOpcode(vm.OpLt)
}

func (p *parser) callonOpcode98() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode98()
}

func (c *current) onOpcode100() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIndex)
}

func (p *parser) callonOpcode100() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode100()
}

func (c *current) onOpcode102() (interface{}, error) {
	return newUnitaryOpcode(vm.OpLen)
}

func (p *parser) callonOpcode102() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode102()
}

func (c *current) onOpcode104() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAppend)
}

func (p *parser) callonOpcode104() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode104()
}

func (c *current) onOpcode106() (interface{}, error) {
	return newUnitaryOpcode(vm.OpExtend)
}

func (p *parser) callonOpcode106() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode106()
}

func (c *current) onOpcode108() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSlice)
}

func (p *parser) callonOpcode108() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode108()
}

func (c *current) onOpcode110() (interface{}, error) {
	return newUnitaryOpcode(vm.OpFieldL)
}

func (p *parser) callonOpcode110() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode110()
}

func (c *current) onOpcode112() (interface{}, error) {
	return newUnitaryOpcode(vm.OpField)
}

func (p *parser) callonOpcode112() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode112()
}

func (c *current) onOpcode114() (interface{}, error) {
	return newUnitaryOpcode(vm.OpChoice)
}

func (p *parser) callonOpcode114() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode114()
}

func (c *current) onOpcode116() (interface{}, error) {
	return newUnitaryOpcode(vm.OpWChoice)
}

func (p *parser) callonOpcode116() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode116()
}

func (c *current) onOpcode118() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSort)
}

func (p *parser) callonOpcode118() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode118()
}

func (c *current) onOpcode120() (interface{}, error) {
	return newUnitaryOpcode(vm.OpLookup)
}

func (p *parser) callonOpcode120() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode120()
}

func (c *current) onOpcode122(n, b interface{}) (interface{}, error) {
	return newCallOpcode(vm.OpCall, n.(string), b.(string))
}

func (p *parser) callonOpcode122() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode122(stack["n"], stack["b"])
}

func (c *current) onOpcode131() (interface{}, error) {
	return newUnitaryOpcode(vm.OpEndDef)
}

func (p *parser) callonOpcode131() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode131()
}

func (c *current) onOpcode133() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIfz)
}

func (p *parser) callonOpcode133() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode133()
}

func (c *current) onOpcode135() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIfnz)
}

func (p *parser) callonOpcode135() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode135()
}

func (c *current) onOpcode137() (interface{}, error) {
	return newUnitaryOpcode(vm.OpElse)
}

func (p *parser) callonOpcode137() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode137()
}

func (c *current) onOpcode139() (interface{}, error) {
	return newUnitaryOpcode(vm.OpEndIf)
}

func (p *parser) callonOpcode139() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode139()
}

func (c *current) onOpcode141() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSum)
}

func (p *parser) callonOpcode141() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode141()
}

func (c *current) onOpcode143() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAvg)
}

func (p *parser) callonOpcode143() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode143()
}

func (c *current) onOpcode145() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMax)
}

func (p *parser) callonOpcode145() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode145()
}

func (c *current) onOpcode147() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMin)
}

func (p *parser) callonOpcode147() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode147()
}

func (c *current) onTimestamp1(ts interface{}) (interface{}, error) {
	return ts.(string), nil
}

func (p *parser) callonTimestamp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestamp1(stack["ts"])
}

func (c *current) onContextConstant2() (interface{}, error) {
	return vm.CtxTest, nil
}

func (p *parser) callonContextConstant2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant2()
}

func (c *current) onContextConstant4() (interface{}, error) {
	return vm.CtxNodePayout, nil
}

func (p *parser) callonContextConstant4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant4()
}

func (c *current) onContextConstant6() (interface{}, error) {
	return vm.CtxEaiTiming, nil
}

func (p *parser) callonContextConstant6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant6()
}

func (c *current) onContextConstant8() (interface{}, error) {
	return vm.CtxNodeQuality, nil
}

func (p *parser) callonContextConstant8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant8()
}

func (c *current) onContextConstant10() (interface{}, error) {
	return vm.CtxMarketPrice, nil
}

func (p *parser) callonContextConstant10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant10()
}

func (c *current) onConstantRef1(k interface{}) (interface{}, error) {
	cm := c.globalStore["constants"].(map[string]string)
	return cm[k.(string)], nil
}

func (p *parser) callonConstantRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstantRef1(stack["k"])
}

func (c *current) onInteger1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

func (c *current) onBytes2(b interface{}) (interface{}, error) {
	return b, nil
}

func (p *parser) callonBytes2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBytes2(stack["b"])
}

func (c *current) onBytes6(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonBytes6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBytes6(stack["s"])
}

func (c *current) onAddress1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddress1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddress1()
}

func (c *current) onByte2() (interface{}, error) {
	return strings.TrimSpace(string(c.text)), nil
}

func (p *parser) callonByte2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onByte2()
}

func (c *current) onByte10() (interface{}, error) {
	return strings.TrimSpace(string(c.text)), nil
}

func (p *parser) callonByte10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onByte10()
}

func (c *current) onConstant1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstant1()
}

func (c *current) onFunctionName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonFunctionName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionName1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
