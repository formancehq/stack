// Code generated from NumScript.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // NumScript

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type NumScriptParser struct {
	*antlr.BaseParser
}

var numscriptParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func numscriptParserInit() {
	staticData := &numscriptParserStaticData
	staticData.literalNames = []string{
		"", "'*'", "'allowing overdraft up to'", "'allowing unbounded overdraft'",
		"','", "", "", "", "", "'vars'", "'meta'", "'set_tx_meta'", "'set_account_meta'",
		"'fail'", "'send'", "'source'", "'from'", "'max'", "'destination'",
		"'to'", "'allocate'", "'+'", "'-'", "'('", "')'", "'['", "']'", "'{'",
		"'}'", "'='", "'account'", "'asset'", "'number'", "'monetary'", "'portion'",
		"'string'", "", "", "'remaining'", "'kept'", "'balance'", "", "'%'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "NEWLINE", "WHITESPACE", "MULTILINE_COMMENT", "LINE_COMMENT",
		"VARS", "META", "SET_TX_META", "SET_ACCOUNT_META", "FAIL", "SEND", "SOURCE",
		"FROM", "MAX", "DESTINATION", "TO", "ALLOCATE", "OP_ADD", "OP_SUB",
		"LPAREN", "RPAREN", "LBRACK", "RBRACK", "LBRACE", "RBRACE", "EQ", "TY_ACCOUNT",
		"TY_ASSET", "TY_NUMBER", "TY_MONETARY", "TY_PORTION", "TY_STRING", "STRING",
		"PORTION", "REMAINING", "KEPT", "BALANCE", "NUMBER", "PERCENT", "VARIABLE_NAME",
		"ACCOUNT", "ASSET",
	}
	staticData.ruleNames = []string{
		"monetary", "monetaryAll", "literal", "variable", "expression", "allotmentPortion",
		"destinationInOrder", "destinationAllotment", "keptOrDestination", "destination",
		"sourceAccountOverdraft", "sourceAccount", "sourceInOrder", "sourceMaxed",
		"source", "sourceAllotment", "valueAwareSource", "monetaryArithmetic",
		"monetaryExpression", "statement", "type_", "origin", "varDecl", "varListDecl",
		"script",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 45, 287, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 3, 2, 67, 8, 2, 1, 3, 1, 3, 1, 4, 1, 4, 3, 4, 73, 8, 4, 1, 5, 1,
		5, 1, 5, 3, 5, 78, 8, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 4, 6,
		87, 8, 6, 11, 6, 12, 6, 88, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1,
		7, 1, 7, 1, 7, 1, 7, 4, 7, 102, 8, 7, 11, 7, 12, 7, 103, 1, 7, 1, 7, 1,
		8, 1, 8, 1, 8, 3, 8, 111, 8, 8, 1, 9, 1, 9, 1, 9, 3, 9, 116, 8, 9, 1, 10,
		1, 10, 1, 10, 3, 10, 121, 8, 10, 1, 11, 1, 11, 3, 11, 125, 8, 11, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 12, 4, 12, 132, 8, 12, 11, 12, 12, 12, 133, 1,
		12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 3, 14,
		146, 8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 4, 15, 155,
		8, 15, 11, 15, 12, 15, 156, 1, 15, 1, 15, 1, 16, 1, 16, 3, 16, 163, 8,
		16, 1, 17, 1, 17, 1, 17, 4, 17, 168, 8, 17, 11, 17, 12, 17, 169, 3, 17,
		172, 8, 17, 1, 18, 1, 18, 3, 18, 176, 8, 18, 1, 19, 1, 19, 1, 19, 1, 19,
		1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1,
		19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19,
		1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1,
		19, 1, 19, 3, 19, 215, 8, 19, 1, 19, 1, 19, 1, 19, 3, 19, 220, 8, 19, 1,
		20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21,
		1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 238, 8, 21, 1, 22, 1, 22, 1,
		22, 1, 22, 3, 22, 244, 8, 22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 4, 23,
		251, 8, 23, 11, 23, 12, 23, 252, 4, 23, 255, 8, 23, 11, 23, 12, 23, 256,
		1, 23, 1, 23, 1, 23, 1, 24, 5, 24, 263, 8, 24, 10, 24, 12, 24, 266, 9,
		24, 1, 24, 3, 24, 269, 8, 24, 1, 24, 1, 24, 1, 24, 5, 24, 274, 8, 24, 10,
		24, 12, 24, 277, 9, 24, 1, 24, 5, 24, 280, 8, 24, 10, 24, 12, 24, 283,
		9, 24, 1, 24, 1, 24, 1, 24, 0, 0, 25, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18,
		20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 0, 2, 1, 0,
		21, 22, 1, 0, 30, 35, 296, 0, 50, 1, 0, 0, 0, 2, 55, 1, 0, 0, 0, 4, 66,
		1, 0, 0, 0, 6, 68, 1, 0, 0, 0, 8, 72, 1, 0, 0, 0, 10, 77, 1, 0, 0, 0, 12,
		79, 1, 0, 0, 0, 14, 95, 1, 0, 0, 0, 16, 110, 1, 0, 0, 0, 18, 115, 1, 0,
		0, 0, 20, 120, 1, 0, 0, 0, 22, 122, 1, 0, 0, 0, 24, 126, 1, 0, 0, 0, 26,
		137, 1, 0, 0, 0, 28, 145, 1, 0, 0, 0, 30, 147, 1, 0, 0, 0, 32, 162, 1,
		0, 0, 0, 34, 164, 1, 0, 0, 0, 36, 175, 1, 0, 0, 0, 38, 219, 1, 0, 0, 0,
		40, 221, 1, 0, 0, 0, 42, 237, 1, 0, 0, 0, 44, 239, 1, 0, 0, 0, 46, 245,
		1, 0, 0, 0, 48, 264, 1, 0, 0, 0, 50, 51, 5, 25, 0, 0, 51, 52, 3, 8, 4,
		0, 52, 53, 5, 41, 0, 0, 53, 54, 5, 26, 0, 0, 54, 1, 1, 0, 0, 0, 55, 56,
		5, 25, 0, 0, 56, 57, 3, 8, 4, 0, 57, 58, 5, 1, 0, 0, 58, 59, 5, 26, 0,
		0, 59, 3, 1, 0, 0, 0, 60, 67, 5, 44, 0, 0, 61, 67, 5, 45, 0, 0, 62, 67,
		5, 41, 0, 0, 63, 67, 5, 36, 0, 0, 64, 67, 5, 37, 0, 0, 65, 67, 3, 0, 0,
		0, 66, 60, 1, 0, 0, 0, 66, 61, 1, 0, 0, 0, 66, 62, 1, 0, 0, 0, 66, 63,
		1, 0, 0, 0, 66, 64, 1, 0, 0, 0, 66, 65, 1, 0, 0, 0, 67, 5, 1, 0, 0, 0,
		68, 69, 5, 43, 0, 0, 69, 7, 1, 0, 0, 0, 70, 73, 3, 4, 2, 0, 71, 73, 3,
		6, 3, 0, 72, 70, 1, 0, 0, 0, 72, 71, 1, 0, 0, 0, 73, 9, 1, 0, 0, 0, 74,
		78, 5, 37, 0, 0, 75, 78, 3, 6, 3, 0, 76, 78, 5, 38, 0, 0, 77, 74, 1, 0,
		0, 0, 77, 75, 1, 0, 0, 0, 77, 76, 1, 0, 0, 0, 78, 11, 1, 0, 0, 0, 79, 80,
		5, 27, 0, 0, 80, 86, 5, 5, 0, 0, 81, 82, 5, 17, 0, 0, 82, 83, 3, 8, 4,
		0, 83, 84, 3, 16, 8, 0, 84, 85, 5, 5, 0, 0, 85, 87, 1, 0, 0, 0, 86, 81,
		1, 0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 86, 1, 0, 0, 0, 88, 89, 1, 0, 0, 0,
		89, 90, 1, 0, 0, 0, 90, 91, 5, 38, 0, 0, 91, 92, 3, 16, 8, 0, 92, 93, 5,
		5, 0, 0, 93, 94, 5, 28, 0, 0, 94, 13, 1, 0, 0, 0, 95, 96, 5, 27, 0, 0,
		96, 101, 5, 5, 0, 0, 97, 98, 3, 10, 5, 0, 98, 99, 3, 16, 8, 0, 99, 100,
		5, 5, 0, 0, 100, 102, 1, 0, 0, 0, 101, 97, 1, 0, 0, 0, 102, 103, 1, 0,
		0, 0, 103, 101, 1, 0, 0, 0, 103, 104, 1, 0, 0, 0, 104, 105, 1, 0, 0, 0,
		105, 106, 5, 28, 0, 0, 106, 15, 1, 0, 0, 0, 107, 108, 5, 19, 0, 0, 108,
		111, 3, 18, 9, 0, 109, 111, 5, 39, 0, 0, 110, 107, 1, 0, 0, 0, 110, 109,
		1, 0, 0, 0, 111, 17, 1, 0, 0, 0, 112, 116, 3, 8, 4, 0, 113, 116, 3, 12,
		6, 0, 114, 116, 3, 14, 7, 0, 115, 112, 1, 0, 0, 0, 115, 113, 1, 0, 0, 0,
		115, 114, 1, 0, 0, 0, 116, 19, 1, 0, 0, 0, 117, 118, 5, 2, 0, 0, 118, 121,
		3, 8, 4, 0, 119, 121, 5, 3, 0, 0, 120, 117, 1, 0, 0, 0, 120, 119, 1, 0,
		0, 0, 121, 21, 1, 0, 0, 0, 122, 124, 3, 8, 4, 0, 123, 125, 3, 20, 10, 0,
		124, 123, 1, 0, 0, 0, 124, 125, 1, 0, 0, 0, 125, 23, 1, 0, 0, 0, 126, 127,
		5, 27, 0, 0, 127, 131, 5, 5, 0, 0, 128, 129, 3, 28, 14, 0, 129, 130, 5,
		5, 0, 0, 130, 132, 1, 0, 0, 0, 131, 128, 1, 0, 0, 0, 132, 133, 1, 0, 0,
		0, 133, 131, 1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 135, 1, 0, 0, 0, 135,
		136, 5, 28, 0, 0, 136, 25, 1, 0, 0, 0, 137, 138, 5, 17, 0, 0, 138, 139,
		3, 8, 4, 0, 139, 140, 5, 16, 0, 0, 140, 141, 3, 28, 14, 0, 141, 27, 1,
		0, 0, 0, 142, 146, 3, 22, 11, 0, 143, 146, 3, 26, 13, 0, 144, 146, 3, 24,
		12, 0, 145, 142, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0, 145, 144, 1, 0, 0, 0,
		146, 29, 1, 0, 0, 0, 147, 148, 5, 27, 0, 0, 148, 154, 5, 5, 0, 0, 149,
		150, 3, 10, 5, 0, 150, 151, 5, 16, 0, 0, 151, 152, 3, 28, 14, 0, 152, 153,
		5, 5, 0, 0, 153, 155, 1, 0, 0, 0, 154, 149, 1, 0, 0, 0, 155, 156, 1, 0,
		0, 0, 156, 154, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157, 158, 1, 0, 0, 0,
		158, 159, 5, 28, 0, 0, 159, 31, 1, 0, 0, 0, 160, 163, 3, 28, 14, 0, 161,
		163, 3, 30, 15, 0, 162, 160, 1, 0, 0, 0, 162, 161, 1, 0, 0, 0, 163, 33,
		1, 0, 0, 0, 164, 171, 3, 8, 4, 0, 165, 166, 7, 0, 0, 0, 166, 168, 3, 8,
		4, 0, 167, 165, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0,
		169, 170, 1, 0, 0, 0, 170, 172, 1, 0, 0, 0, 171, 167, 1, 0, 0, 0, 171,
		172, 1, 0, 0, 0, 172, 35, 1, 0, 0, 0, 173, 176, 3, 2, 1, 0, 174, 176, 3,
		34, 17, 0, 175, 173, 1, 0, 0, 0, 175, 174, 1, 0, 0, 0, 176, 37, 1, 0, 0,
		0, 177, 178, 5, 11, 0, 0, 178, 179, 5, 23, 0, 0, 179, 180, 5, 36, 0, 0,
		180, 181, 5, 4, 0, 0, 181, 182, 3, 8, 4, 0, 182, 183, 5, 24, 0, 0, 183,
		220, 1, 0, 0, 0, 184, 185, 5, 12, 0, 0, 185, 186, 5, 23, 0, 0, 186, 187,
		3, 8, 4, 0, 187, 188, 5, 4, 0, 0, 188, 189, 5, 36, 0, 0, 189, 190, 5, 4,
		0, 0, 190, 191, 3, 8, 4, 0, 191, 192, 5, 24, 0, 0, 192, 220, 1, 0, 0, 0,
		193, 220, 5, 13, 0, 0, 194, 195, 5, 14, 0, 0, 195, 196, 3, 36, 18, 0, 196,
		197, 5, 23, 0, 0, 197, 214, 5, 5, 0, 0, 198, 199, 5, 15, 0, 0, 199, 200,
		5, 29, 0, 0, 200, 201, 3, 32, 16, 0, 201, 202, 5, 5, 0, 0, 202, 203, 5,
		18, 0, 0, 203, 204, 5, 29, 0, 0, 204, 205, 3, 18, 9, 0, 205, 215, 1, 0,
		0, 0, 206, 207, 5, 18, 0, 0, 207, 208, 5, 29, 0, 0, 208, 209, 3, 18, 9,
		0, 209, 210, 5, 5, 0, 0, 210, 211, 5, 15, 0, 0, 211, 212, 5, 29, 0, 0,
		212, 213, 3, 32, 16, 0, 213, 215, 1, 0, 0, 0, 214, 198, 1, 0, 0, 0, 214,
		206, 1, 0, 0, 0, 215, 216, 1, 0, 0, 0, 216, 217, 5, 5, 0, 0, 217, 218,
		5, 24, 0, 0, 218, 220, 1, 0, 0, 0, 219, 177, 1, 0, 0, 0, 219, 184, 1, 0,
		0, 0, 219, 193, 1, 0, 0, 0, 219, 194, 1, 0, 0, 0, 220, 39, 1, 0, 0, 0,
		221, 222, 7, 1, 0, 0, 222, 41, 1, 0, 0, 0, 223, 224, 5, 10, 0, 0, 224,
		225, 5, 23, 0, 0, 225, 226, 3, 8, 4, 0, 226, 227, 5, 4, 0, 0, 227, 228,
		5, 36, 0, 0, 228, 229, 5, 24, 0, 0, 229, 238, 1, 0, 0, 0, 230, 231, 5,
		40, 0, 0, 231, 232, 5, 23, 0, 0, 232, 233, 3, 8, 4, 0, 233, 234, 5, 4,
		0, 0, 234, 235, 3, 8, 4, 0, 235, 236, 5, 24, 0, 0, 236, 238, 1, 0, 0, 0,
		237, 223, 1, 0, 0, 0, 237, 230, 1, 0, 0, 0, 238, 43, 1, 0, 0, 0, 239, 240,
		3, 40, 20, 0, 240, 243, 3, 6, 3, 0, 241, 242, 5, 29, 0, 0, 242, 244, 3,
		42, 21, 0, 243, 241, 1, 0, 0, 0, 243, 244, 1, 0, 0, 0, 244, 45, 1, 0, 0,
		0, 245, 246, 5, 9, 0, 0, 246, 247, 5, 27, 0, 0, 247, 254, 5, 5, 0, 0, 248,
		250, 3, 44, 22, 0, 249, 251, 5, 5, 0, 0, 250, 249, 1, 0, 0, 0, 251, 252,
		1, 0, 0, 0, 252, 250, 1, 0, 0, 0, 252, 253, 1, 0, 0, 0, 253, 255, 1, 0,
		0, 0, 254, 248, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 254, 1, 0, 0, 0,
		256, 257, 1, 0, 0, 0, 257, 258, 1, 0, 0, 0, 258, 259, 5, 28, 0, 0, 259,
		260, 5, 5, 0, 0, 260, 47, 1, 0, 0, 0, 261, 263, 5, 5, 0, 0, 262, 261, 1,
		0, 0, 0, 263, 266, 1, 0, 0, 0, 264, 262, 1, 0, 0, 0, 264, 265, 1, 0, 0,
		0, 265, 268, 1, 0, 0, 0, 266, 264, 1, 0, 0, 0, 267, 269, 3, 46, 23, 0,
		268, 267, 1, 0, 0, 0, 268, 269, 1, 0, 0, 0, 269, 270, 1, 0, 0, 0, 270,
		275, 3, 38, 19, 0, 271, 272, 5, 5, 0, 0, 272, 274, 3, 38, 19, 0, 273, 271,
		1, 0, 0, 0, 274, 277, 1, 0, 0, 0, 275, 273, 1, 0, 0, 0, 275, 276, 1, 0,
		0, 0, 276, 281, 1, 0, 0, 0, 277, 275, 1, 0, 0, 0, 278, 280, 5, 5, 0, 0,
		279, 278, 1, 0, 0, 0, 280, 283, 1, 0, 0, 0, 281, 279, 1, 0, 0, 0, 281,
		282, 1, 0, 0, 0, 282, 284, 1, 0, 0, 0, 283, 281, 1, 0, 0, 0, 284, 285,
		5, 0, 0, 1, 285, 49, 1, 0, 0, 0, 26, 66, 72, 77, 88, 103, 110, 115, 120,
		124, 133, 145, 156, 162, 169, 171, 175, 214, 219, 237, 243, 252, 256, 264,
		268, 275, 281,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// NumScriptParserInit initializes any static state used to implement NumScriptParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewNumScriptParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func NumScriptParserInit() {
	staticData := &numscriptParserStaticData
	staticData.once.Do(numscriptParserInit)
}

// NewNumScriptParser produces a new parser instance for the optional input antlr.TokenStream.
func NewNumScriptParser(input antlr.TokenStream) *NumScriptParser {
	NumScriptParserInit()
	this := new(NumScriptParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &numscriptParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "NumScript.g4"

	return this
}

// NumScriptParser tokens.
const (
	NumScriptParserEOF               = antlr.TokenEOF
	NumScriptParserT__0              = 1
	NumScriptParserT__1              = 2
	NumScriptParserT__2              = 3
	NumScriptParserT__3              = 4
	NumScriptParserNEWLINE           = 5
	NumScriptParserWHITESPACE        = 6
	NumScriptParserMULTILINE_COMMENT = 7
	NumScriptParserLINE_COMMENT      = 8
	NumScriptParserVARS              = 9
	NumScriptParserMETA              = 10
	NumScriptParserSET_TX_META       = 11
	NumScriptParserSET_ACCOUNT_META  = 12
	NumScriptParserFAIL              = 13
	NumScriptParserSEND              = 14
	NumScriptParserSOURCE            = 15
	NumScriptParserFROM              = 16
	NumScriptParserMAX               = 17
	NumScriptParserDESTINATION       = 18
	NumScriptParserTO                = 19
	NumScriptParserALLOCATE          = 20
	NumScriptParserOP_ADD            = 21
	NumScriptParserOP_SUB            = 22
	NumScriptParserLPAREN            = 23
	NumScriptParserRPAREN            = 24
	NumScriptParserLBRACK            = 25
	NumScriptParserRBRACK            = 26
	NumScriptParserLBRACE            = 27
	NumScriptParserRBRACE            = 28
	NumScriptParserEQ                = 29
	NumScriptParserTY_ACCOUNT        = 30
	NumScriptParserTY_ASSET          = 31
	NumScriptParserTY_NUMBER         = 32
	NumScriptParserTY_MONETARY       = 33
	NumScriptParserTY_PORTION        = 34
	NumScriptParserTY_STRING         = 35
	NumScriptParserSTRING            = 36
	NumScriptParserPORTION           = 37
	NumScriptParserREMAINING         = 38
	NumScriptParserKEPT              = 39
	NumScriptParserBALANCE           = 40
	NumScriptParserNUMBER            = 41
	NumScriptParserPERCENT           = 42
	NumScriptParserVARIABLE_NAME     = 43
	NumScriptParserACCOUNT           = 44
	NumScriptParserASSET             = 45
)

// NumScriptParser rules.
const (
	NumScriptParserRULE_monetary               = 0
	NumScriptParserRULE_monetaryAll            = 1
	NumScriptParserRULE_literal                = 2
	NumScriptParserRULE_variable               = 3
	NumScriptParserRULE_expression             = 4
	NumScriptParserRULE_allotmentPortion       = 5
	NumScriptParserRULE_destinationInOrder     = 6
	NumScriptParserRULE_destinationAllotment   = 7
	NumScriptParserRULE_keptOrDestination      = 8
	NumScriptParserRULE_destination            = 9
	NumScriptParserRULE_sourceAccountOverdraft = 10
	NumScriptParserRULE_sourceAccount          = 11
	NumScriptParserRULE_sourceInOrder          = 12
	NumScriptParserRULE_sourceMaxed            = 13
	NumScriptParserRULE_source                 = 14
	NumScriptParserRULE_sourceAllotment        = 15
	NumScriptParserRULE_valueAwareSource       = 16
	NumScriptParserRULE_monetaryArithmetic     = 17
	NumScriptParserRULE_monetaryExpression     = 18
	NumScriptParserRULE_statement              = 19
	NumScriptParserRULE_type_                  = 20
	NumScriptParserRULE_origin                 = 21
	NumScriptParserRULE_varDecl                = 22
	NumScriptParserRULE_varListDecl            = 23
	NumScriptParserRULE_script                 = 24
)

// IMonetaryContext is an interface to support dynamic dispatch.
type IMonetaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAmt returns the amt token.
	GetAmt() antlr.Token

	// SetAmt sets the amt token.
	SetAmt(antlr.Token)

	// GetAsset returns the asset rule contexts.
	GetAsset() IExpressionContext

	// SetAsset sets the asset rule contexts.
	SetAsset(IExpressionContext)

	// IsMonetaryContext differentiates from other interfaces.
	IsMonetaryContext()
}

type MonetaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	asset  IExpressionContext
	amt    antlr.Token
}

