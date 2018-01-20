package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/golang/glog"
)

//ParseExcelNewCompany 用来解析导入的excel文件
func ParseExcelNewCompany(fileName string) (string, string, error) {
	glog.Infof("ParseExcelNewCompany")
	var errorLines string
	var errorStr string
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		errmsg := fmt.Sprintf("cannot open file %s, err %s", fileName, err)
		glog.Errorf(errmsg)
		return "", "", errors.New(errmsg)
	}

	var townName, countryName, companyName string

	rows := f.GetRows("Sheet1")
OUT:
	for index, row := range rows {
		tx := db.Begin()
		if index == 0 {
			glog.Infof("line %d: session rollback", index+1)
			tx.Rollback()
			continue
		}

		if index > 100 {
			break
		}

		townName = row[0]
		if townName == "" {
			errmsg := fmt.Sprintf("line %d: town name %s is empty", index+1, countryName)
			if errorLines == "" {
				errorLines = fmt.Sprintf("%d", index+1)
				errorStr = errmsg
			} else {
				errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
				errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
			}
			glog.Error(errmsg)
			glog.Infof("line %d: session rollback", index+1)
			tx.Rollback()
			continue
		}

		town := Town{}
		tx.Debug().Where("name = ?", townName).First(&town)
		if town.ID == 0 {
			glog.Infof("line %d: town %s does not exist", index+1, townName)
			town.Name = townName
			tx.Debug().Create(&town)
			if town.ID == 0 {
				errmsg := fmt.Sprintf("line %d, town %s created failed", index+1, townName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
				}
				glog.Error(errmsg)
				glog.Infof("line %d: session rollback", index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("line %d: town %s created successfully, id %d", index+1, town.Name, town.ID)
			}
		}

		countryName = row[1]
		if countryName == "" {
			glog.Infof("line %d: country name is empty", index+1)
			glog.Infof("line %d: session commit", index+1)
			tx.Commit()
			continue
		}

		country := Country{}
		countries := make([]Country, 0)
		tx.Debug().Where(" town_id = ?", town.ID).Find(&countries)
		for _, c := range countries {
			if countryName == c.Name {
				country = c
			}
		}

		if country.ID == 0 {
			glog.Infof("line %d: country %s is not in town %s", index+1, countryName, townName)
			country.Name = countryName
			country.TownId = town.ID
			tx.Debug().Create(&country)
			if country.ID == 0 {
				errmsg := fmt.Sprintf("line %d: country %s created failed", index+1, countryName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
				}
				glog.Error(errmsg)
				glog.Infof("line %d: session rollback", index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("line %d: country %s created successfully, id %d", index+1, country.Name, country.ID)
			}
		}

		companyName = row[2]
		if companyName == "" {
			glog.Infof("line %d: company name %s is empty", index+1)
			glog.Infof("line %d: session commit", index+1)
			tx.Commit()
			continue
		}
		company := Company{}
		tx.Debug().Where("name = ?", companyName).First(&company)
		if company.ID == 0 {
			glog.Infof("line %d: company %s (country %s) does not exist", index+1, companyName, company.CountryName)
			//create new company
			company.Name = companyName
			company.CountryName = country.Name
			company.CountryId = country.ID
			company.Address = row[3]
			tx.Debug().Create(&company)
			if company.ID == 0 {
				errmsg := fmt.Sprintf("line %d: company %s cannot created", index+1, companyName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
				}
				glog.Error(errmsg)
				glog.Infof("line %d: session rollback", index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("line %d: company %s is created successfully, id %d", index+1, companyName, company.ID)
			}
		}

		if company.CountryId != country.ID {
			glog.Errorf("line %d: company %s belongs to country %s, not country %s", index+1, company.Name, company.CountryName, country.Name)
			glog.Infof("line %d: session rollback", index+1)
			tx.Rollback()
			continue
		}

		column := 4
		rowLen := len(row)
		for {
			if column > 12 || column+2 >= rowLen {
				break
			}
			user := User{}
			hashCode := md5.New()
			user.Name = row[column]
			if user.Name == "" {
				errmsg := fmt.Sprintf("line %d: user cannot created, empty user name", index+1)
				glog.Error(errmsg)
				glog.Infof("line %d: session commit", index+1)
				tx.Commit()
				continue OUT
			}
			user.Phone = row[column+1]
			if len(user.Phone) != 11 {
				errmsg := fmt.Sprintf("line %d: user %s cannot created, phone %s len %d", index+1, user.Name, user.Phone, len(user.Phone))
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
				}
				glog.Error(errmsg)
				glog.Infof("line %d: session rollback", index+1)
				tx.Rollback()
				continue OUT
			}

			samePhoneUser := User{}
			tx.Debug().Where("phone = ?", user.Phone).First(&samePhoneUser)
			if samePhoneUser.ID != 0 {
				if samePhoneUser.Name != user.Name {
					errmsg := fmt.Sprintf("line %d: user %s cannot created, phone %s confict with user %s(company %s)", index+1, user.Name, user.Phone, samePhoneUser.Name, samePhoneUser.CompanyName)
					if errorLines == "" {
						errorLines = fmt.Sprintf("%d", index+1)
						errorStr = errmsg
					} else {
						errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
						errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
					}
					glog.Error(errmsg)
					glog.Infof("line %d: session rollback", index+1)
					tx.Rollback()
					continue OUT
				}
			} else {
				user.Job = row[column+2]
				user.CompanyId = company.ID
				io.WriteString(hashCode, user.Phone[5:])
				user.Password = fmt.Sprintf("%x", hashCode.Sum(nil))
				tx.Debug().Create(&user)
				if user.ID == 0 {
					errmsg := fmt.Sprintf("line %d: user %s cannot created", index+1, user.Name)
					if errorLines == "" {
						errorLines = fmt.Sprintf("%d", index+1)
						errorStr = errmsg
					} else {
						errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
						errorStr = fmt.Sprintf("%s|%s", errorStr, errmsg)
					}
					glog.Error(errmsg)
					glog.Infof("line %d: session rollback", index+1)
					tx.Rollback()
					continue OUT
				} else {
					glog.Infof("line %d: create user %s, id %d", index+1, user.Name, user.ID)
				}
			}
			column += 3
		}
		glog.Infof("line %d, session commit", index+1)
		tx.Commit()
	}

	return errorLines, errorStr, nil
}
