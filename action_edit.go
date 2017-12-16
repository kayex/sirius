package sirius

import (
	"github.com/kayex/sirius/text"
)

// TextEditAction represents a series of
// modifications to the message Text property
type TextEditAction struct {
	mutations []text.Mutation
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

func (edit *TextEditAction) Set(txt string) *TextEditAction {
	edit.add(&text.Set{
		Text: txt,
	})

	return edit
}

func (edit *TextEditAction) Substitute(search string, sub string) *TextEditAction {
	edit.add(&text.Sub{
		Search: search,
		Sub:    sub,
	})

	return edit
}

func (edit *TextEditAction) SubstituteQuery(q text.Query, sub string) *TextEditAction {
	edit.add(&text.SubQuery{
		Search: q,
		Sub:    sub,
	})

	return edit
}

func (edit *TextEditAction) Append(app string) *TextEditAction {
	edit.add(&text.Append{
		Appendix: app,
	})

	return edit
}

func (edit *TextEditAction) Prepend(pre string) *TextEditAction {
	edit.add(&text.Prepend{
		Prefix: pre,
	})

	return edit
}

func (edit *TextEditAction) add(m text.Mutation) {
	edit.mutations = append(edit.mutations, m)
}
