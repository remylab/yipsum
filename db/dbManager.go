package db

import (
)

type DbManager interface {
    Close() error

    IsNewUri(s string) (bool,error)
    ValidateUriKey(ipsum string, key string) (bool,error)

    GetIpsum(s string) (map[string]string, error)
    CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error)
    AddText(ipsumId int, text string) (sqlRes, error)
    UpdateText(dataId int, text string) (sqlRes, error)
    DeleteText(dataId int) (sqlRes, error)
}
