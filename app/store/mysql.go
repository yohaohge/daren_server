package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"server.com/daren/config"
	"server.com/daren/def"
	"server.com/daren/pkg/crypto"
	"strconv"
	"strings"
	"time"
)

type SqlClient struct {
	c *sqlx.DB
}

var MC *SqlClient
var MC_Dblog *SqlClient

// user 表
// uid_generator 表

func ConnectMysql() {
	// business
	db, err := sqlx.Connect("mysql", config.GetMainMysqlDsn())
	if err != nil {
		log.Fatal("mysql connect error: " + err.Error())
		panic(err)
	}
	db.SetMaxOpenConns(100)
	MC = new(SqlClient)
	MC.c = db

	// log
	dbLog, err := sqlx.Connect("mysql", config.GetLogMysqlDsn())
	if err != nil {
		log.Fatal("mysql connect error: " + err.Error())
		panic(err)
	}
	db.SetMaxOpenConns(100)
	MC_Dblog = new(SqlClient)
	MC_Dblog.c = dbLog
}

func (sc *SqlClient) _GetUserByOpenId(openId string) (*User, error) {
	user := &User{}
	stmt, err := sc.c.Prepare("SELECT * FROM user WHERE openId = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(openId)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.OpenId,
			&user.Uid,
			&user.Name,
			&user.VipTime,
			&user.Version,
			&user.Password,
			&user.Data,
			&user.ClientIP,
			&user.CreateTime,
		)
		if err != nil {
			return user, err
		}
	} else {
		return nil, nil
	}
	return user, nil
}

func (sc *SqlClient) GetUserByOpenId(openId string) (*User, error) {
	user := &User{}
	query := "SELECT * FROM user WHERE openId = ?"
	err := sc.c.Get(user, query, openId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return user, nil
}

func (sc *SqlClient) GetUserByUid(uid int) (*User, error) {
	user := &User{}
	query := "SELECT * FROM user WHERE uid = ?"
	err := sc.c.Get(user, query, uid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return user, nil
}

func (sc *SqlClient) SaveUser(user *User) (int64, error) {
	if user == nil {
		return 0, errors.New("param empty")
	}
	query := "UPDATE user SET name = ?, vipTime = ?, version = ?, password = ?, data = ?, clientIp = ?, createTime = ?, device = ?"
	ret, err := sc.c.Exec(query, user.Name, user.VipTime, user.Version, user.Password, user.Data, user.ClientIP, user.CreateTime, user.Device)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

func (sc *SqlClient) AddUser(user *User) (int64, error) {
	if user == nil {
		return 0, errors.New("param empty")
	}
	query := "INSERT INTO user VALUES (?,?,?,?,?,?,?,?,?,?)"
	ret, err := sc.c.Exec(query, user.OpenId, user.Uid, user.Name, user.VipTime, user.Version, user.Password, user.Data, user.ClientIP, user.CreateTime)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

func (sc *SqlClient) UpdateUser(uid int) {
}

func (sc *SqlClient) GetNewUid() (int64, error) {
	query := "INSERT INTO uid_generator VALUES ()"
	result, err := sc.c.Exec(query)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (sc *SqlClient) SaveSession(uid int, ip string) (string, error) {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(uid))
	sb.WriteString("-")
	sb.WriteString(ip)
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(time.Now().UnixMilli())))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(rand.Int()))

	expire := time.Now().Unix() + def.SessionExpireTime
	sessionId := crypto.Md5(sb.String())
	query := "INSERT INTO user_session (uid, sessionId, expire) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE sessionId = ?, expire = ?"
	_, err := sc.c.Exec(query, uid, sessionId, expire, sessionId, expire)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (sc *SqlClient) GetSession(sid string) *Session {
	sess := &Session{}
	query := "SELECT * from user_session WHERE sessionId = ?"
	err := sc.c.Get(sess, query, sid)

	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	if sess.Expire < time.Now().Unix() {
		query = "DELETE FROM user_session WHERE  sessionId = ?"
		_, err = sc.c.Exec(query, sid)

		return nil
	}

	return sess
}

func (sc *SqlClient) IsSessionValid(sid string, uid int) bool {
	sess := &Session{}
	query := "SELECT * from user_session WHERE uid = ? and sessionId = ?"
	err := sc.c.Get(sess, query, uid, sid)

	if err != nil {
		if err != sql.ErrNoRows {
			logrus.Errorln(err)
		}
		return false
	}
	if sess.Expire < time.Now().Unix() {
		query = "DELETE FROM user_session WHERE  uid = ?"
		_, err = sc.c.Exec(query, uid)
		if err != nil {
			logrus.Errorln(err)
		}
		return false
	}
	return true
}

func (sc *SqlClient) UpdateSession(sid, uid string) bool {
	expire := time.Now().Unix() + def.SessionExpireTime
	query := "UPDATE user_session SET expire = ? WHERE sessionId = ? and uid = ?"

	_, err := sc.c.Exec(query, expire, sid, uid)
	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

func (sc *SqlClient) GetVideListTotal() int {
	var count int
	query := "SELECT COUNT(*) FROM video_content"
	err := sc.c.Get(&count, query)

	if err != nil {
		logrus.Errorln(err)
		return 0
	}

	return count
}

func (sc *SqlClient) GetVideoListWithOutData(pageNo, pageCount int) []*VideoInfo {
	list := make([]*VideoInfo, 0)
	offset := (pageNo - 1) * pageCount
	query := "SELECT id, name, cover, total, `desc`, label FROM video_content LIMIT ? OFFSET ?"

	rows, err := sc.c.Query(query, pageCount, offset)
	if err != nil {
		logrus.Errorln(err)
		return list
	}

	for rows.Next() {
		vi := &VideoInfo{}
		err = rows.Scan(&vi.Id, &vi.Name, &vi.Cover, &vi.Total, &vi.Desc, &vi.Label)
		if err != nil {
			logrus.Errorln(err)
			break
		}
		list = append(list, vi)
	}
	if err = rows.Err(); err != nil {
		logrus.Errorln(err)
	}

	return list
}

func (sc *SqlClient) getVideoDetail(id int) *VideoDetailInfo {
	item := new(VideoDetailInfo)
	query := "SELECT id, name ,data, cover, total, `desc`, label FROM video_content WHERE id = ?"

	vi := &VideoInfo{}
	err := sc.c.Get(vi, query, id)
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.Errorln(err)
		}
		return nil
	}

	if err = json.Unmarshal([]byte(vi.Data), &item.VData); err != nil {
		return nil
	}

	vi.Data = "" // 此项不需要后续优化结构
	item.VideoInfo = vi

	return item
}

func (sc *SqlClient) SearchVideo(name string, pageNo, pageCount int) (int, []*VideoInfo) {
	list := make([]*VideoInfo, 0)
	count := 0
	query := "SELECT COUNT(*) FROM video_content WHERE name like ?"
	err := sc.c.Get(&count, query, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		logrus.Errorln(err)
		return count, list
	}
	if count == 0 {
		return count, list
	}

	offset := (pageNo - 1) * pageCount
	query = "SELECT id, name, cover, total, `desc`, label FROM video_content WHERE name like ? LIMIT ? OFFSET ?"

	rows, err := sc.c.Query(query, fmt.Sprintf("%%%s%%", name), pageCount, offset)
	if err != nil {
		logrus.Errorln(err)
		return count, list
	}

	for rows.Next() {
		vi := &VideoInfo{}
		err = rows.Scan(&vi.Id, &vi.Name, &vi.Cover, &vi.Total, &vi.Desc, &vi.Label)
		if err != nil {
			logrus.Errorln(err)
			break
		}
		list = append(list, vi)
	}
	if err = rows.Err(); err != nil {
		logrus.Errorln(err)
	}

	return count, list
}

func (sc *SqlClient) IsInvalidCDKEY(cdkey string) bool {
	var cdkeyDb string
	query := "SELECT cdkey FROM cdkey_config WHERE cdkey = ? AND (num > 0 OR cdkeyType = ?) AND expireTime > ?"
	err := sc.c.Get(&cdkeyDb, query, cdkey, def.CDKEY_TYPE_ANYONE, time.Now().Unix())
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.Errorln(err)
		}
		return false
	}

	return cdkeyDb != ""
}