func NewEmptyMonetaryContext() *MonetaryContext {
	var p = new(MonetaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_monetary
	return p
}

func (*MonetaryContext) IsMonetaryContext() {}

func NewMonetaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MonetaryContext {
	var p = new(MonetaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_monetary

	return p
}

func (s *MonetaryContext) GetParser() antlr.Parser { return s.parser }

func (s *MonetaryContext) GetAmt() antlr.Token { return s.amt }

func (s *MonetaryContext) SetAmt(v antlr.Token) { s.amt = v }

func (s *MonetaryContext) GetAsset() IExpressionContext { return s.asset }

func (s *MonetaryContext) SetAsset(v IExpressionContext) { s.asset = v }

func (s *MonetaryContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACK, 0)
}

func (s *MonetaryContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACK, 0)
}

func (s *MonetaryContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MonetaryContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(NumScriptParserNUMBER, 0)
}

func (s *MonetaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MonetaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterMonetary(s)
	}
}

func (s *MonetaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitMonetary(s)
	}
}

func (p *NumScriptParser) Monetary() (localctx IMonetaryContext) {
	this := p
	_ = this

	localctx = NewMonetaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, NumScriptParserRULE_monetary)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(50)
		p.Match(NumScriptParserLBRACK)
	}
	{
		p.SetState(51)

		var _x = p.Expression()

		localctx.(*MonetaryContext).asset = _x
	}
	{
		p.SetState(52)

		var _m = p.Match(NumScriptParserNUMBER)

		localctx.(*MonetaryContext).amt = _m
	}
	{
		p.SetState(53)
		p.Match(NumScriptParserRBRACK)
	}

	return localctx
}

// IMonetaryAllContext is an interface to support dynamic dispatch.
type IMonetaryAllContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAsset returns the asset rule contexts.
	GetAsset() IExpressionContext

	// SetAsset sets the asset rule contexts.
	SetAsset(IExpressionContext)

	// IsMonetaryAllContext differentiates from other interfaces.
	IsMonetaryAllContext()
}

type MonetaryAllContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	asset  IExpressionContext
}

func NewEmptyMonetaryAllContext() *MonetaryAllContext {
	var p = new(MonetaryAllContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_monetaryAll
	return p
}

func (*MonetaryAllContext) IsMonetaryAllContext() {}

func NewMonetaryAllContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MonetaryAllContext {
	var p = new(MonetaryAllContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_monetaryAll

	return p
}

func (s *MonetaryAllContext) GetParser() antlr.Parser { return s.parser }

func (s *MonetaryAllContext) GetAsset() IExpressionContext { return s.asset }

func (s *MonetaryAllContext) SetAsset(v IExpressionContext) { s.asset = v }

func (s *MonetaryAllContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACK, 0)
}

func (s *MonetaryAllContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACK, 0)
}

func (s *MonetaryAllContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MonetaryAllContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryAllContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MonetaryAllContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterMonetaryAll(s)
	}
}

func (s *MonetaryAllContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitMonetaryAll(s)
	}
}

func (p *NumScriptParser) MonetaryAll() (localctx IMonetaryAllContext) {
	this := p
	_ = this

	localctx = NewMonetaryAllContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, NumScriptParserRULE_monetaryAll)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		p.Match(NumScriptParserLBRACK)
	}
	{
		p.SetState(56)

		var _x = p.Expression()

		localctx.(*MonetaryAllContext).asset = _x
	}
	{
		p.SetState(57)
		p.Match(NumScriptParserT__0)
	}
	{
		p.SetState(58)
		p.Match(NumScriptParserRBRACK)
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyFrom(ctx *LiteralContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LitPortionContext struct {
	*LiteralContext
}

func NewLitPortionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitPortionContext {
	var p = new(LitPortionContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitPortionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitPortionContext) PORTION() antlr.TerminalNode {
	return s.GetToken(NumScriptParserPORTION, 0)
}

func (s *LitPortionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitPortion(s)
	}
}

func (s *LitPortionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitPortion(s)
	}
}

type LitStringContext struct {
	*LiteralContext
}

func NewLitStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitStringContext {
	var p = new(LitStringContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitStringContext) STRING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSTRING, 0)
}

