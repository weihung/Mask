package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

func main() {
	defer CloseDB()

	if v := os.Getenv("PORT"); len(v) > 0 {
		port = v
	}

	if v := os.Getenv("API_PATH"); len(v) > 0 {
		path = "/" + v
	}

	log.Println("Mask server Start")
	for i := range os.Args {
		arg := os.Args[i]
		if arg == "-s" {
			file := os.Args[i+1]
			UpdateStore(file)
		}
	}

	go FetchNewData()

	mux := http.NewServeMux()
	mux.HandleFunc(path, ApiMask)
	log.Println("Mask API server is running")
	handler := cors.Default().Handler(mux)
	log.Println(http.ListenAndServe(":"+port, handler))

	log.Println("Server shut down.")
}

// ApiMask: 回應口罩資料
func ApiMask(w http.ResponseWriter, r *http.Request) {
	lat, ok := r.URL.Query()["lat"]

	if !ok || len(lat[0]) < 1 {
		log.Println("Url Param 'lat' is missing")
		return
	}

	lng, ok := r.URL.Query()["lng"]

	if !ok || len(lng[0]) < 1 {
		log.Println("Url Param 'lng' is missing")
		return
	}

	rs := GetData(lat[0], lng[0])
	ResponseWithJson(w, http.StatusOK, rs)
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetData: 從 SQL 取回資料
func GetData(lat, lng string) []Result {
	var rs []Result
	SqlDb.Table("stores").Select("lefts.m_id, stores.name, stores.tel, stores.address, stores.lat, stores.lng, lefts.adult, lefts.child, lefts.update_time").Joins("RIGHT JOIN lefts ON lefts.m_id = stores.m_id").Where("DISTANCE(lat, lng, ?, ?, 'KM' ) < 1", lat, lng).Find(&rs)
	return rs
}

// FetchNewData 重複取回資料時間
func FetchNewData() {
	d := time.Minute * 5
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
		t, err2 := time.ParseInLocation("2006/01/02 15:04:05", strings.TrimSpace(data[6]), time.Local)
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
