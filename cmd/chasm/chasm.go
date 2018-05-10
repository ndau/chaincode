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
			pos:  position{line: 10, col: 1, offset: 192},
			expr: &actionExpr{
				pos: position{line: 10, col: 10, offset: 201},
				run: (*parser).callonCode1,
				expr: &seqExpr{
					pos: position{line: 10, col: 10, offset: 201},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 10, col: 10, offset: 201},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 10, offset: 201},
								name: "EOL",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 10, col: 15, offset: 206},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 15, offset: 206},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 10, col: 18, offset: 209},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 10, col: 22, offset: 213},
							label: "s",
							expr: &oneOrMoreExpr{
								pos: position{line: 10, col: 24, offset: 215},
								expr: &ruleRefExpr{
									pos:  position{line: 10, col: 24, offset: 215},
									name: "Line",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 10, col: 30, offset: 221},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 30, offset: 221},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 10, col: 33, offset: 224},
							val:        "}",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 10, col: 37, offset: 228},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 37, offset: 228},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 12, col: 1, offset: 252},
			expr: &choiceExpr{
				pos: position{line: 13, col: 7, offset: 266},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 13, col: 7, offset: 266},
						run: (*parser).callonLine2,
						expr: &seqExpr{
							pos: position{line: 13, col: 7, offset: 266},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 13, col: 7, offset: 266},
									expr: &ruleRefExpr{
										pos:  position{line: 13, col: 7, offset: 266},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 13, col: 10, offset: 269},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 13, col: 13, offset: 272},
										name: "Operation",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 13, col: 23, offset: 282},
									name: "EOL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 14, col: 7, offset: 311},
						run: (*parser).callonLine9,
						expr: &ruleRefExpr{
							pos:  position{line: 14, col: 7, offset: 311},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Operation",
			pos:  position{line: 17, col: 1, offset: 342},
			expr: &choiceExpr{
				pos: position{line: 18, col: 7, offset: 361},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 18, col: 7, offset: 361},
						name: "ConstDef",
					},
					&ruleRefExpr{
						pos:  position{line: 19, col: 7, offset: 376},
						name: "Opcode",
					},
				},
			},
		},
		{
			name: "ConstDef",
			pos:  position{line: 22, col: 1, offset: 390},
			expr: &actionExpr{
				pos: position{line: 23, col: 7, offset: 408},
				run: (*parser).callonConstDef1,
				expr: &seqExpr{
					pos: position{line: 23, col: 7, offset: 408},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 23, col: 7, offset: 408},
							label: "k",
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 9, offset: 410},
								name: "Constant",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 23, col: 18, offset: 419},
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 18, offset: 419},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 23, col: 21, offset: 422},
							val:        "=",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 23, col: 25, offset: 426},
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 25, offset: 426},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 23, col: 28, offset: 429},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 30, offset: 431},
								name: "Value",
							},
						},
					},
				},
			},
		},
		{
			name: "Opcode",
			pos:  position{line: 35, col: 1, offset: 800},
			expr: &choiceExpr{
				pos: position{line: 36, col: 7, offset: 815},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 36, col: 7, offset: 815},
						run: (*parser).callonOpcode2,
						expr: &litMatcher{
							pos:        position{line: 36, col: 7, offset: 815},
							val:        "nop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 37, col: 7, offset: 862},
						run: (*parser).callonOpcode4,
						expr: &litMatcher{
							pos:        position{line: 37, col: 7, offset: 862},
							val:        "drop2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 38, col: 7, offset: 913},
						run: (*parser).callonOpcode6,
						expr: &litMatcher{
							pos:        position{line: 38, col: 7, offset: 913},
							val:        "drop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 39, col: 7, offset: 962},
						run: (*parser).callonOpcode8,
						expr: &litMatcher{
							pos:        position{line: 39, col: 7, offset: 962},
							val:        "dup2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 40, col: 7, offset: 1011},
						run: (*parser).callonOpcode10,
						expr: &litMatcher{
							pos:        position{line: 40, col: 7, offset: 1011},
							val:        "dup",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 41, col: 7, offset: 1058},
						run: (*parser).callonOpcode12,
						expr: &litMatcher{
							pos:        position{line: 41, col: 7, offset: 1058},
							val:        "swap",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 42, col: 7, offset: 1107},
						run: (*parser).callonOpcode14,
						expr: &litMatcher{
							pos:        position{line: 42, col: 7, offset: 1107},
							val:        "over",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 43, col: 7, offset: 1156},
						run: (*parser).callonOpcode16,
						expr: &litMatcher{
							pos:        position{line: 43, col: 7, offset: 1156},
							val:        "ret",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 44, col: 7, offset: 1203},
						run: (*parser).callonOpcode18,
						expr: &litMatcher{
							pos:        position{line: 44, col: 7, offset: 1203},
							val:        "fail",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 45, col: 7, offset: 1252},
						run: (*parser).callonOpcode20,
						expr: &litMatcher{
							pos:        position{line: 45, col: 7, offset: 1252},
							val:        "zero",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 46, col: 7, offset: 1301},
						run: (*parser).callonOpcode22,
						expr: &litMatcher{
							pos:        position{line: 46, col: 7, offset: 1301},
							val:        "false",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 47, col: 7, offset: 1352},
						run: (*parser).callonOpcode24,
						expr: &seqExpr{
							pos: position{line: 47, col: 7, offset: 1352},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 47, col: 7, offset: 1352},
									val:        "push64",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 47, col: 16, offset: 1361},
									expr: &ruleRefExpr{
										pos:  position{line: 47, col: 16, offset: 1361},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 47, col: 19, offset: 1364},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 47, col: 21, offset: 1366},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 48, col: 7, offset: 1412},
						run: (*parser).callonOpcode31,
						expr: &seqExpr{
							pos: position{line: 48, col: 7, offset: 1412},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 48, col: 7, offset: 1412},
									val:        "push",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 48, col: 14, offset: 1419},
									expr: &ruleRefExpr{
										pos:  position{line: 48, col: 14, offset: 1419},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 48, col: 17, offset: 1422},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 48, col: 19, offset: 1424},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 7, offset: 1474},
						run: (*parser).callonOpcode38,
						expr: &litMatcher{
							pos:        position{line: 49, col: 7, offset: 1474},
							val:        "one",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 50, col: 7, offset: 1521},
						run: (*parser).callonOpcode40,
						expr: &litMatcher{
							pos:        position{line: 50, col: 7, offset: 1521},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 51, col: 7, offset: 1570},
						run: (*parser).callonOpcode42,
						expr: &litMatcher{
							pos:        position{line: 51, col: 7, offset: 1570},
							val:        "neg1",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 52, col: 7, offset: 1619},
						run: (*parser).callonOpcode44,
						expr: &seqExpr{
							pos: position{line: 52, col: 7, offset: 1619},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 52, col: 7, offset: 1619},
									val:        "pusht",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 52, col: 15, offset: 1627},
									expr: &ruleRefExpr{
										pos:  position{line: 52, col: 15, offset: 1627},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 52, col: 18, offset: 1630},
									label: "ts",
									expr: &ruleRefExpr{
										pos:  position{line: 52, col: 21, offset: 1633},
										name: "Timestamp",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 53, col: 7, offset: 1691},
						run: (*parser).callonOpcode51,
						expr: &litMatcher{
							pos:        position{line: 53, col: 7, offset: 1691},
							val:        "now",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 54, col: 7, offset: 1738},
						run: (*parser).callonOpcode53,
						expr: &litMatcher{
							pos:        position{line: 54, col: 7, offset: 1738},
							val:        "rand",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 55, col: 7, offset: 1787},
						run: (*parser).callonOpcode55,
						expr: &litMatcher{
							pos:        position{line: 55, col: 7, offset: 1787},
							val:        "add",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 56, col: 7, offset: 1834},
						run: (*parser).callonOpcode57,
						expr: &litMatcher{
							pos:        position{line: 56, col: 7, offset: 1834},
							val:        "sub",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 57, col: 7, offset: 1881},
						run: (*parser).callonOpcode59,
						expr: &litMatcher{
							pos:        position{line: 57, col: 7, offset: 1881},
							val:        "mul",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 58, col: 7, offset: 1928},
						run: (*parser).callonOpcode61,
						expr: &litMatcher{
							pos:        position{line: 58, col: 7, offset: 1928},
							val:        "div",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 59, col: 7, offset: 1975},
						run: (*parser).callonOpcode63,
						expr: &litMatcher{
							pos:        position{line: 59, col: 7, offset: 1975},
							val:        "mod",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 60, col: 7, offset: 2022},
						run: (*parser).callonOpcode65,
						expr: &litMatcher{
							pos:        position{line: 60, col: 7, offset: 2022},
							val:        "not",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 61, col: 7, offset: 2069},
						run: (*parser).callonOpcode67,
						expr: &litMatcher{
							pos:        position{line: 61, col: 7, offset: 2069},
							val:        "neg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 62, col: 7, offset: 2116},
						run: (*parser).callonOpcode69,
						expr: &litMatcher{
							pos:        position{line: 62, col: 7, offset: 2116},
							val:        "inc",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 7, offset: 2163},
						run: (*parser).callonOpcode71,
						expr: &litMatcher{
							pos:        position{line: 63, col: 7, offset: 2163},
							val:        "dec",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 64, col: 7, offset: 2210},
						run: (*parser).callonOpcode73,
						expr: &litMatcher{
							pos:        position{line: 64, col: 7, offset: 2210},
							val:        "index",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 7, offset: 2261},
						run: (*parser).callonOpcode75,
						expr: &litMatcher{
							pos:        position{line: 65, col: 7, offset: 2261},
							val:        "len",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 66, col: 7, offset: 2308},
						run: (*parser).callonOpcode77,
						expr: &litMatcher{
							pos:        position{line: 66, col: 7, offset: 2308},
							val:        "append",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 67, col: 7, offset: 2361},
						run: (*parser).callonOpcode79,
						expr: &litMatcher{
							pos:        position{line: 67, col: 7, offset: 2361},
							val:        "extend",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 68, col: 7, offset: 2414},
						run: (*parser).callonOpcode81,
						expr: &litMatcher{
							pos:        position{line: 68, col: 7, offset: 2414},
							val:        "slice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 69, col: 7, offset: 2465},
						run: (*parser).callonOpcode83,
						expr: &litMatcher{
							pos:        position{line: 69, col: 7, offset: 2465},
							val:        "fieldl",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 70, col: 7, offset: 2518},
						run: (*parser).callonOpcode85,
						expr: &litMatcher{
							pos:        position{line: 70, col: 7, offset: 2518},
							val:        "field",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 71, col: 7, offset: 2569},
						run: (*parser).callonOpcode87,
						expr: &litMatcher{
							pos:        position{line: 71, col: 7, offset: 2569},
							val:        "choice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 72, col: 7, offset: 2622},
						run: (*parser).callonOpcode89,
						expr: &litMatcher{
							pos:        position{line: 72, col: 7, offset: 2622},
							val:        "wchoice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 73, col: 7, offset: 2677},
						run: (*parser).callonOpcode91,
						expr: &litMatcher{
							pos:        position{line: 73, col: 7, offset: 2677},
							val:        "sort",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 74, col: 7, offset: 2726},
						run: (*parser).callonOpcode93,
						expr: &litMatcher{
							pos:        position{line: 74, col: 7, offset: 2726},
							val:        "lookup",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 75, col: 7, offset: 2779},
						run: (*parser).callonOpcode95,
						expr: &litMatcher{
							pos:        position{line: 75, col: 7, offset: 2779},
							val:        "ifz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 76, col: 7, offset: 2826},
						run: (*parser).callonOpcode97,
						expr: &litMatcher{
							pos:        position{line: 76, col: 7, offset: 2826},
							val:        "ifnz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 77, col: 7, offset: 2875},
						run: (*parser).callonOpcode99,
						expr: &litMatcher{
							pos:        position{line: 77, col: 7, offset: 2875},
							val:        "else",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 78, col: 7, offset: 2924},
						run: (*parser).callonOpcode101,
						expr: &litMatcher{
							pos:        position{line: 78, col: 7, offset: 2924},
							val:        "end",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 79, col: 7, offset: 2971},
						run: (*parser).callonOpcode103,
						expr: &litMatcher{
							pos:        position{line: 79, col: 7, offset: 2971},
							val:        "sum",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 80, col: 7, offset: 3018},
						run: (*parser).callonOpcode105,
						expr: &litMatcher{
							pos:        position{line: 80, col: 7, offset: 3018},
							val:        "avg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 81, col: 7, offset: 3065},
						run: (*parser).callonOpcode107,
						expr: &litMatcher{
							pos:        position{line: 81, col: 7, offset: 3065},
							val:        "max",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 82, col: 7, offset: 3112},
						run: (*parser).callonOpcode109,
						expr: &litMatcher{
							pos:        position{line: 82, col: 7, offset: 3112},
							val:        "min",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Timestamp",
			pos:  position{line: 85, col: 1, offset: 3160},
			expr: &actionExpr{
				pos: position{line: 85, col: 14, offset: 3173},
				run: (*parser).callonTimestamp1,
				expr: &labeledExpr{
					pos:   position{line: 85, col: 14, offset: 3173},
					label: "ts",
					expr: &seqExpr{
						pos: position{line: 85, col: 18, offset: 3177},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 85, col: 18, offset: 3177},
								name: "Date",
							},
							&litMatcher{
								pos:        position{line: 85, col: 23, offset: 3182},
								val:        "T",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 27, offset: 3186},
								name: "Time",
							},
							&litMatcher{
								pos:        position{line: 85, col: 32, offset: 3191},
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
			pos:  position{line: 86, col: 1, offset: 3224},
			expr: &seqExpr{
				pos: position{line: 86, col: 9, offset: 3232},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 86, col: 9, offset: 3232},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 15, offset: 3238},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 21, offset: 3244},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 27, offset: 3250},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 86, col: 33, offset: 3256},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 37, offset: 3260},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 43, offset: 3266},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 86, col: 49, offset: 3272},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 53, offset: 3276},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 86, col: 59, offset: 3282},
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
			pos:  position{line: 87, col: 1, offset: 3288},
			expr: &seqExpr{
				pos: position{line: 87, col: 10, offset: 3297},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 87, col: 10, offset: 3297},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 16, offset: 3303},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 87, col: 22, offset: 3309},
						val:        ":",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 26, offset: 3313},
						val:        "[0:9]",
						chars:      []rune{'0', ':', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 32, offset: 3319},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 87, col: 38, offset: 3325},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 42, offset: 3329},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 87, col: 48, offset: 3335},
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
			pos:  position{line: 89, col: 1, offset: 3342},
			expr: &choiceExpr{
				pos: position{line: 90, col: 7, offset: 3367},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 90, col: 7, offset: 3367},
						run: (*parser).callonContextConstant2,
						expr: &litMatcher{
							pos:        position{line: 90, col: 7, offset: 3367},
							val:        "TEST",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 91, col: 7, offset: 3404},
						run: (*parser).callonContextConstant4,
						expr: &litMatcher{
							pos:        position{line: 91, col: 7, offset: 3404},
							val:        "NODE_PAYOUT",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 92, col: 7, offset: 3454},
						run: (*parser).callonContextConstant6,
						expr: &litMatcher{
							pos:        position{line: 92, col: 7, offset: 3454},
							val:        "EAI_TIMING",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 93, col: 7, offset: 3502},
						run: (*parser).callonContextConstant8,
						expr: &litMatcher{
							pos:        position{line: 93, col: 7, offset: 3502},
							val:        "NODE_QUALITY",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 94, col: 7, offset: 3554},
						run: (*parser).callonContextConstant10,
						expr: &litMatcher{
							pos:        position{line: 94, col: 7, offset: 3554},
							val:        "MARKET_PRICE",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 97, col: 1, offset: 3607},
			expr: &choiceExpr{
				pos: position{line: 98, col: 7, offset: 3621},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 98, col: 7, offset: 3621},
						name: "Integer",
					},
					&ruleRefExpr{
						pos:  position{line: 99, col: 7, offset: 3635},
						name: "ConstantRef",
					},
				},
			},
		},
		{
			name: "ConstantRef",
			pos:  position{line: 102, col: 1, offset: 3654},
			expr: &actionExpr{
				pos: position{line: 102, col: 16, offset: 3669},
				run: (*parser).callonConstantRef1,
				expr: &labeledExpr{
					pos:   position{line: 102, col: 16, offset: 3669},
					label: "k",
					expr: &ruleRefExpr{
						pos:  position{line: 102, col: 18, offset: 3671},
						name: "Constant",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 103, col: 1, offset: 3740},
			expr: &actionExpr{
				pos: position{line: 103, col: 12, offset: 3751},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 103, col: 12, offset: 3751},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 103, col: 12, offset: 3751},
							expr: &litMatcher{
								pos:        position{line: 103, col: 12, offset: 3751},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 103, col: 17, offset: 3756},
							expr: &charClassMatcher{
								pos:        position{line: 103, col: 17, offset: 3756},
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
			name: "Constant",
			pos:  position{line: 104, col: 1, offset: 3821},
			expr: &actionExpr{
				pos: position{line: 104, col: 13, offset: 3833},
				run: (*parser).callonConstant1,
				expr: &seqExpr{
					pos: position{line: 104, col: 13, offset: 3833},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 104, col: 13, offset: 3833},
							val:        "[A-Z]",
							ranges:     []rune{'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 104, col: 19, offset: 3839},
							expr: &charClassMatcher{
								pos:        position{line: 104, col: 19, offset: 3839},
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
			name: "_",
			pos:  position{line: 106, col: 1, offset: 3903},
			expr: &oneOrMoreExpr{
				pos: position{line: 106, col: 6, offset: 3908},
				expr: &charClassMatcher{
					pos:        position{line: 106, col: 6, offset: 3908},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 108, col: 1, offset: 3916},
			expr: &seqExpr{
				pos: position{line: 108, col: 8, offset: 3923},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 108, col: 8, offset: 3923},
						expr: &ruleRefExpr{
							pos:  position{line: 108, col: 8, offset: 3923},
							name: "_",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 108, col: 11, offset: 3926},
						expr: &ruleRefExpr{
							pos:  position{line: 108, col: 11, offset: 3926},
							name: "Comment",
						},
					},
					&choiceExpr{
						pos: position{line: 108, col: 21, offset: 3936},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 108, col: 21, offset: 3936},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 108, col: 30, offset: 3945},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 108, col: 39, offset: 3954},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 108, col: 46, offset: 3961},
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
			pos:  position{line: 110, col: 1, offset: 3969},
			expr: &seqExpr{
				pos: position{line: 110, col: 12, offset: 3980},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 110, col: 12, offset: 3980},
						val:        ";",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 110, col: 16, offset: 3984},
						expr: &charClassMatcher{
							pos:        position{line: 110, col: 16, offset: 3984},
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
			pos:  position{line: 112, col: 1, offset: 3994},
			expr: &notExpr{
				pos: position{line: 112, col: 8, offset: 4001},
				expr: &anyMatcher{
					line: 112, col: 9, offset: 4002,
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
	return newPreambleNode(cc.(byte))
}

func (p *parser) callonPreamble1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreamble1(stack["cc"])
}

func (c *current) onCode1(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonCode1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCode1(stack["s"])
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
		c.state[key] = v.(string)
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
	return newUnitaryOpcode(OpNop)
}

func (p *parser) callonOpcode2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode2()
}

func (c *current) onOpcode4() (interface{}, error) {
	return newUnitaryOpcode(OpDrop2)
}

func (p *parser) callonOpcode4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode4()
}

func (c *current) onOpcode6() (interface{}, error) {
	return newUnitaryOpcode(OpDrop)
}

func (p *parser) callonOpcode6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode6()
}

