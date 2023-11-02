package caskdb

func validateKV(key string, value []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}

	if len(key) > MaxKeySize {
		return ErrLargeKey
	}

	if len(value) > MaxValueSize {
		return ErrLargeValue
	}

	return nil
}
