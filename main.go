package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return

}
func main() {
	c := getClient()

	t, err := Search(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	FololowAndRetweetIfNeed(c, t)
}

func getClient() *twitter.Client {
	loadEnv()
	clientID := os.Getenv("APIKEY")
	clientSecret := os.Getenv("APIKEY_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(clientID, clientSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func query() string {
	return "フォロー＆RT -裏垢"
}

func Search(c *twitter.Client) (*twitter.Search, error) {
	// Search Tweets
	search, _, err := c.Search.Tweets(&twitter.SearchTweetParams{
		Query: "フォロー＆RT",
		Count: 10,
	})
	return search, err
}

func FololowAndRetweetIfNeed(c *twitter.Client, s *twitter.Search) {
	count := 0
	for _, t := range s.Statuses {
		count++
		subject, err := isSubject(&t)
		if err != nil {
			fmt.Println(err)
		}
		if !subject {
			continue
		}
		err = FollowAndRetweet(c, &t)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Minute)
		if count > 5 {
			break
		}
	}
}

func FollowAndRetweet(c *twitter.Client, t *twitter.Tweet) error {
	_, _, err := c.Friendships.Create(&twitter.FriendshipCreateParams{
		UserID: t.RetweetedStatus.User.ID,
	})
	time.Sleep(10 * time.Second)
	fmt.Printf("%v\n", t)
	fmt.Println(t.Text)
	fmt.Println("----------------------------------------")
	fmt.Printf("user-iduser:%s \n", t.User.Name)
	fmt.Printf("rt-iduser:%s \n", t.RetweetedStatus.User.Name)
	if err != nil {
		return err
	}
	_, _, err = c.Statuses.Retweet(t.ID, &twitter.StatusRetweetParams{})
	if err != nil {
		return err
	}
	fmt.Println("ok!!!!")
	return nil
}

var (
	DbName    = "history.sql"
	TableName = "followRTs"
)

func isSubject(t *twitter.Tweet) (bool, error) {
	db, err := sql.Open("sqlite3", DbName)
	defer db.Close()

	if err != nil {
		return false, err
	}
	if _, err = db.Exec(fmt.Sprintf("select count(*) from %s", TableName)); err != nil {
		sqlStmt := fmt.Sprintf("create table %s (id integer not null primary key AUTOINCREMENT, TweetID text)", TableName)
		fmt.Printf("create tablel: %s\n", sqlStmt)
		if _, err = db.Exec(sqlStmt); err != nil {
			return false, err
		}
	}
	// already add DB
	query := fmt.Sprintf("select count(*) from %s where TweetID = %d", TableName, t.ID)
	if _, err := db.Exec(query); err != nil {
		fmt.Printf("query error:%s", query)
		return false, err
	}

	// add DB
	db.Begin()
	if _, err := db.Exec(fmt.Sprintf("insert into %s(TweetID) values(%d)", TableName, t.ID)); err != nil {
		return false, err
	}
	return true, nil

}