func (c *current) onOpcode8() (interface{}, error) {
	return newUnitaryOpcode(OpDup2)
}

func (p *parser) callonOpcode8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode8()
}

func (c *current) onOpcode10() (interface{}, error) {
	return newUnitaryOpcode(OpDup)
}

func (p *parser) callonOpcode10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode10()
}

func (c *current) onOpcode12() (interface{}, error) {
	return newUnitaryOpcode(OpSwap)
}

func (p *parser) callonOpcode12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode12()
}

func (c *current) onOpcode14() (interface{}, error) {
	return newUnitaryOpcode(OpOver)
}

func (p *parser) callonOpcode14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode14()
}

func (c *current) onOpcode16() (interface{}, error) {
	return newUnitaryOpcode(OpRet)
}

func (p *parser) callonOpcode16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode16()
}

func (c *current) onOpcode18() (interface{}, error) {
	return newUnitaryOpcode(OpFail)
}

func (p *parser) callonOpcode18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode18()
}

func (c *current) onOpcode20() (interface{}, error) {
	return newUnitaryOpcode(OpZero)
}

func (p *parser) callonOpcode20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode20()
}

func (c *current) onOpcode22() (interface{}, error) {
	return newUnitaryOpcode(OpFalse)
}

func (p *parser) callonOpcode22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode22()
}

