# sirius
Sirius is a standalone Slack extension runner written in Go. It enables you to write and run simple extensions that analyze and alter your outgoing messages in realtime.

## What does it do?
Sirius aims to improve the basic Slack messaging experience with a series of extensions that add powerful editing capabilities. Extensions examine every message you send, and can perform inline substitutions or add data.

For example, the `thumbs_up` extension automatically swaps all ocurrences of `(y)` in your messages to `:+1:` (thumbs up emojii).

## How does it work?
Sirius connects to the Slack RTM API using your Slack OAuth token. Once logged in, it actively monitors your outgoing messages, executing extensions which read and modify your messages.

## Wait, does this mean that Sirius can read all my messages?
Yes. Any message sent or received by your Slack account while Sirius is running will be intercepted via the RTM API and processed by the enabled extensions. However, Sirius does not store any messages or message metadata, and does not collect any logs. Messages are only kept in memory for the duration of the extensions' execution.

## Creating a new extension
Each extension is run concurrently and has a generous execution time limit. In addition to this, extensions may perform any type of I/O, including network requests. Message updates are batched on a fixed time interval, which allows quick executing extensions to send their message modifications even though there are slower extensions that haven't yet completed.

Creating a new extension is only a matter of implementing the `Plugin` interface:
```go
type Plugin interface {
	Run(model.Message) []Transformation
}
```

Every plugin invokation must return a slice of zero or more `Transformation`s. The transformations will be applied to the message `Text` property, which will then be broadcasted as a message update via the RTM API.

## Bundled plugins

### thumbs_up
Converts `(y)` to the `:+1:` (thumbs up) emojii in all outgoing messages.

Before
```
kayex: Awesome (y)
```

After
```
kayex: Awesome :+1: (edited)
```

### ripperino
Adds a random ending to any outgoing messages that contain the phrase *ripperino* and nothing else.

Before
```
kayex: ripperino
```

After
```
kayex: ripperino casino (edited)
```
