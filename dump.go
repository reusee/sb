package sb

func dumpStream(
	stream Stream,
	prefix string,
) {
	for {
		token, err := stream.Next()
		if err != nil {
			panic(err)
		}
		if token == nil {
			break
		}
		pt("%s%+v\n", prefix, token)
	}
}
