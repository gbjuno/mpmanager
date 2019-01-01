package main

import (
	//	"github.com/emicklei/go-restful"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	//全部
	TotalDaysAll               int     `gorm:"column:total_days_all;default:0;" json:"total_days_all"`
	RelaxDaysAll               int     `gorm:"column:relax_days_all;default:0;" json:"relax_days_all"`
	MaxContinuousFinishDaysAll int     `gorm:"column:max_continuous_finish_days_all;default:0;" json:"max_continuous_finish_days_all"`
	ContinuousFinishDaysAll    int     `gorm:"column:continuous_finish_days_all;default:0;" json:"continuous_finish_days_all"`
	FinishDaysAll              int     `gorm:"column:finish_days_all;default:0;" json:"finish_days_all"`
	FinishPercentageAll        float64 `gorm:"column:finish_percentage_all;default:0;" json:"finish_percentage_all"`
	//最近365天
	TotalDaysInLast365days               int     `gorm:"column:total_days_last_365_days;default:0;" json:"total_days_last_365_days"`
	RelaxDaysInLast365days               int     `gorm:"column:relax_days_last_365_days;default:0;" json:"relax_days_last_365_days"`
	MaxContinuousFinishDaysInLast365days int     `gorm:"column:max_continuous_finish_days_last_365_days;default:0;" json:"max_continuous_finish_days_last_365_days"`
	ContinuousFinishDaysInLast365days    int     `gorm:"column:continuous_finish_days_last_365_days;default:0;" json:"continuous_finish_days_last_365_days"`
	FinishDaysInLast365days              int     `gorm:"column:finish_days_last_365_days;default:0;" json:"finish_days_last_365_days"`
	FinishPercentageInLast365days        float64 `gorm:"column:finish_percentage_last_365_days;default:0;" json:"finish_percentage_last_365_days"`
	//最近182天
	TotalDaysInLast182days               int     `gorm:"column:total_days_last_182_days;default:0;" json:"total_days_last_182_days"`
	RelaxDaysInLast182days               int     `gorm:"column:relax_days_last_182_days;default:0;" json:"relax_days_last_182_days"`
	MaxContinuousFinishDaysInLast182days int     `gorm:"column:max_continuous_finish_days_last_182_days;default:0;" json:"max_continuous_finish_days_last_182_days"`
	ContinuousFinishDaysInLast182days    int     `gorm:"column:continuous_finish_days_last_182_days;default:0;" json:"continuous_finish_days_last_182_days"`
	FinishDaysInLast182days              int     `gorm:"column:finish_days_last_182_days;default:0;" json:"finish_days_last_182_days"`
	FinishPercentageInLast182days        float64 `gorm:"column:finish_percentage_last_182_days;default:0;" json:"finish_percentage_last_182_days"`
	//最近90天
	TotalDaysInLast90days               int     `gorm:"column:total_days_last_90_days;default:0;" json:"total_days_last_90_days"`
	RelaxDaysInLast90days               int     `gorm:"column:relax_days_last_90_days;default:0;" json:"relax_days_last_90_days"`
	MaxContinuousFinishDaysInLast90days int     `gorm:"column:max_continuous_finish_days_last_90_days;default:0;" json:"max_continuous_finish_days_last_90_days"`
	ContinuousFinishDaysInLast90days    int     `gorm:"column:continuous_finish_days_last_90_days;default:0;" json:"continuous_finish_days_last_90_days"`
	FinishDaysInLast90days              int     `gorm:"column:finish_days_last_90_days;default:0;" json:"finish_days_last_90_days"`
	FinishPercentageInLast90days        float64 `gorm:"column:finish_percentage_last_90_days;default:0;" json:"finish_percentage_last_90_days"`
	//最近30天
	TotalDaysInLast30days               int     `gorm:"column:total_days_last_30_days;default:0;" json:"total_days_last_30_days"`
	RelaxDaysInLast30days               int     `gorm:"column:relax_days_last_30_days;default:0;" json:"relax_days_last_30_days"`
	MaxContinuousFinishDaysInLast30days int     `gorm:"column:max_continuous_finish_days_last_30_days;default:0;" json:"max_continuous_finish_days_last_30_days"`
	ContinuousFinishDaysInLast30days    int     `gorm:"column:continuous_finish_days_last_30_days;default:0;" json:"continuous_finish_days_last_30_days"`
	FinishDaysInLast30days              int     `gorm:"column:finish_days_last_30_days;default:0;" json:"finish_days_last_30_days"`
	FinishPercentageInLast30days        float64 `gorm:"column:finish_percentage_last_30_days;default:0;" json:"finish_percentage_last_30_days"`
}

func (Company) TableName() string {
	return "companies"
}

