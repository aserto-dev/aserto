package pb

import (
	"bytes"
	"encoding/json"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoToBuf marshal protomessage to buffer.
func ProtoToBuf(w io.Writer, msg proto.Message) error {
	b, err := protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

// BufToProto unmarshal buffer to protomessage.
func BufToProto(r io.Reader, msg proto.Message) error {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(buf.Bytes(), msg)
}

// UnmarshalNext unmarshal next protomessage in stream.
func UnmarshalNext(d *json.Decoder, m proto.Message) error {
	var b json.RawMessage
	if err := d.Decode(&b); err != nil {
		return err
	}
	return protojson.Unmarshal(b, m)
}

// ProtoToStr marshal protomessage to string.
func ProtoToStr(msg proto.Message) string {
	return protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}.Format(msg)
}
