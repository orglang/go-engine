package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/adt/contkey"

	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func TestExpRecRoundTripPreservesConnectionIDs(t *testing.T) {
	connID := identity.New()
	contID := identity.New()
	valueID := identity.New()
	commChnl := symbol.ADT("comm")
	contChnl := symbol.ADT("cont")
	newChnl := symbol.ADT("new")
	label := uniqsym.New(symbol.ADT("label"))

	tests := []termexp.ExpRec{
		termexp.SendRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			ValChnlID:  valueID,
			ValExpCK:   contkey.ADT(7),
			ConnID:     connID,
		},
		termexp.RecvRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			NewChnlPH:  newChnl,
			ContExp:    termexp.CloseSpec{ContChnlPH: contChnl},
		},
		termexp.LabRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			ValExpCK:   contkey.ADT(9),
			ValLabQN:   label,
		},
		termexp.CaseRec{
			CommChnlPH: commChnl,
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