func (c *current) onOpcode24(v interface{}) (interface{}, error) {
	return newPush64(v.(string))
}

func (p *parser) callonOpcode24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode24(stack["v"])
}

func (c *current) onOpcode31(v interface{}) (interface{}, error) {
	return newPushOpcode(v.(string))
}

func (p *parser) callonOpcode31() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode31(stack["v"])
}

func (c *current) onOpcode38() (interface{}, error) {
	return newUnitaryOpcode(OpOne)
}

func (p *parser) callonOpcode38() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode38()
}

func (c *current) onOpcode40() (interface{}, error) {
	return newUnitaryOpcode(OpTrue)
}

func (p *parser) callonOpcode40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode40()
}

func (c *current) onOpcode42() (interface{}, error) {
	return newUnitaryOpcode(OpNeg1)
}

func (p *parser) callonOpcode42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode42()
}

func (c *current) onOpcode44(ts interface{}) (interface{}, error) {
	return newPushTimestamp(ts.(string))
}

func (p *parser) callonOpcode44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode44(stack["ts"])
}

func (c *current) onOpcode51() (interface{}, error) {
	return newUnitaryOpcode(OpNow)
}

func (p *parser) callonOpcode51() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode51()
}

func (c *current) onOpcode53() (interface{}, error) {
	return newUnitaryOpcode(OpRand)
}