func (s *LitStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitString(s)
	}
}

func (s *LitStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitString(s)
	}
}

type LitAccountContext struct {
	*LiteralContext
}

func NewLitAccountContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitAccountContext {
	var p = new(LitAccountContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitAccountContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitAccountContext) ACCOUNT() antlr.TerminalNode {
	return s.GetToken(NumScriptParserACCOUNT, 0)
}

func (s *LitAccountContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitAccount(s)
	}
}

func (s *LitAccountContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitAccount(s)
	}
}

type LitAssetContext struct {
	*LiteralContext
}

func NewLitAssetContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitAssetContext {
	var p = new(LitAssetContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitAssetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitAssetContext) ASSET() antlr.TerminalNode {
	return s.GetToken(NumScriptParserASSET, 0)
}

func (s *LitAssetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitAsset(s)
	}
}

func (s *LitAssetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitAsset(s)
	}
}

type LitMonetaryContext struct {
	*LiteralContext
}

func NewLitMonetaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitMonetaryContext {
	var p = new(LitMonetaryContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitMonetaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitMonetaryContext) Monetary() IMonetaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMonetaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMonetaryContext)
}

func (s *LitMonetaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitMonetary(s)
	}
}

func (s *LitMonetaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitMonetary(s)
	}
}

type LitNumberContext struct {
	*LiteralContext
}

func NewLitNumberContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LitNumberContext {
	var p = new(LitNumberContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *LitNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LitNumberContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(NumScriptParserNUMBER, 0)
}

func (s *LitNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterLitNumber(s)
	}
}

func (s *LitNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitLitNumber(s)
	}
}

func (p *NumScriptParser) Literal() (localctx ILiteralContext) {
	this := p
	_ = this

	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, NumScriptParserRULE_literal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(66)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserACCOUNT:
		localctx = NewLitAccountContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(60)
			p.Match(NumScriptParserACCOUNT)
		}

	case NumScriptParserASSET:
		localctx = NewLitAssetContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(61)
			p.Match(NumScriptParserASSET)
		}

	case NumScriptParserNUMBER:
		localctx = NewLitNumberContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(62)
			p.Match(NumScriptParserNUMBER)
		}

	case NumScriptParserSTRING:
		localctx = NewLitStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(63)
			p.Match(NumScriptParserSTRING)
		}

	case NumScriptParserPORTION:
		localctx = NewLitPortionContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(64)
			p.Match(NumScriptParserPORTION)
		}

	case NumScriptParserLBRACK:
		localctx = NewLitMonetaryContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(65)
			p.Monetary()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_variable
	return p
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) VARIABLE_NAME() antlr.TerminalNode {
	return s.GetToken(NumScriptParserVARIABLE_NAME, 0)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (p *NumScriptParser) Variable() (localctx IVariableContext) {
	this := p
	_ = this

	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, NumScriptParserRULE_variable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(68)
		p.Match(NumScriptParserVARIABLE_NAME)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ExprLiteralContext struct {
	*ExpressionContext
	lit ILiteralContext
}

func NewExprLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprLiteralContext {
	var p = new(ExprLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ExprLiteralContext) GetLit() ILiteralContext { return s.lit }

func (s *ExprLiteralContext) SetLit(v ILiteralContext) { s.lit = v }

func (s *ExprLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprLiteralContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ExprLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterExprLiteral(s)
	}
}

func (s *ExprLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitExprLiteral(s)
	}
}

type ExprVariableContext struct {
	*ExpressionContext
	var_ IVariableContext
}

func NewExprVariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprVariableContext {
	var p = new(ExprVariableContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *ExprVariableContext) GetVar_() IVariableContext { return s.var_ }

func (s *ExprVariableContext) SetVar_(v IVariableContext) { s.var_ = v }

func (s *ExprVariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprVariableContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *ExprVariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterExprVariable(s)
	}
}

func (s *ExprVariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitExprVariable(s)
	}
}

func (p *NumScriptParser) Expression() (localctx IExpressionContext) {
	this := p
	_ = this

	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, NumScriptParserRULE_expression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(72)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserLBRACK, NumScriptParserSTRING, NumScriptParserPORTION, NumScriptParserNUMBER, NumScriptParserACCOUNT, NumScriptParserASSET:
		localctx = NewExprLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(70)

			var _x = p.Literal()

			localctx.(*ExprLiteralContext).lit = _x
		}

	case NumScriptParserVARIABLE_NAME:
		localctx = NewExprVariableContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(71)

			var _x = p.Variable()

			localctx.(*ExprVariableContext).var_ = _x
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAllotmentPortionContext is an interface to support dynamic dispatch.
type IAllotmentPortionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAllotmentPortionContext differentiates from other interfaces.
	IsAllotmentPortionContext()
}

type AllotmentPortionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAllotmentPortionContext() *AllotmentPortionContext {
	var p = new(AllotmentPortionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_allotmentPortion
	return p
}

func (*AllotmentPortionContext) IsAllotmentPortionContext() {}

func NewAllotmentPortionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AllotmentPortionContext {
	var p = new(AllotmentPortionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_allotmentPortion

	return p
}

func (s *AllotmentPortionContext) GetParser() antlr.Parser { return s.parser }

func (s *AllotmentPortionContext) CopyFrom(ctx *AllotmentPortionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *AllotmentPortionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllotmentPortionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AllotmentPortionRemainingContext struct {
	*AllotmentPortionContext
}

func NewAllotmentPortionRemainingContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AllotmentPortionRemainingContext {
	var p = new(AllotmentPortionRemainingContext)

	p.AllotmentPortionContext = NewEmptyAllotmentPortionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AllotmentPortionContext))

	return p
}

func (s *AllotmentPortionRemainingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllotmentPortionRemainingContext) REMAINING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserREMAINING, 0)
}

func (s *AllotmentPortionRemainingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterAllotmentPortionRemaining(s)
	}
}

func (s *AllotmentPortionRemainingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitAllotmentPortionRemaining(s)
	}
}

type AllotmentPortionVarContext struct {
	*AllotmentPortionContext
	por IVariableContext
}

func NewAllotmentPortionVarContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AllotmentPortionVarContext {
	var p = new(AllotmentPortionVarContext)

	p.AllotmentPortionContext = NewEmptyAllotmentPortionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AllotmentPortionContext))

	return p
}

func (s *AllotmentPortionVarContext) GetPor() IVariableContext { return s.por }

func (s *AllotmentPortionVarContext) SetPor(v IVariableContext) { s.por = v }

func (s *AllotmentPortionVarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllotmentPortionVarContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *AllotmentPortionVarContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterAllotmentPortionVar(s)
	}
}

func (s *AllotmentPortionVarContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitAllotmentPortionVar(s)
	}
}

type AllotmentPortionConstContext struct {
	*AllotmentPortionContext
}

func NewAllotmentPortionConstContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AllotmentPortionConstContext {
	var p = new(AllotmentPortionConstContext)

	p.AllotmentPortionContext = NewEmptyAllotmentPortionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AllotmentPortionContext))

	return p
}

func (s *AllotmentPortionConstContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllotmentPortionConstContext) PORTION() antlr.TerminalNode {
	return s.GetToken(NumScriptParserPORTION, 0)
}

func (s *AllotmentPortionConstContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterAllotmentPortionConst(s)
	}
}

func (s *AllotmentPortionConstContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitAllotmentPortionConst(s)
	}
}

func (p *NumScriptParser) AllotmentPortion() (localctx IAllotmentPortionContext) {
	this := p
	_ = this

	localctx = NewAllotmentPortionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, NumScriptParserRULE_allotmentPortion)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(77)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserPORTION:
		localctx = NewAllotmentPortionConstContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(74)
			p.Match(NumScriptParserPORTION)
		}

	case NumScriptParserVARIABLE_NAME:
		localctx = NewAllotmentPortionVarContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(75)

			var _x = p.Variable()

			localctx.(*AllotmentPortionVarContext).por = _x
		}

	case NumScriptParserREMAINING:
		localctx = NewAllotmentPortionRemainingContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(76)
			p.Match(NumScriptParserREMAINING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDestinationInOrderContext is an interface to support dynamic dispatch.
type IDestinationInOrderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_expression returns the _expression rule contexts.
	Get_expression() IExpressionContext

	// Get_keptOrDestination returns the _keptOrDestination rule contexts.
	Get_keptOrDestination() IKeptOrDestinationContext

	// GetRemainingDest returns the remainingDest rule contexts.
	GetRemainingDest() IKeptOrDestinationContext

	// Set_expression sets the _expression rule contexts.
	Set_expression(IExpressionContext)

	// Set_keptOrDestination sets the _keptOrDestination rule contexts.
	Set_keptOrDestination(IKeptOrDestinationContext)

	// SetRemainingDest sets the remainingDest rule contexts.
	SetRemainingDest(IKeptOrDestinationContext)

	// GetAmounts returns the amounts rule context list.
	GetAmounts() []IExpressionContext

	// GetDests returns the dests rule context list.
	GetDests() []IKeptOrDestinationContext

	// SetAmounts sets the amounts rule context list.
	SetAmounts([]IExpressionContext)

	// SetDests sets the dests rule context list.
	SetDests([]IKeptOrDestinationContext)

	// IsDestinationInOrderContext differentiates from other interfaces.
	IsDestinationInOrderContext()
}

type DestinationInOrderContext struct {
	*antlr.BaseParserRuleContext
	parser             antlr.Parser
	_expression        IExpressionContext
	amounts            []IExpressionContext
	_keptOrDestination IKeptOrDestinationContext
	dests              []IKeptOrDestinationContext
	remainingDest      IKeptOrDestinationContext
}

func NewEmptyDestinationInOrderContext() *DestinationInOrderContext {
	var p = new(DestinationInOrderContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_destinationInOrder
	return p
}

func (*DestinationInOrderContext) IsDestinationInOrderContext() {}

func NewDestinationInOrderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DestinationInOrderContext {
	var p = new(DestinationInOrderContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_destinationInOrder

	return p
}

func (s *DestinationInOrderContext) GetParser() antlr.Parser { return s.parser }

func (s *DestinationInOrderContext) Get_expression() IExpressionContext { return s._expression }

func (s *DestinationInOrderContext) Get_keptOrDestination() IKeptOrDestinationContext {
	return s._keptOrDestination
}

func (s *DestinationInOrderContext) GetRemainingDest() IKeptOrDestinationContext {
	return s.remainingDest
}

func (s *DestinationInOrderContext) Set_expression(v IExpressionContext) { s._expression = v }

func (s *DestinationInOrderContext) Set_keptOrDestination(v IKeptOrDestinationContext) {
	s._keptOrDestination = v
}

func (s *DestinationInOrderContext) SetRemainingDest(v IKeptOrDestinationContext) {
	s.remainingDest = v
}

func (s *DestinationInOrderContext) GetAmounts() []IExpressionContext { return s.amounts }

func (s *DestinationInOrderContext) GetDests() []IKeptOrDestinationContext { return s.dests }

func (s *DestinationInOrderContext) SetAmounts(v []IExpressionContext) { s.amounts = v }

func (s *DestinationInOrderContext) SetDests(v []IKeptOrDestinationContext) { s.dests = v }

func (s *DestinationInOrderContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACE, 0)
}

func (s *DestinationInOrderContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *DestinationInOrderContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *DestinationInOrderContext) REMAINING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserREMAINING, 0)
}

func (s *DestinationInOrderContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACE, 0)
}

func (s *DestinationInOrderContext) AllKeptOrDestination() []IKeptOrDestinationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IKeptOrDestinationContext); ok {
			len++
		}
	}

	tst := make([]IKeptOrDestinationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IKeptOrDestinationContext); ok {
			tst[i] = t.(IKeptOrDestinationContext)
			i++
		}
	}

	return tst
}

func (s *DestinationInOrderContext) KeptOrDestination(i int) IKeptOrDestinationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeptOrDestinationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeptOrDestinationContext)
}

func (s *DestinationInOrderContext) AllMAX() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserMAX)
}

func (s *DestinationInOrderContext) MAX(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserMAX, i)
}

func (s *DestinationInOrderContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *DestinationInOrderContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *DestinationInOrderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestinationInOrderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DestinationInOrderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterDestinationInOrder(s)
	}
}

func (s *DestinationInOrderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitDestinationInOrder(s)
	}
}

