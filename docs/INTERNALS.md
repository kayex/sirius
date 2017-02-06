# The internals of Sirius
This page aims to serve as a reference for the implementation and overall architecture of sirius.

## Domain entities

### SlackID
A globally unique user identifier consisting of a combination of the user `UserID` and `TeamID` properties.

### User
A single Slack user account.

### Extension
A single, runnable extension.

### ExtensionConfig
The extension configuration object.

### Connection
A single connection to the RTM API, serving a single user (token).

### Client
Manages a single Connection and its respective User. Responsible for loading and executing extensions.

#### Methods of Interest
```go
func NewClient(user *User) *Client
func (c *Client) Start(ctx context.Context)
```

### MessageAction
Represents an action that an extension wishes to perform on the
current message after execution has finished.
