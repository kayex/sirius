package sirius

import (
	"strings"
)

/*
TextEditAction represents a series of
edits to the message Text property

Usage:

var msg Message

edit := msg.EditText()
edit.Substitute("foo", "bar")
edit.Append("-ending")

 */
type TextEditAction struct {
	mutations []TextMutation
}

func (*Message) EditText() *TextEditAction {
	return &TextEditAction{}
}

func (edit *TextEditAction) Perform(msg *Message) error {
	for _, m := range edit.mutations {
		msg.Text = m.Apply(msg.Text)
	}

	return nil
}

func (edit *TextEditAction) Substitute(search string, sub string) *TextEditAction {
	edit.add(&SubMutation{
		Search: search,
		Sub:    sub,
	})

	return edit
}

func (edit *TextEditAction) Append(app string) *TextEditAction {
	edit.add(&AppendMutation{
		Appendix: app,
	})

	return edit
}

func (edit *TextEditAction) add(m TextMutation) {
	edit.mutations = append(edit.mutations, m)
}

/*
TextTransform represents a string mutation
*/
type TextMutation interface {
	Apply(text string) string
}

type SubMutation struct {
	Search string
	Sub    string
}

type AppendMutation struct {
	Appendix string
}

func (st *SubMutation) Apply(text string) string {
	return strings.Replace(text, st.Search, st.Sub, -1)
}

func (at *AppendMutation) Apply(text string) string {
	if len(at.Appendix) == 0 {
		return text
	}

	return text + at.Appendix
}
