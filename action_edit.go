package sirius

import (
	"strings"
)

type TextEditAction struct {
	trans []TextTransform
}

func (edit *TextEditAction) Get() []TextTransform {
	return edit.trans
}

func (m *Message) EditText() *TextEditAction {
	return &TextEditAction{}
}

func (edit *TextEditAction) Substitute(search string, sub string) *TextEditAction {
	edit.add(&SubTransform{
		Search: search,
		Sub:    sub,
	})

	return edit
}

func (edit *TextEditAction) Append(app string) *TextEditAction {
	edit.add(&AppendTransform{
		Appendix: app,
	})

	return edit
}

func (edit *TextEditAction) Perform(msg *Message) error {
	for _, t := range edit.trans {
		msg.Text = t.Apply(msg.Text)
	}

	return nil
}

func (edit *TextEditAction) add(t TextTransform) {
	edit.trans = append(edit.trans, t)
}

/*
TextTransform represents an alteration of the message Text property.
*/
type TextTransform interface {
	Apply(text string) string
}

type SubTransform struct {
	Search string
	Sub    string
}

type AppendTransform struct {
	Appendix string
}

func (st *SubTransform) Apply(text string) string {
	return strings.Replace(text, st.Search, st.Sub, -1)
}

func (at *AppendTransform) Apply(text string) string {
	if len(at.Appendix) == 0 {
		return text
	}

	return text + at.Appendix
}
