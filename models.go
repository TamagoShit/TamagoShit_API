package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserId             int        `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Name               string     `json:"name" gorm:"not null"`
	LastName           string     `json:"last_name" gorm:"not null"`
	UserName           string     `json:"user_name" gorm:"unique;not null"`
	Email              string     `json:"email" gorm:"unique;not null"`
	Password           string     `json:"-" gorm:"not null"` // Hashed password, hidden from JSON
	ProfilPicture      *string    `json:"profil_picture"`
	GamingTime         int        `json:"gaming_time" gorm:"default:0"`
	CreationDate       time.Time  `json:"creation_date" gorm:"default:CURRENT_TIMESTAMP"`
	LastConnectionDate *time.Time `json:"last_connection_date"`
	Role               string     `json:"role" gorm:"default:'general'"` // Added for authorization
}

func (User) TableName() string {
	return "Users"
}

type Tama struct {
	TamaId      int        `json:"tama_id" gorm:"primaryKey;autoIncrement"`
	UserId      int        `json:"user_id" gorm:"not null"`
	TamaStatsID int        `json:"tama_stats_id" gorm:"not null"`
	Name        string     `json:"name" gorm:"not null"`
	Sexe        bool       `json:"sexe"`
	Race        string     `json:"race" gorm:"not null"`
	Sickness    *string    `json:"sickness"`
	Birthday    *time.Time `json:"birthday"`
	DeathDay    *time.Time `json:"death_day"`
	Traits      *string    `json:"traits"`
	LifeChoices *string    `json:"life_choices"`
	User        User       `gorm:"foreignKey:UserId"`
	TamaStats   TamaStats  `gorm:"foreignKey:TamaStatsID"`
}

func (Tama) TableName() string {
	return "tamas"
}

type TamaStats struct {
	TamaStatId    int     `json:"tama_stat_id" gorm:"primaryKey;autoIncrement"`
	Food          int     `json:"food" gorm:"default:0"`
	Play          int     `json:"play" gorm:"default:0"`
	Cleaned       int     `json:"cleaned" gorm:"default:0"`
	CarAccident   int     `json:"car_accident" gorm:"default:0"`
	WorkAccident  int     `json:"work_accident" gorm:"default:0"`
	SocialSatis   float64 `json:"social_satis" gorm:"default:0"`
	WorkSatis     float64 `json:"work_satis" gorm:"default:0"`
	PersonalSatis float64 `json:"personal_satis" gorm:"default:0"`
}

func (TamaStats) TableName() string {
	return "Tama_stats"
}

type Friend struct {
	UserID            int       `json:"user_id" gorm:"primaryKey"`
	FriendID          int       `json:"friend_id" gorm:"primaryKey"`
	DateBecameFriends time.Time `json:"date_became_friends" gorm:"default:CURRENT_DATE"`
	User              User      `gorm:"foreignKey:UserID"`
	Friend            User      `gorm:"foreignKey:FriendID"`
}

func (Friend) TableName() string {
	return "Friends"
}

type Sponsor struct {
	SponsorId     int       `json:"sponsor_id" gorm:"primaryKey"`
	SponsoredId   int       `json:"sponsored_id" gorm:"primaryKey"`
	DateOfSponsor time.Time `json:"date_of_sponsor" gorm:"default:CURRENT_DATE"`
	SponsorUser   User      `gorm:"foreignKey:SponsorId"`
	SponsoredUser User      `gorm:"foreignKey:SponsoredId"`
}

func (Sponsor) TableName() string {
	return "Sponsor"
}

type Race struct {
	RaceId int     `json:"race_id" gorm:"primaryKey;autoIncrement"`
	Name   string  `json:"name" gorm:"unique;not null"`
	Desc   *string `json:"desc"`
	Bonus  *string `json:"bonus"`
	Malus  *string `json:"malus"`
}

func (Race) TableName() string {
	return "Race"
}

type Sickness struct {
	SicknessId     int     `json:"sickness_id" gorm:"primaryKey;autoIncrement"`
	Name           string  `json:"name" gorm:"not null"`
	Desc           *string `json:"desc"`
	ExpirationDays *int    `json:"expiration_days"`
	Bonus          *string `json:"bonus"`
	Malus          *string `json:"malus"`
}

func (Sickness) TableName() string {
	return "Sickness"
}

type Trait struct {
	TraitId int     `json:"trait_id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name" gorm:"not null"`
	Desc    *string `json:"desc"`
	Bonus   *string `json:"bonus"`
	Malus   *string `json:"malus"`
}

func (Trait) TableName() string {
	return "Trait"
}

type Bonus struct {
	BonusId int     `json:"bonus_id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name" gorm:"not null"`
	Desc    *string `json:"desc"`
	Effect  string  `json:"effect"`
}

func (Bonus) TableName() string {
	return "Bonus"
}

type Malus struct {
	MalusId int     `json:"malus_id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name" gorm:"not null"`
	Desc    *string `json:"desc"`
	Effect  string  `json:"effect"`
}

func (Malus) TableName() string {
	return "Malus"
}

type Event struct {
	EventId int     `json:"event_id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name" gorm:"not null"`
	Desc    *string `json:"desc"`
	Bonus   *string `json:"bonus"`
	Malus   *string `json:"malus"`
}

func (Event) TableName() string {
	return "Event"
}

type LifeChoice struct {
	LifeChoicesId int     `json:"life_choices_id" gorm:"primaryKey;autoIncrement"`
	Name          string  `json:"name" gorm:"not null"`
	Desc          *string `json:"desc"`
	Traits        *string `json:"traits"`
}

func (LifeChoice) TableName() string {
	return "LifeChoices"
}

var DB *gorm.DB
