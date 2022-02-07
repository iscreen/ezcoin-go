package model

type JwtBlacklist struct {
	EZ_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