func (p *parser) callonOpcode53() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode53()
}

func (c *current) onOpcode55() (interface{}, error) {
	return newUnitaryOpcode(OpAdd)
}

func (p *parser) callonOpcode55() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode55()
}

func (c *current) onOpcode57() (interface{}, error) {
	return newUnitaryOpcode(OpSub)
}

func (p *parser) callonOpcode57() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode57()
}

func (c *current) onOpcode59() (interface{}, error) {
	return newUnitaryOpcode(OpMul)
}

func (p *parser) callonOpcode59() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode59()
}

func (c *current) onOpcode61() (interface{}, error) {
	return newUnitaryOpcode(OpDiv)
}

func (p *parser) callonOpcode61() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode61()
}

func (c *current) onOpcode63() (interface{}, error) {
	return newUnitaryOpcode(OpMod)
}

func (p *parser) callonOpcode63() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode63()
}

func (c *current) onOpcode65() (interface{}, error) {
	return newUnitaryOpcode(OpNot)
}

func (p *parser) callonOpcode65() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode65()
}

func (c *current) onOpcode67() (interface{}, error) {
	return newUnitaryOpcode(OpNeg)
}

func (p *parser) callonOpcode67() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode67()
}

func (c *current) onOpcode69() (interface{}, error) {
	return newUnitaryOpcode(OpInc)
}

