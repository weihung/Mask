package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Left struct {
	gorm.Model
	MId        string `gorm:"unique;not null"`
	Name       string
	Address    string
	Tel        string
	Adult      int64
	Child      int64
	UpdateTime *time.Time
}

type Store struct {
	gorm.Model
	MId        string `gorm:"unique;not null"`
	Name       string
	Tel        string
	Address    string
	MapAddress string
	Lat        float64
	Lng        float64
	Remark     string
}

type CardStatus struct {
	gorm.Model
	MId    string
	Status int
}

type Result struct {
	MId        string  `json:"mId"`
	Name       string  `json:"name"`
	Tel        string  `json:"tel"`
	Address    string  `json:"address"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Adult      int     `json:"adult"`
	Child      int     `json:"child"`
	UpdateTime string  `json:"time"`
	StatusTime string  `json:"cardUpdateTime"`
	Status     int     `json:"cardStatus"`
}
