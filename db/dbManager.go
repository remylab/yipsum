package db

import (
)

type DbManager interface {
    Close() error

    IsNewUri(s string) (bool,error)
    ValidateUriKey(ipsum string, key string) (bool,error)

    GetIpsum(s string) (map[string]string, error)
    GetIpsumTextsForPage(ipsumId int64, pageNum int64, resByPage int64) ([]map[string]string, error)
    GetTotalIpsumTexts(ipsumId int64) (int,error)

    CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error)
    AddText(ipsumId int64, text string) (sqlRes, error)
    UpdateText(ipsumId int64, dataId int64, text string) (sqlRes, error)
    DeleteText(ipsumId int64, dataId int64) (sqlRes, error)
}