//每月的完成率报告
type CompanyMonthStat struct {
	UpdateAt                         time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	CompanyID                        int       `gorm:"column:company_id;primary_key;not null" json:"company_id"`
	Date                             time.Time `gorm:"column:date;primary_key;not null;" json:"date"`
	CompanyName                      string    `gorm:"-" json:"company_name"`
	TotalDays                        int       `gorm:"column:total_days;default:0;" json:"total_days"`
	RelaxDays                        int       `gorm:"column:relax_days;default:0;" json:"relax_days"`
	MaxContinuousFinishDaysThisMonth int       `gorm:"column:max_continuous_finish_days_this_month;default:0;" json:"max_continuous_finish_days_this_month"`
	ContinuousFinishDaysThisMonth    int       `gorm:"column:continuous_finish_days_this_month;default:0;" json:"continuous_finish_days_this_month"`
	FinishDaysThisMonth              int       `gorm:"column:finish_days_this_month;default:0;" json:"finish_days_this_month"`
	FinishPercentageThisMonth        float64   `gorm:"column:finish_percentage_this_month;default:0;" json:"finish_percentage_this_month"`
}

func (CompanyMonthStat) TableName() string {
	return "company_month_stat"
}

//每年的完成率报告
type CompanyYearStat struct {
	UpdateAt                        time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	CompanyID                       int       `gorm:"column:company_id;primary_key;not null" json:"company_id"`
	Date                            time.Time `gorm:"column:date;primary_key;not null;" json:"date"`
	CompanyName                     string    `gorm:"-" json:"company_name"`
	TotalDays                       int       `gorm:"column:total_days;default:0;" json:"total_days"`
	RelaxDays                       int       `gorm:"column:relax_days;default:0;" json:"relax_days"`
	MaxContinuousFinishDaysThisYear int       `gorm:"column:max_continuous_finish_days_this_year;default:0;" json:"max_continuous_finish_days_this_year"`
	ContinuousFinishDaysThisYear    int       `gorm:"column:continuous_finish_days_this_year;default:0;" json:"continuous_finish_days_this_year"`
	FinishDaysThisYear              int       `gorm:"column:finish_days_this_year;default:0;" json:"finish_days_this_year"`
	FinishPercentageThisYear        float64   `gorm:"column:finish_percentage_this_year;default:0;" json:"finish_percentage_this_year`
}

func (CompanyYearStat) TableName() string {
	return "company_year_stat"
}

//全局公共假期
type GlobalRelaxPeriod struct {
	ID       int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	StartAt  time.Time `gorm:"column:start_at;not null;" json:"start_at"`
	EndAt    time.Time `gorm:"column:end_at;not null;" json:"end_at"`
}

func (GlobalRelaxPeriod) TableName() string {
	return "global_relax_period"
}

//公司休息假期的时间
type CompanyRelaxPeriod struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt    time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt    time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	CompanyId   int       `gorm:"column:company_id;not null" json:"company_id"`
	CompanyName string    `gorm:"-" json:"company_name"`
	StartAt     time.Time `gorm:"column:start_at;not null;" json:"start_at"`
	EndAt       time.Time `gorm:"column:end_at;not null;" json:"end_at"`
}

func (CompanyRelaxPeriod) TableName() string {
	return "company_relax_period"
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
	WxOpenId    *string   `gorm:"column:wx_openid;size:50;unique_index" json:"wx_openid"`
	Enable      string    `gorm:"column:enable;size:1;default:'T';not null" json:"enable"`
	Admin       string    `gorm:"column:admin;size:1;default:'F';not null" json:"admin"`
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
	ID              int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	CreateAt        time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	UpdateAt        time.Time `gorm:"column:update_at;not null;" json:"update_at"`
	Name            string    `gorm:"column:name;size:20;not null" json:"name"`
	CompanyId       int       `gorm:"column:company_id;not null" json:"company_id"`
	CompanyName     string    `gorm:"-" json:"company_name"`
	MonitorTypeId   int       `gorm:"column:monitor_type_id;index" json:"monitor_type_id"`
	MonitorTypeName string    `gorm:"-" json:"monitor_type_name"`
	QrcodePath      string    `gorm:"column:qrcode_path" json:"qrcode_path"`
	QrcodeURI       string    `gorm:"column:qrcode_uri" json:"qrcode_uri"`
	Pictures        []Picture `gorm:"ForeignKey:MonitorPlaceId" json:"-"`
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
	Judgement      string    `gorm:"column:judgement;size:1;not null;default:'T'" json:"judgement"`
	JudgeComment   string    `gorm:"column:comment" json:"judgecomment"`
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
	RelaxDay    string    `gorm:"column:relax_day;default:'F';" json:"relax_day"`
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
	IsUpload         string    `gorm:"column:is_upload;string;size:1;not null;default:'F'" json:"is_upload"`
	Judgement        string    `gorm:"column:judgement;size:1;not null;default:'T'" json:"judgement"`
	EverJudge        string    `gorm:"column:ever_judge;string;size:1;not null;default:'F'" json:"ever_judge"`
}

