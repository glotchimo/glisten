# gleam

A simple Go framework for building command-driven Twitch chat bots.

## Usage

A full example can be found in `example/main.go`, but we will walk through
the essentials here.

The first step is creating a new bot with `gleam.NewBot`. This method requires
some options for authenticating with Twitch.

```go
bot, err := gleam.NewBot(&gleam.BotOptions{
    Channel: "some channel",
    Username: "some username",
    Password: "some password",
})
if err != nil {
    log.Fatal("error creating new bot: ", err)
}
```

Note that you can make a bot password using [this](https://twitchapps.com/tmi/)
implicit grant flow implementation.

Once the bot's ready, the next step is handler registration using. Handler
functions use a specific signature (`func (m gleam.Message) gleam.Event`) and
are registered using the `bot.AddHandler` method.

```go
bot.AddHandler("!timer", func(m gleam.Message) gleam.Event {
    // ... do some stuff
    return Event{ ... }
}
```

After we've added all our handlers, we go ahead and start the connection and
listening process. Note that the `Bot.Connect` method is blocking and should be
launched in a goroutine.

```go
go bot.Connect()
```

Once we see a confirmation message in the browser tab and the program's logs,
we know we're ready to start listening for events. There are two important
channels to watch: `bot.Events` and `bot.Errors`.

```go
for {
    select {
    case event := <-bot.Events:
        // ... handle events
    case error := <-bot.Errors:
        // ... handle errors
    }
}
```

