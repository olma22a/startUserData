package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type UserData struct {
	UserID    int64  `bson:"user_id"`
	Username  string `bson:"username"`
	FirstName string `bson:"first_name"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	collection := client.Database("telegram").Collection("user_data")

	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я тебя запомнил чуркабес")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}

				user := UserData{
					UserID:    int64(update.Message.Chat.ID),
					Username:  update.Message.Chat.UserName,
					FirstName: update.Message.Chat.FirstName,
				}

				_, err := collection.InsertOne(context.Background(), user)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
