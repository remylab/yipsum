package common

import (
    "os"
    "time"
    "math/rand"
)

func GetRootPath() string {
    return os.Getenv("yip_root")
}

func RandomString(strlen int) string {
    rand.Seed(time.Now().UTC().UnixNano())
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
    result := make([]byte, strlen)
    for i := 0; i < strlen; i++ {
        result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}