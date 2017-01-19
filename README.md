# sirius

Sirius is a standalone Slack extensions runner written in Go.

## How does it work?
Sirius connects to the Slack RTM API using your Slack OAuth token. Once logged in, it actively monitors your outgoing messages for triggers, i.e. words or syntaxes that trigger the execution of an extension. An extension may optionally modify the message that triggered it, which will push a message edit via the RTM API.

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
