# sirius ⚡
![](https://travis-ci.org/kayex/sirius.svg?branch=develop)

An extension server for Slack.

Extensions are small pieces of user-created functionality that enhance the default Slack experience. For example, the `thumbs_up` extension converts `(y)` to `👍` in all outgoing messages.

Sirius is an extension server that allows you to run these extensions as a standalone service, and use them without having to install anything on the devices you use Slack from.

## How does it work?
Sirius connects to the [Slack Real Time Messaging API](https://api.slack.com/rtm) using your Slack OAuth credentials. Once logged in, it monitors your active channels for messages and processes them using the enabled extensions. Only messages sent by you are processed. Modifications to the message text body are instantly propagated using a regular message edit.

Sirius does not store or forward any message data, including metadata.

## Bundled extensions
These extensions come included with sirius by default, and can be enabled immediately using their EIDs (names).

### thumbs_up
Converts `(y)` to `👍` in all outgoing messages.

>**kayex**: Awesome (y)

⚡

>**kayex**: Awesome 👍

--

### geocode
Type `!address` followed by any sort of geographical location, and the `geocode` extension will fetch the exact address and coordinates for you.

>**kayex**: !address Empire State Building

⚡

>**kayex**: **350 5th Ave, New York, NY 10118, USA**  
`(40.748441, -73.985664)`

--

### ip_lookup
Type `!ip` followed by an IP address to fetch related geolocation information.

>**kayex**: !ip 8.8.8.8  
>**kayex**: !ip 2001:4860:4860::8888 // IPv6

⚡

>**kayex**: `8.8.8.8`  
Mountain View, United States (`US`)  
Google  
>**kayex**: `2001:4860:4860::8888`  
Chicago, United States (`US`)  
Google

--

### quotes
Avoids breaking blockquotes that contain newlines.

>**kayex**: >This is  
           a multi-line  
	   quote.

⚡

>**kayex**: >This is  
           >a multi-line  
	   >quote.

--

## Getting started
Sirius is available as a free, hosted service at http://adsa.se/sirius.

You can also run Sirius yourself by following the instructions below.

### Building
```
$ go build github.com/kayex/sirius/cmd/sirius-local
```

### Running
Before starting the service, a `users.json` file needs to be created in the same directory as the executable will run from. The file should consist of a single JSON array containing OAuth tokens for the Slack accounts that Sirius should be enabled for:

**users.json**
```json
[
	"xoxp-234234234-23234234-234234324234234-2342343242433"
]
```

You can then start the main service by simply running the `sirius-local` executable:
```
$ ./sirius-local
```

## Can I request a new extension?
Of course! Just [submit a new issue](https://github.com/kayex/sirius/issues/new) and make sure to tag it with the `extension` label. You can also submit your own extension by creating a pull request.

## Creating a new extension
Creating a new extension is only a matter of implementing the `Extension` interface:
```go
package sirius

type Extension interface {
	Run(Message, ExtensionConfig) (MessageAction, error)
}
```

The `Run` function is called with every outgoing message captured via the RTM API, and should return either an `error` or a `MessageAction`. Each extension is passed an `ExtensionConfig` with every invocation, which is a read-only key/value configuration store. The `ExtensionConfig` is unique per user and extension, and is used to access all extension-specific configuration values.

`MessageAction`s are returned by extensions to describe changes that should be made to the processed message. This includes things such as editing the message text or deleting the message entirely. Actions are accumulated by the extension runner and broadcasted via the RTM API in timed batches.

Extensions that do not need to modify the message in any way can simply `return NoAction(), nil`.

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
func (*ThumbsUp) Run(m Message, cfg ExtensionConfig) (MessageAction, error) {
	edit := m.EditText()
	
	edit.Substitute("(y)", ":+1:")
	edit.Substitute("(Y)", ":+1:")

	return edit, nil
}
```
