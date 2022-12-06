package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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
	p := twitter.FriendListParams{ScreenName: "tak3zaki3"}
	list, _, err := c.Friends.List(&p)
	for _, v := range list.Users {
		fmt.Println(v.Name)
	}

	return
	t, err := Search(c)
	if err != nil {
		log.Println(err)
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
func containsNGWord(text string, list []string) bool {
	for _, v := range list {
		if strings.Contains(text, v) {
			return true
		}
	}
	return false
}

func FololowAndRetweetIfNeed(c *twitter.Client, s *twitter.Search) {
	ngList := getNGList()
	count := 0
	for _, t := range s.Statuses {
		if containsNGWord(t.Text, ngList) {
			continue
		}

		subject, err := isSubject(&t)
		if err != nil {
			log.Println(err)
		}
		if !subject {
			continue
		}
		err = FollowAndRetweet(c, &t)
		if err != nil {
			Notify(os.Getenv("LINE_TOKEN"), err.Error())
			log.Println(err)
			break
		}
		log.Println("ok")
		time.Sleep(5 * time.Minute)
		if count > 5 {
			break
		}
		count++
	}
}
func isFriends(c *twitter.Client, t *twitter.User) bool {
	return false
}

func FollowAndRetweet(c *twitter.Client, t *twitter.Tweet) error {
	if !isFriends(c, t.RetweetedStatus.User) {
		_, _, err := c.Friendships.Create(&twitter.FriendshipCreateParams{
			UserID: t.RetweetedStatus.User.ID,
		})
		time.Sleep(10 * time.Second)
		if err != nil {
			log.Println(err)
		}
	}
	_, _, err := c.Statuses.Retweet(t.ID, &twitter.StatusRetweetParams{})
	if err != nil {
		return err
	}
	return nil
}

var (
	DbName      = "history.sql"
	TableName   = "followRTs"
	NGlistTable = "NGList"
)

func getNGList() []string {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		fmt.Println("open")
		return []string{}
	}
	defer db.Close()

	query := fmt.Sprintf("create table %s (id integer not null primary key AUTOINCREMENT, NGText text)", NGlistTable)
	err = tableCreateIfNeed(db, NGlistTable, query)
	if err != nil {
		fmt.Println("create")
		return []string{}
	}

	db.Begin()
	rows, err := db.Query(fmt.Sprintf("select NGText from %s;", NGlistTable))
	if err != nil {
		fmt.Println("ng")
		return []string{}
	}
	defer rows.Close()

	list := make([]string, 0, 10)
	for rows.Next() {
		var text string
		rows.Scan(&text)
		list = append(list, text)
	}
	return list
}

func tableCreateIfNeed(db *sql.DB, name string, tableQuery string) error {
	if _, err := db.Exec(fmt.Sprintf("select count(*) from %s", name)); err != nil {
		log.Printf("create tablel: %s\n", &tableQuery)
		if _, err = db.Exec(tableQuery); err != nil {
			return err
		}
	}
	return nil
}

func isSubject(t *twitter.Tweet) (bool, error) {
	db, err := sql.Open("sqlite3", DbName)
	defer db.Close()

	if err != nil {
		return false, err
	}
	if _, err = db.Exec(fmt.Sprintf("select count(*) from %s", TableName)); err != nil {
		sqlStmt := fmt.Sprintf("create table %s (id integer not null primary key AUTOINCREMENT, TweetID text)", TableName)
		log.Printf("create tablel: %s\n", sqlStmt)
		if _, err = db.Exec(sqlStmt); err != nil {
			return false, err
		}
	}
	// already add DB
	query := fmt.Sprintf("select count(*) from %s where TweetID = %d", TableName, t.ID)
	v, err := db.Exec(query)
	if err != nil {
		return false, err
	}
	count, _ := v.LastInsertId()
	if count != 0 {
		log.Printf("same id:%d,%d\n", t.ID, count)
		return false, nil
	}
	// add DB
	db.Begin()
	if _, err := db.Exec(fmt.Sprintf("insert into %s(TweetID) values(%d)", TableName, t.ID)); err != nil {
		return false, err
	}
	return true, nil

}
