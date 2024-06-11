package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"savetabs/storage"
)

const (
	DBFile = "/Users/mikeschinkel/Projects/savetabs/daemon/data/savetabs.db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Must pass an integer offset for links to send")
		return
	}
	offset := os.Args[1]
	n, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid value '%s'; must pass an integer offset for links to send", offset)
	}
	slog.Info("Preparing to send 100 links", "offset", n)
	produce(ctx, cancel, n)
	consume(ctx)
}

func consume(ctx context.Context) {
	conf := ReadConfig()
	// sets the consumer group ID and offset
	conf["group.id"] = "go-group-1"
	conf["auto.offset.reset"] = "earliest"

	slog.Info("Creating new consumer")
	// creates a new consumer and subscribes to your topic
	consumer, _ := kafka.NewConsumer(&conf)
	slog.Info("Subscribing consumer to links")
	err := consumer.SubscribeTopics([]string{"links"}, nil)
	if err != nil {
		slog.Error("Failed to subscribe to topic", "topic", "links", "error", err)
		os.Exit(3)
	}
	counter := 0
	run := true
	slog.Info("Starting to poll for events")
	for run {
		// consumes messages from the subscribed topic and prints them to the console
		e := consumer.Poll(1000)
		slog.Info("Event received", "event", e, "counter", counter)
		switch ev := e.(type) {
		case *kafka.Message:
			eventInfoMsg(ev, "Consumed event from topic")
			counter++
		case kafka.Error:
			slog.Error("Consuming event failed", "event", ev)
			run = false
			counter++
		}
		if counter >= 100 {
			slog.Info("All recent events consumed")
			break
		}
	}

	// closes the consumer connection
	must(consumer.Close())
}

type Link struct {
	Id        int64  `json:"-"`
	URL       string `json:"url"`
	Timestamp int32  `json:"timestamp"`
}

func (l Link) Key() []byte {
	return []byte(fmt.Sprintf("link-%d", l.Id))
}

func (l Link) Bytes() []byte {
	b, err := json.Marshal(l)
	if err != nil {
		slog.Error("Failed to marshal link to JSON", "link", l)
		return nil
	}
	return b
}

func newLink(u string, ts int32) Link {
	return Link{
		URL:       u,
		Timestamp: ts,
	}
}

func produce(ctx context.Context, cancel context.CancelFunc, offset int64) {
	var liteLinks []storage.LinkLite

	// creates a new producer instance
	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)
	topic := "links"

	wg := sync.WaitGroup{}

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			slog.Info("Received events", "event", e)
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					slog.Error("Failed to deliver message", "event", ev)
				} else {
					eventInfoMsg(ev, "Produced event to topic")
				}
				wg.Done()
			}
		}
	}()
	err := storage.Initialize(ctx, DBFile)
	if err != nil {
		cancel()
		slog.Error("Failed to initialize DB", "error", err)
		goto end
	}

	slog.Info("Loading links")
	liteLinks, err = storage.ListLinksLite(ctx, storage.ListLinksLiteArgs{
		ID:            offset,
		LinksArchived: storage.NotArchived,
		LinksDeleted:  storage.NotDeleted,
	})
	if err != nil {
		slog.Error("Failed to load links", "error", err)
		goto end
	}
	slog.Info("Links loaded", "num_links", len(liteLinks))
	for _, ll := range liteLinks {
		// produces a sample message to the user-created topic
		link := Link{
			Id:        ll.Id,
			URL:       ll.URL,
			Timestamp: int32(ll.Visited),
		}
		km := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            link.Key(),
			Value:          link.Bytes(),
		}
		wg.Add(1)
		slog.Info("Sending key-value", "key", string(link.Key()), "value", string(link.Bytes()))
		err = p.Produce(km, nil)
		if err != nil {
			slog.Error("Failed to produce to topic", "topic", topic, "error", err)
		}
	}
	// send any outstanding or buffered messages to the Kafka broker and close the connection
	p.Flush(15 * 1000)
	p.Close()
	slog.Info("Waiting for confirmation all messages were sent")
	wg.Wait()
	slog.Info("All messages sent")
end:
}

func ReadConfig() kafka.ConfigMap {
	config := map[string]kafka.ConfigValue{
		"bootstrap.servers":  "pkc-p11xm.us-east-1.aws.confluent.cloud:9092",
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"sasl.username":      strings.TrimSpace(os.Getenv("KAFKA_API_KEY")),
		"sasl.password":      strings.TrimSpace(os.Getenv("KAFKA_API_SECRET")),
		"session.timeout.ms": 45000,
	}
	if len(config["sasl.username"].(string)) == 0 {
		slog.Error("ERROR: KAFKA_API_KEY environment variable not set.")
		os.Exit(1)
	}
	if len(config["sasl.password"].(string)) == 0 {
		slog.Error("ERROR: KAFKA_API_SECRET environment variable not set.")
		os.Exit(2)
	}
	return config
}

func must(err error) {
	if err != nil {
		slog.Warn("Failed to close writer", "error", err)
	}
}

func eventInfoMsg(ev *kafka.Message, msg string) {
	slog.Info(msg,
		"topic", *ev.TopicPartition.Topic,
		"key", string(ev.Key),
		"value", string(ev.Value),
	)
}
