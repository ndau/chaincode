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
			name: "Grammar",
			pos:  position{line: 5, col: 1, offset: 18},
			expr: &actionExpr{
				pos: position{line: 5, col: 11, offset: 30},
				run: (*parser).callonGrammar1,
				expr: &seqExpr{
					pos: position{line: 5, col: 11, offset: 30},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 5, col: 11, offset: 30},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 5, col: 14, offset: 33},
							label: "initializer",
							expr: &zeroOrOneExpr{
								pos: position{line: 5, col: 26, offset: 45},
								expr: &seqExpr{
									pos: position{line: 5, col: 28, offset: 47},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 28, offset: 47},
											name: "Initializer",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 40, offset: 59},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 5, col: 46, offset: 65},
							label: "rules",
							expr: &oneOrMoreExpr{
								pos: position{line: 5, col: 52, offset: 71},
								expr: &seqExpr{
									pos: position{line: 5, col: 54, offset: 73},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 5, col: 54, offset: 73},
											name: "Rule",
										},
										&ruleRefExpr{
											pos:  position{line: 5, col: 59, offset: 78},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 5, col: 65, offset: 84},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Initializer",
			pos:  position{line: 7, col: 1, offset: 112},
			expr: &seqExpr{
				pos: position{line: 7, col: 15, offset: 128},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 7, col: 15, offset: 128},
						label: "code",
						expr: &ruleRefExpr{
							pos:  position{line: 7, col: 20, offset: 133},
							name: "CodeBlock",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 7, col: 30, offset: 143},
						name: "EOS",
					},
				},
			},
		},
		{
			name: "Rule",
			pos:  position{line: 9, col: 1, offset: 148},
			expr: &seqExpr{
				pos: position{line: 9, col: 8, offset: 157},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 9, col: 8, offset: 157},
						label: "name",
						expr: &ruleRefExpr{
							pos:  position{line: 9, col: 13, offset: 162},
							name: "IdentifierName",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 9, col: 28, offset: 177},
						name: "__",
					},
					&labeledExpr{
						pos:   position{line: 9, col: 31, offset: 180},
						label: "display",
						expr: &zeroOrOneExpr{
							pos: position{line: 9, col: 39, offset: 188},
							expr: &seqExpr{
								pos: position{line: 9, col: 41, offset: 190},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 9, col: 41, offset: 190},
										name: "StringLiteral",
									},
									&ruleRefExpr{
										pos:  position{line: 9, col: 55, offset: 204},
										name: "__",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 9, col: 61, offset: 210},
						name: "RuleDefOp",
					},
					&ruleRefExpr{
						pos:  position{line: 9, col: 71, offset: 220},
						name: "__",
					},
					&labeledExpr{
						pos:   position{line: 9, col: 74, offset: 223},
						label: "expr",
						expr: &ruleRefExpr{
							pos:  position{line: 9, col: 79, offset: 228},
							name: "Expression",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 9, col: 90, offset: 239},
						name: "EOS",
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 11, col: 1, offset: 244},
			expr: &ruleRefExpr{
				pos:  position{line: 11, col: 14, offset: 259},
				name: "RecoveryExpr",
			},
		},
		{
			name: "RecoveryExpr",
			pos:  position{line: 13, col: 1, offset: 273},
			expr: &seqExpr{
				pos: position{line: 13, col: 16, offset: 290},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 13, col: 16, offset: 290},
						label: "expr",
						expr: &ruleRefExpr{
							pos:  position{line: 13, col: 21, offset: 295},
							name: "ChoiceExpr",
						},
					},
					&labeledExpr{
						pos:   position{line: 13, col: 32, offset: 306},
						label: "recoverExprs",
						expr: &zeroOrMoreExpr{
							pos: position{line: 13, col: 45, offset: 319},
							expr: &seqExpr{
								pos: position{line: 13, col: 47, offset: 321},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 13, col: 47, offset: 321},
										name: "__",
									},
									&litMatcher{
										pos:        position{line: 13, col: 50, offset: 324},
										val:        "//{",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 13, col: 56, offset: 330},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 13, col: 59, offset: 333},
										name: "Labels",
									},
									&ruleRefExpr{
										pos:  position{line: 13, col: 66, offset: 340},
										name: "__",
									},
									&litMatcher{
										pos:        position{line: 13, col: 69, offset: 343},
										val:        "}",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 13, col: 73, offset: 347},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 13, col: 76, offset: 350},
										name: "ChoiceExpr",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Labels",
			pos:  position{line: 15, col: 1, offset: 365},
			expr: &seqExpr{
				pos: position{line: 15, col: 10, offset: 376},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 15, col: 10, offset: 376},
						label: "label",
						expr: &ruleRefExpr{
							pos:  position{line: 15, col: 16, offset: 382},
							name: "IdentifierName",
						},
					},
					&labeledExpr{
						pos:   position{line: 15, col: 31, offset: 397},
						label: "labels",
						expr: &zeroOrMoreExpr{
							pos: position{line: 15, col: 38, offset: 404},
							expr: &seqExpr{
								pos: position{line: 15, col: 40, offset: 406},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 15, col: 40, offset: 406},
										name: "__",
									},
									&litMatcher{
										pos:        position{line: 15, col: 43, offset: 409},
										val:        ",",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 47, offset: 413},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 15, col: 50, offset: 416},
										name: "IdentifierName",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ChoiceExpr",
			pos:  position{line: 17, col: 1, offset: 434},
			expr: &seqExpr{
				pos: position{line: 17, col: 14, offset: 449},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 17, col: 14, offset: 449},
						label: "first",
						expr: &ruleRefExpr{
							pos:  position{line: 17, col: 20, offset: 455},
							name: "ActionExpr",
						},
					},
					&labeledExpr{
						pos:   position{line: 17, col: 31, offset: 466},
						label: "rest",
						expr: &zeroOrMoreExpr{
							pos: position{line: 17, col: 36, offset: 471},
							expr: &seqExpr{
								pos: position{line: 17, col: 38, offset: 473},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 17, col: 38, offset: 473},
										name: "__",
									},
									&litMatcher{
										pos:        position{line: 17, col: 41, offset: 476},
										val:        "/",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 17, col: 45, offset: 480},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 17, col: 48, offset: 483},
										name: "ActionExpr",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ActionExpr",
			pos:  position{line: 19, col: 1, offset: 498},
			expr: &seqExpr{
				pos: position{line: 19, col: 14, offset: 513},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 19, col: 14, offset: 513},
						label: "expr",
						expr: &ruleRefExpr{
							pos:  position{line: 19, col: 19, offset: 518},
							name: "SeqExpr",
						},
					},
					&labeledExpr{
						pos:   position{line: 19, col: 27, offset: 526},
						label: "code",
						expr: &zeroOrOneExpr{
							pos: position{line: 19, col: 32, offset: 531},
							expr: &seqExpr{
								pos: position{line: 19, col: 34, offset: 533},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 19, col: 34, offset: 533},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 19, col: 37, offset: 536},
										name: "CodeBlock",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SeqExpr",
			pos:  position{line: 21, col: 1, offset: 550},
			expr: &seqExpr{
				pos: position{line: 21, col: 11, offset: 562},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 21, col: 11, offset: 562},
						label: "first",
						expr: &ruleRefExpr{
							pos:  position{line: 21, col: 17, offset: 568},
							name: "LabeledExpr",
						},
					},
					&labeledExpr{
						pos:   position{line: 21, col: 29, offset: 580},
						label: "rest",
						expr: &zeroOrMoreExpr{
							pos: position{line: 21, col: 34, offset: 585},
							expr: &seqExpr{
								pos: position{line: 21, col: 36, offset: 587},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 21, col: 36, offset: 587},
										name: "__",
									},
									&ruleRefExpr{
										pos:  position{line: 21, col: 39, offset: 590},
										name: "LabeledExpr",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LabeledExpr",
			pos:  position{line: 23, col: 1, offset: 606},
			expr: &choiceExpr{
				pos: position{line: 23, col: 15, offset: 622},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 23, col: 15, offset: 622},
						exprs: []interface{}{
							&labeledExpr{
								pos:   position{line: 23, col: 15, offset: 622},
								label: "label",
								expr: &ruleRefExpr{
									pos:  position{line: 23, col: 21, offset: 628},
									name: "Identifier",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 23, col: 32, offset: 639},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 23, col: 35, offset: 642},
								val:        ":",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 23, col: 39, offset: 646},
								name: "__",
							},
							&labeledExpr{
								pos:   position{line: 23, col: 42, offset: 649},
								label: "expr",
								expr: &ruleRefExpr{
									pos:  position{line: 23, col: 47, offset: 654},
									name: "PrefixedExpr",
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 23, col: 63, offset: 670},
						name: "PrefixedExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 23, col: 78, offset: 685},
						name: "ThrowExpr",
					},
				},
			},
		},
		{
			name: "PrefixedExpr",
			pos:  position{line: 25, col: 1, offset: 696},
			expr: &choiceExpr{
				pos: position{line: 25, col: 16, offset: 713},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 25, col: 16, offset: 713},
						exprs: []interface{}{
							&labeledExpr{
								pos:   position{line: 25, col: 16, offset: 713},
								label: "op",
								expr: &ruleRefExpr{
									pos:  position{line: 25, col: 19, offset: 716},
									name: "PrefixedOp",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 25, col: 30, offset: 727},
								name: "__",
							},
							&labeledExpr{
								pos:   position{line: 25, col: 33, offset: 730},
								label: "expr",
								expr: &ruleRefExpr{
									pos:  position{line: 25, col: 38, offset: 735},
									name: "SuffixedExpr",
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 25, col: 54, offset: 751},
						name: "SuffixedExpr",
					},
				},
			},
		},
		{
			name: "PrefixedOp",
			pos:  position{line: 27, col: 1, offset: 765},
			expr: &choiceExpr{
				pos: position{line: 27, col: 16, offset: 782},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 27, col: 16, offset: 782},
						val:        "&",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 27, col: 22, offset: 788},
						val:        "!",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SuffixedExpr",
			pos:  position{line: 29, col: 1, offset: 795},
			expr: &choiceExpr{
				pos: position{line: 29, col: 16, offset: 812},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 29, col: 16, offset: 812},
						exprs: []interface{}{
							&labeledExpr{
								pos:   position{line: 29, col: 16, offset: 812},
								label: "expr",
								expr: &ruleRefExpr{
									pos:  position{line: 29, col: 21, offset: 817},
									name: "PrimaryExpr",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 29, col: 33, offset: 829},
								name: "__",
							},
							&labeledExpr{
								pos:   position{line: 29, col: 36, offset: 832},
								label: "op",
								expr: &ruleRefExpr{
									pos:  position{line: 29, col: 39, offset: 835},
									name: "SuffixedOp",
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 29, col: 52, offset: 848},
						name: "PrimaryExpr",
					},
				},
			},
		},
		{
			name: "SuffixedOp",
			pos:  position{line: 31, col: 1, offset: 861},
			expr: &choiceExpr{
				pos: position{line: 31, col: 16, offset: 878},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 31, col: 16, offset: 878},
						val:        "?",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 31, col: 22, offset: 884},
						val:        "*",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 31, col: 28, offset: 890},
						val:        "+",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "PrimaryExpr",
			pos:  position{line: 33, col: 1, offset: 897},
			expr: &choiceExpr{
				pos: position{line: 33, col: 15, offset: 913},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 33, col: 15, offset: 913},
						name: "LitMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 28, offset: 926},
						name: "CharClassMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 47, offset: 945},
						name: "AnyMatcher",
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 60, offset: 958},
						name: "RuleRefExpr",
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 74, offset: 972},
						name: "SemanticPredExpr",
					},
					&seqExpr{
						pos: position{line: 33, col: 93, offset: 991},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 33, col: 93, offset: 991},
								val:        "(",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 33, col: 97, offset: 995},
								name: "__",
							},
							&labeledExpr{
								pos:   position{line: 33, col: 100, offset: 998},
								label: "expr",
								expr: &ruleRefExpr{
									pos:  position{line: 33, col: 105, offset: 1003},
									name: "Expression",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 33, col: 116, offset: 1014},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 33, col: 119, offset: 1017},
								val:        ")",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "RuleRefExpr",
			pos:  position{line: 35, col: 1, offset: 1022},
			expr: &seqExpr{
				pos: position{line: 35, col: 15, offset: 1038},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 35, col: 15, offset: 1038},
						label: "name",
						expr: &ruleRefExpr{
							pos:  position{line: 35, col: 20, offset: 1043},
							name: "IdentifierName",
						},
					},
					&notExpr{
						pos: position{line: 35, col: 35, offset: 1058},
						expr: &seqExpr{
							pos: position{line: 35, col: 38, offset: 1061},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 35, col: 38, offset: 1061},
									name: "__",
								},
								&zeroOrOneExpr{
									pos: position{line: 35, col: 41, offset: 1064},
									expr: &seqExpr{
										pos: position{line: 35, col: 43, offset: 1066},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 35, col: 43, offset: 1066},
												name: "StringLiteral",
											},
											&ruleRefExpr{
												pos:  position{line: 35, col: 57, offset: 1080},
												name: "__",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 35, col: 63, offset: 1086},
									name: "RuleDefOp",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SemanticPredExpr",
			pos:  position{line: 37, col: 1, offset: 1099},
			expr: &seqExpr{
				pos: position{line: 37, col: 20, offset: 1120},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 37, col: 20, offset: 1120},
						label: "op",
						expr: &ruleRefExpr{
							pos:  position{line: 37, col: 23, offset: 1123},
							name: "SemanticPredOp",
						},
					},
					&ruleRefExpr{
						pos:  position{line: 37, col: 38, offset: 1138},
						name: "__",
					},
					&labeledExpr{
						pos:   position{line: 37, col: 41, offset: 1141},
						label: "code",
						expr: &ruleRefExpr{
							pos:  position{line: 37, col: 46, offset: 1146},
							name: "CodeBlock",
						},
					},
				},
			},
		},
		{
			name: "SemanticPredOp",
			pos:  position{line: 39, col: 1, offset: 1157},
			expr: &choiceExpr{
				pos: position{line: 39, col: 20, offset: 1178},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 39, col: 20, offset: 1178},
						val:        "#",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 39, col: 26, offset: 1184},
						val:        "&",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 39, col: 32, offset: 1190},
						val:        "!",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "RuleDefOp",
			pos:  position{line: 41, col: 1, offset: 1197},
			expr: &choiceExpr{
				pos: position{line: 41, col: 13, offset: 1211},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 41, col: 13, offset: 1211},
						val:        "=",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 41, col: 19, offset: 1217},
						val:        "<-",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 41, col: 26, offset: 1224},
						val:        "←",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 41, col: 37, offset: 1235},
						val:        "⟵",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 43, col: 1, offset: 1245},
			expr: &anyMatcher{
				line: 43, col: 14, offset: 1260,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 44, col: 1, offset: 1262},
			expr: &choiceExpr{
				pos: position{line: 44, col: 11, offset: 1274},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 44, col: 11, offset: 1274},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 44, col: 30, offset: 1293},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 45, col: 1, offset: 1311},
			expr: &seqExpr{
				pos: position{line: 45, col: 20, offset: 1332},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 45, col: 20, offset: 1332},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 45, col: 25, offset: 1337},
						expr: &seqExpr{
							pos: position{line: 45, col: 27, offset: 1339},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 45, col: 27, offset: 1339},
									expr: &litMatcher{
										pos:        position{line: 45, col: 28, offset: 1340},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 45, col: 33, offset: 1345},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 45, col: 47, offset: 1359},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 46, col: 1, offset: 1364},
			expr: &seqExpr{
				pos: position{line: 46, col: 36, offset: 1401},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 46, col: 36, offset: 1401},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 46, col: 41, offset: 1406},
						expr: &seqExpr{
							pos: position{line: 46, col: 43, offset: 1408},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 46, col: 43, offset: 1408},
									expr: &choiceExpr{
										pos: position{line: 46, col: 46, offset: 1411},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 46, col: 46, offset: 1411},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 46, col: 53, offset: 1418},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 46, col: 59, offset: 1424},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 46, col: 73, offset: 1438},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 47, col: 1, offset: 1443},
			expr: &seqExpr{
				pos: position{line: 47, col: 21, offset: 1465},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 47, col: 21, offset: 1465},
						expr: &litMatcher{
							pos:        position{line: 47, col: 23, offset: 1467},
							val:        "//{",
							ignoreCase: false,
						},
					},
					&litMatcher{
						pos:        position{line: 47, col: 30, offset: 1474},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 47, col: 35, offset: 1479},
						expr: &seqExpr{
							pos: position{line: 47, col: 37, offset: 1481},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 47, col: 37, offset: 1481},
									expr: &ruleRefExpr{
										pos:  position{line: 47, col: 38, offset: 1482},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 47, col: 42, offset: 1486},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 49, col: 1, offset: 1501},
			expr: &labeledExpr{
				pos:   position{line: 49, col: 14, offset: 1516},
				label: "ident",
				expr: &ruleRefExpr{
					pos:  position{line: 49, col: 20, offset: 1522},
					name: "IdentifierName",
				},
			},
		},
		{
			name: "IdentifierName",
			pos:  position{line: 51, col: 1, offset: 1538},
			expr: &seqExpr{
				pos: position{line: 51, col: 18, offset: 1557},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 51, col: 18, offset: 1557},
						name: "IdentifierStart",
					},
					&zeroOrMoreExpr{
						pos: position{line: 51, col: 34, offset: 1573},
						expr: &ruleRefExpr{
							pos:  position{line: 51, col: 34, offset: 1573},
							name: "IdentifierPart",
						},
					},
				},
			},
		},
		{
			name: "IdentifierStart",
			pos:  position{line: 52, col: 1, offset: 1589},
			expr: &charClassMatcher{
				pos:        position{line: 52, col: 19, offset: 1609},
				val:        "[\\pL_]",
				chars:      []rune{'_'},
				classes:    []*unicode.RangeTable{rangeTable("L")},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IdentifierPart",
			pos:  position{line: 53, col: 1, offset: 1616},
			expr: &choiceExpr{
				pos: position{line: 53, col: 18, offset: 1635},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 53, col: 18, offset: 1635},
						name: "IdentifierStart",
					},
					&charClassMatcher{
						pos:        position{line: 53, col: 36, offset: 1653},
						val:        "[\\p{Nd}]",
						classes:    []*unicode.RangeTable{rangeTable("Nd")},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "LitMatcher",
			pos:  position{line: 55, col: 1, offset: 1663},
			expr: &seqExpr{
				pos: position{line: 55, col: 14, offset: 1678},
				exprs: []interface{}{
					&labeledExpr{
						pos:   position{line: 55, col: 14, offset: 1678},
						label: "lit",
						expr: &ruleRefExpr{
							pos:  position{line: 55, col: 18, offset: 1682},
							name: "StringLiteral",
						},
					},
					&labeledExpr{
						pos:   position{line: 55, col: 32, offset: 1696},
						label: "ignore",
						expr: &zeroOrOneExpr{
							pos: position{line: 55, col: 39, offset: 1703},
							expr: &litMatcher{
								pos:        position{line: 55, col: 39, offset: 1703},
								val:        "i",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 57, col: 1, offset: 1709},
			expr: &choiceExpr{
				pos: position{line: 57, col: 17, offset: 1727},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 57, col: 19, offset: 1729},
						alternatives: []interface{}{
							&seqExpr{
								pos: position{line: 57, col: 19, offset: 1729},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 57, col: 19, offset: 1729},
										val:        "\"",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 57, col: 23, offset: 1733},
										expr: &ruleRefExpr{
											pos:  position{line: 57, col: 23, offset: 1733},
											name: "DoubleStringChar",
										},
									},
									&litMatcher{
										pos:        position{line: 57, col: 41, offset: 1751},
										val:        "\"",
										ignoreCase: false,
									},
								},
							},
							&seqExpr{
								pos: position{line: 57, col: 47, offset: 1757},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 57, col: 47, offset: 1757},
										val:        "'",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 57, col: 51, offset: 1761},
										name: "SingleStringChar",
									},
									&litMatcher{
										pos:        position{line: 57, col: 68, offset: 1778},
										val:        "'",
										ignoreCase: false,
									},
								},
							},
							&seqExpr{
								pos: position{line: 57, col: 74, offset: 1784},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 57, col: 74, offset: 1784},
										val:        "`",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 57, col: 78, offset: 1788},
										expr: &ruleRefExpr{
											pos:  position{line: 57, col: 78, offset: 1788},
											name: "RawStringChar",
										},
									},
									&litMatcher{
										pos:        position{line: 57, col: 93, offset: 1803},
										val:        "`",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&choiceExpr{
						pos: position{line: 58, col: 6, offset: 1814},
						alternatives: []interface{}{
							&seqExpr{
								pos: position{line: 58, col: 8, offset: 1816},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 58, col: 8, offset: 1816},
										val:        "\"",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 58, col: 12, offset: 1820},
										expr: &ruleRefExpr{
											pos:  position{line: 58, col: 12, offset: 1820},
											name: "DoubleStringChar",
										},
									},
									&choiceExpr{
										pos: position{line: 58, col: 32, offset: 1840},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 58, col: 32, offset: 1840},
												name: "EOL",
											},
											&ruleRefExpr{
												pos:  position{line: 58, col: 38, offset: 1846},
												name: "EOF",
											},
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 58, col: 50, offset: 1858},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 58, col: 50, offset: 1858},
										val:        "'",
										ignoreCase: false,
									},
									&zeroOrOneExpr{
										pos: position{line: 58, col: 54, offset: 1862},
										expr: &ruleRefExpr{
											pos:  position{line: 58, col: 54, offset: 1862},
											name: "SingleStringChar",
										},
									},
									&choiceExpr{
										pos: position{line: 58, col: 74, offset: 1882},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 58, col: 74, offset: 1882},
												name: "EOL",
											},
											&ruleRefExpr{
												pos:  position{line: 58, col: 80, offset: 1888},
												name: "EOF",
											},
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 58, col: 90, offset: 1898},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 58, col: 90, offset: 1898},
										val:        "`",
										ignoreCase: false,
									},
									&zeroOrMoreExpr{
										pos: position{line: 58, col: 94, offset: 1902},
										expr: &ruleRefExpr{
											pos:  position{line: 58, col: 94, offset: 1902},
											name: "RawStringChar",
										},
									},
									&ruleRefExpr{
										pos:  position{line: 58, col: 109, offset: 1917},
										name: "EOF",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleStringChar",
			pos:  position{line: 61, col: 1, offset: 1925},
			expr: &choiceExpr{
				pos: position{line: 61, col: 20, offset: 1946},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 61, col: 20, offset: 1946},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 61, col: 20, offset: 1946},
								expr: &choiceExpr{
									pos: position{line: 61, col: 23, offset: 1949},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 61, col: 23, offset: 1949},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 61, col: 29, offset: 1955},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 61, col: 36, offset: 1962},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 61, col: 42, offset: 1968},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 61, col: 55, offset: 1981},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 61, col: 55, offset: 1981},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 61, col: 60, offset: 1986},
								name: "DoubleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringChar",
			pos:  position{line: 62, col: 1, offset: 2005},
			expr: &choiceExpr{
				pos: position{line: 62, col: 20, offset: 2026},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 62, col: 20, offset: 2026},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 62, col: 20, offset: 2026},
								expr: &choiceExpr{
									pos: position{line: 62, col: 23, offset: 2029},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 62, col: 23, offset: 2029},
											val:        "'",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 62, col: 29, offset: 2035},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 62, col: 36, offset: 2042},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 62, col: 42, offset: 2048},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 62, col: 55, offset: 2061},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 62, col: 55, offset: 2061},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 62, col: 60, offset: 2066},
								name: "SingleStringEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "RawStringChar",
			pos:  position{line: 63, col: 1, offset: 2085},
			expr: &seqExpr{
				pos: position{line: 63, col: 17, offset: 2103},
				exprs: []interface{}{
					&notExpr{
						pos: position{line: 63, col: 17, offset: 2103},
						expr: &litMatcher{
							pos:        position{line: 63, col: 18, offset: 2104},
							val:        "`",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 22, offset: 2108},
						name: "SourceChar",
					},
				},
			},
		},
		{
			name: "DoubleStringEscape",
			pos:  position{line: 65, col: 1, offset: 2120},
			expr: &choiceExpr{
				pos: position{line: 65, col: 22, offset: 2143},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 65, col: 24, offset: 2145},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 65, col: 24, offset: 2145},
								val:        "\"",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 65, col: 30, offset: 2151},
								name: "CommonEscapeSequence",
							},
						},
					},
					&choiceExpr{
						pos: position{line: 66, col: 9, offset: 2182},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 66, col: 9, offset: 2182},
								name: "SourceChar",
							},
							&ruleRefExpr{
								pos:  position{line: 66, col: 22, offset: 2195},
								name: "EOL",
							},
							&ruleRefExpr{
								pos:  position{line: 66, col: 28, offset: 2201},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "SingleStringEscape",
			pos:  position{line: 68, col: 1, offset: 2208},
			expr: &choiceExpr{
				pos: position{line: 68, col: 22, offset: 2231},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 68, col: 24, offset: 2233},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 68, col: 24, offset: 2233},
								val:        "'",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 68, col: 30, offset: 2239},
								name: "CommonEscapeSequence",
							},
						},
					},
					&choiceExpr{
						pos: position{line: 69, col: 9, offset: 2270},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 69, col: 9, offset: 2270},
								name: "SourceChar",
							},
							&ruleRefExpr{
								pos:  position{line: 69, col: 22, offset: 2283},
								name: "EOL",
							},
							&ruleRefExpr{
								pos:  position{line: 69, col: 28, offset: 2289},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "CommonEscapeSequence",
			pos:  position{line: 72, col: 1, offset: 2297},
			expr: &choiceExpr{
				pos: position{line: 72, col: 24, offset: 2322},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 72, col: 24, offset: 2322},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 43, offset: 2341},
						name: "OctalEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 57, offset: 2355},
						name: "HexEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 69, offset: 2367},
						name: "LongUnicodeEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 89, offset: 2387},
						name: "ShortUnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 73, col: 1, offset: 2406},
			expr: &choiceExpr{
				pos: position{line: 73, col: 20, offset: 2427},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 73, col: 20, offset: 2427},
						val:        "a",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 26, offset: 2433},
						val:        "b",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 32, offset: 2439},
						val:        "n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 38, offset: 2445},
						val:        "f",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 44, offset: 2451},
						val:        "r",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 50, offset: 2457},
						val:        "t",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 56, offset: 2463},
						val:        "v",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 73, col: 62, offset: 2469},
						val:        "\\",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "OctalEscape",
			pos:  position{line: 74, col: 1, offset: 2474},
			expr: &choiceExpr{
				pos: position{line: 74, col: 15, offset: 2490},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 74, col: 15, offset: 2490},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 74, col: 15, offset: 2490},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 74, col: 26, offset: 2501},
								name: "OctalDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 74, col: 37, offset: 2512},
								name: "OctalDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 75, col: 7, offset: 2529},
						run: (*parser).callonOctalEscape6,
						expr: &seqExpr{
							pos: position{line: 75, col: 7, offset: 2529},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 75, col: 7, offset: 2529},
									name: "OctalDigit",
								},
								&choiceExpr{
									pos: position{line: 75, col: 20, offset: 2542},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 75, col: 20, offset: 2542},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 75, col: 33, offset: 2555},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 75, col: 39, offset: 2561},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "HexEscape",
			pos:  position{line: 78, col: 1, offset: 2622},
			expr: &choiceExpr{
				pos: position{line: 78, col: 13, offset: 2636},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 78, col: 13, offset: 2636},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 78, col: 13, offset: 2636},
								val:        "x",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 78, col: 17, offset: 2640},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 78, col: 26, offset: 2649},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 79, col: 7, offset: 2664},
						run: (*parser).callonHexEscape6,
						expr: &seqExpr{
							pos: position{line: 79, col: 7, offset: 2664},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 79, col: 7, offset: 2664},
									val:        "x",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 79, col: 13, offset: 2670},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 79, col: 13, offset: 2670},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 79, col: 26, offset: 2683},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 79, col: 32, offset: 2689},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LongUnicodeEscape",
			pos:  position{line: 82, col: 1, offset: 2756},
			expr: &choiceExpr{
				pos: position{line: 83, col: 5, offset: 2782},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 83, col: 5, offset: 2782},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 83, col: 5, offset: 2782},
								val:        "U",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 9, offset: 2786},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 18, offset: 2795},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 27, offset: 2804},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 36, offset: 2813},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 45, offset: 2822},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 54, offset: 2831},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 63, offset: 2840},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 83, col: 72, offset: 2849},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 7, offset: 2864},
						run: (*parser).callonLongUnicodeEscape12,
						expr: &seqExpr{
							pos: position{line: 84, col: 7, offset: 2864},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 84, col: 7, offset: 2864},
									val:        "U",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 84, col: 13, offset: 2870},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 84, col: 13, offset: 2870},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 84, col: 26, offset: 2883},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 84, col: 32, offset: 2889},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ShortUnicodeEscape",
			pos:  position{line: 87, col: 1, offset: 2952},
			expr: &choiceExpr{
				pos: position{line: 88, col: 5, offset: 2979},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 88, col: 5, offset: 2979},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 88, col: 5, offset: 2979},
								val:        "u",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 88, col: 9, offset: 2983},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 88, col: 18, offset: 2992},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 88, col: 27, offset: 3001},
								name: "HexDigit",
							},
							&ruleRefExpr{
								pos:  position{line: 88, col: 36, offset: 3010},
								name: "HexDigit",
							},
						},
					},
					&actionExpr{
						pos: position{line: 89, col: 7, offset: 3025},
						run: (*parser).callonShortUnicodeEscape8,
						expr: &seqExpr{
							pos: position{line: 89, col: 7, offset: 3025},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 89, col: 7, offset: 3025},
									val:        "u",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 89, col: 13, offset: 3031},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 89, col: 13, offset: 3031},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 89, col: 26, offset: 3044},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 89, col: 32, offset: 3050},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OctalDigit",
			pos:  position{line: 93, col: 1, offset: 3114},
			expr: &charClassMatcher{
				pos:        position{line: 93, col: 14, offset: 3129},
				val:        "[0-7]",
				ranges:     []rune{'0', '7'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 94, col: 1, offset: 3135},
			expr: &charClassMatcher{
				pos:        position{line: 94, col: 16, offset: 3152},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 95, col: 1, offset: 3158},
			expr: &charClassMatcher{
				pos:        position{line: 95, col: 12, offset: 3171},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "CharClassMatcher",
			pos:  position{line: 97, col: 1, offset: 3182},
			expr: &choiceExpr{
				pos: position{line: 97, col: 20, offset: 3203},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 97, col: 20, offset: 3203},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 97, col: 20, offset: 3203},
								val:        "[",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 97, col: 24, offset: 3207},
								expr: &choiceExpr{
									pos: position{line: 97, col: 26, offset: 3209},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 97, col: 26, offset: 3209},
											name: "ClassCharRange",
										},
										&ruleRefExpr{
											pos:  position{line: 97, col: 43, offset: 3226},
											name: "ClassChar",
										},
										&seqExpr{
											pos: position{line: 97, col: 55, offset: 3238},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 97, col: 55, offset: 3238},
													val:        "\\",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 97, col: 60, offset: 3243},
													name: "UnicodeClassEscape",
												},
											},
										},
									},
								},
							},
							&litMatcher{
								pos:        position{line: 97, col: 82, offset: 3265},
								val:        "]",
								ignoreCase: false,
							},
							&zeroOrOneExpr{
								pos: position{line: 97, col: 86, offset: 3269},
								expr: &litMatcher{
									pos:        position{line: 97, col: 86, offset: 3269},
									val:        "i",
									ignoreCase: false,
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 98, col: 4, offset: 3277},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 98, col: 4, offset: 3277},
								val:        "[",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 98, col: 8, offset: 3281},
								expr: &seqExpr{
									pos: position{line: 98, col: 10, offset: 3283},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 98, col: 10, offset: 3283},
											expr: &ruleRefExpr{
												pos:  position{line: 98, col: 13, offset: 3286},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 98, col: 19, offset: 3292},
											name: "SourceChar",
										},
									},
								},
							},
							&choiceExpr{
								pos: position{line: 98, col: 35, offset: 3308},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 98, col: 35, offset: 3308},
										name: "EOL",
									},
									&ruleRefExpr{
										pos:  position{line: 98, col: 41, offset: 3314},
										name: "EOF",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ClassCharRange",
			pos:  position{line: 101, col: 1, offset: 3322},
			expr: &seqExpr{
				pos: position{line: 101, col: 18, offset: 3341},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 101, col: 18, offset: 3341},
						name: "ClassChar",
					},
					&litMatcher{
						pos:        position{line: 101, col: 28, offset: 3351},
						val:        "-",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 101, col: 32, offset: 3355},
						name: "ClassChar",
					},
				},
			},
		},
		{
			name: "ClassChar",
			pos:  position{line: 102, col: 1, offset: 3365},
			expr: &choiceExpr{
				pos: position{line: 102, col: 13, offset: 3379},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 102, col: 13, offset: 3379},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 102, col: 13, offset: 3379},
								expr: &choiceExpr{
									pos: position{line: 102, col: 16, offset: 3382},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 102, col: 16, offset: 3382},
											val:        "]",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 102, col: 22, offset: 3388},
											val:        "\\",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 102, col: 29, offset: 3395},
											name: "EOL",
										},
									},
								},
							},
							&ruleRefExpr{
								pos:  position{line: 102, col: 35, offset: 3401},
								name: "SourceChar",
							},
						},
					},
					&seqExpr{
						pos: position{line: 102, col: 48, offset: 3414},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 102, col: 48, offset: 3414},
								val:        "\\",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 102, col: 53, offset: 3419},
								name: "CharClassEscape",
							},
						},
					},
				},
			},
		},
		{
			name: "CharClassEscape",
			pos:  position{line: 103, col: 1, offset: 3435},
			expr: &choiceExpr{
				pos: position{line: 103, col: 19, offset: 3455},
				alternatives: []interface{}{
					&choiceExpr{
						pos: position{line: 103, col: 21, offset: 3457},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 103, col: 21, offset: 3457},
								val:        "]",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 103, col: 27, offset: 3463},
								name: "CommonEscapeSequence",
							},
						},
					},
					&actionExpr{
						pos: position{line: 104, col: 7, offset: 3492},
						run: (*parser).callonCharClassEscape5,
						expr: &seqExpr{
							pos: position{line: 104, col: 7, offset: 3492},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 104, col: 7, offset: 3492},
									expr: &litMatcher{
										pos:        position{line: 104, col: 8, offset: 3493},
										val:        "p",
										ignoreCase: false,
									},
								},
								&choiceExpr{
									pos: position{line: 104, col: 14, offset: 3499},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 104, col: 14, offset: 3499},
											name: "SourceChar",
										},
										&ruleRefExpr{
											pos:  position{line: 104, col: 27, offset: 3512},
											name: "EOL",
										},
										&ruleRefExpr{
											pos:  position{line: 104, col: 33, offset: 3518},
											name: "EOF",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "UnicodeClassEscape",
			pos:  position{line: 108, col: 1, offset: 3584},
			expr: &seqExpr{
				pos: position{line: 108, col: 22, offset: 3607},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 108, col: 22, offset: 3607},
						val:        "p",
						ignoreCase: false,
					},
					&choiceExpr{
						pos: position{line: 109, col: 7, offset: 3619},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 109, col: 7, offset: 3619},
								name: "SingleCharUnicodeClass",
							},
							&actionExpr{
								pos: position{line: 110, col: 7, offset: 3648},
								run: (*parser).callonUnicodeClassEscape5,
								expr: &seqExpr{
									pos: position{line: 110, col: 7, offset: 3648},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 110, col: 7, offset: 3648},
											expr: &litMatcher{
												pos:        position{line: 110, col: 8, offset: 3649},
												val:        "{",
												ignoreCase: false,
											},
										},
										&choiceExpr{
											pos: position{line: 110, col: 14, offset: 3655},
											alternatives: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 110, col: 14, offset: 3655},
													name: "SourceChar",
												},
												&ruleRefExpr{
													pos:  position{line: 110, col: 27, offset: 3668},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 110, col: 33, offset: 3674},
													name: "EOF",
												},
											},
										},
									},
								},
							},
							&seqExpr{
								pos: position{line: 111, col: 7, offset: 3745},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 111, col: 7, offset: 3745},
										val:        "{",
										ignoreCase: false,
									},
									&labeledExpr{
										pos:   position{line: 111, col: 11, offset: 3749},
										label: "ident",
										expr: &ruleRefExpr{
											pos:  position{line: 111, col: 17, offset: 3755},
											name: "IdentifierName",
										},
									},
									&litMatcher{
										pos:        position{line: 111, col: 32, offset: 3770},
										val:        "}",
										ignoreCase: false,
									},
								},
							},
							&actionExpr{
								pos: position{line: 112, col: 7, offset: 3780},
								run: (*parser).callonUnicodeClassEscape18,
								expr: &seqExpr{
									pos: position{line: 112, col: 7, offset: 3780},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 112, col: 7, offset: 3780},
											val:        "{",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 112, col: 11, offset: 3784},
											name: "IdentifierName",
										},
										&choiceExpr{
											pos: position{line: 112, col: 28, offset: 3801},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 112, col: 28, offset: 3801},
													val:        "]",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 112, col: 34, offset: 3807},
													name: "EOL",
												},
												&ruleRefExpr{
													pos:  position{line: 112, col: 40, offset: 3813},
													name: "EOF",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SingleCharUnicodeClass",
			pos:  position{line: 116, col: 1, offset: 3896},
			expr: &charClassMatcher{
				pos:        position{line: 116, col: 26, offset: 3923},
				val:        "[LMNCPZS]",
				chars:      []rune{'L', 'M', 'N', 'C', 'P', 'Z', 'S'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "AnyMatcher",
			pos:  position{line: 118, col: 1, offset: 3934},
			expr: &litMatcher{
				pos:        position{line: 118, col: 14, offset: 3949},
				val:        ".",
				ignoreCase: false,
			},
		},
		{
			name: "ThrowExpr",
			pos:  position{line: 120, col: 1, offset: 3954},
			expr: &choiceExpr{
				pos: position{line: 120, col: 13, offset: 3968},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 120, col: 13, offset: 3968},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 120, col: 13, offset: 3968},
								val:        "%",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 120, col: 17, offset: 3972},
								val:        "{",
								ignoreCase: false,
							},
							&labeledExpr{
								pos:   position{line: 120, col: 21, offset: 3976},
								label: "label",
								expr: &ruleRefExpr{
									pos:  position{line: 120, col: 27, offset: 3982},
									name: "IdentifierName",
								},
							},
							&litMatcher{
								pos:        position{line: 120, col: 42, offset: 3997},
								val:        "}",
								ignoreCase: false,
							},
						},
					},
					&actionExpr{
						pos: position{line: 121, col: 4, offset: 4004},
						run: (*parser).callonThrowExpr8,
						expr: &seqExpr{
							pos: position{line: 121, col: 4, offset: 4004},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 121, col: 4, offset: 4004},
									val:        "%",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 121, col: 8, offset: 4008},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 121, col: 12, offset: 4012},
									name: "IdentifierName",
								},
								&ruleRefExpr{
									pos:  position{line: 121, col: 27, offset: 4027},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 125, col: 1, offset: 4098},
			expr: &choiceExpr{
				pos: position{line: 125, col: 13, offset: 4112},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 125, col: 13, offset: 4112},
						run: (*parser).callonCodeBlock2,
						expr: &seqExpr{
							pos: position{line: 125, col: 13, offset: 4112},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 125, col: 13, offset: 4112},
									val:        "{",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 125, col: 17, offset: 4116},
									label: "code",
									expr: &ruleRefExpr{
										pos:  position{line: 125, col: 22, offset: 4121},
										name: "Code",
									},
								},
								&litMatcher{
									pos:        position{line: 125, col: 27, offset: 4126},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 128, col: 4, offset: 4166},
						run: (*parser).callonCodeBlock8,
						expr: &seqExpr{
							pos: position{line: 128, col: 4, offset: 4166},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 128, col: 4, offset: 4166},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 8, offset: 4170},
									name: "Code",
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 13, offset: 4175},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 132, col: 1, offset: 4240},
			expr: &zeroOrMoreExpr{
				pos: position{line: 132, col: 8, offset: 4249},
				expr: &choiceExpr{
					pos: position{line: 132, col: 10, offset: 4251},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 132, col: 10, offset: 4251},
							expr: &seqExpr{
								pos: position{line: 132, col: 12, offset: 4253},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 132, col: 12, offset: 4253},
										expr: &charClassMatcher{
											pos:        position{line: 132, col: 13, offset: 4254},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&ruleRefExpr{
										pos:  position{line: 132, col: 18, offset: 4259},
										name: "SourceChar",
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 132, col: 34, offset: 4275},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 132, col: 34, offset: 4275},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 132, col: 38, offset: 4279},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 132, col: 43, offset: 4284},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 134, col: 1, offset: 4292},
			expr: &zeroOrMoreExpr{
				pos: position{line: 134, col: 6, offset: 4299},
				expr: &choiceExpr{
					pos: position{line: 134, col: 8, offset: 4301},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 134, col: 8, offset: 4301},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 134, col: 21, offset: 4314},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 134, col: 27, offset: 4320},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 135, col: 1, offset: 4331},
			expr: &zeroOrMoreExpr{
				pos: position{line: 135, col: 5, offset: 4337},
				expr: &choiceExpr{
					pos: position{line: 135, col: 7, offset: 4339},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 135, col: 7, offset: 4339},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 20, offset: 4352},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 137, col: 1, offset: 4389},
			expr: &charClassMatcher{
				pos:        position{line: 137, col: 14, offset: 4404},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 138, col: 1, offset: 4412},
			expr: &litMatcher{
				pos:        position{line: 138, col: 7, offset: 4420},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 139, col: 1, offset: 4425},
			expr: &choiceExpr{
				pos: position{line: 139, col: 7, offset: 4433},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 139, col: 7, offset: 4433},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 139, col: 7, offset: 4433},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 139, col: 10, offset: 4436},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 139, col: 16, offset: 4442},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 139, col: 16, offset: 4442},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 139, col: 18, offset: 4444},
								expr: &ruleRefExpr{
									pos:  position{line: 139, col: 18, offset: 4444},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 139, col: 37, offset: 4463},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 139, col: 43, offset: 4469},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 139, col: 43, offset: 4469},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 139, col: 46, offset: 4472},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 141, col: 1, offset: 4477},
			expr: &notExpr{
				pos: position{line: 141, col: 7, offset: 4485},
				expr: &anyMatcher{
					line: 141, col: 8, offset: 4486,
				},
			},
		},
	},
}

