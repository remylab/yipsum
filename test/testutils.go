package test

import (
    "os"
    "os/exec"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    
    "github.com/remylab/yipsum/common"
)

func GetEcho() *echo.Echo {

    e := echo.New()
    e.Pre(middleware.RemoveTrailingSlash())
    e.SetRenderer(common.GetTemplate())
    return e
}

func ImportData(targetDb string, script string) error {

    var err error
    // cat conf/evol/createdb.sql | sqlite3 work/yipsum.db 
    // sqlite3 work/yipsum.db < conf/evol/createdb.sql
    c1 := exec.Command("cat", common.GetRootPath()+script)
    c2 := exec.Command("sqlite3", targetDb)
    c2.Stdin, err = c1.StdoutPipe()
    if (err != nil) { return err}
    c2.Stdout = os.Stdout
    err = c2.Start()
    if (err != nil) { return err}
    err = c1.Run()
    if (err != nil) { return err}
    err = c2.Wait()
    return err
}

func LoadTestData(targetDb string, dataPath string) {

    _ = ImportData(targetDb,"/conf/evol/fulldb.sql")

    if (dataPath != "") {
        c1 := exec.Command("cat", dataPath)
        c2 := exec.Command("sqlite3", targetDb)
        c2.Stdin, _ = c1.StdoutPipe()
        c2.Stdout = os.Stdout
        _ = c2.Start()
        _ = c1.Run()
        _ = c2.Wait()
    }

}
