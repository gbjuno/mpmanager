package tools

import (
	//	"github.com/emicklei/go-restful"

	"time"

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
