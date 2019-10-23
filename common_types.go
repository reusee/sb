package sb

import (
	"reflect"
	"time"
)

var commonTokenizers = map[reflect.Type]func(any) []Token{
	// time.Time
	reflect.TypeOf((*time.Time)(nil)).Elem(): func(v any) []Token {
		bin, err := v.(time.Time).MarshalBinary()
		if err != nil {
			panic(err)
		}
		return []Token{
			{KindString, string(bin)},
		}
	},
}

var commonDetokenizers = map[reflect.Type]func(Stream, any) (Token, error){
	// *time.Time
	reflect.TypeOf((**time.Time)(nil)).Elem(): func(stream Stream, ptr any) (token Token, err error) {
		p := stream.Next()
		if p == nil {
			return
		}
		token = *p
		if token.Kind != KindString {
			return
		}
		if err = ptr.(*time.Time).UnmarshalBinary([]byte(token.Value.(string))); err != nil {
			return
		}
		return
	},
}
