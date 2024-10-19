package database

import (
	"encoding/json"
	"log"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"unique"`
	Password    string `gorm:"unique"`
	Secret_word string `gorm:"unique"`
}

func CreateUser(login string, password string, secret string) error {
	result := DB.Create(&User{Login: login, Password: password, Secret_word: secret})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID(id uint) (User, error) {
	var user User
	result := DB.First(&user, id)
	return user, result.Error
}

func UpdateUser(id uint, payload string) error {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		log.Fatal(err)
	}
	if password, exists := data["password"]; exists {
		DB.Model(&user).Update("Password", password)
	}
	if login, exists := data["login"]; exists {
		DB.Model(&user).Update("Login", login)
	}
	if secret, exists := data["secret"]; exists {
		DB.Model(&user).Update("Secret_word", secret)
	}
	return nil
}

func DeleteUser(id uint) error {
	var user User
	result := DB.Delete(&user, id)
	return result.Error
}
