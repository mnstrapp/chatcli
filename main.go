package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mnstrapp/chatcli/gossip"
)

var (
	host string
	name string
)

func init() {
	flag.StringVar(&host, "host", "", "Gossip host")
	flag.StringVar(&name, "name", "", "Nickname")
	flag.Parse()
}

func main() {
	if host == "" {
		flag.Usage()
		os.Exit(1)
	}

	user := &gossip.User{Id: name}

	client := gossip.NewClient(host, name)

	go func() {
		stream, err := client.SubscribeToEvents(context.Background(), user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error subscribing to events: %s", err)
			os.Exit(1)
		}
		for {
			event, err := stream.Recv()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Fprintf(os.Stderr, "error listening for events: %s", err)
					os.Exit(1)
				}
				continue
			}

			timestamp := time.Unix(int64(event.Timestamp), 0)
			switch event.Type {
			case gossip.EventType_UserEventType:
				fmt.Fprintf(os.Stdout, "[%s:%s] %s\n", *event.FromId, event.GetUser().Status, timestamp)
			case gossip.EventType_MessageEventType:
				fmt.Fprintf(os.Stdout, "[%s] %s - %s\n", *event.FromId, timestamp, event.GetMessage().Content)
			default:
				fmt.Fprintf(os.Stdout, "[%s:%s] %s\n", *event.FromId, event.Type, timestamp)
			}

		}
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		content, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %s\n", err)
		}

		if strings.Contains(content, "/quit") {
			client.UnsubscribeFromEvents(context.Background(), user)
			return
		}

		event := &gossip.Event{
			Type:      gossip.EventType_MessageEventType,
			FromId:    &user.Id,
			Timestamp: uint64(time.Now().Unix()),
			Body:      &gossip.Event_Message{Message: &gossip.MessageEvent{Content: []byte(content)}},
		}
		_, err = client.SendEvent(context.Background(), event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %s", err)
		}
	}
}
