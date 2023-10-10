package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Name          string
	From          string
	IPv4          string
	IPv6          string
	Received      int64
	Sent          int64
	TotalReceived int64
	TotalSent     int64
	Status        bool
	Last          string
}

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(Env.GetString("db")))
	if err != nil {
		Logger.Fatalln(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		Logger.Fatalln(err)
	}
}

func Add(cli User) (bool, User) {
	var user User
	res := db.Model(&User{}).Where("name = ?", cli.Name).Limit(1).Find(&user)
	if res.RowsAffected != 0 {
		Logger.Warning(res.Error)
		return false, User{}
	} else {
		cli.Last = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 PM")
		if err := db.Model(&User{}).Create(&cli).Error; err != nil {
			Logger.Error(err)
			return false, User{}
		}
		return true, cli
	}
}

func Update(cli User) (bool, User) {
	var user User
	res := db.Model(&User{}).Where("name = ?", cli.Name).Limit(1).Find(&user)
	if res.RowsAffected != 1 {
		ok, added := Add(cli)
		if !ok {
			Logger.Warning("Add Client Error")
			return false, User{}
		}
		return true, added
	} else {
		cli.Last = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05 PM")
		cli.TotalReceived = user.TotalReceived
		cli.TotalSent = user.TotalSent
		res = db.Model(&User{}).Where("name = ?", cli.Name).Limit(1).Select("*").Updates(cli)
		if res.RowsAffected != 1 {
			Logger.Warning("Update Client Error")
			return false, User{}
		}
		res = db.Model(&User{}).Where("name = ?", cli.Name).Limit(1).Find(&cli)
		if res.RowsAffected != 1 {
			Logger.Warning("Update Client Error")
			return false, User{}
		}
		return true, cli
	}
}

func Offline(online []string) bool {
	_ = db.Model(&User{}).Not(online).Update("status", false)
	return true
}

func Sum() {
	var users []User
	_ = db.Model(&User{}).Where("status = ?", false).Find(&users)
	if len(users) != 0 {
		for k, user := range users {
			users[k].TotalSent += user.Sent
			users[k].TotalReceived += user.Received
			users[k].Sent = 0
			users[k].Received = 0
			res := db.Model(&User{}).Where("name = ?", user.Name).Limit(1).Select("*").Updates(users[k])
			if res.RowsAffected != 1 {
				Logger.Warning("Update Client Error")
			}
		}
	}
}

func GetData() (bool, []User) {
	var users []User
	res := db.Model(&User{}).Find(&users)
	if res.Error != nil {
		Logger.Warning(res.Error)
		return false, []User{}
	}
	return true, users
}
