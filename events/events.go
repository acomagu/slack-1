package events

import (
	"encoding/json"
	"fmt"

	"github.com/lestrrat-go/slack/objects"
	"github.com/pkg/errors"
)

func (e *Event) UnmarshalJSON(data []byte) error {
	// Event contains an "item" field whose structure is unknown
	// until we actually receive the event. We could just
	// let json.Unmarshal to deserialize it into a map[string]interface{}
	// But in Go, that's just ugly.
	//
	// So we use a proxy that stores json.RawMessage, and
	// decouple the deserialization and composition of the
	// Event struct
	fmt.Println("OPPAIs222")
	var p eventUnmarshalProxy
	if err := json.Unmarshal(data, &p); err != nil {
		return errors.Wrap(err, `failed to unmarshal payload`)
	}

	return p.Populate(e)
}

func (p *eventUnmarshalProxy) Populate(e *Event) error {
	var item interface{}
	fmt.Println("OPPAI2")

	// TODO: Add more types
	switch p.Type {
	case MessageChannelsType, MessageGroupsType, MessageImType, MessageMpimType:
		item = &objects.Message{}
	case URLVerificationType:
		fmt.Println("OPPAI")
		item = &objects.Challenge{}
	default:
		return errors.Errorf("unknown event type: %s", p.Type)
	}

	fmt.Println("data:", string(p.Item))
	if err := json.Unmarshal(p.Item, item); err != nil {
		return errors.Wrap(err, `failed to unmarshal event item`)
	}
	*e = Event{}
	e.EventTimestamp = p.EventTimestamp
	e.Item = item
	e.Timestamp = p.Timestamp
	e.Type = p.Type
	e.User = p.User

	return nil
}
