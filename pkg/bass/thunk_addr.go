package bass

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/vito/bass/pkg/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type ThunkAddr struct {
	Thunk  Thunk
	Port   string
	Format string
}

var _ Value = ThunkAddr{}

func (value ThunkAddr) String() string {
	return fmt.Sprintf("%s:%s", value.Thunk, value.Port)
}

func (value ThunkAddr) Equal(other Value) bool {
	var o ThunkAddr
	return other.Decode(&o) == nil &&
		value.Thunk.Equal(o.Thunk) &&
		value.Port == o.Port &&
		value.Format == o.Format
}

func (value ThunkAddr) Decode(dest any) error {
	switch x := dest.(type) {
	case *ThunkAddr:
		*x = value
		return nil
	case *Value:
		*x = value
		return nil
	case Decodable:
		return x.FromValue(value)
	default:
		return DecodeError{
			Source:      value,
			Destination: dest,
		}
	}
}

// Eval returns the value.
func (value ThunkAddr) Eval(_ context.Context, _ *Scope, cont Cont) ReadyCont {
	return cont.Call(value, nil)
}

func (value ThunkAddr) Render(scope *Scope) (string, error) {
	var errs error
	rendered := os.Expand(value.Format, func(name string) string {
		bnd := Symbol(name)

		val, found := scope.Get(bnd)
		if !found {
			errs = multierror.Append(errs, UnboundError{
				Symbol: bnd,
				Scope:  scope,
			})
		}

		var str string
		if err := val.Decode(&str); err != nil {
			str = val.String()
		}

		return str
	})
	if errs != nil {
		return "", errs
	}

	return rendered, nil
}

var _ ProtoMarshaler = ThunkAddr{}

func (value ThunkAddr) MarshalProto() (proto.Message, error) {
	t, err := value.Thunk.MarshalProto()
	if err != nil {
		return nil, fmt.Errorf("thunk: %w", err)
	}

	return &proto.ThunkAddr{
		Thunk:  t.(*proto.Thunk),
		Port:   value.Port,
		Format: value.Format,
	}, nil
}

func (value *ThunkAddr) UnmarshalProto(msg proto.Message) error {
	p, ok := msg.(*proto.ThunkAddr)
	if !ok {
		return fmt.Errorf("unmarshal proto: have %T, want %T", msg, p)
	}

	if err := value.Thunk.UnmarshalProto(p.Thunk); err != nil {
		return err
	}

	value.Port = p.GetPort()
	value.Format = p.GetFormat()

	return nil
}

func (value ThunkAddr) MarshalJSON() ([]byte, error) {
	msg, err := value.MarshalProto()
	if err != nil {
		return nil, err
	}

	return protojson.Marshal(msg)
}

func (value *ThunkAddr) UnmarshalJSON(b []byte) error {
	msg := &proto.ThunkAddr{}
	err := protojson.Unmarshal(b, msg)
	if err != nil {
		return err
	}

	return value.UnmarshalProto(msg)
}
