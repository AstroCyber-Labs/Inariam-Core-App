package config

import (
	"errors"
	"fmt"
	"gitea/pcp-inariam/inariam/pkgs/log"
	"github.com/go-playground/validator/v10"
)

// Validate checks the Config struct for any validation errors based on the struct tags.
// If there are errors, it prints detailed information about each of them.
func (config *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(config)

	if err != nil {
		errString := ""
		for _, err := range err.(validator.ValidationErrors) {
			log.Logger.Infof("%s %s %s %s %s %s %s %s %s %s", err.Namespace(), err.Field(),
				err.StructNamespace(),
				err.StructField(),
				err.Tag(),
				err.ActualTag(),
				err.Kind(),
				err.Type(),
				err.Value(),
				err.Param(),
			)

			errString += fmt.Sprintf("%s %s: %s -  %s\n", err.StructNamespace(), err.StructField(), err.Value(), err.Error())
		}

		return errors.New(errString)
		// from here you can create your own error messages in whatever language you wish
	}

	return nil
}