func (p *NumScriptParser) DestinationInOrder() (localctx IDestinationInOrderContext) {
	this := p
	_ = this

	localctx = NewDestinationInOrderContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, NumScriptParserRULE_destinationInOrder)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.Match(NumScriptParserLBRACE)
	}
	{
		p.SetState(80)
		p.Match(NumScriptParserNEWLINE)
	}
	p.SetState(86)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == NumScriptParserMAX {
		{
			p.SetState(81)
			p.Match(NumScriptParserMAX)
		}
		{
			p.SetState(82)

			var _x = p.Expression()

			localctx.(*DestinationInOrderContext)._expression = _x
		}
		localctx.(*DestinationInOrderContext).amounts = append(localctx.(*DestinationInOrderContext).amounts, localctx.(*DestinationInOrderContext)._expression)
		{
			p.SetState(83)

			var _x = p.KeptOrDestination()

			localctx.(*DestinationInOrderContext)._keptOrDestination = _x
		}
		localctx.(*DestinationInOrderContext).dests = append(localctx.(*DestinationInOrderContext).dests, localctx.(*DestinationInOrderContext)._keptOrDestination)
		{
			p.SetState(84)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(88)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(90)
		p.Match(NumScriptParserREMAINING)
	}
	{
		p.SetState(91)

		var _x = p.KeptOrDestination()

		localctx.(*DestinationInOrderContext).remainingDest = _x
	}
	{
		p.SetState(92)
		p.Match(NumScriptParserNEWLINE)
	}
	{
		p.SetState(93)
		p.Match(NumScriptParserRBRACE)
	}

	return localctx
}

// IDestinationAllotmentContext is an interface to support dynamic dispatch.
type IDestinationAllotmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_allotmentPortion returns the _allotmentPortion rule contexts.
	Get_allotmentPortion() IAllotmentPortionContext

	// Get_keptOrDestination returns the _keptOrDestination rule contexts.
	Get_keptOrDestination() IKeptOrDestinationContext

	// Set_allotmentPortion sets the _allotmentPortion rule contexts.
	Set_allotmentPortion(IAllotmentPortionContext)

	// Set_keptOrDestination sets the _keptOrDestination rule contexts.
	Set_keptOrDestination(IKeptOrDestinationContext)

	// GetPortions returns the portions rule context list.
	GetPortions() []IAllotmentPortionContext

	// GetDests returns the dests rule context list.
	GetDests() []IKeptOrDestinationContext

	// SetPortions sets the portions rule context list.
	SetPortions([]IAllotmentPortionContext)

	// SetDests sets the dests rule context list.
	SetDests([]IKeptOrDestinationContext)

	// IsDestinationAllotmentContext differentiates from other interfaces.
	IsDestinationAllotmentContext()
}

type DestinationAllotmentContext struct {
	*antlr.BaseParserRuleContext
	parser             antlr.Parser
	_allotmentPortion  IAllotmentPortionContext
	portions           []IAllotmentPortionContext
	_keptOrDestination IKeptOrDestinationContext
	dests              []IKeptOrDestinationContext
}

func NewEmptyDestinationAllotmentContext() *DestinationAllotmentContext {
	var p = new(DestinationAllotmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_destinationAllotment
	return p
}

func (*DestinationAllotmentContext) IsDestinationAllotmentContext() {}

func NewDestinationAllotmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DestinationAllotmentContext {
	var p = new(DestinationAllotmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_destinationAllotment

	return p
}

func (s *DestinationAllotmentContext) GetParser() antlr.Parser { return s.parser }

func (s *DestinationAllotmentContext) Get_allotmentPortion() IAllotmentPortionContext {
	return s._allotmentPortion
}

func (s *DestinationAllotmentContext) Get_keptOrDestination() IKeptOrDestinationContext {
	return s._keptOrDestination
}

func (s *DestinationAllotmentContext) Set_allotmentPortion(v IAllotmentPortionContext) {
	s._allotmentPortion = v
}

func (s *DestinationAllotmentContext) Set_keptOrDestination(v IKeptOrDestinationContext) {
	s._keptOrDestination = v
}

func (s *DestinationAllotmentContext) GetPortions() []IAllotmentPortionContext { return s.portions }

func (s *DestinationAllotmentContext) GetDests() []IKeptOrDestinationContext { return s.dests }

func (s *DestinationAllotmentContext) SetPortions(v []IAllotmentPortionContext) { s.portions = v }

func (s *DestinationAllotmentContext) SetDests(v []IKeptOrDestinationContext) { s.dests = v }

func (s *DestinationAllotmentContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACE, 0)
}

func (s *DestinationAllotmentContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *DestinationAllotmentContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *DestinationAllotmentContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACE, 0)
}

func (s *DestinationAllotmentContext) AllAllotmentPortion() []IAllotmentPortionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAllotmentPortionContext); ok {
			len++
		}
	}

	tst := make([]IAllotmentPortionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAllotmentPortionContext); ok {
			tst[i] = t.(IAllotmentPortionContext)
			i++
		}
	}

	return tst
}

func (s *DestinationAllotmentContext) AllotmentPortion(i int) IAllotmentPortionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAllotmentPortionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAllotmentPortionContext)
}

func (s *DestinationAllotmentContext) AllKeptOrDestination() []IKeptOrDestinationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IKeptOrDestinationContext); ok {
			len++
		}
	}

	tst := make([]IKeptOrDestinationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IKeptOrDestinationContext); ok {
			tst[i] = t.(IKeptOrDestinationContext)
			i++
		}
	}

	return tst
}

func (s *DestinationAllotmentContext) KeptOrDestination(i int) IKeptOrDestinationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeptOrDestinationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeptOrDestinationContext)
}

func (s *DestinationAllotmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestinationAllotmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DestinationAllotmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterDestinationAllotment(s)
	}
}

func (s *DestinationAllotmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitDestinationAllotment(s)
	}
}

func (p *NumScriptParser) DestinationAllotment() (localctx IDestinationAllotmentContext) {
	this := p
	_ = this

	localctx = NewDestinationAllotmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, NumScriptParserRULE_destinationAllotment)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		p.Match(NumScriptParserLBRACE)
	}
	{
		p.SetState(96)
		p.Match(NumScriptParserNEWLINE)
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-37)&-(0x1f+1)) == 0 && ((1<<uint((_la-37)))&((1<<(NumScriptParserPORTION-37))|(1<<(NumScriptParserREMAINING-37))|(1<<(NumScriptParserVARIABLE_NAME-37)))) != 0) {
		{
			p.SetState(97)

			var _x = p.AllotmentPortion()

			localctx.(*DestinationAllotmentContext)._allotmentPortion = _x
		}
		localctx.(*DestinationAllotmentContext).portions = append(localctx.(*DestinationAllotmentContext).portions, localctx.(*DestinationAllotmentContext)._allotmentPortion)
		{
			p.SetState(98)

			var _x = p.KeptOrDestination()

			localctx.(*DestinationAllotmentContext)._keptOrDestination = _x
		}
		localctx.(*DestinationAllotmentContext).dests = append(localctx.(*DestinationAllotmentContext).dests, localctx.(*DestinationAllotmentContext)._keptOrDestination)
		{
			p.SetState(99)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(103)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(105)
		p.Match(NumScriptParserRBRACE)
	}

	return localctx
}

// IKeptOrDestinationContext is an interface to support dynamic dispatch.
type IKeptOrDestinationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsKeptOrDestinationContext differentiates from other interfaces.
	IsKeptOrDestinationContext()
}

type KeptOrDestinationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKeptOrDestinationContext() *KeptOrDestinationContext {
	var p = new(KeptOrDestinationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_keptOrDestination
	return p
}

func (*KeptOrDestinationContext) IsKeptOrDestinationContext() {}

func NewKeptOrDestinationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KeptOrDestinationContext {
	var p = new(KeptOrDestinationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_keptOrDestination

	return p
}

func (s *KeptOrDestinationContext) GetParser() antlr.Parser { return s.parser }

func (s *KeptOrDestinationContext) CopyFrom(ctx *KeptOrDestinationContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *KeptOrDestinationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KeptOrDestinationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type IsKeptContext struct {
	*KeptOrDestinationContext
}

func NewIsKeptContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IsKeptContext {
	var p = new(IsKeptContext)

	p.KeptOrDestinationContext = NewEmptyKeptOrDestinationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*KeptOrDestinationContext))

	return p
}

func (s *IsKeptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IsKeptContext) KEPT() antlr.TerminalNode {
	return s.GetToken(NumScriptParserKEPT, 0)
}

func (s *IsKeptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterIsKept(s)
	}
}

func (s *IsKeptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitIsKept(s)
	}
}

type IsDestinationContext struct {
	*KeptOrDestinationContext
}

func NewIsDestinationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IsDestinationContext {
	var p = new(IsDestinationContext)

	p.KeptOrDestinationContext = NewEmptyKeptOrDestinationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*KeptOrDestinationContext))

	return p
}

func (s *IsDestinationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IsDestinationContext) TO() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTO, 0)
}

func (s *IsDestinationContext) Destination() IDestinationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDestinationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDestinationContext)
}

func (s *IsDestinationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterIsDestination(s)
	}
}

func (s *IsDestinationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitIsDestination(s)
	}
}

func (p *NumScriptParser) KeptOrDestination() (localctx IKeptOrDestinationContext) {
	this := p
	_ = this

	localctx = NewKeptOrDestinationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, NumScriptParserRULE_keptOrDestination)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(110)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserTO:
		localctx = NewIsDestinationContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(107)
			p.Match(NumScriptParserTO)
		}
		{
			p.SetState(108)
			p.Destination()
		}

	case NumScriptParserKEPT:
		localctx = NewIsKeptContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(109)
			p.Match(NumScriptParserKEPT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDestinationContext is an interface to support dynamic dispatch.
type IDestinationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDestinationContext differentiates from other interfaces.
	IsDestinationContext()
}

type DestinationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDestinationContext() *DestinationContext {
	var p = new(DestinationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_destination
	return p
}

func (*DestinationContext) IsDestinationContext() {}

func NewDestinationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DestinationContext {
	var p = new(DestinationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_destination

	return p
}

func (s *DestinationContext) GetParser() antlr.Parser { return s.parser }

func (s *DestinationContext) CopyFrom(ctx *DestinationContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *DestinationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestinationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DestAccountContext struct {
	*DestinationContext
}

func NewDestAccountContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DestAccountContext {
	var p = new(DestAccountContext)

	p.DestinationContext = NewEmptyDestinationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DestinationContext))

	return p
}

func (s *DestAccountContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestAccountContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *DestAccountContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterDestAccount(s)
	}
}

func (s *DestAccountContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitDestAccount(s)
	}
}

type DestAllotmentContext struct {
	*DestinationContext
}

func NewDestAllotmentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DestAllotmentContext {
	var p = new(DestAllotmentContext)

	p.DestinationContext = NewEmptyDestinationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DestinationContext))

	return p
}

func (s *DestAllotmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestAllotmentContext) DestinationAllotment() IDestinationAllotmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDestinationAllotmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDestinationAllotmentContext)
}

func (s *DestAllotmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterDestAllotment(s)
	}
}

func (s *DestAllotmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitDestAllotment(s)
	}
}

type DestInOrderContext struct {
	*DestinationContext
}

func NewDestInOrderContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DestInOrderContext {
	var p = new(DestInOrderContext)

	p.DestinationContext = NewEmptyDestinationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DestinationContext))

	return p
}

func (s *DestInOrderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DestInOrderContext) DestinationInOrder() IDestinationInOrderContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDestinationInOrderContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDestinationInOrderContext)
}

func (s *DestInOrderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterDestInOrder(s)
	}
}

func (s *DestInOrderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitDestInOrder(s)
	}
}

func (p *NumScriptParser) Destination() (localctx IDestinationContext) {
	this := p
	_ = this

	localctx = NewDestinationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, NumScriptParserRULE_destination)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(115)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDestAccountContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(112)
			p.Expression()
		}

	case 2:
		localctx = NewDestInOrderContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(113)
			p.DestinationInOrder()
		}

	case 3:
		localctx = NewDestAllotmentContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(114)
			p.DestinationAllotment()
		}

	}

	return localctx
}

// ISourceAccountOverdraftContext is an interface to support dynamic dispatch.
type ISourceAccountOverdraftContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSourceAccountOverdraftContext differentiates from other interfaces.
	IsSourceAccountOverdraftContext()
}

type SourceAccountOverdraftContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySourceAccountOverdraftContext() *SourceAccountOverdraftContext {
	var p = new(SourceAccountOverdraftContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_sourceAccountOverdraft
	return p
}

func (*SourceAccountOverdraftContext) IsSourceAccountOverdraftContext() {}

func NewSourceAccountOverdraftContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceAccountOverdraftContext {
	var p = new(SourceAccountOverdraftContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_sourceAccountOverdraft

	return p
}

func (s *SourceAccountOverdraftContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceAccountOverdraftContext) CopyFrom(ctx *SourceAccountOverdraftContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SourceAccountOverdraftContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceAccountOverdraftContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SrcAccountOverdraftSpecificContext struct {
	*SourceAccountOverdraftContext
	specific IExpressionContext
}

func NewSrcAccountOverdraftSpecificContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcAccountOverdraftSpecificContext {
	var p = new(SrcAccountOverdraftSpecificContext)

	p.SourceAccountOverdraftContext = NewEmptySourceAccountOverdraftContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SourceAccountOverdraftContext))

	return p
}

func (s *SrcAccountOverdraftSpecificContext) GetSpecific() IExpressionContext { return s.specific }

func (s *SrcAccountOverdraftSpecificContext) SetSpecific(v IExpressionContext) { s.specific = v }

func (s *SrcAccountOverdraftSpecificContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcAccountOverdraftSpecificContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SrcAccountOverdraftSpecificContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcAccountOverdraftSpecific(s)
	}
}

func (s *SrcAccountOverdraftSpecificContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcAccountOverdraftSpecific(s)
	}
}

type SrcAccountOverdraftUnboundedContext struct {
	*SourceAccountOverdraftContext
}

func NewSrcAccountOverdraftUnboundedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcAccountOverdraftUnboundedContext {
	var p = new(SrcAccountOverdraftUnboundedContext)

	p.SourceAccountOverdraftContext = NewEmptySourceAccountOverdraftContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SourceAccountOverdraftContext))

	return p
}

func (s *SrcAccountOverdraftUnboundedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcAccountOverdraftUnboundedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcAccountOverdraftUnbounded(s)
	}
}

func (s *SrcAccountOverdraftUnboundedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcAccountOverdraftUnbounded(s)
	}
}

func (p *NumScriptParser) SourceAccountOverdraft() (localctx ISourceAccountOverdraftContext) {
	this := p
	_ = this

	localctx = NewSourceAccountOverdraftContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, NumScriptParserRULE_sourceAccountOverdraft)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(120)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserT__1:
		localctx = NewSrcAccountOverdraftSpecificContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(117)
			p.Match(NumScriptParserT__1)
		}
		{
			p.SetState(118)

			var _x = p.Expression()

			localctx.(*SrcAccountOverdraftSpecificContext).specific = _x
		}

	case NumScriptParserT__2:
		localctx = NewSrcAccountOverdraftUnboundedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(119)
			p.Match(NumScriptParserT__2)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISourceAccountContext is an interface to support dynamic dispatch.
type ISourceAccountContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAccount returns the account rule contexts.
	GetAccount() IExpressionContext

	// GetOverdraft returns the overdraft rule contexts.
	GetOverdraft() ISourceAccountOverdraftContext

	// SetAccount sets the account rule contexts.
	SetAccount(IExpressionContext)

	// SetOverdraft sets the overdraft rule contexts.
	SetOverdraft(ISourceAccountOverdraftContext)

	// IsSourceAccountContext differentiates from other interfaces.
	IsSourceAccountContext()
}

type SourceAccountContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	account   IExpressionContext
	overdraft ISourceAccountOverdraftContext
}

func NewEmptySourceAccountContext() *SourceAccountContext {
	var p = new(SourceAccountContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_sourceAccount
	return p
}

func (*SourceAccountContext) IsSourceAccountContext() {}

func NewSourceAccountContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceAccountContext {
	var p = new(SourceAccountContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_sourceAccount

	return p
}

func (s *SourceAccountContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceAccountContext) GetAccount() IExpressionContext { return s.account }

func (s *SourceAccountContext) GetOverdraft() ISourceAccountOverdraftContext { return s.overdraft }

func (s *SourceAccountContext) SetAccount(v IExpressionContext) { s.account = v }

func (s *SourceAccountContext) SetOverdraft(v ISourceAccountOverdraftContext) { s.overdraft = v }

func (s *SourceAccountContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SourceAccountContext) SourceAccountOverdraft() ISourceAccountOverdraftContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceAccountOverdraftContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceAccountOverdraftContext)
}

func (s *SourceAccountContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceAccountContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SourceAccountContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSourceAccount(s)
	}
}

func (s *SourceAccountContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSourceAccount(s)
	}
}

func (p *NumScriptParser) SourceAccount() (localctx ISourceAccountContext) {
	this := p
	_ = this

	localctx = NewSourceAccountContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, NumScriptParserRULE_sourceAccount)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(122)

		var _x = p.Expression()

		localctx.(*SourceAccountContext).account = _x
	}
	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == NumScriptParserT__1 || _la == NumScriptParserT__2 {
		{
			p.SetState(123)

			var _x = p.SourceAccountOverdraft()

			localctx.(*SourceAccountContext).overdraft = _x
		}

	}

	return localctx
}

// ISourceInOrderContext is an interface to support dynamic dispatch.
type ISourceInOrderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_source returns the _source rule contexts.
	Get_source() ISourceContext

	// Set_source sets the _source rule contexts.
	Set_source(ISourceContext)

	// GetSources returns the sources rule context list.
	GetSources() []ISourceContext

	// SetSources sets the sources rule context list.
	SetSources([]ISourceContext)

	// IsSourceInOrderContext differentiates from other interfaces.
	IsSourceInOrderContext()
}

type SourceInOrderContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	_source ISourceContext
	sources []ISourceContext
}

func NewEmptySourceInOrderContext() *SourceInOrderContext {
	var p = new(SourceInOrderContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_sourceInOrder
	return p
}

func (*SourceInOrderContext) IsSourceInOrderContext() {}

func NewSourceInOrderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceInOrderContext {
	var p = new(SourceInOrderContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_sourceInOrder

	return p
}

func (s *SourceInOrderContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceInOrderContext) Get_source() ISourceContext { return s._source }

func (s *SourceInOrderContext) Set_source(v ISourceContext) { s._source = v }

func (s *SourceInOrderContext) GetSources() []ISourceContext { return s.sources }

func (s *SourceInOrderContext) SetSources(v []ISourceContext) { s.sources = v }

func (s *SourceInOrderContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACE, 0)
}

func (s *SourceInOrderContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *SourceInOrderContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *SourceInOrderContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACE, 0)
}

func (s *SourceInOrderContext) AllSource() []ISourceContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISourceContext); ok {
			len++
		}
	}

	tst := make([]ISourceContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISourceContext); ok {
			tst[i] = t.(ISourceContext)
			i++
		}
	}

	return tst
}

func (s *SourceInOrderContext) Source(i int) ISourceContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *SourceInOrderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceInOrderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SourceInOrderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSourceInOrder(s)
	}
}

func (s *SourceInOrderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSourceInOrder(s)
	}
}

func (p *NumScriptParser) SourceInOrder() (localctx ISourceInOrderContext) {
	this := p
	_ = this

	localctx = NewSourceInOrderContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, NumScriptParserRULE_sourceInOrder)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(126)
		p.Match(NumScriptParserLBRACE)
	}
	{
		p.SetState(127)
		p.Match(NumScriptParserNEWLINE)
	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-17)&-(0x1f+1)) == 0 && ((1<<uint((_la-17)))&((1<<(NumScriptParserMAX-17))|(1<<(NumScriptParserLBRACK-17))|(1<<(NumScriptParserLBRACE-17))|(1<<(NumScriptParserSTRING-17))|(1<<(NumScriptParserPORTION-17))|(1<<(NumScriptParserNUMBER-17))|(1<<(NumScriptParserVARIABLE_NAME-17))|(1<<(NumScriptParserACCOUNT-17))|(1<<(NumScriptParserASSET-17)))) != 0) {
		{
			p.SetState(128)

			var _x = p.Source()

			localctx.(*SourceInOrderContext)._source = _x
		}
		localctx.(*SourceInOrderContext).sources = append(localctx.(*SourceInOrderContext).sources, localctx.(*SourceInOrderContext)._source)
		{
			p.SetState(129)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(135)
		p.Match(NumScriptParserRBRACE)
	}

	return localctx
}

// ISourceMaxedContext is an interface to support dynamic dispatch.
type ISourceMaxedContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetMax returns the max rule contexts.
	GetMax() IExpressionContext

	// GetSrc returns the src rule contexts.
	GetSrc() ISourceContext

	// SetMax sets the max rule contexts.
	SetMax(IExpressionContext)

	// SetSrc sets the src rule contexts.
	SetSrc(ISourceContext)

	// IsSourceMaxedContext differentiates from other interfaces.
	IsSourceMaxedContext()
}

type SourceMaxedContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	max    IExpressionContext
	src    ISourceContext
}

func NewEmptySourceMaxedContext() *SourceMaxedContext {
	var p = new(SourceMaxedContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_sourceMaxed
	return p
}

func (*SourceMaxedContext) IsSourceMaxedContext() {}

func NewSourceMaxedContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceMaxedContext {
	var p = new(SourceMaxedContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_sourceMaxed

	return p
}

func (s *SourceMaxedContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceMaxedContext) GetMax() IExpressionContext { return s.max }

func (s *SourceMaxedContext) GetSrc() ISourceContext { return s.src }

func (s *SourceMaxedContext) SetMax(v IExpressionContext) { s.max = v }

func (s *SourceMaxedContext) SetSrc(v ISourceContext) { s.src = v }

func (s *SourceMaxedContext) MAX() antlr.TerminalNode {
	return s.GetToken(NumScriptParserMAX, 0)
}

func (s *SourceMaxedContext) FROM() antlr.TerminalNode {
	return s.GetToken(NumScriptParserFROM, 0)
}

func (s *SourceMaxedContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SourceMaxedContext) Source() ISourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *SourceMaxedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceMaxedContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SourceMaxedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSourceMaxed(s)
	}
}

func (s *SourceMaxedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSourceMaxed(s)
	}
}

func (p *NumScriptParser) SourceMaxed() (localctx ISourceMaxedContext) {
	this := p
	_ = this

	localctx = NewSourceMaxedContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, NumScriptParserRULE_sourceMaxed)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(137)
		p.Match(NumScriptParserMAX)
	}
	{
		p.SetState(138)

		var _x = p.Expression()

		localctx.(*SourceMaxedContext).max = _x
	}
	{
		p.SetState(139)
		p.Match(NumScriptParserFROM)
	}
	{
		p.SetState(140)

		var _x = p.Source()

		localctx.(*SourceMaxedContext).src = _x
	}

	return localctx
}

// ISourceContext is an interface to support dynamic dispatch.
type ISourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSourceContext differentiates from other interfaces.
	IsSourceContext()
}

type SourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySourceContext() *SourceContext {
	var p = new(SourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_source
	return p
}

func (*SourceContext) IsSourceContext() {}

func NewSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceContext {
	var p = new(SourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_source

	return p
}

func (s *SourceContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceContext) CopyFrom(ctx *SourceContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SrcAccountContext struct {
	*SourceContext
}

func NewSrcAccountContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcAccountContext {
	var p = new(SrcAccountContext)

	p.SourceContext = NewEmptySourceContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SourceContext))

	return p
}

func (s *SrcAccountContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcAccountContext) SourceAccount() ISourceAccountContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceAccountContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceAccountContext)
}

func (s *SrcAccountContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcAccount(s)
	}
}

func (s *SrcAccountContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcAccount(s)
	}
}

type SrcMaxedContext struct {
	*SourceContext
}

func NewSrcMaxedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcMaxedContext {
	var p = new(SrcMaxedContext)

	p.SourceContext = NewEmptySourceContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SourceContext))

	return p
}

func (s *SrcMaxedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcMaxedContext) SourceMaxed() ISourceMaxedContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceMaxedContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceMaxedContext)
}

func (s *SrcMaxedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcMaxed(s)
	}
}

func (s *SrcMaxedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcMaxed(s)
	}
}

type SrcInOrderContext struct {
	*SourceContext
}

func NewSrcInOrderContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcInOrderContext {
	var p = new(SrcInOrderContext)

	p.SourceContext = NewEmptySourceContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SourceContext))

	return p
}

func (s *SrcInOrderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcInOrderContext) SourceInOrder() ISourceInOrderContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceInOrderContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceInOrderContext)
}

func (s *SrcInOrderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcInOrder(s)
	}
}

func (s *SrcInOrderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcInOrder(s)
	}
}

func (p *NumScriptParser) Source() (localctx ISourceContext) {
	this := p
	_ = this

	localctx = NewSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, NumScriptParserRULE_source)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(145)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserLBRACK, NumScriptParserSTRING, NumScriptParserPORTION, NumScriptParserNUMBER, NumScriptParserVARIABLE_NAME, NumScriptParserACCOUNT, NumScriptParserASSET:
		localctx = NewSrcAccountContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(142)
			p.SourceAccount()
		}

	case NumScriptParserMAX:
		localctx = NewSrcMaxedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(143)
			p.SourceMaxed()
		}

	case NumScriptParserLBRACE:
		localctx = NewSrcInOrderContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(144)
			p.SourceInOrder()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISourceAllotmentContext is an interface to support dynamic dispatch.
type ISourceAllotmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_allotmentPortion returns the _allotmentPortion rule contexts.
	Get_allotmentPortion() IAllotmentPortionContext

	// Get_source returns the _source rule contexts.
	Get_source() ISourceContext

	// Set_allotmentPortion sets the _allotmentPortion rule contexts.
	Set_allotmentPortion(IAllotmentPortionContext)

	// Set_source sets the _source rule contexts.
	Set_source(ISourceContext)

	// GetPortions returns the portions rule context list.
	GetPortions() []IAllotmentPortionContext

	// GetSources returns the sources rule context list.
	GetSources() []ISourceContext

	// SetPortions sets the portions rule context list.
	SetPortions([]IAllotmentPortionContext)

	// SetSources sets the sources rule context list.
	SetSources([]ISourceContext)

	// IsSourceAllotmentContext differentiates from other interfaces.
	IsSourceAllotmentContext()
}