func (sc *SqlClient) GetCDKEYInfo(cdkey string) (bool, *CDKEYInfo) {
	ret := &CDKEYInfo{}
	query := "SELECT * FROM cdkey_config WHERE cdkey = ? AND (num > 0 OR cdkeyType = ?) AND expireTime > ?"
	err := sc.c.Get(ret, query, cdkey, def.CDKEY_TYPE_ANYONE, time.Now().Unix())
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.Errorln(err)
		}
		return false, ret
	}
	if ret.Cdkey == "" {
		return false, ret
	}
	return true, ret
}

func (sc *SqlClient) IsUserUsed(uid int, cdkey string, cdkeyType int) bool {
	var cdkeyDb string
	var query string
	var err error
	if cdkeyType == def.CDKEY_TYPE_NORMAL {
		query = "SELECT cdkey FROM user_cdkey_use WHERE cdkey = ?"
		err = sc.c.Get(&cdkeyDb, query, cdkey)
	} else if cdkeyType == def.CDKEY_TYPE_ANYONE {
		query = "SELECT cdkey FROM user_cdkey_use WHERE uid = ? AND cdkey = ?"
		err = sc.c.Get(&cdkeyDb, query, uid, cdkey)
	} else {
		return true
	}
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.Errorln(err)
			return true
		} else {
			return false
		}
	}
	return cdkeyDb != ""
}

func UseCDKEY(uid int, cdkey string) bool {
	query := "UPDATE cdkey_config SET num = num - 1 WHERE cdkey = ? AND num > 0 AND cdkeyType != ?"
	result, err := MC.c.Exec(query, cdkey, def.CDKEY_TYPE_ANYONE)
	if err != nil {
		logrus.Errorln(err)
		return false
	}

	count, err := result.RowsAffected()
	if err != nil {
		logrus.Errorln(err)
		return false
	}

	// 使用记录
	query = "INSERT INTO user_cdkey_use (uid, cdkey, createTime) VALUES (?,?,?)"
	_, err = MC.c.Exec(query, uid, cdkey, time.Now().Unix())
	if err != nil {
		logrus.Errorln(err)
	}

	return count > 0
}
