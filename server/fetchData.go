package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// FetchNewData 重複取回資料時間
func FetchNewData() {
	d := time.Minute * 1
	t := time.NewTicker(d)
	defer t.Stop()

	FetchData()
	for {
		<-t.C
		FetchData()
	}
}

// FetchData 從 open data 取回資料
func FetchData() {
	log.Println("Fetch new data")
	resp, err := http.Get(DataURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("Analysis data")
	reader := bufio.NewReader(bytes.NewReader(body))
	reader.ReadSlice('\n')
	tx := SqlDb.Begin()
	i := 0
	for {
		i++
		line, err := reader.ReadSlice('\n')
		if err != nil {
			log.Println(err)
			break
		}
		data := strings.Split(string(line), ",")
		loc, _ := time.LoadLocation("Asia/Taipei")
		t, err2 := time.ParseInLocation("2006/01/02 15:04:05", strings.TrimSpace(data[6]), loc)
		if err2 != nil {
			log.Println(err2)
		}
		adult, _ := strconv.ParseInt(data[4], 10, 64)
		child, _ := strconv.ParseInt(data[5], 10, 64)

		left := Left{
			MId:        data[0],
			Name:       data[1],
			Address:    data[2],
			Tel:        data[3],
			Adult:      adult,
			Child:      child,
			UpdateTime: &t}
		if err3 := tx.Where("m_id = ?", data[0]).First(&Left{}).Error; err3 != nil {
			if gorm.IsRecordNotFoundError(err3) {
				tx.Create(&left)
			}
		} else {
			tx.Model(&Left{}).Where("m_id = ?", data[0]).Update(&left)
		}
	}
	tx.Commit()
	log.Printf("Found %d records", i)
}
