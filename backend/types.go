package main

import (
	//	"github.com/emicklei/go-restful"
	"fmt"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

type Town struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt  time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt  time.Time `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Name      string    `gorm:"column:name;size:20;not null;unique_index" json:"name"`
	Countries []Country `gorm:"ForeignKey:TownId" json:"-"`
}

func (Town) TableName() string {
	return "towns"
}

type Country struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt  time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt  time.Time `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Name      string    `gorm:"column:name;size:20;not null" json:"name"`
	TownId    int       `gorm:"column:town_id" json:"town_id"`
	Companies []Company `gorm:"ForeignKey:CountryId" json:"-"`
}

func (Country) TableName() string {
	return "countries"
}

type Company struct {
	ID            int            `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt      time.Time      `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt      time.Time      `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Name          string         `gorm:"column:name;size:60;not null;unique_index" json:"name"`
	Address       string         `gorm:"column:address;size:100;not null" json:"address"`
	CountryId     int            `gorm:"column:country_id" json:"country_id"`
	Users         []User         `gorm:"ForeignKey:CompanyId" json:"-"`
	MonitorPlaces []MonitorPlace `gorm:"ForeignKey:CompanyId" json:"-"`
	Summaries     []Summary      `gorm:"ForeignKey:CompanyId" json:"-"`
}

func (Company) TableName() string {
	return "companies"
}

type User struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt  time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt  time.Time `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Phone     string    `gorm:"column:phone;size:11;not null;unique_index" json:"phone"`
	Name      string    `gorm:"column:name;size:20;not null" json:"name"`
	Password  string    `gorm:"column:password;size:30;not null"`
	Job       string    `gorm:"column:job;size:20" json:"job"`
	CompanyId int       `gorm:"column:company_id;not null" json:"company_id"`
	WxOpenId  string    `gorm:"column:wx_openid;size:50" json:"wx_openid"`
	Enable    string    `gorm:"column:enable;size:1;not null" json:"enable"`
	Pictures  []Picture `gorm:"ForeignKey:UserId" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type MonitorType struct {
	ID            int            `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt      time.Time      `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt      time.Time      `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Name          string         `gorm:"column:name;size:20;not null;unique_index" json:"name"`
	Comment       string         `gorm:"column:comment" json:"comment"`
	MonitorPlaces []MonitorPlace `gorm:"ForeignKey:MonitorTypeId" json:"-"`
}

func (MonitorType) TableName() string {
	return "monitor_types"
}

type MonitorPlace struct {
	ID            int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt      time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt      time.Time `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	Name          string    `gorm:"column:name;size:20;not null" json:"name"`
	CompanyId     int       `gorm:"column:company_id;not null;index" json:"company_id"`
	MonitorTypeId int       `gorm:"column:monitor_type_id;index" json:"monitor_type_id"`
	Qrcode        string    `gorm:"column:qrcode" json:"-"`
	QrcodePath    string    `gorm:"-" json:"qrcode_path"`
	Pictures      []Picture `gorm:"ForeignKey:MonitorPlaceId" json:"-"`
}

func (MonitorPlace) TableName() string {
	return "monitor_places"
}

type Picture struct {
	ID             int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt       time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt       time.Time `gorm:"column:update_at;not null;default:NOW()" json:"update_at"`
	MonitorPlaceId int       `gorm:"column:monitor_place_id;not null;index" json:"monitor_place_id"`
	Thumb          string    `gorm:"column:thumb" json:"-"`
	Full           string    `gorm:"column:full" json:"-"`
	ThumbPath      string    `gorm:"-" json:"thumb_path"`
	FullPath       string    `gorm:"-" json:"full_path"`
	Corrective     string    `gorm:"column:corrective;size:1;not null" json:"corrective"`
	UserId         int       `gorm:"column:user_id;index" json:"user_id"`
}

func (Picture) TableName() string {
	return "pictures"
}

type Summary struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Day         time.Time `gorm:"column:day;not null;index" json:"day"`
	CompanyId   int       `gorm:"column:company_id;not null;index" json:"company_id"`
	Finish      string    `gorm:"column:string;size:1;not null" json:"finish"`
	UnfinishIds string    `gorm:"column:unfinish_ids" json:"unfinish_ids"`
}

func (Summary) TableName() string {
	return "summary"
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func InitializeDB() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/mpmanager?charset=utf8&parseTime=True&loc=Local", dbuser, dbpass, dbip, dbport))
	if err != nil {
		glog.Fatalf("cannot initialize database connection, err %s", err)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Town{}, &Country{}, &Company{}, &User{}, &MonitorType{}, &MonitorPlace{}, &Picture{}, &Summary{})
	db.Model(&Country{}).AddForeignKey("town_id", "towns(id)", "SET NULL", "SET NULL")
	db.Model(&Company{}).AddForeignKey("country_id", "countries(id)", "SET NULL", "SET NULL")
	db.Model(&User{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Model(&MonitorPlace{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Model(&MonitorPlace{}).AddForeignKey("monitor_type_id", "monitor_types(id)", "SET NULL", "CASCADE")
	db.Model(&Picture{}).AddForeignKey("user_id", "users(id)", "SET NULL", "SET NULL")
	db.Model(&Picture{}).AddForeignKey("monitor_place_id", "monitor_places(id)", "CASCADE", "CASCADE")
	db.Model(&Summary{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
}
