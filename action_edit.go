package sirius

import (
	"strings"
)

// TextEditAction represents a series of
// modifications to the message Text property
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

func (edit *TextEditAction) ReplaceWith(replacement string) *TextEditAction {
	edit.add(&ReplaceMutation{
		Replacement: replacement,
	})

	return edit
}

func (edit *TextEditAction) Substitute(search string, sub string) *TextEditAction {
	edit.add(&SubMutation{
		Search: search,
		Sub:    sub,
	})

	return edit
}

func (edit *TextEditAction) SubstituteWord(search string, sub string) *TextEditAction {
	edit.add(&SubWordMutation{
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

func (edit *TextEditAction) Prepend(pre string) *TextEditAction {
	edit.add(&PrependMutation{
		Prefix: pre,
	})

	return edit
}

func (edit *TextEditAction) add(m TextMutation) {
	edit.mutations = append(edit.mutations, m)
}

type TextMutation interface {
	Apply(text string) string
}

type ReplaceMutation struct {
	Replacement string
}

type SubMutation struct {
	Search string
	Sub    string
}

type SubWordMutation struct {
	Search string
	Sub    string
}

type AppendMutation struct {
	Appendix string
}

type PrependMutation struct {
	Prefix string
}

func (rm *ReplaceMutation) Apply(text string) string {
	return rm.Replacement
}

func (sm *SubMutation) Apply(text string) string {
	return strings.Replace(text, sm.Search, sm.Sub, -1)
}

func (sm *SubWordMutation) Apply(text string) string {
	if text == sm.Search {
		return strings.Replace(text, sm.Search, sm.Sub, -1)
	}

	if strings.HasPrefix(text, sm.Search+" ") {
		text = sm.Sub + text[len(sm.Search):]
	}

	if strings.HasSuffix(text, " "+sm.Search) {
		text = text[:len(text)-len(sm.Search)] + sm.Sub
	}

	return text
}

func (am *AppendMutation) Apply(text string) string {
	if len(am.Appendix) == 0 {
		return text
	}

	return text + am.Appendix
}

func (pm *PrependMutation) Apply(text string) string {
	if len(pm.Prefix) == 0 {
		return text
	}

	return pm.Prefix + text
}