func (p *parser) callonOpcode69() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode69()
}

func (c *current) onOpcode71() (interface{}, error) {
	return newUnitaryOpcode(OpDec)
}

func (p *parser) callonOpcode71() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode71()
}

func (c *current) onOpcode73() (interface{}, error) {
	return newUnitaryOpcode(OpIndex)
}

func (p *parser) callonOpcode73() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode73()
}

func (c *current) onOpcode75() (interface{}, error) {
	return newUnitaryOpcode(OpLen)
}

func (p *parser) callonOpcode75() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode75()
}

func (c *current) onOpcode77() (interface{}, error) {
	return newUnitaryOpcode(OpAppend)
}

func (p *parser) callonOpcode77() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode77()
}

func (c *current) onOpcode79() (interface{}, error) {
	return newUnitaryOpcode(OpExtend)
}

func (p *parser) callonOpcode79() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode79()
}

func (c *current) onOpcode81() (interface{}, error) {
	return newUnitaryOpcode(OpSlice)
}

func (p *parser) callonOpcode81() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode81()
}

func (c *current) onOpcode83() (interface{}, error) {
	return newUnitaryOpcode(OpFieldL)
}

func (p *parser) callonOpcode83() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode83()
}

func (c *current) onOpcode85() (interface{}, error) {
	return newUnitaryOpcode(OpField)
}