type SourceAllotmentContext struct {
	*antlr.BaseParserRuleContext
	parser            antlr.Parser
	_allotmentPortion IAllotmentPortionContext
	portions          []IAllotmentPortionContext
	_source           ISourceContext
	sources           []ISourceContext
}

func NewEmptySourceAllotmentContext() *SourceAllotmentContext {
	var p = new(SourceAllotmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_sourceAllotment
	return p
}

func (*SourceAllotmentContext) IsSourceAllotmentContext() {}

func NewSourceAllotmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceAllotmentContext {
	var p = new(SourceAllotmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_sourceAllotment

	return p
}

func (s *SourceAllotmentContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceAllotmentContext) Get_allotmentPortion() IAllotmentPortionContext {
	return s._allotmentPortion
}

func (s *SourceAllotmentContext) Get_source() ISourceContext { return s._source }

func (s *SourceAllotmentContext) Set_allotmentPortion(v IAllotmentPortionContext) {
	s._allotmentPortion = v
}

func (s *SourceAllotmentContext) Set_source(v ISourceContext) { s._source = v }

func (s *SourceAllotmentContext) GetPortions() []IAllotmentPortionContext { return s.portions }

func (s *SourceAllotmentContext) GetSources() []ISourceContext { return s.sources }

func (s *SourceAllotmentContext) SetPortions(v []IAllotmentPortionContext) { s.portions = v }

func (s *SourceAllotmentContext) SetSources(v []ISourceContext) { s.sources = v }

func (s *SourceAllotmentContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACE, 0)
}

func (s *SourceAllotmentContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *SourceAllotmentContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *SourceAllotmentContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACE, 0)
}

func (s *SourceAllotmentContext) AllFROM() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserFROM)
}

func (s *SourceAllotmentContext) FROM(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserFROM, i)
}

func (s *SourceAllotmentContext) AllAllotmentPortion() []IAllotmentPortionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAllotmentPortionContext); ok {
			len++
		}
	}

	tst := make([]IAllotmentPortionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAllotmentPortionContext); ok {
			tst[i] = t.(IAllotmentPortionContext)
			i++
		}
	}

	return tst
}

func (s *SourceAllotmentContext) AllotmentPortion(i int) IAllotmentPortionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAllotmentPortionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAllotmentPortionContext)
}

func (s *SourceAllotmentContext) AllSource() []ISourceContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISourceContext); ok {
			len++
		}
	}

	tst := make([]ISourceContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISourceContext); ok {
			tst[i] = t.(ISourceContext)
			i++
		}
	}

	return tst
}

func (s *SourceAllotmentContext) Source(i int) ISourceContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *SourceAllotmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceAllotmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SourceAllotmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSourceAllotment(s)
	}
}

func (s *SourceAllotmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSourceAllotment(s)
	}
}

func (p *NumScriptParser) SourceAllotment() (localctx ISourceAllotmentContext) {
	this := p
	_ = this

	localctx = NewSourceAllotmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, NumScriptParserRULE_sourceAllotment)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(147)
		p.Match(NumScriptParserLBRACE)
	}
	{
		p.SetState(148)
		p.Match(NumScriptParserNEWLINE)
	}
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-37)&-(0x1f+1)) == 0 && ((1<<uint((_la-37)))&((1<<(NumScriptParserPORTION-37))|(1<<(NumScriptParserREMAINING-37))|(1<<(NumScriptParserVARIABLE_NAME-37)))) != 0) {
		{
			p.SetState(149)

			var _x = p.AllotmentPortion()

			localctx.(*SourceAllotmentContext)._allotmentPortion = _x
		}
		localctx.(*SourceAllotmentContext).portions = append(localctx.(*SourceAllotmentContext).portions, localctx.(*SourceAllotmentContext)._allotmentPortion)
		{
			p.SetState(150)
			p.Match(NumScriptParserFROM)
		}
		{
			p.SetState(151)

			var _x = p.Source()

			localctx.(*SourceAllotmentContext)._source = _x
		}
		localctx.(*SourceAllotmentContext).sources = append(localctx.(*SourceAllotmentContext).sources, localctx.(*SourceAllotmentContext)._source)
		{
			p.SetState(152)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(156)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(158)
		p.Match(NumScriptParserRBRACE)
	}

	return localctx
}

// IValueAwareSourceContext is an interface to support dynamic dispatch.
type IValueAwareSourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueAwareSourceContext differentiates from other interfaces.
	IsValueAwareSourceContext()
}

type ValueAwareSourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueAwareSourceContext() *ValueAwareSourceContext {
	var p = new(ValueAwareSourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_valueAwareSource
	return p
}

func (*ValueAwareSourceContext) IsValueAwareSourceContext() {}

func NewValueAwareSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueAwareSourceContext {
	var p = new(ValueAwareSourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_valueAwareSource

	return p
}

func (s *ValueAwareSourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueAwareSourceContext) CopyFrom(ctx *ValueAwareSourceContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ValueAwareSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueAwareSourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SrcContext struct {
	*ValueAwareSourceContext
}

func NewSrcContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcContext {
	var p = new(SrcContext)

	p.ValueAwareSourceContext = NewEmptyValueAwareSourceContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ValueAwareSourceContext))

	return p
}

func (s *SrcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcContext) Source() ISourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *SrcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrc(s)
	}
}

func (s *SrcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrc(s)
	}
}

type SrcAllotmentContext struct {
	*ValueAwareSourceContext
}

func NewSrcAllotmentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SrcAllotmentContext {
	var p = new(SrcAllotmentContext)

	p.ValueAwareSourceContext = NewEmptyValueAwareSourceContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ValueAwareSourceContext))

	return p
}

func (s *SrcAllotmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SrcAllotmentContext) SourceAllotment() ISourceAllotmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceAllotmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceAllotmentContext)
}

func (s *SrcAllotmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSrcAllotment(s)
	}
}

func (s *SrcAllotmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSrcAllotment(s)
	}
}

func (p *NumScriptParser) ValueAwareSource() (localctx IValueAwareSourceContext) {
	this := p
	_ = this

	localctx = NewValueAwareSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, NumScriptParserRULE_valueAwareSource)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(162)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSrcContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(160)
			p.Source()
		}

	case 2:
		localctx = NewSrcAllotmentContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(161)
			p.SourceAllotment()
		}

	}

	return localctx
}

// IMonetaryArithmeticContext is an interface to support dynamic dispatch.
type IMonetaryArithmeticContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_OP_ADD returns the _OP_ADD token.
	Get_OP_ADD() antlr.Token

	// Get_OP_SUB returns the _OP_SUB token.
	Get_OP_SUB() antlr.Token

	// Get_tset745 returns the _tset745 token.
	Get_tset745() antlr.Token

	// Set_OP_ADD sets the _OP_ADD token.
	Set_OP_ADD(antlr.Token)

	// Set_OP_SUB sets the _OP_SUB token.
	Set_OP_SUB(antlr.Token)

	// Set_tset745 sets the _tset745 token.
	Set_tset745(antlr.Token)

	// GetOperators returns the operators token list.
	GetOperators() []antlr.Token

	// SetOperators sets the operators token list.
	SetOperators([]antlr.Token)

	// GetBaseOperand returns the baseOperand rule contexts.
	GetBaseOperand() IExpressionContext

	// Get_expression returns the _expression rule contexts.
	Get_expression() IExpressionContext

	// SetBaseOperand sets the baseOperand rule contexts.
	SetBaseOperand(IExpressionContext)

	// Set_expression sets the _expression rule contexts.
	Set_expression(IExpressionContext)

	// GetOperands returns the operands rule context list.
	GetOperands() []IExpressionContext

	// SetOperands sets the operands rule context list.
	SetOperands([]IExpressionContext)

	// IsMonetaryArithmeticContext differentiates from other interfaces.
	IsMonetaryArithmeticContext()
}

type MonetaryArithmeticContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	baseOperand IExpressionContext
	_OP_ADD     antlr.Token
	operators   []antlr.Token
	_OP_SUB     antlr.Token
	_tset745    antlr.Token
	_expression IExpressionContext
	operands    []IExpressionContext
}

func NewEmptyMonetaryArithmeticContext() *MonetaryArithmeticContext {
	var p = new(MonetaryArithmeticContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_monetaryArithmetic
	return p
}

func (*MonetaryArithmeticContext) IsMonetaryArithmeticContext() {}

func NewMonetaryArithmeticContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MonetaryArithmeticContext {
	var p = new(MonetaryArithmeticContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_monetaryArithmetic

	return p
}

func (s *MonetaryArithmeticContext) GetParser() antlr.Parser { return s.parser }

func (s *MonetaryArithmeticContext) Get_OP_ADD() antlr.Token { return s._OP_ADD }

func (s *MonetaryArithmeticContext) Get_OP_SUB() antlr.Token { return s._OP_SUB }

func (s *MonetaryArithmeticContext) Get_tset745() antlr.Token { return s._tset745 }

func (s *MonetaryArithmeticContext) Set_OP_ADD(v antlr.Token) { s._OP_ADD = v }

func (s *MonetaryArithmeticContext) Set_OP_SUB(v antlr.Token) { s._OP_SUB = v }

func (s *MonetaryArithmeticContext) Set_tset745(v antlr.Token) { s._tset745 = v }

func (s *MonetaryArithmeticContext) GetOperators() []antlr.Token { return s.operators }

func (s *MonetaryArithmeticContext) SetOperators(v []antlr.Token) { s.operators = v }

func (s *MonetaryArithmeticContext) GetBaseOperand() IExpressionContext { return s.baseOperand }

func (s *MonetaryArithmeticContext) Get_expression() IExpressionContext { return s._expression }

func (s *MonetaryArithmeticContext) SetBaseOperand(v IExpressionContext) { s.baseOperand = v }

func (s *MonetaryArithmeticContext) Set_expression(v IExpressionContext) { s._expression = v }

func (s *MonetaryArithmeticContext) GetOperands() []IExpressionContext { return s.operands }

func (s *MonetaryArithmeticContext) SetOperands(v []IExpressionContext) { s.operands = v }

func (s *MonetaryArithmeticContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MonetaryArithmeticContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MonetaryArithmeticContext) AllOP_ADD() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserOP_ADD)
}

func (s *MonetaryArithmeticContext) OP_ADD(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserOP_ADD, i)
}

func (s *MonetaryArithmeticContext) AllOP_SUB() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserOP_SUB)
}

func (s *MonetaryArithmeticContext) OP_SUB(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserOP_SUB, i)
}

func (s *MonetaryArithmeticContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryArithmeticContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MonetaryArithmeticContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterMonetaryArithmetic(s)
	}
}

func (s *MonetaryArithmeticContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitMonetaryArithmetic(s)
	}
}

func (p *NumScriptParser) MonetaryArithmetic() (localctx IMonetaryArithmeticContext) {
	this := p
	_ = this

	localctx = NewMonetaryArithmeticContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, NumScriptParserRULE_monetaryArithmetic)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(164)

		var _x = p.Expression()

		localctx.(*MonetaryArithmeticContext).baseOperand = _x
	}
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == NumScriptParserOP_ADD || _la == NumScriptParserOP_SUB {
		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == NumScriptParserOP_ADD || _la == NumScriptParserOP_SUB {
			{
				p.SetState(165)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*MonetaryArithmeticContext)._tset745 = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == NumScriptParserOP_ADD || _la == NumScriptParserOP_SUB) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*MonetaryArithmeticContext)._tset745 = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			localctx.(*MonetaryArithmeticContext).operators = append(localctx.(*MonetaryArithmeticContext).operators, localctx.(*MonetaryArithmeticContext)._tset745)
			{
				p.SetState(166)

				var _x = p.Expression()

				localctx.(*MonetaryArithmeticContext)._expression = _x
			}
			localctx.(*MonetaryArithmeticContext).operands = append(localctx.(*MonetaryArithmeticContext).operands, localctx.(*MonetaryArithmeticContext)._expression)

			p.SetState(169)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}

	return localctx
}

// IMonetaryExpressionContext is an interface to support dynamic dispatch.
type IMonetaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMonetaryExpressionContext differentiates from other interfaces.
	IsMonetaryExpressionContext()
}

type MonetaryExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMonetaryExpressionContext() *MonetaryExpressionContext {
	var p = new(MonetaryExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_monetaryExpression
	return p
}

func (*MonetaryExpressionContext) IsMonetaryExpressionContext() {}

func NewMonetaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MonetaryExpressionContext {
	var p = new(MonetaryExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_monetaryExpression

	return p
}

func (s *MonetaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *MonetaryExpressionContext) CopyFrom(ctx *MonetaryExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *MonetaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type MonetaryExpressionAllContext struct {
	*MonetaryExpressionContext
}

func NewMonetaryExpressionAllContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MonetaryExpressionAllContext {
	var p = new(MonetaryExpressionAllContext)

	p.MonetaryExpressionContext = NewEmptyMonetaryExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MonetaryExpressionContext))

	return p
}

func (s *MonetaryExpressionAllContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryExpressionAllContext) MonetaryAll() IMonetaryAllContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMonetaryAllContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMonetaryAllContext)
}

