package main

import (
	"context"

	"github.com/MaximkaSha/gophkeeper/internal/client"
	"github.com/MaximkaSha/gophkeeper/internal/ui"
)

var (
	BuildVersion string = "N/A"
	BuildTime    string = "N/A"
)

func main() {
	ctx := context.Background()
	client := client.NewClient(BuildVersion, BuildTime)

	ui.UI(ctx, *client)

	/*
			rand.Seed(time.Now().UnixNano())
		user := models.User{
			Email:    "test1233111@qqqqq.ru" + fmt.Sprintf("%v", (rand.Intn(10000))),
			Password: "Passqword",
		}
		err := client.UserRegister(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		user.HashPassword()

		err = client.UserLogin(ctx, user)
		if err != nil {
			log.Fatal(err)
		}
		err = client.GetAllDataFromDB(ctx)
		if err != nil {
			log.Println(err)
		}
		log.Println("No data found: good !")
		data := models.Password{
			Login:    "111111",
			Password: "22222222",
			Tag:      "333333333",
		}
		err = client.AddData(ctx, data)
		if err != nil {
			log.Fatal(err)
		}
		err = client.GetAllDataFromDB(ctx)
		if err != nil {
			log.Println("No data found: bad !")
			log.Println(err)
		}
		client.PrinStorage()
		for _, val := range client.LocalStorage.PasswordStorage {
			val.Password = "new passw"
			client.AddData(ctx, val)
		}
		client.PrinStorage()
		for _, val := range client.LocalStorage.PasswordStorage {
			client.DelData(ctx, val)
		}
		client.PrinStorage()
	*/
}
