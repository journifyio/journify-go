# Journify SDK for Go

## Installation
```
go get github.com/journifyio/journify-go-sdk
```

## Documentation

You can find more detailed documentation at the following link: https://docs.journify.io/sources/go.

## Usage

```go
package main

import (
    "os"
	"fmt"

    "github.com/journifyio/journify-go-sdk"
)

func main() {
    // Instantiates a client to use send messages to the Journify API.
    client, err := journify.NewWithConfig(os.Getenv("JOURNIFY_WRITE_KEY"), journify.Config{
		BatchSize: 1,
	})
	
	if err != nil {
		fmt.Println("could not initialize journify client", err)
		os.Exit(1)
    }
	
    defer client.Close()
	
    _ = client.Enqueue(journify.Identify{
        UserId: "test-user-12345",
        Traits: map[string]any{
            "email": "joe@doe.com",
            "firstname": "Joe",
            "lastname": "Doe",
        },
    })
	
    _ = client.Enqueue(journify.Track{
        UserId: "test-user-12345",
        Event:  "test-snippet",
    })

    _ = client.Enqueue(journify.Group{
        UserId: "test-user-12345",
        GroupId: "test-group-12345",
        Traits:  map[string]any{
            "name":        "Initech",
            "description": "Accounting Software",
        },
    })

    _ = client.Enqueue(journify.Page{
        UserId:     "test-user-12345",
        Name:       "Home page",
        Properties: journify.NewProperties().
            SetURL("https://journify.io"),
    })
}
```

## License

The library is released under the [MIT license](License.md).
