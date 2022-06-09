package bass

import (
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/vito/bass/pkg/proto"
)

var Secrets = NewEmptyScope()

func init() {
	Ground.Set("mask",
		Func("mask", "[secret name]", func(val String, name Symbol) Secret {
			return NewSecret(name.String(), []byte(val))
		}),
		`shrouds a string in secrecy`,
		`Prevents the string from being revealed when the value is displayed.`,
		`Prevents the string from being revealed in a serialized thunk or thunk path.`,
		`Does NOT currently prevent the string's value from being displayed in log output; you still have to be careful there.`,
		`=> (mask "super secret" :github-token)`)
}

type Secret struct {
	Name string `json:"secret"`

	// private to guard against accidentally revealing it when encoding to JSON
	// or something
	secret []byte
}

func NewSecret(name string, inner []byte) Secret {
	return Secret{
		Name:   name,
		secret: inner,
	}
}

func (secret Secret) Reveal() []byte {
	return secret.secret
}

var _ Value = Secret{}

func (secret Secret) String() string {
	return fmt.Sprintf("<secret: %s (%d bytes)>", secret.Name, len(secret.secret))
}

// Eval does nothing and returns the secret.
func (secret Secret) Eval(ctx context.Context, scope *Scope, cont Cont) ReadyCont {
	return cont.Call(secret, nil)
}

// Equal returns false; secrets cannot be compared.
func (secret Secret) Equal(other Value) bool {
	var o Secret
	return other.Decode(&o) == nil &&
		subtle.ConstantTimeCompare(secret.secret, o.secret) == 1
}

// Decode only supports decoding into a Secret or Value; it will not reveal the
// inner secret.
func (value Secret) Decode(dest any) error {
	switch x := dest.(type) {
	case *Secret:
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

var _ Decodable = (*Secret)(nil)

func (value *Secret) UnmarshalProto(msg proto.Message) error {
	p, ok := msg.(*proto.Secret)
	if !ok {
		return fmt.Errorf("unmarshal proto: %w", DecodeError{msg, value})
	}

	value.Name = p.Name
	value.secret = p.Value

	return nil
}

func (value *Secret) FromValue(val Value) error {
	var obj *Scope
	if err := val.Decode(&obj); err != nil {
		return fmt.Errorf("%T.FromValue: %w", value, err)
	}

	return decodeStruct(obj, value)
}
