package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_Query(t *testing.T) {
	now := time.Now()
	mugong := MonitorType{CreateAt: now, UpdateAt: now, Name: "mugong", Comment: "mugong"}
	db.Debug().Create(&mugong)

	lecong := Town{CreateAt: now, UpdateAt: now, Name: "lecong"}
	db.Debug().Create(&lecong)

	daluocun := Country{CreateAt: now, UpdateAt: now, Name: "daluocun", TownId: lecong.ID}
	db.Debug().Create(&daluocun)

	company := Company{CreateAt: now, UpdateAt: now, Name: "a", Address: "1 street", CountryId: daluocun.ID}
	db.Debug().Create(&company)

	user := User{CreateAt: now, UpdateAt: now, Phone: "13333333333", Name: "userA", Password: "123456", Job: "boss", CompanyId: company.ID, Enable: "Y"}
	db.Debug().Create(&user)

	place := MonitorPlace{CreateAt: now, UpdateAt: now, Name: "mugong1", CompanyId: company.ID, MonitorTypeId: mugong.ID}
	db.Debug().Create(&place)
	place.Qrcode = fmt.Sprintf("qrcode/%d/%d", place.CompanyId, place.ID)
	place.QrcodePath = fmt.Sprintf("/monitorplace/%d/qrcode", place.ID)
	db.Debug().Save(&place)

	picture := Picture{CreateAt: now, UpdateAt: now, MonitorPlaceId: place.ID, Corrective: "N", UserId: user.ID}
	db.Debug().Create(&picture)
	monitorplace := MonitorPlace{}
	db.Debug().First(&monitorplace, picture.MonitorPlaceId)
	year, month, day := picture.CreateAt.Date()
	createDate := fmt.Sprintf("%d%d%d", year, month, day)
	picture.Thumb = fmt.Sprintf("picture/%s/%d/%d_%d_thumb", createDate, monitorplace.CompanyId, monitorplace.ID, picture.CreateAt.Unix())
	picture.Full = fmt.Sprintf("picture/%s/%d/%d_%d_full", createDate, monitorplace.CompanyId, monitorplace.ID, picture.CreateAt.Unix())
	picture.ThumbPath = fmt.Sprintf("/picture/%d/thumb", picture.ID)
	picture.FullPath = fmt.Sprintf("picture/%d/full", picture.ID)
	db.Debug().Save(&picture)

	place1 := MonitorPlace{CreateAt: now, UpdateAt: now, Name: "mugong2", CompanyId: company.ID, MonitorTypeId: mugong.ID}
	db.Debug().Create(&place1)
	place1.Qrcode = fmt.Sprintf("qrcode/%d/%d", place1.CompanyId, place1.ID)
	place1.QrcodePath = fmt.Sprintf("/monitorplace/%d/qrcode", place1.ID)
	db.Debug().Save(&place1)

	place2 := MonitorPlace{CreateAt: now, UpdateAt: now, Name: "mugong3", CompanyId: company.ID, MonitorTypeId: mugong.ID}
	db.Debug().Create(&place2)
	place2.Qrcode = fmt.Sprintf("qrcode/%d/%d", place2.CompanyId, place2.ID)
	place2.QrcodePath = fmt.Sprintf("/monitorplace/%d/qrcode", place2.ID)
	db.Debug().Save(&place2)

	places := make([]MonitorPlace, 2)
	db.Debug().Model(&company).Related(&places).Where("id <> ?", "1")

	db.Debug().Delete(&picture)
	db.Debug().Delete(&place2)
	db.Debug().Delete(&place1)
	db.Debug().Delete(&place)
	db.Debug().Delete(&user)
	db.Debug().Delete(&company)
	db.Debug().Delete(&daluocun)
	db.Debug().Delete(&mugong)
	db.Debug().Delete(&lecong)
}
