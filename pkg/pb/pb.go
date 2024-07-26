package pb

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Message[T any] interface {
	proto.Message
	*T
}

func opts() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		Multiline:         true,
		Indent:            "  ",
		AllowPartial:      true,
		UseProtoNames:     true,
		UseEnumNumbers:    false,
		EmitUnpopulated:   true,
		EmitDefaultValues: true,
	}
}

// WriteMsg - write proto message to io.Writer as JSON object.
func WriteMsg[T any, M Message[T]](w io.Writer, msg M) error {
	if buf, err := opts().Marshal(msg); err == nil {
		if _, err := w.Write(buf); err != nil {
			return err
		}
	} else {
		return err
	}

	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

// WriteMsgArray - write array of proto messages to io.Writer as a JSON array.
func WriteMsgArray[T any, M Message[T]](w io.Writer, msgs []M) error {
	f := false

	if _, err := w.Write([]byte("[")); err != nil {
		return err
	}

	for _, msg := range msgs {
		if f {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		}

		buf, err := opts().Marshal(msg)
		if err != nil {
			return err
		}

		if _, err := w.Write(buf); err != nil {
			return err
		}

		if !f {
			f = true
		}
	}

	if _, err := w.Write([]byte("]\n")); err != nil {
		return err
	}

	return nil
}
