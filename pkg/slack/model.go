package slack

import "encoding/json"

func UnmarshalSlack(data []byte) (Model, error) {
	var r Model
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Model) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Model struct {
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Color  string  `json:"color"`
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type   string  `json:"type"`
	Text   *Text   `json:"text,omitempty"`
	Fields *[]Text `json:"fields,omitempty"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
