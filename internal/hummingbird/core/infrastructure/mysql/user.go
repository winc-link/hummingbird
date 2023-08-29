package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"

	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//
	//"gitlab.com/tedge/edgex/internal/models"
	"gorm.io/gorm"
)

func (c *Client) GetUserByUserName(username string) (models.User, error) {
	user, edgeXErr := userByUserName(c, username)
	if edgeXErr != nil {
		return user, edgeXErr
	}
	return user, nil
}

func (c *Client) GetAllUser() ([]models.User, error) {
	return getAllUser(c)
}

func (c *Client) AddUsers(users []models.User) error {
	return addUsers(c, users)
}

func (c *Client) AddUser(u models.User) (models.User, error) {
	return addUser(c, u)
}

func (c *Client) UpdateUser(u models.User) error {
	return updateUser(c, u)
}

func userByUserName(c *Client, username string) (models.User, error) {
	user := models.User{}
	err := c.client.GetObject(&models.User{Username: username}, &user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, errort.NewCommonEdgeX(errort.AppPasswordError, fmt.Sprintf("fail to query username %s", username), err)
		} else {
			return user, err
		}
	}
	return user, nil
}

func getAllUser(c *Client) ([]models.User, error) {
	var users []models.User
	if err := c.Pool.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func addUsers(c *Client, users []models.User) error {
	if len(users) <= 0 {
		return nil
	}
	return c.Pool.Create(users).Error
}

func updateUser(c *Client, u models.User) error {
	err := c.Pool.Table(u.TableName()).Where(&models.User{
		Username: u.Username,
	}).Save(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func userExist(c *Client, username string) (bool, error) {
	exists, err := c.client.ExistObject(&models.User{Username: username})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func addUser(c *Client, u models.User) (models.User, error) {
	exists, edgeXErr := userExist(c, u.Username)
	if edgeXErr != nil {
		return u, edgeXErr
	} else if exists {
		return u, errort.NewCommonEdgeX(errort.DefaultNameRepeat, fmt.Sprintf("username %s exists", u.Username), edgeXErr)
	}
	if err := c.client.CreateObject(&u); err != nil {
		return u, errort.NewCommonEdgeX(errort.DefaultSystemError, "user creation failed", err)
	}
	return u, nil
}
