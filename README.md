# KwikDesk Partner Platform Go Client

This is a client to interact with the KwikDesk Partner Platform using the Go Language. 

## Partner Documentation

The partner documentation can be found at [partners.kwikdesk.com](https://partners.kwikdesk.com).

## Usage

This package is a barebone straightforward package that allows you to complete the lifecycle of a
message on Kwikdesk from [token](https://partners.kwikdesk.com/#token) creation to [message](https://partners.kwikdesk.com/#message) creation to [search](https://partners.kwikdesk.com/#search) and [secure channel](https://partners.kwikdesk.com/#channel) communications.

### Create a Token

When creating a token, we need to pass either an **application name** or if we wish to be contacted at some point if needs be, an *email address*.

```go
import "github.com/kwikdesk/kwikdesk"

client    := kwikdesk.NewClient("") // We pass an empty token here.
token, _  := kwikdesk.CreateToken("MyApplicationName")
```

**Make sure not to loose your token as it is not retrievable for privacy reasons**


The `Token` API documentation can be found [here](https://partners.kwikdesk.com/#token).

### Create Message

Let's assume our token is `token-token-token`. 


First we are going to create a message that is searchable (Not private — **third parameter**).

```go
import "github.com/kwikdesk/kwikdesk"

client    := kwikdesk.NewClient("token-token-token")
message, _  := kwikdesk.Messages("Content of my message #testing", 1440, false)

fmt.Printf("%v", message)
```

This created a public message (public but only associated with your account)

Then we want to create a message that will only be retrievable via **secure-channels**:


```go
import "github.com/kwikdesk/kwikdesk"

client      := kwikdesk.NewClient("token-token-token")
private, _  := kwikdesk.Messages("Private message. No Hashtags Required.", 1440, true)

fmt.Printf("%v", private)
```

The `Messages` API documentation page can be found [here](https://partners.kwikdesk.com/#create).

### Search For Messages

Now the search allows you to search for messages associated with your application token and
that have the **private** flag set to false. You can execute a search like this:

```go
import "github.com/kwikdesk/kwikdesk"

client     := kwikdesk.NewClient("token-token-token")
search, _  := kwikdesk.Search("testing")

fmt.Printf("%v", search)
```

You can then iterate over the map of interfaces.

The `Search` API documentation page can be found [here](https://partners.kwikdesk.com/#search).

### Secure Communication Channel

Now that you've created messages that can't be searched for, you can only retrieve them by using 
the channels. All  you have to do is:

```go
import "github.com/kwikdesk/kwikdesk"

client      := kwikdesk.NewClient("token-token-token")
channel, _  := kwikdesk.Channel()

fmt.Printf("%v", channel)
```

The `Channel` API documentation page can be found [here](https://partners.kwikdesk.com/#channel).

# Example

If you are interested in an example that is fully runnable from A to Z, you can check out the [examples/main.go](examples/main.go) and run it. 

It will **generate** a token, **create** a searchable message, **create** a private message, **search** for the public message, and **retrieve the channel** messages.
# License

Copyright (c) KwikDesk. All rights reserved. See the [LICENSE](LICENSE) file.