func (TodaySummary) TableName() string {
	return "today_summary"
}

type MaterialPicture struct {
	ID      int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Day     time.Time `gorm:"column:day;not null;unique_index:day_company_place;default:NOW()" json:"day"`
	MediaId string    `gorm:"column:media_id;not null;unique_index" json:"media_id"`
	Url     string    `gorm:"column:url;not null" json:"url"`
}

func (MaterialPicture) TableName() string {
	return "material_picture"
}

type MaterialVideo struct {
	ID           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Day          time.Time `gorm:"column:day;not null;unique_index:day_company_place;default:NOW()" json:"day"`
	Title        string    `gorm:"column:title;not nul" json:"title"`
	Introduction string    `gorm:"column:introduction" json:"introduction"`
	MediaId      string    `gorm:"column:media_id" json:"media_id"`
	Url          string    `gorm:"column:url;not null" json:"url"`
}

func (MaterialVideo) TableName() string {
	return "material_video"
}

type MaterialAudio struct {
	ID           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Day          time.Time `gorm:"column:day;not null;unique_index:day_company_place;default:NOW()" json:"day"`
	Title        string    `gorm:"column:title;not nul" json:"title"`
	Introduction string    `gorm:"column:introduction" json:"introduction"`
	MediaId      string    `gorm:"column:media_id" json:"media_id"`
	Url          string    `gorm:"column:url;not null" json:"url"`
}

func (MaterialAudio) TableName() string {
	return "material_audio"
}

type MediaPicture struct {
	Url string `gorm:"-" json:"url"`
}

type Chapter struct {
	ID               int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	NewsId           int    `gorm:"column:news_id" json:"news_id"`
	Title            string `gorm:"column:title;not null" json:"title"`
	ThumbMediaId     string `gorm:"column:thumb_media_id" json:"thumb_media_id"`
	ThumbUrl         string `gorm:"column:thumb_url;not null" json:"thumb_url"`
	ShowCoverPic     int    `gorm:"column:show_cover_pic" json:"show_cover_pic"`
	Author           string `gorm:"column:author" json:"author"`
	Digest           string `gorm:"column:digest" json:"digest"`
	Content          string `gorm:"column:content" json:"content"`
	Url              string `gorm:"column:url" json:"url"`
	ContentSourceUrl string `gorm:"column:content_source_url" json:"content_source_url"`
	TemplatePageIds  string `gorm:"column:templatepageids" json:"templatepageids"`
}

func (Chapter) TableName() string {
	return "chapter"
}

type News struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	MediaId     string    `gorm:"column:media_id;not null;unique_index" json:"media_id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	ChapterIds  string    `gorm:"column:chapterids" json:"chapterids"`
	ChapterList []Chapter `gorm:"-" json:"chapter_list"`
}

func (News) TableName() string {
	return "news"
}

type GroupSend struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	MediaId   string    `gorm:"column:media_id;not null;unique_index" json:"media_id"`
	CreateAt  time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	NewsID    int       `gorm:"column:news_id" json:"news_id"`
	NewsName  string    `gorm:"column:news_name;index" json:"news_name"`
	Errcode   int       `gorm:"errcode" json:"errcode"`
	Errmsg    string    `gorm:"columen:errmsg" json:"errmsg"`
	MsgId     int64     `gorm:"msg_id" json:"msg_id"`
	MsgDataId int64     `gorm:"msg_data_id" json:"msg_data_id"`
}

type TemplatePage struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	CreateAt    time.Time `gorm:"column:create_at;not null;default:NOW()" json:"create_at"`
	ChapterIds  string    `gorm:"column:chapterids" json:"chapterids"`
	ChapterList []Chapter `gorm:"-" json:"chapter_list"`
	URL         string    `gorm:"column:url" json:"url"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func InitializeDB(dbuser, dbpass, dbip, dbport, dbname string) {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbuser, dbpass, dbip, dbport, dbname))
	if err != nil {
		glog.Fatalf("cannot initialize database connection, err %s", err)
	}
	db.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&Town{}, &Country{}, &Company{}, &User{}, &MonitorType{}, &MonitorPlace{}, &Picture{}, &Summary{}, &TodaySummary{}, &MaterialPicture{}, &MaterialVideo{}, &MaterialAudio{}, &Chapter{}, &News{}, &GroupSend{}, &TemplatePage{}, &GlobalRelaxPeriod{}, &CompanyRelaxPeriod{}, &CompanyMonthStat{}, &CompanyYearStat{})
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
	db.Debug().Model(&CompanyRelaxPeriod{}).AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE")
}