func (s *MonetaryExpressionAllContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterMonetaryExpressionAll(s)
	}
}

func (s *MonetaryExpressionAllContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitMonetaryExpressionAll(s)
	}
}

type MonetaryExpressionArithmeticContext struct {
	*MonetaryExpressionContext
}

func NewMonetaryExpressionArithmeticContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MonetaryExpressionArithmeticContext {
	var p = new(MonetaryExpressionArithmeticContext)

	p.MonetaryExpressionContext = NewEmptyMonetaryExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MonetaryExpressionContext))

	return p
}

func (s *MonetaryExpressionArithmeticContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MonetaryExpressionArithmeticContext) MonetaryArithmetic() IMonetaryArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMonetaryArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMonetaryArithmeticContext)
}

func (s *MonetaryExpressionArithmeticContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterMonetaryExpressionArithmetic(s)
	}
}

func (s *MonetaryExpressionArithmeticContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitMonetaryExpressionArithmetic(s)
	}
}

func (p *NumScriptParser) MonetaryExpression() (localctx IMonetaryExpressionContext) {
	this := p
	_ = this

	localctx = NewMonetaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, NumScriptParserRULE_monetaryExpression)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(175)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMonetaryExpressionAllContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(173)
			p.MonetaryAll()
		}

	case 2:
		localctx = NewMonetaryExpressionArithmeticContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(174)
			p.MonetaryArithmetic()
		}

	}

	return localctx
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) CopyFrom(ctx *StatementContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SetTxMetaContext struct {
	*StatementContext
	key   antlr.Token
	value IExpressionContext
}

func NewSetTxMetaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SetTxMetaContext {
	var p = new(SetTxMetaContext)

	p.StatementContext = NewEmptyStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*StatementContext))

	return p
}

func (s *SetTxMetaContext) GetKey() antlr.Token { return s.key }

func (s *SetTxMetaContext) SetKey(v antlr.Token) { s.key = v }

func (s *SetTxMetaContext) GetValue() IExpressionContext { return s.value }

func (s *SetTxMetaContext) SetValue(v IExpressionContext) { s.value = v }

func (s *SetTxMetaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetTxMetaContext) SET_TX_META() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSET_TX_META, 0)
}

func (s *SetTxMetaContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLPAREN, 0)
}

func (s *SetTxMetaContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRPAREN, 0)
}

func (s *SetTxMetaContext) STRING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSTRING, 0)
}

func (s *SetTxMetaContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SetTxMetaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSetTxMeta(s)
	}
}

func (s *SetTxMetaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSetTxMeta(s)
	}
}

type SetAccountMetaContext struct {
	*StatementContext
	acc   IExpressionContext
	key   antlr.Token
	value IExpressionContext
}

func NewSetAccountMetaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SetAccountMetaContext {
	var p = new(SetAccountMetaContext)

	p.StatementContext = NewEmptyStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*StatementContext))

	return p
}

func (s *SetAccountMetaContext) GetKey() antlr.Token { return s.key }

func (s *SetAccountMetaContext) SetKey(v antlr.Token) { s.key = v }

func (s *SetAccountMetaContext) GetAcc() IExpressionContext { return s.acc }

func (s *SetAccountMetaContext) GetValue() IExpressionContext { return s.value }

func (s *SetAccountMetaContext) SetAcc(v IExpressionContext) { s.acc = v }

func (s *SetAccountMetaContext) SetValue(v IExpressionContext) { s.value = v }

func (s *SetAccountMetaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetAccountMetaContext) SET_ACCOUNT_META() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSET_ACCOUNT_META, 0)
}

func (s *SetAccountMetaContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLPAREN, 0)
}

func (s *SetAccountMetaContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRPAREN, 0)
}

func (s *SetAccountMetaContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *SetAccountMetaContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SetAccountMetaContext) STRING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSTRING, 0)
}

func (s *SetAccountMetaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSetAccountMeta(s)
	}
}

func (s *SetAccountMetaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSetAccountMeta(s)
	}
}

type FailContext struct {
	*StatementContext
}

func NewFailContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FailContext {
	var p = new(FailContext)

	p.StatementContext = NewEmptyStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*StatementContext))

	return p
}

func (s *FailContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FailContext) FAIL() antlr.TerminalNode {
	return s.GetToken(NumScriptParserFAIL, 0)
}

func (s *FailContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterFail(s)
	}
}

func (s *FailContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitFail(s)
	}
}

type SendContext struct {
	*StatementContext
	mon  IMonetaryExpressionContext
	src  IValueAwareSourceContext
	dest IDestinationContext
}

func NewSendContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SendContext {
	var p = new(SendContext)

	p.StatementContext = NewEmptyStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*StatementContext))

	return p
}

func (s *SendContext) GetMon() IMonetaryExpressionContext { return s.mon }

func (s *SendContext) GetSrc() IValueAwareSourceContext { return s.src }

func (s *SendContext) GetDest() IDestinationContext { return s.dest }

func (s *SendContext) SetMon(v IMonetaryExpressionContext) { s.mon = v }

func (s *SendContext) SetSrc(v IValueAwareSourceContext) { s.src = v }

func (s *SendContext) SetDest(v IDestinationContext) { s.dest = v }

func (s *SendContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SendContext) SEND() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSEND, 0)
}

func (s *SendContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLPAREN, 0)
}

func (s *SendContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *SendContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *SendContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRPAREN, 0)
}

func (s *SendContext) MonetaryExpression() IMonetaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMonetaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMonetaryExpressionContext)
}

func (s *SendContext) SOURCE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSOURCE, 0)
}

func (s *SendContext) AllEQ() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserEQ)
}

func (s *SendContext) EQ(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserEQ, i)
}

func (s *SendContext) DESTINATION() antlr.TerminalNode {
	return s.GetToken(NumScriptParserDESTINATION, 0)
}

func (s *SendContext) ValueAwareSource() IValueAwareSourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueAwareSourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueAwareSourceContext)
}

func (s *SendContext) Destination() IDestinationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDestinationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDestinationContext)
}

func (s *SendContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterSend(s)
	}
}

func (s *SendContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitSend(s)
	}
}

func (p *NumScriptParser) Statement() (localctx IStatementContext) {
	this := p
	_ = this

	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, NumScriptParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(219)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserSET_TX_META:
		localctx = NewSetTxMetaContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(177)
			p.Match(NumScriptParserSET_TX_META)
		}
		{
			p.SetState(178)
			p.Match(NumScriptParserLPAREN)
		}
		{
			p.SetState(179)

			var _m = p.Match(NumScriptParserSTRING)

			localctx.(*SetTxMetaContext).key = _m
		}
		{
			p.SetState(180)
			p.Match(NumScriptParserT__3)
		}
		{
			p.SetState(181)

			var _x = p.Expression()

			localctx.(*SetTxMetaContext).value = _x
		}
		{
			p.SetState(182)
			p.Match(NumScriptParserRPAREN)
		}

	case NumScriptParserSET_ACCOUNT_META:
		localctx = NewSetAccountMetaContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(184)
			p.Match(NumScriptParserSET_ACCOUNT_META)
		}
		{
			p.SetState(185)
			p.Match(NumScriptParserLPAREN)
		}
		{
			p.SetState(186)

			var _x = p.Expression()

			localctx.(*SetAccountMetaContext).acc = _x
		}
		{
			p.SetState(187)
			p.Match(NumScriptParserT__3)
		}
		{
			p.SetState(188)

			var _m = p.Match(NumScriptParserSTRING)

			localctx.(*SetAccountMetaContext).key = _m
		}
		{
			p.SetState(189)
			p.Match(NumScriptParserT__3)
		}
		{
			p.SetState(190)

			var _x = p.Expression()

			localctx.(*SetAccountMetaContext).value = _x
		}
		{
			p.SetState(191)
			p.Match(NumScriptParserRPAREN)
		}

	case NumScriptParserFAIL:
		localctx = NewFailContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(193)
			p.Match(NumScriptParserFAIL)
		}

	case NumScriptParserSEND:
		localctx = NewSendContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(194)
			p.Match(NumScriptParserSEND)
		}
		{
			p.SetState(195)

			var _x = p.MonetaryExpression()

			localctx.(*SendContext).mon = _x
		}
		{
			p.SetState(196)
			p.Match(NumScriptParserLPAREN)
		}
		{
			p.SetState(197)
			p.Match(NumScriptParserNEWLINE)
		}
		p.SetState(214)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case NumScriptParserSOURCE:
			{
				p.SetState(198)
				p.Match(NumScriptParserSOURCE)
			}
			{
				p.SetState(199)
				p.Match(NumScriptParserEQ)
			}
			{
				p.SetState(200)

				var _x = p.ValueAwareSource()

				localctx.(*SendContext).src = _x
			}
			{
				p.SetState(201)
				p.Match(NumScriptParserNEWLINE)
			}
			{
				p.SetState(202)
				p.Match(NumScriptParserDESTINATION)
			}
			{
				p.SetState(203)
				p.Match(NumScriptParserEQ)
			}
			{
				p.SetState(204)

				var _x = p.Destination()

				localctx.(*SendContext).dest = _x
			}

		case NumScriptParserDESTINATION:
			{
				p.SetState(206)
				p.Match(NumScriptParserDESTINATION)
			}
			{
				p.SetState(207)
				p.Match(NumScriptParserEQ)
			}
			{
				p.SetState(208)

				var _x = p.Destination()

				localctx.(*SendContext).dest = _x
			}
			{
				p.SetState(209)
				p.Match(NumScriptParserNEWLINE)
			}
			{
				p.SetState(210)
				p.Match(NumScriptParserSOURCE)
			}
			{
				p.SetState(211)
				p.Match(NumScriptParserEQ)
			}
			{
				p.SetState(212)

				var _x = p.ValueAwareSource()

				localctx.(*SendContext).src = _x
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(216)
			p.Match(NumScriptParserNEWLINE)
		}
		{
			p.SetState(217)
			p.Match(NumScriptParserRPAREN)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IType_Context is an interface to support dynamic dispatch.
type IType_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsType_Context differentiates from other interfaces.
	IsType_Context()
}

type Type_Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyType_Context() *Type_Context {
	var p = new(Type_Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_type_
	return p
}

func (*Type_Context) IsType_Context() {}

func NewType_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Type_Context {
	var p = new(Type_Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_type_

	return p
}

func (s *Type_Context) GetParser() antlr.Parser { return s.parser }

func (s *Type_Context) TY_ACCOUNT() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_ACCOUNT, 0)
}

func (s *Type_Context) TY_ASSET() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_ASSET, 0)
}

func (s *Type_Context) TY_NUMBER() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_NUMBER, 0)
}

func (s *Type_Context) TY_STRING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_STRING, 0)
}

func (s *Type_Context) TY_MONETARY() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_MONETARY, 0)
}

func (s *Type_Context) TY_PORTION() antlr.TerminalNode {
	return s.GetToken(NumScriptParserTY_PORTION, 0)
}

func (s *Type_Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Type_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Type_Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterType_(s)
	}
}

func (s *Type_Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitType_(s)
	}
}

func (p *NumScriptParser) Type_() (localctx IType_Context) {
	this := p
	_ = this

	localctx = NewType_Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, NumScriptParserRULE_type_)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(221)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-30)&-(0x1f+1)) == 0 && ((1<<uint((_la-30)))&((1<<(NumScriptParserTY_ACCOUNT-30))|(1<<(NumScriptParserTY_ASSET-30))|(1<<(NumScriptParserTY_NUMBER-30))|(1<<(NumScriptParserTY_MONETARY-30))|(1<<(NumScriptParserTY_PORTION-30))|(1<<(NumScriptParserTY_STRING-30)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IOriginContext is an interface to support dynamic dispatch.
type IOriginContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOriginContext differentiates from other interfaces.
	IsOriginContext()
}

type OriginContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOriginContext() *OriginContext {
	var p = new(OriginContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_origin
	return p
}

func (*OriginContext) IsOriginContext() {}

func NewOriginContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OriginContext {
	var p = new(OriginContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_origin

	return p
}

func (s *OriginContext) GetParser() antlr.Parser { return s.parser }

func (s *OriginContext) CopyFrom(ctx *OriginContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *OriginContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OriginContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type OriginAccountBalanceContext struct {
	*OriginContext
	account IExpressionContext
	asset   IExpressionContext
}

func NewOriginAccountBalanceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OriginAccountBalanceContext {
	var p = new(OriginAccountBalanceContext)

	p.OriginContext = NewEmptyOriginContext()
	p.parser = parser
	p.CopyFrom(ctx.(*OriginContext))

	return p
}

func (s *OriginAccountBalanceContext) GetAccount() IExpressionContext { return s.account }

func (s *OriginAccountBalanceContext) GetAsset() IExpressionContext { return s.asset }

func (s *OriginAccountBalanceContext) SetAccount(v IExpressionContext) { s.account = v }

func (s *OriginAccountBalanceContext) SetAsset(v IExpressionContext) { s.asset = v }

func (s *OriginAccountBalanceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OriginAccountBalanceContext) BALANCE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserBALANCE, 0)
}

func (s *OriginAccountBalanceContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLPAREN, 0)
}

func (s *OriginAccountBalanceContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRPAREN, 0)
}

func (s *OriginAccountBalanceContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *OriginAccountBalanceContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *OriginAccountBalanceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterOriginAccountBalance(s)
	}
}

