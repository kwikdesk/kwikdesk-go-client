package main

import (
    "github.com/kwikdesk/kwikdesk-go-client/kwikdesk"
    "crypto/rand"
    "fmt"
)

func Load() (client *kwikdesk.Client) {
    client = kwikdesk.NewClient("")
    return client
}

func main() {
    // Here we create an empty Client Object and use it to create things.
    var (
        client = Load()
        privateString = randString(10)
        publicMessage = fmt.Sprintf("Token-Level Searchable Content %v #test", privateString)
        privateMessage = fmt.Sprintf("Private Channel-Only Content %v #test", privateString)
    )


    fmt.Println("Initiating Workflow Example")
    fmt.Println("=====================================================================")
    fmt.Println("")

    // 0. Create a token to create messages 
    fmt.Println("Create a token to use in future message and channel creation:")
    fmt.Println("-------------------------------------------------------------")
    token, _ := client.CreateToken("YourAppName")
    fmt.Printf("%v", token)

    fmt.Println("")
    fmt.Println("")

    // 1. Post a message with your token. This message is searchable.
    fmt.Println("Posting a new message that is searchable and only associated to your token:")
    fmt.Println("---------------------------------------------------------------------------")
    post, _ := client.Messages(publicMessage, 100000, false)
    fmt.Printf("%v", post)

    fmt.Println("")
    fmt.Println("")

    // 2. Execute a search for the newly created message associated with your token.
    fmt.Println("Searching for the created message:")
    fmt.Println("----------------------------------")
    search, _ := client.Search("test")
    fmt.Printf("%v", search)

    fmt.Println("")
    fmt.Println("")

    // 3. Post a message that is private. Not searchable but only retrievable
    //    via the `Channel` call. Third parameter is the private flag.

    fmt.Println("Adding a private message that should not appear in the search:")
    fmt.Println("--------------------------------------------------------------")
    private, _ := client.Messages(privateMessage, 100000, true)
    fmt.Printf("%v", private)

    fmt.Println("")
    fmt.Println("")

    // 4. Do another search to show that #test doesn't appear when it's private.
    fmt.Println("Prove the last assertion by not seeing the private message we created.")
    fmt.Println("----------------------------------------------------------------------")
    search_privacy, _ := client.Search("test")
    fmt.Printf("%v", search_privacy)

    fmt.Println("")
    fmt.Println("")

    // 5. Retrieve the channel messages to show the private message is there.
    fmt.Println("Now we retrieve the private message through a channel")
    fmt.Println("-----------------------------------------------------")
    channel, _ := client.Channel()
    fmt.Printf("%v", channel)

    fmt.Println("")
    fmt.Println("")
}

// For the purpose of this example, we have a random generator.
func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}
