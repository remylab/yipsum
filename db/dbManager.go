package db

import (
)

type DbManager interface {
    Close() error
    CheckUri(s string) (bool,error)
    CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error)
}
