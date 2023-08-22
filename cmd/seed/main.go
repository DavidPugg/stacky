package main

import (
	"fmt"
	"math/rand"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/spf13/viper"
)

var images = []string{
	"https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
	"https://cdn.pixabay.com/photo/2023/07/20/04/45/leva-8138344_640.jpg",
	"https://cdn.pixabay.com/photo/2023/06/05/13/56/cat-8042342_640.jpg",
	"https://cdn.pixabay.com/photo/2019/12/07/14/57/rubber-4679464_640.png",
	"https://cdn.pixabay.com/photo/2023/07/27/13/37/cottage-8153413_640.jpg",
	"https://cdn.pixabay.com/photo/2023/06/25/08/46/woman-8086721_640.jpg",
	"https://cdn.pixabay.com/photo/2023/07/08/12/48/butterfly-8114483_640.jpg",
	"https://cdn.pixabay.com/photo/2023/05/24/06/53/bird-8014191_640.jpg",
}

var descriptions = []string{
	"Do you like this photo?",
	"Check out this photo!",
	"First post!",
	"Been a while since I posted",
	"Testing out the app",
	"lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore magna aliqua",
}

var posts = []*data.DBPost{
	{
		Image:       "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Description: "This is a test post",
	},
	{
		Image:       "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Description: "This is a test post number 2",
	},
}

var users = []*data.DBUser{
	{
		Avatar:   "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Username: "Dejv",
		Email:    "tdest@test.com",
		Password: "password",
	},
	{
		Avatar:   "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg",
		Username: "Bine",
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

	for _, user := range users {
		id, err := data.CreateUser(user.Avatar, user.Username, user.Email, user.Password)
		if err != nil {
			fmt.Println(err)
		}

		random := rand.Intn(6) + 2

		for i := 0; i < random; i++ {
			randImage := rand.Intn(len(images))
			randDescription := rand.Intn(len(descriptions))
			err := data.CreatePost(id, images[randImage], descriptions[randDescription])
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
