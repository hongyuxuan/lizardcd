package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddUser(username, password string, db *gorm.DB) (err error) {
	var hashPwd []byte
	if hashPwd, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost); err != nil {
		return
	}
	user := types.User{
		Username: username,
		Password: string(hashPwd),
		UpdateAt: time.Now(),
	}
	return db.Create(&user).Error
}

func ValidatedUser(username, password string, db *gorm.DB) (err error) {
	var user types.User
	if err = db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewDefaultError(fmt.Sprintf("user \"%s\" not found", username))
		}
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		err = errorx.NewError(http.StatusUnauthorized, "wrong username or password", err)
	}
	return
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func ModifyPassword(username, oldPassword, newPassword string, db *gorm.DB) (err error) {
	if err = ValidatedUser(username, oldPassword, db); err != nil {
		return
	}
	var hashPwd []byte
	if hashPwd, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost); err != nil {
		return
	}
	return db.Where("username = ?", username).Table("user").Update("password", string(hashPwd)).Error
}
