package live

import (
	"encoding/base64"
	"encoding/json"
)

// EventConfig configures an event.
type EventConfig func(e *Event) error

const (
	// EventError indicates an error has occured.
	EventError = "err"
	// EventPatch a patch event containing a diff.
	EventPatch = "patch"
	// EventAck sent when an event is ackknowledged.
	EventAck = "ack"
	// EventConnect sent as soon as the server accepts the
	// WS connection.
	EventConnect = "connect"
	// EventParams sent for a URL parameter update. Can be
	// sent both directions.
	EventParams = "params"
	// EventRedirect sent in order to trigger a browser
	// redirect.
	EventRedirect = "redirect"

	EventUpload = "allow_upload"
)

// Event messages that are sent and received by the
// socket.
type Event struct {
	T        string          `json:"t"`
	ID       int             `json:"i,omitempty"`
	Data     json.RawMessage `json:"d,omitempty"`
	SelfData interface{}     `json:"-"`
}

// Params extract params from inbound message.
func (e Event) Params() (Params, error) {
	if e.Data == nil {
		return Params{}, nil
	}
	var p Params
	if err := json.Unmarshal(e.Data, &p); err != nil {
		return nil, ErrMessageMalformed
	}
	return p, nil
}

type FileTest struct {
	File struct {
		LastModified     string `json:"-"`
		LastModifiedDate string `json:"-"`
		Name             string `json:"name"`
		Size             int    `json:"Size"`
		Type             string `json:"-"`
	} `json:"file"`
	Field string `json:"field"`
	Chunk string `json:"chunk"`
}

type FileMeta struct {
	LastModified     string
	LastModifiedDate string
	Name             string
	Size             int
	Type             string
}

type FileTest2 struct {
	File  FileMeta
	Field string
	Chunk []byte
}

// Params extract data from inbound message.
func (e Event) File() (*FileTest2, error) {
	if e.Data == nil {
		return nil, nil
	}
	// p := []byte(e.Data)

	var p FileTest
	err := json.Unmarshal(e.Data, &p)
	if err != nil {
		return nil, err
	}

	b, err := base64.StdEncoding.DecodeString(p.Chunk)
	if err != nil {
		panic(err)
	}
	// p = []byte(dst)

	d := &FileTest2{
		Chunk: b,
		File: FileMeta{
			Name: p.File.Name,
			Size: p.File.Size,
		},
		Field: p.Field,
	}

	// dst := make([]byte, len(p)*len(p)/base64.StdEncoding.DecodedLen(len(p)))
	// _, err := base64.StdEncoding.Decode(dst, []byte(p))
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return []byte{}, nil
	// }

	return d, nil
}

// WithID sets an ID on an event.
func WithID(ID int) EventConfig {
	return func(e *Event) error {
		e.ID = ID
		return nil
	}
}

type ErrorEvent struct {
	Source Event  `json:"source"`
	Err    string `json:"err"`
}
