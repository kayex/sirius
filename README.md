# sirius
Sirius is a standalone Slack extension runner written in Go. It enables you to write and run simple extensions that can read and alter your outgoing messages in realtime.

For example, the `thumbs_up` extension automatically swaps all ocurrences of `(y)` in your messages to `👍` (thumbs up emojii).

## How does it work?
Sirius connects to the [Slack Real Time Messaging API](https://api.slack.com/rtm) using your Slack OAuth token. Once logged in, it monitors your active conversations, making intelligent edits to your messages based on their contents. The gist of the functionality is provided through so-called *extensions*, which are small, stateless functions that are executed with every message you send. Extensions can be enabled and disabled individually.

Sirius is run as a standalone service, and does not have to be run on the same device that you are messaging from. Multiple Slack accounts are supported within the same running instance.

For beta access to the cloud version, please contact the author of this repository.

## Wait, does this mean that Sirius can read all my messages?
Yes. Any message sent or received by your Slack account while Sirius is running will be intercepted via the RTM API and processed by the enabled extensions. However, Sirius does not store any messages or message metadata, and does not collect any message content in its logs. Messages are only kept in memory while the extensions are actively executing.

## Setup
Before starting the service, you need to create a `users.json` file in the same directory as the executable. The file should consist of a single JSON array containing the OAuth tokens for the Slack accounts you wish to enable sirius for.

**users.json**
```json
[
	"xoxp-234234234-23234234-234234324234234-2342343242433"
]
```

## Bundled extensions

### thumbs_up
Converts `(y)` to `👍` (thumbs up emojii) in all outgoing messages.

**kayex** Awesome (y)  
**kayex** Awesome 👍 (edited)

### ripperino
Adds a random ending to any outgoing messages that contain the phrase *ripperino* and nothing else.

**kayex** ripperino  
**kayex** ripperino casino (edited)


## Can I request a new extension?
Of course! Just [submit a new issue](https://github.com/kayex/sirius/issues/new) and make sure to tag it with the `extension` label. You can also submit your own extension for inclusion in the set of default extensions by submitting it as a pull request.

## Creating a new extension
Creating a new extension is only a matter of implementing the `Extension` interface:
```go
type Extension interface {
	Run(Message) (error, MessageAction)
}
```

The `Run` function is called with every outgoing message captured via the RTM API, and should return either an `error` or a `MessageAction`.

`MessageAction`s are returned by extensions to describe changes that should be made to the processed message. This includes things such as editing the message text or deleting it entirely. These changes are accumulated by the extension runner and are broadcasted via the RTM API in timed batches.

Extensions that do not need to modify the message in any way can simply `return NoAction()`.

An extension has exactly **200 ms** to finish execution if it wishes to provide a `MessageAction` other than the `EmptyAction` (as returned by `NoAction()`). Any extensions executing beyond this point will be allowed to finish, but none of the message actions they return will be applied or broadcasted.

### MessageActions
Modifications to the message text are easily described using `TextEditAction`.
```go
func (tu *HelloWorld) Run(m Message) (error, MessageAction) {
	edit := TextEdit()
	
	edit.Substitute("(y)", ":+1:")
	edit.Substitute("(Y)", ":+1:")

	return nil, edit
}
```
