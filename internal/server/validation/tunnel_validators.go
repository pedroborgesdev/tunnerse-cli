package validation

import (
	"errors"
	"regexp"
)

var (
	ErrBannedName = errors.New("tunnel name contais inavalid characters")
)

type TunnelValidator struct {
	nameRegex *regexp.Regexp
}

func NewTunnelValidator() *TunnelValidator {
	return &TunnelValidator{
		nameRegex: regexp.MustCompile(`^[a-z-]{1,20}$`),
	}
}

func (v *TunnelValidator) ValidateTunnelRegister(name string) error {




	return nil
}







