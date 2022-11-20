package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sivchari/gotwtr"
)

func loadEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("YOUR_TWITTER_BEARER_TOKEN")
}

func main() {
	token := loadEnv()
	client := gotwtr.New(token)

	opt := gotwtr.SearchTweetsOption{
		TweetFields: []gotwtr.TweetField{
			gotwtr.TweetFieldAuthorID,
			gotwtr.TweetFieldAttachments,
		},
		MaxResults: 10,
	}
	t, err := client.SearchRecentTweets(context.Background(), "フォロー＆RT", &opt)
	if err != nil {
		panic(err)
	}
	for _, v := range t.Tweets {
		fmt.Println("------------------------------")
		fmt.Println(v.Text)
		fmt.Println("===")
		fmt.Println(v.ID)
	}
}