func (p *parser) callonOpcode85() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode85()
}

func (c *current) onOpcode87() (interface{}, error) {
	return newUnitaryOpcode(OpChoice)
}

func (p *parser) callonOpcode87() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode87()
}

func (c *current) onOpcode89() (interface{}, error) {
	return newUnitaryOpcode(OpWChoice)
}

func (p *parser) callonOpcode89() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode89()
}

func (c *current) onOpcode91() (interface{}, error) {
	return newUnitaryOpcode(OpSort)
}

func (p *parser) callonOpcode91() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode91()
}

func (c *current) onOpcode93() (interface{}, error) {
	return newUnitaryOpcode(OpLookup)
}

func (p *parser) callonOpcode93() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode93()
}

func (c *current) onOpcode95() (interface{}, error) {
	return newUnitaryOpcode(OpIfz)
}

func (p *parser) callonOpcode95() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode95()
}

func (c *current) onOpcode97() (interface{}, error) {
	return newUnitaryOpcode(OpIfnz)
}

func (p *parser) callonOpcode97() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode97()
}

func (c *current) onOpcode99() (interface{}, error) {
	return newUnitaryOpcode(OpElse)
}

func (p *parser) callonOpcode99() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode99()
}

func (c *current) onOpcode101() (interface{}, error) {
	return newUnitaryOpcode(OpEnd)
}

