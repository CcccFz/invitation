package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thoas/go-funk"
)

const (
	formatDate = "2006-01-02"
)

var (
	db *gorm.DB
)

type dbConfig struct {
	Product  string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func initDB(c *dbConfig) {
	var err error
	var config string

	switch c.Product {
	case "postgres":
		config = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name, c.Password,
		)
	case "mysql":
		config = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Name,
		)
	default:
		panic("不支持的数据库: " + c.Product)
	}

	if db, err = gorm.Open(c.Product, config); err != nil {
		panic(err)
	}

	db.AutoMigrate(&invitation{})

	return
}

const (
	TypeOut    = "赶礼"
	TypeIn     = "收礼"
	TypeOutRed = "发红包"
	TypeInRed  = "收红包"
	TypeNotOut = "未赶礼"
	TypeNotIn  = "未收礼"

	CategoryWedding    = "结婚"
	CategoryBaby       = "满月"
	CategoryGraduation = "升学"
	CategoryFuneral    = "葬礼"
	CategoryBirthDay   = "生日"
	CategoryHouse      = "新房"
	CategoryIll        = "探病"
)

const (
	IdxName = iota
	IdxMoney
	IdxNote
	IDxType
	IdxCategory
	IdxAt
)

var (
	categories = []string{
		CategoryWedding, CategoryBaby, CategoryGraduation,
		CategoryFuneral, CategoryBirthDay, CategoryHouse, CategoryIll,
	}
)

// 邀请
type invitation struct {
	Name  string `gorm:"primary_key;type:varchar(10);not null"` // 姓名
	Money int    `gorm:"type:integer;not null"`                 // 礼金。0表示，拒绝邀请；大于0，表示收礼；小于0，表示赶礼

	Note     string    `gorm:"type:varchar(30);not null"`            // 备注
	Type     string    `gorm:"primary_key;type:varchar(5);not null"` // 类型
	Category string    `gorm:"primary_key;type:varchar(5);not null"` // 类别
	At       time.Time `gorm:"type:timestamp with time zone"`        // 时间

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func newInvitation(values []string) *invitation {
	name := values[IdxName]
	money, _ := strconv.Atoi(values[IdxMoney])
	note := values[IdxNote]
	typ := values[IDxType]
	category := values[IdxCategory]
	at := values[IdxAt]

	if strings.Contains(name, " ") {
		panic(fmt.Errorf("name不能有空格: %s", name))
	}
	if strings.Contains(note, " ") {
		panic(fmt.Errorf("note不能有空格: %s", note))
	}
	if !funk.ContainsString(categories, category) {
		panic(fmt.Errorf("不支持的类别: %s", category))
	}
	switch typ {
	case TypeOut, TypeIn, TypeInRed, TypeOutRed:
		if money <= 0 {
			panic(fmt.Errorf("%s: %d元", typ, money))
		}
	case TypeNotIn, TypeNotOut:
		if money != 0 {
			panic(fmt.Errorf("%s: %d元", typ, money))
		}
	default:
		panic(fmt.Errorf("不支持的类型: %s", typ))
	}
	_at, err := time.ParseInLocation(formatDate, at, time.Local)
	if err != nil {
		panic(err)
	}

	v := new(invitation)
	v.Name = name
	v.Money = money
	v.Note = note
	v.Type = typ
	v.Category = category
	v.At = _at

	if v.Type == TypeOut || v.Type == TypeOutRed {
		v.Money = -v.Money
	}

	return v
}

func creates(invitations []*invitation) {
	var err error

	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	if err = tx.Error; err != nil {
		panic(err)
	}

	for i := range invitations {
		if invitations[i] == nil {
			continue
		}

		if err = tx.Create(invitations[i]).Error; err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err = tx.Commit().Error; err != nil {
		panic(err)
	}
}

func findAll() (invitations []invitation) {
	err := db.Find(&invitations).Error
	if err != nil {
		panic(err)
	}
	return
}
