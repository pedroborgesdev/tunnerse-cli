package validators

import (
	"errors"
	"regexp"
	"tunnerse/dto"
)

var (
	ErrInvalidTunnelID = errors.New("tunnel id contais invalid characters")
	ErrInvalidAddress  = errors.New("address contains invalid characters or format")
)

type ArgsValidator struct {
	regexTunnelID *regexp.Regexp
	regexAddress  *regexp.Regexp
}

// NewArgsValidator creates and returns a new ArgsValidator with compiled regex patterns.
func NewArgsValidator() *ArgsValidator {
	return &ArgsValidator{
		regexTunnelID: regexp.MustCompile(`^[a-z-]{1,20}$`),
		regexAddress:  regexp.MustCompile(`^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([1-9][0-9]{0,3})|0)$`),
	}
}

// ValidateUsageArgs checks if the command-line arguments have the expected length.
func (v *ArgsValidator) ValidateUsageArgs(args []string) string {
	if len(args) > 3 {
		return dto.Invalid
	}

	if len(args) > 2 {
		if args[1] == "help" || args[1] == "version" {
			return dto.Invalid
		}
	}

	if len(args) == 2 {
		switch args[1] {
		case "help":
			return dto.Help
		case "version":
			return dto.Info
		default:
			return dto.Invalid
		}
	}

	if len(args) == 3 {
		if err := v.validateTunnelID(args[1]); err != nil {
			return dto.InvalidID
		}

		if err := v.validateAddress(args[2]); err != nil {
			return dto.InvalidPort
		}

		return ""
	}

	return dto.Help
}

// ValidateExposeArgs validates both the tunnel ID and the address provided as arguments.
func (v *ArgsValidator) ValidateExposeArgs(tunnelID, address string) error {
	if err := v.validateTunnelID(tunnelID); err != nil {
		return err
	}

	if err := v.validateAddress(address); err != nil {
		return err
	}

	return nil
}

// validateTunnelID checks whether the tunnel ID matches the expected regex pattern.
func (v *ArgsValidator) validateTunnelID(tunnelID string) error {
	if !v.regexTunnelID.MatchString(tunnelID) {
		return ErrInvalidTunnelID
	}
	return nil
}

// validateAddress checks whether the address matches the expected IP:PORT or localhost:PORT format.
func (v *ArgsValidator) validateAddress(address string) error {
	if !v.regexAddress.MatchString(address) {
		return ErrInvalidAddress
	}
	return nil
}
