package aliases

func shouldError(err error) error {
	switch err.Error() {
	case "input does not match format",
		"expected space in input to match format",
		"EOF",
		"unexpected EOF":
		return nil
	default:
		return err
	}
}
