package configutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Binds the key with viper.bindenv, returns the value in lowercase or nil if unset
func GetEnvString(vpr *viper.Viper, name string) *string {
	if err := vpr.BindEnv(name); err != nil {
		return nil
	}

	res := strings.ToLower(vpr.GetString(name))
	if res == "" {
		return nil
	}

	return &res
}

func GetEnvInt(vpr *viper.Viper, name string) *int {
	str := GetEnvString(vpr, name)
	if str == nil {
		return nil
	}

	num, err := strconv.Atoi(*str)
	if err != nil {
		return nil
	}

	return &num
}

type ValueNotAllowedError struct {
	name     string
	received string
	allowed  []string
}

func (e ValueNotAllowedError) Error() string {
	return fmt.Sprintf("Invalid value for %s; Expected one of %v; Got %s", e.name, e.allowed, e.received)
}

func NewErrValueNotAllowed(varName, received string, allowed []string) ValueNotAllowedError {
	return ValueNotAllowedError{
		name:     varName,
		received: received,
		allowed:  allowed,
	}
}