func (s *OriginAccountBalanceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitOriginAccountBalance(s)
	}
}

type OriginAccountMetaContext struct {
	*OriginContext
	account IExpressionContext
	key     antlr.Token
}

func NewOriginAccountMetaContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OriginAccountMetaContext {
	var p = new(OriginAccountMetaContext)

	p.OriginContext = NewEmptyOriginContext()
	p.parser = parser
	p.CopyFrom(ctx.(*OriginContext))

	return p
}

func (s *OriginAccountMetaContext) GetKey() antlr.Token { return s.key }

func (s *OriginAccountMetaContext) SetKey(v antlr.Token) { s.key = v }

func (s *OriginAccountMetaContext) GetAccount() IExpressionContext { return s.account }

func (s *OriginAccountMetaContext) SetAccount(v IExpressionContext) { s.account = v }

func (s *OriginAccountMetaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OriginAccountMetaContext) META() antlr.TerminalNode {
	return s.GetToken(NumScriptParserMETA, 0)
}

func (s *OriginAccountMetaContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLPAREN, 0)
}

func (s *OriginAccountMetaContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRPAREN, 0)
}

func (s *OriginAccountMetaContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *OriginAccountMetaContext) STRING() antlr.TerminalNode {
	return s.GetToken(NumScriptParserSTRING, 0)
}

func (s *OriginAccountMetaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterOriginAccountMeta(s)
	}
}

func (s *OriginAccountMetaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitOriginAccountMeta(s)
	}
}

func (p *NumScriptParser) Origin() (localctx IOriginContext) {
	this := p
	_ = this

	localctx = NewOriginContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, NumScriptParserRULE_origin)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(237)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case NumScriptParserMETA:
		localctx = NewOriginAccountMetaContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(223)
			p.Match(NumScriptParserMETA)
		}
		{
			p.SetState(224)
			p.Match(NumScriptParserLPAREN)
		}
		{
			p.SetState(225)

			var _x = p.Expression()

			localctx.(*OriginAccountMetaContext).account = _x
		}
		{
			p.SetState(226)
			p.Match(NumScriptParserT__3)
		}
		{
			p.SetState(227)

			var _m = p.Match(NumScriptParserSTRING)

			localctx.(*OriginAccountMetaContext).key = _m
		}
		{
			p.SetState(228)
			p.Match(NumScriptParserRPAREN)
		}

	case NumScriptParserBALANCE:
		localctx = NewOriginAccountBalanceContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(230)
			p.Match(NumScriptParserBALANCE)
		}
		{
			p.SetState(231)
			p.Match(NumScriptParserLPAREN)
		}
		{
			p.SetState(232)

			var _x = p.Expression()

			localctx.(*OriginAccountBalanceContext).account = _x
		}
		{
			p.SetState(233)
			p.Match(NumScriptParserT__3)
		}
		{
			p.SetState(234)

			var _x = p.Expression()

			localctx.(*OriginAccountBalanceContext).asset = _x
		}
		{
			p.SetState(235)
			p.Match(NumScriptParserRPAREN)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IVarDeclContext is an interface to support dynamic dispatch.
type IVarDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTy returns the ty rule contexts.
	GetTy() IType_Context

	// GetName returns the name rule contexts.
	GetName() IVariableContext

	// GetOrig returns the orig rule contexts.
	GetOrig() IOriginContext

	// SetTy sets the ty rule contexts.
	SetTy(IType_Context)

	// SetName sets the name rule contexts.
	SetName(IVariableContext)

	// SetOrig sets the orig rule contexts.
	SetOrig(IOriginContext)

	// IsVarDeclContext differentiates from other interfaces.
	IsVarDeclContext()
}

type VarDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	ty     IType_Context
	name   IVariableContext
	orig   IOriginContext
}

func NewEmptyVarDeclContext() *VarDeclContext {
	var p = new(VarDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_varDecl
	return p
}

func (*VarDeclContext) IsVarDeclContext() {}

func NewVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclContext {
	var p = new(VarDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_varDecl

	return p
}

func (s *VarDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *VarDeclContext) GetTy() IType_Context { return s.ty }

func (s *VarDeclContext) GetName() IVariableContext { return s.name }

func (s *VarDeclContext) GetOrig() IOriginContext { return s.orig }

func (s *VarDeclContext) SetTy(v IType_Context) { s.ty = v }

func (s *VarDeclContext) SetName(v IVariableContext) { s.name = v }

func (s *VarDeclContext) SetOrig(v IOriginContext) { s.orig = v }

func (s *VarDeclContext) Type_() IType_Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IType_Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IType_Context)
}

func (s *VarDeclContext) Variable() IVariableContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *VarDeclContext) EQ() antlr.TerminalNode {
	return s.GetToken(NumScriptParserEQ, 0)
}

func (s *VarDeclContext) Origin() IOriginContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOriginContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOriginContext)
}

func (s *VarDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterVarDecl(s)
	}
}

func (s *VarDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitVarDecl(s)
	}
}

func (p *NumScriptParser) VarDecl() (localctx IVarDeclContext) {
	this := p
	_ = this

	localctx = NewVarDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, NumScriptParserRULE_varDecl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(239)

		var _x = p.Type_()

		localctx.(*VarDeclContext).ty = _x
	}
	{
		p.SetState(240)

		var _x = p.Variable()

		localctx.(*VarDeclContext).name = _x
	}
	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == NumScriptParserEQ {
		{
			p.SetState(241)
			p.Match(NumScriptParserEQ)
		}
		{
			p.SetState(242)

			var _x = p.Origin()

			localctx.(*VarDeclContext).orig = _x
		}

	}

	return localctx
}

// IVarListDeclContext is an interface to support dynamic dispatch.
type IVarListDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_varDecl returns the _varDecl rule contexts.
	Get_varDecl() IVarDeclContext

	// Set_varDecl sets the _varDecl rule contexts.
	Set_varDecl(IVarDeclContext)

	// GetV returns the v rule context list.
	GetV() []IVarDeclContext

	// SetV sets the v rule context list.
	SetV([]IVarDeclContext)

	// IsVarListDeclContext differentiates from other interfaces.
	IsVarListDeclContext()
}

type VarListDeclContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	_varDecl IVarDeclContext
	v        []IVarDeclContext
}

func NewEmptyVarListDeclContext() *VarListDeclContext {
	var p = new(VarListDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_varListDecl
	return p
}

func (*VarListDeclContext) IsVarListDeclContext() {}

func NewVarListDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarListDeclContext {
	var p = new(VarListDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_varListDecl

	return p
}

func (s *VarListDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *VarListDeclContext) Get_varDecl() IVarDeclContext { return s._varDecl }

func (s *VarListDeclContext) Set_varDecl(v IVarDeclContext) { s._varDecl = v }

func (s *VarListDeclContext) GetV() []IVarDeclContext { return s.v }

func (s *VarListDeclContext) SetV(v []IVarDeclContext) { s.v = v }

func (s *VarListDeclContext) VARS() antlr.TerminalNode {
	return s.GetToken(NumScriptParserVARS, 0)
}

func (s *VarListDeclContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserLBRACE, 0)
}

func (s *VarListDeclContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *VarListDeclContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *VarListDeclContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(NumScriptParserRBRACE, 0)
}

func (s *VarListDeclContext) AllVarDecl() []IVarDeclContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IVarDeclContext); ok {
			len++
		}
	}

	tst := make([]IVarDeclContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IVarDeclContext); ok {
			tst[i] = t.(IVarDeclContext)
			i++
		}
	}

	return tst
}

func (s *VarListDeclContext) VarDecl(i int) IVarDeclContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarDeclContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarDeclContext)
}

func (s *VarListDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarListDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarListDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterVarListDecl(s)
	}
}

func (s *VarListDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitVarListDecl(s)
	}
}

func (p *NumScriptParser) VarListDecl() (localctx IVarListDeclContext) {
	this := p
	_ = this

	localctx = NewVarListDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, NumScriptParserRULE_varListDecl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(245)
		p.Match(NumScriptParserVARS)
	}
	{
		p.SetState(246)
		p.Match(NumScriptParserLBRACE)
	}
	{
		p.SetState(247)
		p.Match(NumScriptParserNEWLINE)
	}
	p.SetState(254)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-30)&-(0x1f+1)) == 0 && ((1<<uint((_la-30)))&((1<<(NumScriptParserTY_ACCOUNT-30))|(1<<(NumScriptParserTY_ASSET-30))|(1<<(NumScriptParserTY_NUMBER-30))|(1<<(NumScriptParserTY_MONETARY-30))|(1<<(NumScriptParserTY_PORTION-30))|(1<<(NumScriptParserTY_STRING-30)))) != 0) {
		{
			p.SetState(248)

			var _x = p.VarDecl()

			localctx.(*VarListDeclContext)._varDecl = _x
		}
		localctx.(*VarListDeclContext).v = append(localctx.(*VarListDeclContext).v, localctx.(*VarListDeclContext)._varDecl)
		p.SetState(250)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == NumScriptParserNEWLINE {
			{
				p.SetState(249)
				p.Match(NumScriptParserNEWLINE)
			}

			p.SetState(252)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

		p.SetState(256)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(258)
		p.Match(NumScriptParserRBRACE)
	}
	{
		p.SetState(259)
		p.Match(NumScriptParserNEWLINE)
	}

	return localctx
}

// IScriptContext is an interface to support dynamic dispatch.
type IScriptContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVars returns the vars rule contexts.
	GetVars() IVarListDeclContext

	// Get_statement returns the _statement rule contexts.
	Get_statement() IStatementContext

	// SetVars sets the vars rule contexts.
	SetVars(IVarListDeclContext)

	// Set_statement sets the _statement rule contexts.
	Set_statement(IStatementContext)

	// GetStmts returns the stmts rule context list.
	GetStmts() []IStatementContext

	// SetStmts sets the stmts rule context list.
	SetStmts([]IStatementContext)

	// IsScriptContext differentiates from other interfaces.
	IsScriptContext()
}

type ScriptContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	vars       IVarListDeclContext
	_statement IStatementContext
	stmts      []IStatementContext
}

func NewEmptyScriptContext() *ScriptContext {
	var p = new(ScriptContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = NumScriptParserRULE_script
	return p
}

func (*ScriptContext) IsScriptContext() {}

func NewScriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScriptContext {
	var p = new(ScriptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = NumScriptParserRULE_script

	return p
}

func (s *ScriptContext) GetParser() antlr.Parser { return s.parser }

func (s *ScriptContext) GetVars() IVarListDeclContext { return s.vars }

func (s *ScriptContext) Get_statement() IStatementContext { return s._statement }

func (s *ScriptContext) SetVars(v IVarListDeclContext) { s.vars = v }

func (s *ScriptContext) Set_statement(v IStatementContext) { s._statement = v }

func (s *ScriptContext) GetStmts() []IStatementContext { return s.stmts }

func (s *ScriptContext) SetStmts(v []IStatementContext) { s.stmts = v }

func (s *ScriptContext) EOF() antlr.TerminalNode {
	return s.GetToken(NumScriptParserEOF, 0)
}

func (s *ScriptContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *ScriptContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *ScriptContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(NumScriptParserNEWLINE)
}

func (s *ScriptContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(NumScriptParserNEWLINE, i)
}

func (s *ScriptContext) VarListDecl() IVarListDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarListDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarListDeclContext)
}

func (s *ScriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.EnterScript(s)
	}
}

func (s *ScriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(NumScriptListener); ok {
		listenerT.ExitScript(s)
	}
}

func (p *NumScriptParser) Script() (localctx IScriptContext) {
	this := p
	_ = this

	localctx = NewScriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, NumScriptParserRULE_script)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(264)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == NumScriptParserNEWLINE {
		{
			p.SetState(261)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(266)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(268)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == NumScriptParserVARS {
		{
			p.SetState(267)

			var _x = p.VarListDecl()

			localctx.(*ScriptContext).vars = _x
		}

	}
	{
		p.SetState(270)

		var _x = p.Statement()

		localctx.(*ScriptContext)._statement = _x
	}
	localctx.(*ScriptContext).stmts = append(localctx.(*ScriptContext).stmts, localctx.(*ScriptContext)._statement)
	p.SetState(275)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(271)
				p.Match(NumScriptParserNEWLINE)
			}
			{
				p.SetState(272)

				var _x = p.Statement()

				localctx.(*ScriptContext)._statement = _x
			}
			localctx.(*ScriptContext).stmts = append(localctx.(*ScriptContext).stmts, localctx.(*ScriptContext)._statement)

		}
		p.SetState(277)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())
	}
	p.SetState(281)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == NumScriptParserNEWLINE {
		{
			p.SetState(278)
			p.Match(NumScriptParserNEWLINE)
		}

		p.SetState(283)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(284)
		p.Match(NumScriptParserEOF)
	}

	return localctx
}
