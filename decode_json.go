package sb

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeJson(r io.Reader, cont Proc) *Proc {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var proc Proc
	proc = func(token *Token) (Proc, error) {
		jsonToken, err := decoder.Token()
		if err == io.EOF {
			return cont, nil
		}
		if err != nil {
			return nil, err
		}

		switch jsonToken := jsonToken.(type) {

		case json.Delim:
			switch jsonToken {
			case '[':
				token.Kind = KindArray
				return proc, nil
			case ']':
				token.Kind = KindArrayEnd
				return proc, nil
			case '{':
				token.Kind = KindObject
				return proc, nil
			case '}':
				token.Kind = KindObjectEnd
				return proc, nil
			default:
				return nil, fmt.Errorf("bad delimiter rune: %v", token)
			}

		case bool:
			token.Kind = KindBool
			token.Value = token
			return proc, nil

		case float64:
			token.Kind = KindFloat64
			token.Value = token
			return proc, nil

		case json.Number:
			token.Kind = KindLiteral
			token.Value = string(jsonToken)
			return proc, nil

		case string:
			token.Kind = KindString
			token.Value = token
			return proc, nil

		case nil:
			token.Kind = KindNil
			return proc, nil

		}

		return nil, fmt.Errorf("bad token type: %T", token)
	}

	return &proc
}
