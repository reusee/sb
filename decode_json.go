package sb

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeJson(r io.Reader, cont Proc) *Proc {
	decoder := json.NewDecoder(r)
	var proc Proc
	proc = func() (*Token, Proc, error) {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil, cont, nil
		}
		if err != nil {
			return nil, nil, err
		}

		switch token := token.(type) {

		case json.Delim:
			switch token {
			case '[':
				return &Token{
					Kind: KindArray,
				}, proc, nil
			case ']':
				return &Token{
					Kind: KindArrayEnd,
				}, proc, nil
			case '{':
				return &Token{
					Kind: KindObject,
				}, proc, nil
			case '}':
				return &Token{
					Kind: KindObjectEnd,
				}, proc, nil
			default:
				return nil, nil, fmt.Errorf("bad delimiter rune: %v", token)
			}

		case bool:
			return &Token{
				Kind:  KindBool,
				Value: token,
			}, proc, nil

		case float64:
			return &Token{
				Kind:  KindFloat64,
				Value: token,
			}, proc, nil

		case json.Number:
			return &Token{
				Kind:  KindLiteral,
				Value: string(token),
			}, proc, nil

		case string:
			return &Token{
				Kind:  KindString,
				Value: token,
			}, proc, nil

		case nil:
			return &Token{
				Kind: KindNil,
			}, proc, nil

		}

		return nil, nil, fmt.Errorf("bad token type: %T", token)
	}
	return &proc
}
