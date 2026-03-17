package tongue

import "testing"

func TestEncodeDecode(t *testing.T) {
	encoded := Encode("complete")
	if encoded == "" {
		t.Error("expected encoded output")
	}
	decoded := Decode(encoded)
	if decoded == "" {
		t.Error("expected decoded output")
	}
}

func TestDecodeKnownToken(t *testing.T) {
	decoded := Decode("CMP")
	if decoded != "complete " {
		t.Errorf("expected 'complete ', got %q", decoded)
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("Prishe", "Zeid", "complete", false)
	if msg.From != "Prishe" {
		t.Errorf("expected Prishe, got %s", msg.From)
	}
	if msg.Encoded {
		t.Error("expected not encoded")
	}
}
