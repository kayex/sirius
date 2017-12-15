package sirius

import (
	"testing"
	"time"

	"github.com/kayex/sirius/slack"
)

type TestExtension struct {
	duration time.Duration // How long the extension should execute for
}

func (x *TestExtension) Run(Message, ExtensionConfig) (MessageAction, error) {
	<-time.After(x.duration)

	return NoAction(), nil
}

func TestAsyncRunner_Run_RespectsTimeout(t *testing.T) {
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
	exe := []Execution{
		*NewExecution(fast, nil),
		*NewExecution(slow, nil),
	}
	r := NewAsyncRunner()
	res := make(chan ExecutionResult, 1)

	r.Run(msg, exe, res, time.Millisecond*1)

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
