package db

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "os"
    "strings"
    "strconv"
    "database/sql"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
    "github.com/remylab/yipsum/common"
    "github.com/labstack/gommon/random"
)

type (
    SqliteManager struct {
        db *sql.DB
    }

    sqlRes struct {
        Ok bool `json:"ok"`
        Msg string `json:"msg"`
    }
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func AfterDbTest(dbm *SqliteManager, dbpath string) func() {
    return func() {
        dbm.Close()
        err := os.Remove(dbpath)
        if err!=nil { fmt.Printf("Cannot remove test db %v : %v\n",dbpath,err) }
    }
}

func NewSqliteManager(dbPath string) (*SqliteManager, error) {

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {return nil, err}

    return &SqliteManager{db}, nil
}

func (m *SqliteManager) Close() error {
    return m.db.Close()
}


func (m *SqliteManager) RemoveResetToken(ipsum string) (error) {

    stmt, err := m.db.Prepare("UPDATE ipsums set resetToken=null, resetTS=0 where uri=?")
    if err != nil {return err}
    defer stmt.Close()

    _ , execErr := stmt.Exec(ipsum)
    return execErr
}

func (m *SqliteManager) UpdateToken(tokenField string, ipsumId int64) (sqlRes, error) {
    ret := sqlRes{false,""}

    // resetToken / resetTS or deleteToken / deleteTS
    stmt, err := m.db.Prepare("UPDATE ipsums set "+tokenField+"Token=?, "+tokenField+"TS=?  where id=?")
    if err != nil {return ret,err}
    defer stmt.Close()

    token := random.String(64)
    created := common.GetTimestamp()

    res, err := stmt.Exec(token, created, ipsumId)
    if err != nil {return ret,err}

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1), token}, err
}
func (m *SqliteManager) DeleteText(ipsumId int64, dataId int64) (sqlRes, error) {

    ret := sqlRes{false,""}

    stmt, err := m.db.Prepare("DELETE from ipsumtext where ipsum_id=? and id= ?")
    if err != nil {return ret,err}
    defer stmt.Close()

    res, err := stmt.Exec(ipsumId, dataId)
    if err != nil {return ret,err}

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),""}, err
}

func (m *SqliteManager) UpdateText(ipsumId int64, dataId int64, text string) (sqlRes, error) {

    ret := sqlRes{false,""}

    stmt, err := m.db.Prepare("UPDATE ipsumtext set data=? where ipsum_id=? and id= ?")
    if err != nil {return ret,err}
    defer stmt.Close()

    escText := template.HTMLEscapeString( strings.TrimSpace(text) )
    res, err := stmt.Exec(escText, ipsumId, dataId)
    if err != nil {return ret,err}

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1), ""}, err
}

func (m *SqliteManager) AddText(ipsumId int64, text string) (sqlRes, error) {

    ret := sqlRes{false,""}

    total, _ := m.GetTotalIpsumTexts(ipsumId) 
    if (total >= 1000) {
        ret.Msg = "too_many"
        return ret, nil
    }

    _, err :=  m.db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {return ret,err}

    stmt, err := m.db.Prepare("INSERT INTO ipsumtext(ipsum_id,data,created) VALUES(?,?,?)")
    if err != nil {return ret,err}
    defer stmt.Close()

    created := common.GetTimestamp()

    escText := template.HTMLEscapeString( strings.TrimSpace(text) )
    res, err := stmt.Exec(ipsumId, escText, created)
    if err != nil {
        sqliteError := err.Error()
        i := strings.Index(sqliteError,"FOREIGN KEY constraint failed")
        if ( i > -1 ) {
            ret.Msg = "unknown"
        }
        return ret, err
    }

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    id, err := res.LastInsertId() 
    if err != nil {return ret,err}
    
    sId := strconv.FormatInt(id, 10)
    return sqlRes{(rowCnt==1), sId}, err
}

func (m *SqliteManager)GetTotalIpsums() (int,error) {

    stmt, err := m.db.Prepare("select count(1) total from ipsums")
    if err != nil {return 0, err}
    defer stmt.Close()

    var count int
    err = stmt.QueryRow().Scan(&count)
    if err != nil {return 0, err}

    return count, nil
}

