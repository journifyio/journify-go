package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/journifyio/journify-go"
	"github.com/segmentio/conf"
)

func main() {
	var config struct {
		WriteKey   string `conf:"writeKey"   help:"The Journify Write Key of the project to send data to"`
		Type       string `conf:"type"       help:"The type of the message to send"`
		UserID     string `conf:"userId"     help:"Unique identifier for the user"`
		GroupID    string `conf:"groupId"    help:"Unique identifier for the group"`
		Traits     string `conf:"traits"     help:"Metadata associated with the user"`
		Event      string `conf:"event"      help:"Name of the track event"`
		Properties string `conf:"properties" help:"Metadata associated with an event or page call"`
		Name       string `conf:"name"       help:"Name of the page"`
	}
	conf.Load(&config)

	callback := callback(make(chan error, 1))

	client, err := journify.NewWithConfig(config.WriteKey, journify.Config{
		BatchSize: 1,
		Callback:  callback,
	})
	if err != nil {
		fmt.Println("could not initialize journify client", err)
		os.Exit(1)
	}

	defer func() {
		_ = client.Close()
	}()

	switch config.Type {
	case "track":
		_ = client.Enqueue(journify.Track{
			UserId:     config.UserID,
			Event:      config.Event,
			Properties: parseJSON(config.Properties),
		})
	case "identify":
		_ = client.Enqueue(journify.Identify{
			UserId: config.UserID,
			Traits: parseJSON(config.Traits),
		})
	case "group":
		_ = client.Enqueue(journify.Group{
			UserId:  config.UserID,
			GroupId: config.GroupID,
			Traits:  parseJSON(config.Traits),
		})
	case "page":
		_ = client.Enqueue(journify.Page{
			UserId:     config.UserID,
			Name:       config.Name,
			Properties: parseJSON(config.Properties),
		})
	}

	if err := <-callback; err != nil {
		os.Exit(1)
	}
}

// parseJSON parses a JSON formatted string into a map.
func parseJSON(v string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		fmt.Println("could not parse json", v)
		fmt.Println("error:", err)
		os.Exit(1)
	}
	return m
}

// callback implements the journify.Callback interface. It is used by the CLI
// to wait for events to be uploaded before exiting.
type callback chan error

func (c callback) Failure(m journify.Message, err error) {
	fmt.Printf("could not upload message %v due to %v\n", m, err)
	c <- err
}

func (c callback) Success(_ journify.Message) {
	c <- nil
}
