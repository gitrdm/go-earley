package scanner

func RunToEnd(scanner Scanner) (bool, error) {
	for !scanner.EndOfStream() {
		ok, err := scanner.Read()
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return scanner.Parser().Accepted(), nil
}
