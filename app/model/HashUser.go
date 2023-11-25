package model

import "github.com/jinzhu/gorm"

type HashUser struct {
	gorm.Model
	UserID                int `json:"UserID" gorm:"primarykey;int;not null"`
	GotDiggCount          int `json:"got_digg_count" gorm:"int;"`
	GotViewCount          int `json:"got_view_count" gorm:"int;"`
	FolloweeCount         int `json:"followee_count" gorm:"int;"`
	FollowerCount         int `json:"follower_count" gorm:"int;"`
	FollowCollectSetCount int `json:"follow_collect_set_count" gorm:"int;"`
	SubscribeTagCount     int `json:"subscribe_tag_count" gorm:"int;"`
}

type HashRequest struct {
	gorm.Model
	ID                    string `json:"ID"`
	Command               string `json:"Command"`
	UserID                string `json:"UserID"`
	GotDiggCount          string `json:"got_digg_count"`
	GotViewCount          string `json:"got_view_count"`
	FolloweeCount         string `json:"followee_count"`
	FollowerCount         string `json:"follower_count"`
	FollowCollectSetCount string `json:"follow_collect_set_count"`
	SubscribeTagCount     string `json:"subscribe_tag_count"`
}
