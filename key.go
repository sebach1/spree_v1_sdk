package spree

type token string

func (t token) validate() error {
	if t == "" {
		return errNilKey
	}
	return nil
}
