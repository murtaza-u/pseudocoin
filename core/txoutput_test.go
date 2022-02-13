package core_test

import (
	"bytes"
	"testing"

	"github.com/murtaza-udaipurwala/core/core"
)

var out = core.TXOutput{
	Value: 100,
}

var outs = core.TXOutputs{
	Outputs: []core.TXOutput{out},
}

func TestLock(t *testing.T) {
	err := out.Lock("17LAQ3wStGiCCqA4RiSLMsA8V5LJpNR5Tn")
	if err != nil {
		t.Error(err)
	}
}

func TestIsLockedWith(t *testing.T) {
	if !out.IsLockedWith(out.PubkeyHash) {
		t.Errorf("Failed to verify ownership of output")
	}
}

func TestSerializeDeserializeOutputs(t *testing.T) {
	data, err := outs.Serialize()
	if err != nil {
		t.Error(err)
	}

	if len(data) == 0 {
		t.Error("Failed to serialize TXOutputs")
	}

	douts, err := core.DeserializeOutputs(data)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(outs.Outputs[0].PubkeyHash, douts.Outputs[0].PubkeyHash) != 0 ||
		outs.Outputs[0].Value != douts.Outputs[0].Value {
		t.Errorf("TXOutputs deserialization failed")
	}
}
