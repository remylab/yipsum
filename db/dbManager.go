package db

import (
)

type DbManager interface {
    Close() error
    IsNewUri(s string) (bool,error)
    CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error)
    GetIpsum(s string) (map[string]string, error)
    ValidateUriKey(ipsum string, key string) (bool,error)
}
