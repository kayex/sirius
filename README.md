# sirius ⚡
Slack extensions that improve your outgoing messages in real-time. Written in Go.

For example, the `thumbs_up` extension automatically swaps all occurrences of `(y)` in your messages for `👍` (the thumbs up emojii).

*The latest release notes can be found [here](https://github.com/kayex/sirius/releases).*

## How does it work?
Sirius runs as a service and connects to the [Slack Real Time Messaging API](https://api.slack.com/rtm) using your Slack OAuth token. Once logged in, it monitors your active conversations and automatically makes intelligent edits to your messages.

## Extensions

### thumbs_up
Converts `(y)` to `👍` in all outgoing messages.

>**kayex**: Awesome (y)

⚡

>**kayex**: Awesome 👍


### geocode
Type `!address` followed by any sort of geographical location, and the `geocode` extension will fetch the exact address and coordinates for you.

>**kayex**: !address Empire State Building

⚡

>**kayex**: **350 5th Ave, New York, NY 10118, USA**  
`(40.748441, -73.985664)`

### ip-lookup
Type `!ip-lookup` followed by an IP address to fetch related geolocation information.

>**kayex**: !ip 8.8.8.8  
>**kayex**: !ip 2001:4860:4860::8888 // IPv6

⚡

>**kayex**: `8.8.8.8`  
Mountain View, United States (`US`)  
Google  
>**kayex**: `2001:4860:4860::8888`  
Chicago, United States (`US`)

### quotes
Avoids breaking blockquotes that contain newlines.

>**kayex**: >This is  
           a multi-line  
	   quote.

⚡

>**kayex**: >This is  
           >a multi-line  
	   >quote.


## Setup
***A cloud version of Sirius is coming soon, which will allow you to use sirius without installing anything.***

### Building
```
$ go build github.com/kayex/sirius/cmd/sirius
```

### Running
Sirius is run as a standalone service, which means it does *not* have to be run on the same device that you are messaging from.

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

## Can I request a new extension?
Of course! Just [submit a new issue](https://github.com/kayex/sirius/issues/new) and make sure to tag it with the `extension` label. You can also submit your own extension by creating a pull request.

## Creating a new extension
Creating a new extension is only a matter of implementing the `Extension` interface:
```go
package sirius

type Extension interface {
	Run(Message, ExtensionConfig) (error, MessageAction)
}
```

The `Run` function is called with every outgoing message captured via the RTM API, and should return either an `error` or a `MessageAction`. Each extension is passed an `ExtensionConfig` with every invocation, which is a read-only key/value configuration store. The `ExtensionConfig` is unique per user and extension, and is used to access all extension-specific configuration values.

`MessageAction`s are returned by extensions to describe changes that should be made to the processed message. This includes things such as editing the message text or deleting the message entirely. Actions are accumulated by the extension runner and broadcasted via the RTM API in timed batches.

Extensions that do not need to modify the message in any way can simply `return nil, NoAction()`.

An extension has exactly **2000 ms** to finish execution if it wishes to provide a `MessageAction` (other than `NoAction()`). Extensions executing past this deadline will be allowed to finish, but any message actions returned will be discarded.

### `MessageAction`s
These are the default `MessageActions`. New actions can be created by implementing the `MessageAction` interface:
```go
type MessageAction interface {
	Perform(*Message) error
}
```

#### TextEditAction
Modifications to the message text are easily done using `(*Message) EditText()`:
```go
func (*ThumbsUp) Run(m Message, cfg ExtensionConfig) (error, MessageAction) {
	edit := m.EditText()
	
	edit.Substitute("(y)", ":+1:")
	edit.Substitute("(Y)", ":+1:")

	return nil, edit
}
```
