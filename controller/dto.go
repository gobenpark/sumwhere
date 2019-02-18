package controllers

import (
	"sumwhere/models"
	"time"
)

const (
	DefaultMaxResultCount = 30
)

type SearchInput struct {
	Sortby         []string `query:"sortby"` // 퀴리스트링용
	Order          []string `query:"order"`
	SkipCount      int      `query:"skipCount"`
	MaxResultCount int      `query:"maxResultCount"`
}

type ValidateInput struct {
	TripId  int64  `query:"id"`
	StartAt string `query:"start"`
	EndAt   string `query:"end"`
}

type ProfileInput struct {
	Age           int                `json:"age" valid:"required"`
	Job           string             `json:"job" valid:"required"`
	TripStyleType string             `json:"tripStyleType"`
	CharacterType []models.Character `json:"characterType"`
	Image1        string             `json:"image1" valid:"required"`
	Image2        string             `json:"image2"`
	Image3        string             `json:"image3"`
	Image4        string             `json:"image4"`
}

type TripInput struct {
	ID          int64  `json:"id" valid:"-"`
	UserId      int64  `json:"userId" valid:"required"`
	MatchTypeID int64  `json:"matchTypeId" valid:"required"`
	GenderType  string `json:"genderType" valid:"required"`
	TripTypeId  int64  `json:"tripTypeId" valid:"required"`
	Region      string `json:"region" valid:"required"`
	Concept     string `json:"concept" valid:"required"`
	StartDate   string `json:"startDate" valid:"required"`
	EndDate     string `json:"endDate" valid:"required"`
}

func (p *ProfileInput) ToModel() (*models.Profile, error) {

	profile := &models.Profile{
		Age:           p.Age,
		Job:           p.Job,
		TripStyleType: p.TripStyleType,
		CharacterType: p.CharacterType,
		Image1:        p.Image1,
		Image2:        p.Image2,
		Image3:        p.Image3,
		Image4:        p.Image4,
	}

	return profile, nil
}

func (t *TripInput) ToModel() (*models.Trip, error) {
	startAt, err := time.Parse("2006-01-02", t.StartDate)
	if err != nil {
		return nil, err
	}

	endAt, err := time.Parse("2006-01-02", t.EndDate)
	if err != nil {
		return nil, err
	}

	trip := &models.Trip{
		UserId:      t.UserId,
		MatchTypeId: t.MatchTypeID,
		TripTypeId:  t.TripTypeId,
		GenderType:  t.GenderType,
		Region:      t.Region,
		Concept:     t.Concept,
		StartDate:   startAt,
		EndDate:     endAt,
	}
	return trip, nil
}
