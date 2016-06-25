package common

import (
	"os"
)

func GetRootPath() string {
	return os.Getenv("yip_root")
}