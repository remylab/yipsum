package common

import (
	"os"
    "os/exec"
)

func createDb(targetDb string) {
	    // cat conf/evol/createdb.sql | sqlite3 work/yipsum.db
 	c1 := exec.Command("cat", os.Getenv("yip_root")+"/conf/evol/fulldb.sql")
    c2 := exec.Command("sqlite3", targetDb)
    c2.Stdin, _ = c1.StdoutPipe()
    c2.Stdout = os.Stdout
    _ = c2.Start()
    _ = c1.Run()
    _ = c2.Wait()
}

func LoadTestData(targetDb string, dataPath string) {

	createDb(targetDb)

 	c1 := exec.Command("cat", dataPath)
    c2 := exec.Command("sqlite3", targetDb)
    c2.Stdin, _ = c1.StdoutPipe()
    c2.Stdout = os.Stdout
    _ = c2.Start()
    _ = c1.Run()
    _ = c2.Wait()

}