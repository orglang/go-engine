package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	connact "orglang/go-engine/internal/proc/compexec"
)

func TestExpRecRoundTripPreservesConnectionIDs(t *testing.T) {
	connID := identity.New()
	contID := identity.New()
	valueID := identity.New()
	connChnl := symbol.ADT("conn")
	contChnl := symbol.ADT("cont")
	newChnl := symbol.ADT("new")
	label := uniqsym.New(symbol.ADT("label"))

	tests := []connact.ExpRec{
		connact.SendRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValChnlID:  valueID,
			ValExpCK:   contkey.ADT(7),
			ConnID:     connID,
		},
		connact.RecvRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			NewChnlPH:  newChnl,
			ContExp:    connact.CloseSpec{ContChnlPH: contChnl},
		},
		connact.LabRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValExpCK:   contkey.ADT(9),
			ValLabQN:   label,
		},
		connact.CaseRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ContExps: map[uniqsym.ADT]connact.ExpSpec{
				label: connact.CloseSpec{ContChnlPH: contChnl},
			},
		},
	}

	for _, want := range tests {
		dto, err := DataFromExpRec(want)
		if err != nil {
			t.Fatalf("DataFromExpRec(%T): %v", want, err)
		}
		got, err := DataToExpRec(dto)
		if err != nil {
			t.Fatalf("DataToExpRec(%T): %v", want, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("round-trip mismatch for %T:\n got: %#v\nwant: %#v", want, got, want)
		}
	}
}
