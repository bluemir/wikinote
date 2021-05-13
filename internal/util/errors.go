package util

type MultipleError struct {
	Causes []error
}

func (errs MultipleError) Error() string {
	str := "multiple error occur. cause:\n"

	for _, err := range errs.Causes {
		str += err.Error() + "\n"
	}

	return str
}

func MergeErrors(errs ...error) error {
	ret := &MultipleError{}
	for _, err := range errs {
		if err == nil {
			continue // skip
		}
		ret.Causes = append(ret.Causes, err)
	}
	if len(ret.Causes) == 0 {
		return nil
	}
	return ret
}
