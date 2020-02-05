package lao_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vectorhacker/lao/pkg/lao"
)

type fakeTokenizer struct {
	tokens   []lao.Token
	position int
}

func (t *fakeTokenizer) Next() bool {
	if t.position < len(t.tokens) {
		t.position++
		return true
	}
	return false
}

func (t *fakeTokenizer) Current() lao.Token {
	if t.position < 0 {
		return lao.Token{}
	}
	if t.position >= len(t.tokens) {
		return lao.Token{
			Kind: lao.KindEnd,
		}
	}
	return t.tokens[t.position]
}

func TestParser(t *testing.T) {
	testCases := []struct {
		desc      string
		tokenizer lao.Tokenizer
		expected  []lao.Node
	}{
		// {
		// 	desc:      "",
		// 	tokenizer: lao.NewTokenizer(strings.NewReader("print\nprint 1.233")),
		// 	expected: []lao.Node{
		// 		lao.PrintStatement{},
		// 		lao.PrintStatement{
		// 			Argumenent: lao.RealNumber{
		// 				Value: "1.233",
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	desc:      "rem statements",
		// 	tokenizer: lao.NewTokenizer(strings.NewReader("rem hello world\nprint\nprint")),
		// 	expected: []lao.Node{
		// 		lao.RemStatement{},
		// 		lao.PrintStatement{},
		// 		lao.PrintStatement{},
		// 	},
		// },
		{
			desc: "read statements",
			tokenizer: lao.NewTokenizer(
				strings.NewReader("read a"),
			),
			expected: []lao.Node{
				lao.ReadStatement{
					Variable: lao.Variable{
						Type: lao.VariableInteger,
						Name: "a",
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			p := lao.NewParser(tC.tokenizer)

			got, err := p.Parse()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.expected, got)
		})
	}
}
