package lao_test

import (
	"github.com/vectorhacker/lao/pkg/lao"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc           string
		input          string
		expectedTokens []lao.Token
	}{
		{
			desc:  "recognize identifiers",
			input: "print then if allow person xyz",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindKeyword,
					Value:  "print",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindKeyword,
					Value:  "then",
					Line:   1,
					Column: 7,
				},
				{
					Kind:   lao.KindKeyword,
					Value:  "if",
					Line:   1,
					Column: 12,
				},
				{
					Kind:   lao.KindIdentifier,
					Value:  "allow",
					Line:   1,
					Column: 15,
				},
				{
					Kind:   lao.KindIdentifier,
					Value:  "person",
					Line:   1,
					Column: 21,
				},
				{
					Kind:   lao.KindIdentifier,
					Value:  "xyz",
					Line:   1,
					Column: 28,
				},
			},
		},
		{
			desc:  "recognize numbers",
			input: "1 1.2 1.22e2.22 1.2E2.5",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindInteger,
					Value:  "1",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindReal,
					Value:  "1.2",
					Line:   1,
					Column: 3,
				},
				{
					Kind:   lao.KindReal,
					Value:  "1.22e2.22",
					Line:   1,
					Column: 7,
				},
				{
					Kind:   lao.KindReal,
					Value:  "1.2E2.5",
					Line:   1,
					Column: 17,
				},
			},
		},
		{
			desc:  "recognize assignment",
			input: "=",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindAssignment,
					Value:  "=",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			desc:  "recognize period",
			input: ".\n.\n.",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindPeriod,
					Value:  ".",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindPeriod,
					Value:  ".",
					Line:   2,
					Column: 1,
				},
				{
					Kind:   lao.KindPeriod,
					Value:  ".",
					Line:   3,
					Column: 1,
				},
			},
		},
		{
			desc:  "recognize arithmetic operators",
			input: ".add. .sub. .div. .mul.",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindArithmeticOperator,
					Value:  ".add.",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindArithmeticOperator,
					Value:  ".sub.",
					Line:   1,
					Column: 7,
				},
				{
					Kind:   lao.KindArithmeticOperator,
					Value:  ".div.",
					Line:   1,
					Column: 13,
				},
				{
					Kind:   lao.KindArithmeticOperator,
					Value:  ".mul.",
					Line:   1,
					Column: 19,
				},
			},
		},
		{
			desc:  "recognize logical operators",
			input: ".or. .and. .not.",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindLogicalOperator,
					Value:  ".or.",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindLogicalOperator,
					Value:  ".and.",
					Line:   1,
					Column: 6,
				},
				{
					Kind:   lao.KindLogicalOperator,
					Value:  ".not.",
					Line:   1,
					Column: 12,
				},
			},
		},
		{
			desc:  "recognize string",
			input: `"HOLA MUNDO"`,
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindString,
					Value:  `"HOLA MUNDO"`,
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			desc:  "recognize relational operators",
			input: ".eq. .lt. .ne. .le. .gt. .ge.",
			expectedTokens: []lao.Token{
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".eq.",
					Line:   1,
					Column: 1,
				},
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".lt.",
					Line:   1,
					Column: 6,
				},
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".ne.",
					Line:   1,
					Column: 11,
				},
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".le.",
					Line:   1,
					Column: 16,
				},
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".gt.",
					Line:   1,
					Column: 21,
				},
				{
					Kind:   lao.KindRelationalOperator,
					Value:  ".ge.",
					Line:   1,
					Column: 26,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			r := strings.NewReader(tC.input)

			tokenizer := lao.NewTokenizer(r)

			got := []lao.Token{}

			for tokenizer.Next() {
				got = append(got, tokenizer.Current())
			}

			assert.Equal(t, tC.expectedTokens, got)
		})
	}
}
