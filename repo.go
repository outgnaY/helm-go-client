package helmclient

import (
	"github.com/pkg/errors"
	"os"
)

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}
