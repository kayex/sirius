# sirius
Sirius is a standalone Slack extension runner written in Go. It enables you to write and run simple extensions that analyze and alter your outgoing messages in realtime.

## What does it do?
Sirius aims to improve the basic Slack messaging experience with a series of extensions that add powerful editing capabilities. Extensions examine every message you send, and can perform inline substitutions or add data.

For example, the `thumbs_up` extension automatically swaps all ocurrences of `(y)` in your messages to `:+1:` (thumbs up emojii).

## How does it work?
Sirius connects to the Slack RTM API using your Slack OAuth token. Once logged in, it actively monitors your outgoing messages, executing extensions which read and modify your messages.

## Wait, does this mean that Sirius can read all my messages?
Yes. Any message sent or received by your Slack account while Sirius is running will be intercepted via the RTM API and processed by the enabled extensions. However, Sirius does not store any messages or message metadata, and does not collect any logs. Messages are only kept in memory while the extensions are actively executing.

## Creating a new extension
Creating a new extension is only a matter of implementing the `Extension` interface:
```go
type Extension interface {
	Run(Message) (error, MessageAction)
}
```

Every extension invokation must return a `MessageAction`, which will be applied to the message. If any modifications are made to the message text property, the updated message is broadcasted via the RTM API. Extensions that do not want to act upon the message in any way can simply return the `EmptyAction`:
```go
return NoAction()
```

Each extension is run concurrently and has a generous execution time limit. In addition to this, extensions may perform any type of I/O, including network requests. Message updates are batched on a fixed time interval, which allows multiple extensions to send near instantaneous updates without exceeding API limits.

## Bundled plugins

### thumbs_up
Converts `(y)` to `👍` (thumbs up emojii) in all outgoing messages.

**kayex** Awesome (y)  
**kayex** Awesome 👍 (edited)

### ripperino
Adds a random ending to any outgoing messages that contain the phrase *ripperino* and nothing else.

**kayex** ripperino  
**kayex** ripperino casino (edited)
