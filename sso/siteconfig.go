package sso

import "github.com/jinzhu/gorm"

type SiteConfig struct {
	gorm.Model
	Name          string
	Providers     []string
	BackgroundURL string
	Title         string
	SubTitle      string
}
