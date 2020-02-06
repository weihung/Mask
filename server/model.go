package main

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Left struct {
	gorm.Model
	MId 		 string `gorm:"unique;not null"`
	Name 		 string 
	Address  string
	Tel 		 string
	Adult 	 int64
	Child 	 int64
	UpdateTime *time.Time
} 

type Store struct {
	gorm.Model
	MId 		 string `gorm:"unique;not null"`
	Name 		 string 
	Tel			 string
	Address  string
	MapAddress string
	Lat			 float64
	Lng			 float64
	Remark   string
}

type Result struct {
	MId string `json:"mId"`
	Name string `json:"name"`
	Tel string `json:"tel"`
	Address string `json:"address"`
	Lat float64  `json:"lat"`
	Lng float64 `json:"lng"`
	Adult int `json:"adult"`
	Child int `json:"child"`
	Update_time string `json:"time"`
}
