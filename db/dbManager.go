package db

import (
)

type DbManager interface {
    Close() error
    CheckUri(s string) (bool,error)
}
