package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"helper/kafkaHelper"
	"helper/twitterHelper"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/namsral/flag"
)

var (
	kafkaBrokerUrl string
	kafkaTopic     string
	kafkaClientId  string
	searchTerm     string
)

func main() {

	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:9092", "Kafka brokers in comma separated value")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.StringVar(&searchTerm, "twitter-search", "", "Twitter search term")

	flag.Parse()

	// TODO: get config vars from etcd cluster
	creds := twitterHelper.Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ApiToken:          os.Getenv("API_TOKEN"),
		ApiTokenSecret:    os.Getenv("API_TOKEN_SECRET"),
	}

	// init twitter
	client, err := twitterHelper.GetClient(&creds)
	errorHandler("Error getting Twitter Client", err)

	// init kafka
	kafkaProducer, err := kafkaHelper.Configure(strings.Split(kafkaBrokerUrl, ","), kafkaClientId, kafkaTopic)
	errorHandler("Error configuring Kafka Producer", err)
	defer kafkaProducer.Close()

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: searchTerm,
	})
	errorHandler("Error searching tweet", err)

	for _, tweet := range search.Statuses {

		err = kafkaHelper.Push(kafkaProducer, context.Background(), nil, []byte(tweet.Text))
		errorHandler("Error Kafka producing message", err)

		fmt.Println("Text: ", tweet.Text)
	}
}

func errorHandler(message string, err error) {

	if err != nil {
		log.Println(message)
		log.Println(err)
		os.Exit(1)
	}
}
