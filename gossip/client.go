package gossip

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cc GossipApiClient

func newConn(peer string) error {
	if cc == nil {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(creds))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error connecting to %s: %s\n", peer, err)
			return err
		}
		cc = NewGossipApiClient(conn)
	}
	return nil
}

type Client struct {
	Peer string
	Name string
}

func NewClient(peer string, name string) *Client {
	return &Client{peer, name}
}

func (c *Client) SubscribeToEvents(ctx context.Context, user *User) (GossipApi_SubscribeToEventsClient, error) {
	_ = newConn(c.Peer)
	return cc.SubscribeToEvents(ctx, user)
}

func (c *Client) SendEvent(ctx context.Context, event *Event) (*Empty, error) {
	_ = newConn(c.Peer)
	return cc.SendEvent(ctx, event)
}

func (c *Client) UnsubscribeFromEvents(ctx context.Context, user *User) (*Empty, error) {
	_ = newConn(c.Peer)
	return cc.UnsubscribeFromEvents(ctx, user)
}
