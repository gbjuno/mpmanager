package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/golang/glog"
)

func getDataHandler(w http.ResponseWriter, r *http.Request) {
}

func getData() (string, error) {
	prefix := fmt.Sprintf("[%s]", "getExcel")
	glog.Info(prefix)

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "镇")
	f.SetCellValue("Sheet1", "B1", "村")
	f.SetCellValue("Sheet1", "C1", "公司")
	f.SetCellValue("Sheet1", "D1", "公司地址")
	f.SetCellValue("Sheet1", "E1", "负责人1")
	f.SetCellValue("Sheet1", "F1", "负责人1的手机")
	f.SetCellValue("Sheet1", "G1", "负责人1的职位")
	f.SetCellValue("Sheet1", "H1", "负责人2")
	f.SetCellValue("Sheet1", "I1", "负责人2的手机")
	f.SetCellValue("Sheet1", "J1", "负责人2的职位")
	f.SetCellValue("Sheet1", "K1", "负责人3")
	f.SetCellValue("Sheet1", "L1", "负责人3的手机")
	f.SetCellValue("Sheet1", "M1", "负责人3的职位")
	f.SetCellValue("Sheet1", "N1", "负责人4")
	f.SetCellValue("Sheet1", "O1", "负责人4的手机")
	f.SetCellValue("Sheet1", "P1", "负责人4的职位")
	f.SetCellValue("Sheet1", "Q1", "负责人5")
	f.SetCellValue("Sheet1", "R1", "负责人5的手机")
	f.SetCellValue("Sheet1", "S1", "负责人5的职位")
	townList := make([]Town, 0)
	db.Debug().Find(&townList)
	line := 2
TOWN:
	for _, town := range townList {
		countryList := make([]Country, 0)
		db.Debug().Where("town_id = ?", town.ID).Find(&countryList)
		if len(countryList) == 0 {
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), town.Name)
			line++
			continue TOWN
		}
	COUNTRY:
		for _, country := range countryList {
			companyList := make([]Company, 0)
			db.Debug().Where("country_id = ?", country.ID).Find(&companyList)
			if len(companyList) == 0 {
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), town.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), country.Name)
				line++
				continue COUNTRY
			}
			for _, company := range companyList {
				userList := make([]User, 0)
				db.Debug().Where("company_id = ?", company.ID).Find(&userList)
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), town.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), country.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), company.Name)
				f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), company.Address)
				init := 'E'
				for index, user := range userList {
					if index == 5 {
						break
					}
					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(init), line), user.Name)
					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(init+1), line), user.Phone)
					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(init+2), line), user.Job)
					init += 3
				}
				line++
			}
		}
	}

	t := time.Now()
	saveFileName := fmt.Sprintf("/tmp/Data_%d%d%d%d%d_%d.xlsx", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Nanosecond())
	err := f.SaveAs(saveFileName)
	if err != nil {
		errmsg := fmt.Sprintf("cannot save file, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return "", errors.New(errmsg)
	}
	glog.Infof("%s save file to %s successfully", prefix, saveFileName)
	return saveFileName, nil
}

//ExcelStatus 是返回的状态结构
type ExcelStatus struct {
	Status     int
	Message    string
	ErrorLines string
	ErrorStr   string
}

//excelHandler 用来处理上传或下载excel文件
func excelHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s]", "setDataHandler")
	status := ExcelStatus{}
	if r.Method == "GET" {
		fileName, err := getData()
		if err != nil {
			glog.Errorf("%s get data failed, err %s", prefix, err)
			status.Status = http.StatusInternalServerError
			status.Message = fmt.Sprintf("get data failed, %s", err)
			returnContent, _ := json.Marshal(status)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(returnContent))
			return
		}
		f, err := os.Open(fileName)
		if err != nil {
			glog.Errorf("%s get data failed, err %s", prefix, err)
			status.Status = http.StatusInternalServerError
			status.Message = fmt.Sprintf("get data failed, cannot open file")
			returnContent, _ := json.Marshal(status)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(returnContent))
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=basedata.xlsx")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		glog.Infof("%s return data file successfully", prefix)
		return
	} else if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			glog.Errorf("%s, cannot read file, err %s", prefix, err)
			status.Status = http.StatusInternalServerError
			status.Message = "parseExcel failed, cannot open file"
			returnContent, _ := json.Marshal(status)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(returnContent))
			return
		}
		defer file.Close()
		saveFileName := fmt.Sprintf("/tmp/%s", handler.Filename)
		f, err := os.OpenFile(saveFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			glog.Errorf("%s, cannot save file to %s, err %s", prefix, saveFileName, err)
			status.Status = http.StatusInternalServerError
			status.Message = "parseExcel failed, cannot open file"
			returnContent, _ := json.Marshal(status)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(returnContent))
			return
		}
		io.Copy(f, file)
		f.Close()
		errorLines, errorStr, err := parseExcel(saveFileName)
		if err != nil {
			glog.Errorf("%s, parseExcel %s failed,  err %s", prefix, saveFileName, err)
			status.Status = http.StatusInternalServerError
			status.Message = "parseExcel failed, cannot open file"
			returnContent, _ := json.Marshal(status)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(returnContent))
		}
		status.Status = http.StatusOK
		status.Message = fmt.Sprintf("OK, upload file is parsed!")
		status.ErrorLines = errorLines
		status.ErrorStr = errorStr
		returnContent, _ := json.Marshal(status)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(returnContent))
		return
	} else {
		status.Status = http.StatusMethodNotAllowed
		status.Message = "Method Not Allowed"
		returnContent, _ := json.Marshal(status)
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, string(returnContent))
		return
	}
}

