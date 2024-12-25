package modelsControlPanel

type User struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Servers  []Server `gorm:"foreignKey:UserID"`
}
