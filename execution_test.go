package sirius

import (
	"testing"
	"time"

	"github.com/kayex/sirius/slack"
)

type TestExtension struct {
	duration time.Duration // How long the extension should process for
}

func (x *TestExtension) Run(Message, ExtensionConfig) (MessageAction, error) {
	<-time.After(x.duration)

	return NoAction(), nil
}

func TestExecutor_Run_RespectsTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping AsyncRunner execution timeout test in short mode.")
	}

	fast := &TestExtension{
		duration: 0,
	}

	slow := &TestExtension{
		duration: time.Millisecond * 10,
	}

	msg := NewMessage(slack.UserID{TeamID: "abc", UserID: "123"}, "Test message.", "#1337", "0")
	exs := []ConfigExtension{
		*NewConfigExtension(fast, nil),
		*NewConfigExtension(slow, nil),
	}
	exe := NewExecutor(nil, time.Millisecond*1)
	res := make(chan ExecutionResult, 1)

	exe.Run(msg, exs, res)

	count := 0

	for range res {
		count++
		if count > 1 {
			t.Fatal("Expected only 1 ExecutionResult but received 2")
		}
	}

	if count == 0 {
		t.Fatal("Expected 1 ExecutionResult but received none")
	}
}
