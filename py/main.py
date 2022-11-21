
import tweepy

client = tweepy.Client("Bearer Token here", user_fields="id")

res = client.search_recent_tweets("フォロー＆RT")

for tweet in res.data:
    print(tweet.text)
	client.retweet(tweet.id)
    client.follow_user