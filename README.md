# sirius

Sirius is a standalone Slack extensions runner written in Go.

## How does it work?
Sirius connects to the Slack RTM API using your Slack OAuth token. Once logged in, it actively monitors your outgoing messages for triggers, i.e. words or syntaxes that trigger the execution of an extension. An extension may optionally modify the message that triggered it, which will push a message edit via the RTM API.

## Wait, does this mean that Sirius can read all my messages?
Yes. Any message sent or received by your Slack account while Sirius is running will be intercepted via the RTM API and processed. However, Sirius does not store any messages or message metadata, and does not collect any logs. Messages are only kept in memory for the duration of the extension execution.

## Extensions

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
