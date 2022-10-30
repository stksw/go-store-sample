package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type User struct {
	Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email" gorm:"unique"`
	Password    []byte `json:"-"`
	IsSeller 		bool   `json:"-"`
	Revenue			*float64 `json:"revenue,omitempty" gorm:"-"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) Name() string {
	return user.FirstName + " " + user.LastName
}

type Seller User

func (seller *Seller) CalculateRevenue(db *gorm.DB) {
	var orders []Order
	var revenue float64 = 0

	db.Preload("OrderItems").Find(&orders, &Order {
		UserId: seller.Id,
		Complete: true,
	})
	
	for _, order := range orders{
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}
	seller.Revenue = &revenue
}

type Admin User 

func (admin *Admin) CalculateRevenue(db *gorm.DB) {
	var orders []Order
	var revenue float64 = 0

	db.Preload("OrderItems").Find(&orders, &Order {
		UserId: admin.Id,
		Complete: true,
	})
	
	for _, order := range orders{
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AdminRevenue
		}
	}
	admin.Revenue = &revenue
}

