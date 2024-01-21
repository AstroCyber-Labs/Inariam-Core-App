package utils

import (
	"fmt"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

func Errorf(msg string, err error) error {
	log.Logger.Errorf(msg, err)
	return fmt.Errorf(msg, err)
}