func (p *parser) callonOpcode101() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode101()
}

func (c *current) onOpcode103() (interface{}, error) {
	return newUnitaryOpcode(OpSum)
}

func (p *parser) callonOpcode103() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode103()
}

func (c *current) onOpcode105() (interface{}, error) {
	return newUnitaryOpcode(OpAvg)
}

func (p *parser) callonOpcode105() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode105()
}

func (c *current) onOpcode107() (interface{}, error) {
	return newUnitaryOpcode(OpMax)
}

func (p *parser) callonOpcode107() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode107()
}

func (c *current) onOpcode109() (interface{}, error) {
	return newUnitaryOpcode(OpMin)
}

func (p *parser) callonOpcode109() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode109()
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
	return CtxTest, nil
}

func (p *parser) callonContextConstant2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant2()
}

func (c *current) onContextConstant4() (interface{}, error) {
	return CtxNodePayout, nil
}

func (p *parser) callonContextConstant4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant4()
}

func (c *current) onContextConstant6() (interface{}, error) {
	return CtxEaiTiming, nil
}

func (p *parser) callonContextConstant6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant6()
}

func (c *current) onContextConstant8() (interface{}, error) {
	return CtxNodeQuality, nil
}

func (p *parser) callonContextConstant8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant8()
}

func (c *current) onContextConstant10() (interface{}, error) {
	return CtxMarketPrice, nil
}

func (p *parser) callonContextConstant10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContextConstant10()
}

func (c *current) onConstantRef1(k interface{}) (interface{}, error) {
	return c.state[k.(string)], nil
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

func (c *current) onConstant1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstant1()
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