//parseExcel 用来解析导入的excel文件
func parseExcel(fileName string) (string, string, error) {
	prefix := fmt.Sprintf("[%s]", "parseExcel")
	glog.Info(prefix)
	var errorLines string
	var errorStr string
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		//errmsg := fmt.Sprintf("cannot open file %s, err %s", fileName, err)
		errmsg := fmt.Sprintf("cannot open file %s, err %s", fileName, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return "", "", errors.New(errmsg)
	}

	var townName, countryName, companyName string
	rows := f.GetRows("Sheet1")
OUT:
	for index, row := range rows {
		tx := db.Begin()
		if index == 0 {
			glog.Infof("%s 第%d行: session rollback", prefix, index+1)
			tx.Rollback()
			continue
		}

		if index > 100 {
			break
		}

		townName = row[0]
		if townName == "" {
			//errmsg := fmt.Sprintf("第%d行: town name %s is empty", index+1, countryName)
			errmsg := fmt.Sprintf("第%d行: 镇名为空", index+1)
			if errorLines == "" {
				errorLines = fmt.Sprintf("%d", index+1)
				errorStr = errmsg
			} else {
				errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
				errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
			}
			glog.Errorf("%s %s", prefix, errmsg)
			glog.Infof("%s 第%d行: session rollback", prefix, index+1)
			tx.Rollback()
			continue
		}

		town := Town{}
		tx.Debug().Where("name = ?", townName).First(&town)
		if town.ID == 0 {
			glog.Infof("%s 第%d行: town %s does not exist", prefix, index+1, townName)
			town.Name = townName
			tx.Debug().Create(&town)
			if town.ID == 0 {
				//errmsg := fmt.Sprintf("第%d行, town %s created failed", index+1, townName)
				errmsg := fmt.Sprintf("第%d行, 镇类型%s创建失败", index+1, townName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
				}
				glog.Errorf("%s %s", prefix, errmsg)
				glog.Infof("%s 第%d行: session rollback", prefix, index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("%s 第%d行: town %s created successfully, id %d", prefix, index+1, town.Name, town.ID)
			}
		}

		countryName = row[1]
		if countryName == "" {
			glog.Infof("%s 第%d行: country name is empty", prefix, index+1)
			glog.Infof("%s 第%d行: session commit", prefix, index+1)
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
			glog.Infof("%s 第%d行: country %s is not in town %s", prefix, index+1, countryName, townName)
			country.Name = countryName
			country.TownId = town.ID
			tx.Debug().Create(&country)
			if country.ID == 0 {
				//errmsg := fmt.Sprintf("第%d行: country %s created failed", index+1, countryName)
				errmsg := fmt.Sprintf("第%d行: 村类型%s创建失败", index+1, countryName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
				}
				glog.Errorf("%s %s", prefix, errmsg)
				glog.Infof("%s 第%d行: session rollback", prefix, index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("%s 第%d行: country %s created successfully, id %d", prefix, index+1, country.Name, country.ID)
			}
		}

		companyName = row[2]
		if companyName == "" {
			glog.Infof("%s 第%d行: company name %s is empty", prefix, index+1)
			glog.Infof("%s 第%d行: session commit", prefix, index+1)
			tx.Commit()
			continue
		}
		company := Company{}
		tx.Debug().Where("name = ?", companyName).First(&company)
		if company.ID == 0 {
			glog.Infof("%s 第%d行: company %s (country %s) does not exist", prefix, index+1, companyName, company.CountryName)
			//create new company
			company.Name = companyName
			company.CountryName = country.Name
			company.CountryId = country.ID
			company.Address = row[3]
			tx.Debug().Create(&company)
			if company.ID == 0 {
				//errmsg := fmt.Sprintf("%s 第%d行: company %s cannot created", prefix, index+1, companyName)
				errmsg := fmt.Sprintf("%s 第%d行: 公司类型%s创建失败", prefix, index+1, companyName)
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
				}
				glog.Errorf("%s %s", prefix, errmsg)
				glog.Infof("%s 第%d行: session rollback", prefix, index+1)
				tx.Rollback()
				continue
			} else {
				glog.Infof("%s 第%d行: company %s is created successfully, id %d", prefix, index+1, companyName, company.ID)
			}
		}

		if company.CountryId != country.ID {
			glog.Errorf("%s 第%d行: company %s belongs to country %s, not country %s", prefix, index+1, company.Name, company.CountryName, country.Name)
			glog.Infof("%s 第%d行: session rollback", prefix, index+1)
			tx.Rollback()
			continue
		}

		column := 4
		rowLen := len(row)
		for {
			if column > 18 || column+2 >= rowLen {
				break
			}
			user := User{}
			hashCode := md5.New()
			user.Name = row[column]
			if user.Name == "" {
				errmsg := fmt.Sprintf("第%d行: user cannot created, empty user name", index+1)
				glog.Errorf("%s %s", prefix, errmsg)
				glog.Infof("%s 第%d行: session commit", prefix, index+1)
				tx.Commit()
				continue OUT
			}
			user.Phone = row[column+1]
			if len(user.Phone) != 11 {
				//errmsg := fmt.Sprintf("第%d行: user %s cannot created, phone %s len %d", index+1, user.Name, user.Phone, len(user.Phone))
				errmsg := fmt.Sprintf("第%d行: 用户%s无法创建成功，电话号码长度为%d", index+1, user.Name, user.Phone, len(user.Phone))
				if errorLines == "" {
					errorLines = fmt.Sprintf("%d", index+1)
					errorStr = errmsg
				} else {
					errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
					errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
				}
				glog.Errorf("%s %s", prefix, errmsg)
				glog.Infof("%s 第%d行: session rollback", prefix, index+1)
				tx.Rollback()
				continue OUT
			}

			samePhoneUser := User{}
			tx.Debug().Where("phone = ?", user.Phone).First(&samePhoneUser)
			if samePhoneUser.ID != 0 {
				if samePhoneUser.Name != user.Name {
					//errmsg := fmt.Sprintf("第%d行: user %s cannot created, phone %s confict with user %s(company %s)", index+1, user.Name, user.Phone, samePhoneUser.Name, samePhoneUser.CompanyName)
					c := Company{}
					db.Debug().Where("id = ?", samePhoneUser.CompanyId).First(&c)
					samePhoneUserCompanyName := c.Name
					errmsg := fmt.Sprintf("第%d行: 用户%s无法创建，手机%s与已有用户%s(公司%s)冲突", index+1, user.Name, user.Phone, samePhoneUser.Name, samePhoneUserCompanyName)
					if errorLines == "" {
						errorLines = fmt.Sprintf("%d", index+1)
						errorStr = errmsg
					} else {
						errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
						errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
					}
					glog.Errorf("%s %s", prefix, errmsg)
					glog.Infof("%s 第%d行: session rollback", prefix, index+1)
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
					errmsg := fmt.Sprintf("第%d行: 用户创建失败，联系管理员", index+1, user.Name)
					if errorLines == "" {
						errorLines = fmt.Sprintf("%d", index+1)
						errorStr = errmsg
					} else {
						errorLines = fmt.Sprintf("%s,%d", errorLines, index+1)
						errorStr = fmt.Sprintf("%s\n%s", errorStr, errmsg)
					}
					glog.Errorf("%s %s", prefix, errmsg)
					glog.Infof("%s 第%d行: session rollback", prefix, index+1)
					tx.Rollback()
					continue OUT
				} else {
					glog.Infof("%s 第%d行: create user %s, id %d", prefix, index+1, user.Name, user.ID)
				}
			}
			column += 3
		}
		glog.Infof("%s 第%d行, session commit", prefix, index+1)
		tx.Commit()
	}

	return errorLines, errorStr, nil
}
