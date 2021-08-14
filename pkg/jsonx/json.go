package jsonx

import (
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func OutputJSON(w io.Writer, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s\n", string(b))

	return nil
}

func OutputJSONPB(w io.Writer, msg proto.Message, opts ...protojson.MarshalOptions) error {
	options := DefaultMarshalOpts()
	if len(opts) == 1 {
		options = opts[0]
	}

	b, err := options.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func EncodeJSONPB(w io.Writer, msg proto.Message, opts ...protojson.MarshalOptions) error {
	options := DefaultMarshalOpts()
	if len(opts) == 1 {
		options = opts[0]
	}

	b, err := options.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func DefaultMarshalOpts() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}
}

func MaskedMarshalOpts() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
}
