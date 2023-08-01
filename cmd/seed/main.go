package main

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/spf13/viper"
)

var posts = []*data.Post_DB{
	{
		UserID:      1,
		Image:       "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Description: "This is a test post",
	},
	{
		UserID:      2,
		Image:       "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Description: "This is a test post number 2",
	},
}

var users = []*data.User_DB{
	{
		Avatar:   "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Username: "dadvid",
		Email:    "tdest@test.com",
		Password: "password",
	},
	{
		Avatar:   "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Username: "bdine",
		Email:    "tedsty@yo.com",
		Password: "password",
	},
}

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	db := data.DBconnect()
	defer db.Close()

	data := data.New(db)

	var uID int

	for _, user := range users {
		id, err := data.CreateUser(user.Avatar, user.Username, user.Email, user.Password)
		if err != nil {
			fmt.Println(err)
		}
		uID = id
	}

	for _, post := range posts {
		err := data.CreatePost(uID, post.Image, post.Description)
		if err != nil {
			fmt.Println(err)
		}
	}
}
