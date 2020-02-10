package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ApiMask : 回應口罩資料
func ApiMask(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)
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

// ApiCardStatus : 號碼牌狀態
func ApiCardStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)
	mId, ok := r.URL.Query()["mId"]
	if !ok || len(mId[0]) < 1 {
		ResponseWithJson(w, http.StatusBadRequest, "Param 'mId' is missing")
		return
	}

	status, ok := r.URL.Query()["status"]
	if !ok || len(status[0]) < 1 {
		ResponseWithJson(w, http.StatusBadRequest, "Param 'status' is missing")
		return
	}
	s, err := strconv.Atoi(status[0])
	if err != nil {
		ResponseWithJson(w, http.StatusBadRequest, "Status error")
		return
	}
	cardStatus := CardStatus{MId: mId[0], Status: s}
	if err := SqlDb.Create(&cardStatus).Error; err != nil {
		ResponseWithJson(w, http.StatusInternalServerError, "Database error")
		return
	}
	ResponseWithJson(w, http.StatusOK, "")
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetData : 從 SQL 取回資料
func GetData(lat, lng string) []Result {
	var rs []Result
	// SqlDb.Table("stores").Select("stores.m_id, stores.name, stores.tel, stores.address, stores.lat, stores.lng, lefts.adult, lefts.child, lefts.update_time, c.updated_at as status_time, c.status").Joins("LEFT JOIN lefts ON lefts.m_id = stores.m_id AND lefts.update_time >= CURDATE()").Joins("LEFT JOIN (SELECT * FROM card_statuses ORDER BY card_statuses.updated_at DESC) as c ON c.m_id = stores.m_id AND c.updated_at >= CURDATE()").Where("DISTANCE(lat, lng, ?, ?, 'KM' ) < 1", lat, lng).Find(&rs)
	SqlDb.Table("stores").Select("stores.m_id, stores.name, stores.tel, stores.address, stores.lat, stores.lng, lefts.adult, lefts.child, lefts.update_time, c.updated_at as status_time, c.status").Joins("LEFT JOIN lefts ON lefts.m_id = stores.m_id").Joins("LEFT JOIN (SELECT * FROM card_statuses ORDER BY card_statuses.updated_at DESC) as c ON c.m_id = stores.m_id AND c.updated_at >= CURDATE()").Where("DISTANCE(lat, lng, ?, ?, 'KM' ) < 1", lat, lng).Find(&rs)
	return rs
}
