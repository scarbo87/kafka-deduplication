package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	clientId string
	groupId  string
)

func init() {
	clientId = os.Getenv("CLIENT_ID")
	if clientId == "" {
		clientId = "default_client_id"
	}

	groupId = os.Getenv("GROUP_ID")
	if groupId == "" {
		groupId = "default_group_id"
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-ch
		cancel()
	}()

	go initConsumerGroup(ctx, clientId, []string{"localhost:9092"}, "test", groupId)

	<-ctx.Done()
	fmt.Println("close application")
	if ctx.Err() != nil && ctx.Err() != context.Canceled {
		log.Fatal("failed to close application:", ctx.Err())
	}
}

func initConsumerGroup(ctx context.Context, clientId string, brokers []string, topic string, groupId string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupId,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
		MinBytes:    10,   // 10B
		MaxBytes:    10e6, // 10MB
		Dialer: &kafka.Dialer{
			ClientID:  clientId,
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	})

	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
