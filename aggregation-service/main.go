package main

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

// load .env file
err := godotenv.Load(".env")
if err != nil {
  log.Fatalf("Error loading .env file")
}

// variables

consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
     "bootstrap.servers":    "host1:9092,host2:9092",
     "group.id":             "foo",
     "auto.offset.reset":    "smallest"})

err = consumer.SubscribeTopics(topics, nil)

for run == true {
	ev := consumer.Poll(100)
	switch e := ev.(type) {
	case *kafka.Message:
		// application-specific processing
	case kafka.Error:
		fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
		run = false
	default:
		fmt.Printf("Ignored %v\n", e)
	}
}

consumer.Close()