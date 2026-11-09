package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func TestExpRecRoundTripPreservesConnectionIDs(t *testing.T) {
	connID := identity.New()
	contID := identity.New()
	valueID := identity.New()
	connChnl := symbol.ADT("conn")
	contChnl := symbol.ADT("cont")
	newChnl := symbol.ADT("new")
	label := uniqsym.New(symbol.ADT("label"))

	tests := []termexp.ExpRec{
		termexp.SendRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValChnlID:  valueID,
			ValExpCK:   contkey.ADT(7),
			ConnID:     connID,
		},
		termexp.RecvRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			NewChnlPH:  newChnl,
			ContExp:    termexp.CloseSpec{ContChnlPH: contChnl},
		},
		termexp.LabRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValExpCK:   contkey.ADT(9),
			ValLabQN:   label,
		},
		termexp.CaseRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ContExps: map[uniqsym.ADT]termexp.ExpSpec{
				label: termexp.CloseSpec{ContChnlPH: contChnl},
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
