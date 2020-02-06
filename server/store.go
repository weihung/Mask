package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
)

func UpdateStore(file string) {
	ReadCsvFile(file)
}

func ReadCsvFile(filePath string) {
	log.Println("Get store file")
	// Load a csv file.
	f, _ := os.Open(filePath)
	defer f.Close()
	// Create a new reader.
	r := csv.NewReader(f)
	tx := SqlDb.Begin()
	i := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		i++

		if err != nil {
			panic(err)
		}
		lat, _ := strconv.ParseFloat(record[5], 64)
		lng, _ := strconv.ParseFloat(record[6], 64)
		store := Store{
			MId:        record[0],
			Name:       record[1],
			Tel:        record[2],
			Address:    record[3],
			MapAddress: record[4],
			Lat:        lat,
			Lng:        lng,
			Remark:     record[7],
		}
		if err3 := tx.Where("m_id = ?", record[0]).First(&Store{}).Error; err3 != nil {
			if gorm.IsRecordNotFoundError(err3) {
				tx.Create(&store)
			}
		} else {
			tx.Model(&Store{}).Where("m_id = ?", record[0]).Update(&store)
		}
	}
	tx.Commit()
	log.Printf("Found %d records", i)
}