func (c *current) onGrammar1(initializer, rules interface{}) (interface{}, error) {
	return c.text, nil
}

func (p *parser) callonGrammar1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammar1(stack["initializer"], stack["rules"])
}

func (c *current) onOctalEscape6() (interface{}, error) {
	return nil, errors.New("invalid octal escape")
}

func (p *parser) callonOctalEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOctalEscape6()
}

func (c *current) onHexEscape6() (interface{}, error) {
	return nil, errors.New("invalid hexadecimal escape")
}

func (p *parser) callonHexEscape6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexEscape6()
}

func (c *current) onLongUnicodeEscape12() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonLongUnicodeEscape12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLongUnicodeEscape12()
}

func (c *current) onShortUnicodeEscape8() (interface{}, error) {
	return nil, errors.New("invalid Unicode escape")
}

func (p *parser) callonShortUnicodeEscape8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onShortUnicodeEscape8()
}

func (c *current) onCharClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid escape character")
}

func (p *parser) callonCharClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCharClassEscape5()
}

func (c *current) onUnicodeClassEscape5() (interface{}, error) {
	return nil, errors.New("invalid Unicode class escape")
}

func (p *parser) callonUnicodeClassEscape5() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape5()
}

func (c *current) onUnicodeClassEscape18() (interface{}, error) {
	return nil, errors.New("Unicode class not terminated")

}

func (p *parser) callonUnicodeClassEscape18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnicodeClassEscape18()
}

func (c *current) onThrowExpr8() (interface{}, error) {
	return nil, errors.New("throw expression not terminated")
}

func (p *parser) callonThrowExpr8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onThrowExpr8()
}

func (c *current) onCodeBlock2(code interface{}) (interface{}, error) {
	return doNothing(c.text)
}

func (p *parser) callonCodeBlock2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock2(stack["code"])
}

func (c *current) onCodeBlock8() (interface{}, error) {
	return nil, errors.New("code block not terminated")
}

func (p *parser) callonCodeBlock8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock8()
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

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
