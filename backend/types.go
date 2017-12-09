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
	UpdateAt  time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	Name      string    `gorm:"column:name;size:20;not null;unique_index" json:"name"`
	Countries []Country `gorm:"ForeignKey:TownId" json:"-"`
}

func (Town) TableName() string {
	return "towns"
}

type Country struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt  time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt  time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	Name      string    `gorm:"column:name;size:20;not null" json:"name"`
	TownId    int       `gorm:"column:town_id" json:"town_id"`
	Companies []Company `gorm:"ForeignKey:CountryId" json:"-"`
}

func (Country) TableName() string {
	return "countries"
}

type Company struct {
	ID             int            `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt       time.Time      `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt       time.Time      `gorm:"column:update_at;not null;" json:"update_at"`
	Name           string         `gorm:"column:name;size:60;not null;unique_index" json:"name"`
	Address        string         `gorm:"column:address;size:100;not null" json:"address"`
	CountryId      int            `gorm:"column:country_id" json:"country_id"`
	CountryName    string         `gorm:"-" json:"country_name"`
	Users          []User         `gorm:"ForeignKey:CompanyId" json:"-"`
	MonitorPlaces  []MonitorPlace `gorm:"ForeignKey:CompanyId" json:"-"`
	Summaries      []Summary      `gorm:"ForeignKey:CompanyId" json:"-"`
	TodaySummaries []TodaySummary `gorm:"ForeignKey:CompanyId" json:"-"`
	Enable         string         `gorm:"column:enable;size:1;default:'T';not null" json:"enable"`
}

func (Company) TableName() string {
	return "companies"
}

type User struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt    time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt    time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	Phone       string    `gorm:"column:phone;size:11;not null;unique_index" json:"phone"`
	Name        string    `gorm:"column:name;size:20;not null" json:"name"`
	Password    string    `gorm:"column:password;size:200;not null"`
	Job         string    `gorm:"column:job;size:20" json:"job"`
	CompanyId   int       `gorm:"column:company_id;not null" json:"company_id"`
	CompanyName string    `gorm:"-" json:"company_name"`
	WxOpenId    string    `gorm:"column:wx_openid;size:50;unique_index" json:"wx_openid"`
	Enable      string    `gorm:"column:enable;size:1;default:'T';not null" json:"enable"`
	Pictures    []Picture `gorm:"ForeignKey:UserId" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type MonitorType struct {
	ID            int            `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt      time.Time      `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt      time.Time      `gorm:"column:update_at;not null;" json:"update_at"`
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
	UpdateAt      time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	Name          string    `gorm:"column:name;size:20;not null" json:"name"`
	CompanyId     int       `gorm:"column:company_id;not null;index" json:"company_id"`
	MonitorTypeId int       `gorm:"column:monitor_type_id;index" json:"monitor_type_id"`
	QrcodePath    string    `gorm:"column:qrcode_path" json:"qrcode_path"`
	QrcodeURI     string    `gorm:"column:qrcode_uri" json:"qrcode_uri"`
	Pictures      []Picture `gorm:"ForeignKey:MonitorPlaceId" json:"-"`
}

func (MonitorPlace) TableName() string {
	return "monitor_places"
}

type Picture struct {
	ID             int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt       time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt       time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	MonitorTypeId  int       `gorm:"-" json:"monitor_type_id"`
	MonitorPlaceId int       `gorm:"column:monitor_place_id;not null;index" json:"monitor_place_id"`
	ThumbPath      string    `gorm:"column:thumb_path" json:"thumb_path"`
	FullPath       string    `gorm:"column:full_path" json:"full_path"`
	ThumbURI       string    `gorm:"column:thumb_uri" json:"thumb_uri"`
	FullURI        string    `gorm:"column:full_uri" json:"full_uri"`
	Corrective     string    `gorm:"column:corrective;size:1;not null" json:"corrective"`
	UserId         int       `gorm:"column:user_id;index" json:"user_id"`
}

func (Picture) TableName() string {
	return "pictures"
}

type Summary struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Day         time.Time `gorm:"column:day;not null;unique_index:day_company;default:NOW()" json:"day"`
	CompanyId   int       `gorm:"column:company_id;not null;unique_index:day_company" json:"company_id"`
	CompanyName string    `gorm:"column:company_name;not null" json:"company_name"`
	IsFinish    string    `gorm:"column:is_finish;string;size:1;not null" json:"finish"`
	UnfinishIds string    `gorm:"column:unfinish_ids" json:"unfinish_ids"`
}

func (Summary) TableName() string {
	return "summary"
}

type TodaySummary struct {
	ID               int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Day              time.Time `gorm:"column:day;not null;unique_index:day_company_place;default:NOW()" json:"day"`
	CompanyId        int       `gorm:"column:company_id;not null;unique_index:day_company_place" json:"company_id"`
	CompanyName      string    `gorm:"column:company_name;not null" json:"company_name"`
	MonitorPlaceId   int       `gorm:"column:monitor_place_id;not null;unique_index:day_company_place" json:"monitor_place_id"`
	MonitorPlaceName string    `gorm:"column:monitor_place_name;not null;" json:"monitor_place_name"`
	IsUpload         string    `gorm:"column:is_upload;string;size:1;not null" json:"is_upload"`
	Corrective       string    `gorm:"column:corrective;string;size:1;not null" json:"corrective"`
	EverCorrective   string    `gorm:"column:ever_corrective;string;size:1;not null" json:"ever_corrective"`
}

func (TodaySummary) TableName() string {
	return "today_summary"
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
	db.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Town{}, &Country{}, &Company{}, &User{}, &MonitorType{}, &MonitorPlace{}, &Picture{}, &Summary{}, &TodaySummary{})
	db.Debug().Model(&Country{}).AddForeignKey("town_id", "towns(id)", "SET NULL", "SET NULL")
	db.Debug().Model(&Company{}).AddForeignKey("country_id", "countries(id)", "SET NULL", "SET NULL")
	db.Debug().Model(&User{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Debug().Model(&MonitorPlace{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Debug().Model(&MonitorPlace{}).AddForeignKey("monitor_type_id", "monitor_types(id)", "SET NULL", "CASCADE")
	db.Debug().Model(&Picture{}).AddForeignKey("user_id", "users(id)", "SET NULL", "SET NULL")
	db.Debug().Model(&Picture{}).AddForeignKey("monitor_place_id", "monitor_places(id)", "CASCADE", "CASCADE")
	db.Debug().Model(&Summary{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Debug().Model(&TodaySummary{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
	db.Debug().Model(&TodaySummary{}).AddForeignKey("monitor_place_id", "monitor_places(id)", "CASCADE", "CASCADE")
}