func (m *SqliteManager)GetIpsumsForPage(pageNum int64, resByPage int64) ([]map[string]string, error) {

    q := "select name, desc, uri from ipsums where exists (select count(1) c from ipsumtext"
    q += " where ipsums.id=ipsumtext.ipsum_id group by ipsumtext.ipsum_id having c > 4 ) order by id desc limit ? offset ?"
    
    stmt, err := m.db.Prepare(q)
    if err != nil {return nil, err}
    defer stmt.Close()

    nbOffset := (pageNum-1) * resByPage
    rows, err := stmt.Query(resByPage, nbOffset)
    if err != nil {return nil, err}
    defer rows.Close()

    var name , desc, uri string

    ipsums := []map[string]string{}
    for rows.Next() {
        err := rows.Scan(&name, &desc, &uri)
        if err != nil {return nil, err}
        t := map[string]string{
            "name": name,
            "desc": desc,
            "url": common.GetDomain() + "/" + uri,
        }
        ipsums = append(ipsums,t)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return ipsums, nil
}

func (m *SqliteManager)GetTotalIpsumTexts(ipsumId int64) (int,error) {

    stmt, err := m.db.Prepare("select count(1) total from ipsumtext where ipsum_id = ?")
    if err != nil {return 0, err}
    defer stmt.Close()

    var count int
    err = stmt.QueryRow(ipsumId).Scan(&count)
    if err != nil {return 0, err}

    return count, nil
}

func (m *SqliteManager)GetIpsumTextsForPage(ipsumId int64, pageNum int64, resByPage int64) ([]map[string]string, error) {

    stmt, err := m.db.Prepare("select id, data from ipsumtext where ipsum_id = ? order by id desc limit ? offset ?;")
    if err != nil {return nil, err}
    defer stmt.Close()

    nbOffset := (pageNum-1) * resByPage
    rows, err := stmt.Query(ipsumId, resByPage, nbOffset)
    if err != nil {return nil, err}
    defer rows.Close()

    var id , data string

    texts := []map[string]string{}
    nbTexts := 0
    for rows.Next() {
        err := rows.Scan(&id, &data)
        if err != nil {return nil, err}
        t := map[string]string{
            "id": id,
            "text": data,
        }
        texts = append(texts,t)
        nbTexts++
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return texts, nil
}

func (m *SqliteManager) ProcessDeleteAction(ipsum string, token string) (sqlRes, error) {

    res := sqlRes{false,"invalid_token"}

    stmt, err := m.db.Prepare("select deleteTS from ipsums where uri = ? and deleteToken = ?")
    if err != nil {return res, err}
    defer stmt.Close()

    var ts int64
    err = stmt.QueryRow(ipsum,token).Scan(&ts)
    if err != nil { return res, err }

    now := common.GetTimestamp()
    deltaSec := now - ts

    if (deltaSec > 3600) {
        res.Msg = "token_expired"
    } else {

        res.Msg = "internal_error"

        ipsumMap, getErr := m.GetIpsum( ipsum )
        if getErr != nil { return res, getErr }
        ipsumId, _ := strconv.ParseInt(ipsumMap["id"], 10, 32)

        stmt, err = m.db.Prepare("Delete from ipsumtext where ipsum_id=?")
        if err != nil {  return res, err} 

        delRes, delErr := stmt.Exec(ipsumId)
        if delErr !=nil { return res, delErr} 

        stmt, err = m.db.Prepare("Delete from ipsums where id=?")
        if err != nil {  return res, err} 

        delRes, delErr = stmt.Exec(ipsumId)
        if delErr !=nil { return res, delErr} 

        rowCount, _ := delRes.RowsAffected()
        if (rowCount == 1) {
            res.Ok = true
        }
    }

    return res, nil
}

func (m *SqliteManager) ProcessResetAction(ipsum string, token string) (sqlRes, error) {

    res := sqlRes{false,"invalid_token"}

    stmt, err := m.db.Prepare("select resetTS from ipsums where uri = ? and resetToken = ?")
    if err != nil {return res, err}
    defer stmt.Close()

    var ts int64
    err = stmt.QueryRow(ipsum,token).Scan(&ts)
    if err != nil { return res, err }

    now := common.GetTimestamp()
    deltaSec := now - ts

    if (deltaSec > 3600) {
        res.Msg = "token_expired"
    } else {

        res.Msg = "internal_error"

        stmt, err = m.db.Prepare("UPDATE ipsums set adminKey=? where uri=?")
        if err != nil {  return res, err} 
        defer stmt.Close()

        newKey := common.RandomString( common.GetAdminKeyLength() )
        upRes, upErr := stmt.Exec(newKey, ipsum)
        if upErr !=nil { return res, upErr} 

        rowCount, _ := upRes.RowsAffected()

        if (rowCount == 1) {
            res.Ok = true
            res.Msg = newKey
        }
    }

    return res, nil
}

func (m *SqliteManager) GetIpsum(s string) (map[string]string, error) {

    ipsumMap := map[string]string{
        "id": "",
        "name": "",
        "desc": "",
        "adminEmail":"",
    }

    stmt, err := m.db.Prepare("select id, name, desc, adminEmail, resetTS, deleteTS from ipsums where uri = ?")
    if err != nil {return ipsumMap, err}
    defer stmt.Close()

    var s1, s2, s3, s4, s5, s6 sql.NullString
    err = stmt.QueryRow(s).Scan(&s1,&s2,&s3,&s4,&s5,&s6)
    if err != nil {return ipsumMap, err}
    
    if ( s1.Valid ) { ipsumMap["id"] = s1.String };
    if ( s2.Valid ) { ipsumMap["name"] = s2.String }; 
    if ( s3.Valid ) { ipsumMap["desc"] = s3.String }; 
    if ( s4.Valid ) { ipsumMap["adminEmail"] = s4.String }; 
    if ( s5.Valid ) { ipsumMap["resetTS"] = s5.String }; 
    if ( s6.Valid ) { ipsumMap["deleteTS"] = s6.String }; 

    return ipsumMap, nil
}

func (m *SqliteManager) ValidateUriKey(ipsum string, key string) (bool,error) {

    stmt, err := m.db.Prepare("select count(1) from ipsums where uri = ? and adminKey = ?")
    if err != nil {return false, err}
    defer stmt.Close()

    var count int
    err = stmt.QueryRow(ipsum,key).Scan(&count)
    if err != nil { return false, err }

    //fmt.Printf("ValidateUriKey %v/%v =%v\n",ipsum, key,(count==1))
    return (count==1), nil
}


func (m *SqliteManager) IsNewUri(s string) (bool,error) {

    stmt, err := m.db.Prepare("select count(1) from ipsums where uri = ?")
    if err != nil {return false, err}
    defer stmt.Close()

    var count int
    err = stmt.QueryRow(s).Scan(&count)
    if err != nil { return false, err }

    //fmt.Printf("counter from db :%v",counter)
    return (count==0), nil
}

func (m *SqliteManager) CreateIpsum(name string, desc string, uri string, adminEmail string) (sqlRes, error) {

    ret := sqlRes{false,""}
    adminKey := common.RandomString( common.GetAdminKeyLength() )

    name = template.HTMLEscapeString(name)
    uri = common.GetUri(uri)
    desc = template.HTMLEscapeString(desc)
    adminEmail = template.HTMLEscapeString(adminEmail)



    stmt, err := m.db.Prepare("INSERT INTO ipsums(name,desc,uri,adminEmail,adminKey,created) VALUES(?,?,?,?,?,?)")
    if err != nil {return ret,err}
    defer stmt.Close()

    created := common.GetTimestamp()
    res, err := stmt.Exec(name,desc,uri,adminEmail,adminKey,created)
    if err != nil {
        sqliteError := err.Error()
        i := strings.Index(sqliteError,"UNIQUE constraint failed")
        if ( i > -1 ) {
            ret.Msg = "taken"
        }
        return ret, err
    }

    rowCnt, err := res.RowsAffected()
    if err != nil {return ret,err}

    return sqlRes{(rowCnt==1),adminKey}, err
}

func (m *SqliteManager)GenerateIpsum(ipsumId int64) ([]string,error) {
    ret := []string{}

    stmt, err := m.db.Prepare("select id, data from ipsumtext where ipsum_id = ? limit 1000;")
    if err != nil {return nil, err}
    defer stmt.Close()

    rows, err := stmt.Query(ipsumId)
    if err != nil {return nil, err}
    defer rows.Close()

    totalLength := 0

    var id , data string
    for rows.Next() {
        err := rows.Scan(&id, &data)
        if err != nil {return ret, err}

        totalLength += len(data)
        sentences := common.GetSentences(data)
        //fmt.Printf(" row data ==%v==\n",data)
        ret = append(ret, sentences...)

        //for _, s := range sentences {
        //    fmt.Printf("\n         line ==%v==\n",s)
        //}
    }

    if err = rows.Err(); err != nil {
        return ret, err
    }

    targetLength := 2000

    if (totalLength <= targetLength) {
        return ret, nil
    }

    d := float64(totalLength) / float64(targetLength)
    chances := int(math.Floor(d)) + 1 
    final := []string{}; finalCount := 0

    done := false
    for i := 1; i <= 3; i++ {

        s1 := rand.NewSource(time.Now().UnixNano())
        r1 := rand.New(s1)

        if ( done ) { break }
        for _, s := range ret {
            randN := r1.Intn(chances)
            yes := randN == 0
            if ( yes ) {
                finalCount += len(s)
                final = append(final, s)
                if ( finalCount > targetLength ) {
                    done = true
                    break
                }
            }
        }
    }

    return final, nil
}
