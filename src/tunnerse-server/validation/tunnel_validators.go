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
	// if err := v.validadeName(name); err != nil {
	// 	return err
	// }

	return nil
}

// func (v *TunnelValidator) validadeName(name string) error {
// 	if !v.nameRegex.MatchString(name) {
// 		return ErrBannedName
// 	}
// 	return nil
// }
