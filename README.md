# sirius
Sirius is a standalone Slack extension runner written in Go.

## What does it do?
Sirius aims to improve the basic Slack messaging experience with a series of extensions that analyze and react to your outgoing messages in realtime. Extensions have the possibility of making instant alterations to messages you send, performing inline substitutions or adding data.

For example, the `thumbs_up` extension automatically swaps all ocurrences of `(y)` in your messages to `:+1:` (thumbs up emojii).

## How does it work?
Sirius connects to the Slack RTM API using your Slack OAuth token. Once logged in, it actively monitors your outgoing messages, executing extensions and modifying your messages based on pattern matching and word triggers.

## Wait, does this mean that Sirius can read all my messages?
Yes. Any message sent or received by your Slack account while Sirius is running will be intercepted via the RTM API and processed by the enabled extensions. However, Sirius does not store any messages or message metadata, and does not collect any logs. Messages are only kept in memory for the duration of the extensions' execution.

## Creating a new extension
Creating a new extension is only a matter of implementing the `Plugin` interface:
```go
type Plugin interface {
	Run(model.Message) []Transformation
}
```

Every plugin invokation must return a series of zero or more `Transformation`s. The transformations will be applied to the message `Text` property, which will then be broadcasted as a message update via the RTM API.

### thumbs_up
Converts `(y)` to the `:+1:` (thumbs up) emojii in all outgoing messages.

Before
```
Awesome (y)
```

After
```
Awesome :+1:
```

### ripperino
Adds a random extension to any outgoing messages consisting of the phrase *ripperino*

Before
```
ripperino
```

After
```
ripperino casino
```
