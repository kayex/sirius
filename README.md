# sirius
Sirius is a standalone Slack extension runner written in Go. It enables you to write and run simple extensions that can read and alter your outgoing messages in realtime.

For example, the `thumbs_up` extension automatically swaps all ocurrences of `(y)` in your messages to `👍` (thumbs up emojii).

*The latest release notes can be found [here](https://github.com/kayex/sirius/releases/tag/0.4.0).*

## How does it work?
Sirius connects to the [Slack Real Time Messaging API](https://api.slack.com/rtm) using your Slack OAuth token. Once logged in, it monitors your active conversations, making intelligent edits to your messages based on their contents.

## Included extensions

### thumbs_up
Converts `(y)` to `👍` (thumbs up emojii) in all outgoing messages.

*before*

>**kayex** Awesome (y)  

*after*

>**kayex** Awesome 👍 (edited)  

### quotes
Avoids breaking blockquotes when the quote contains newlines.

*before*

>**kayex** >This is  
           a multi-paragraph  
	   quote.  
	     
*after*

>**kayex** >This is  
           >a multi-paragraph  
	   >quote. (edited)


## Setup and running
Sirius is run as a standalone service, which means it does *not* have to be run on the same device that you are messaging from. *A cloud version of Sirius is coming soon.*

Before starting the service, you need to create a `users.json` file in the same directory as the executable. The file should consist of a single JSON array containing OAuth tokens for the Slack accounts you wish to enable Sirius for:

**users.json**
```json
[
	"xoxp-234234234-23234234-234234324234234-2342343242433"
]
```

You can then start the main service by simply running the `sirius` executable:
```
$ ./sirius
```

### Users with multiple Slack accounts
Sirius has multi-account support, which means that you can use Sirius with any number of Slack accounts without having to run multiple service instances. Simply add one OAuth token per account that you wish to use with Sirius to the `users.json` file.

## Can I request a new extension?
Of course! Just [submit a new issue](https://github.com/kayex/sirius/issues/new) and make sure to tag it with the `extension` label. You can also submit your own extension for inclusion in the set of default extensions, by submitting it as a pull request.

## Creating a new extension
Creating a new extension is only a matter of implementing either the `Extension` interface:
```go
package sirius

type Extension interface {
	Run(Message, ExtensionConfig) (error, MessageAction)
}
```

The `Run` function is called with every outgoing message captured via the RTM API, and should return either an `error` or a `MessageAction`. It is passed an `ExtensionConfig` with every invocation, which is a read-only key/value configuration store. The `ExtensionConfig` is unique per user and extension.

`MessageAction`s are returned by extensions to describe changes that should be made to the processed message. This includes things such as editing the message text, or deleting the message entirely. These changes are accumulated by the extension runner and broadcasted via the RTM API in timed batches.

Extensions that do not need to modify the message in any way can simply `return NoAction()`.

An extension has exactly **200 ms** to finish execution if it wishes to provide a `MessageAction` other than the `EmptyAction` (as returned by `NoAction()`). Extensions that fail to complete execution before this deadline will be allowed to finish, but none of the message actions they return will be applied to the message or broadcasted via the API.

### Standard MessageActions
These are the default `MessageActions`. New actions can be created by implementing the `MessageAction` interface:
```go
type MessageAction interface {
	Perform(*Message) error
}
```

#### TextEditAction
Modifications to the message text are easily done using `(*Message) EditText()`.
```go
func (*ThumbsUp) Run(m Message, cfg ExtensionConfig) (error, MessageAction) {
	edit := m.EditText()
	
	edit.Substitute("(y)", ":+1:")
	edit.Substitute("(Y)", ":+1:")

	return nil, edit
}
```
