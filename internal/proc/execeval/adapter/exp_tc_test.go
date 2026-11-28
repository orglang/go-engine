package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	commact "orglang/go-engine/internal/proc/commact/core"
)

func TestExpRecRoundTripPreservesCommunicationIDs(t *testing.T) {
	commID := identity.New()
	contID := identity.New()
	valueID := identity.New()
	commChnl := symbol.ADT("conn")
	contChnl := symbol.ADT("cont")
	newChnl := symbol.ADT("new")
	label := uniqsym.New(symbol.ADT("label"))

	tests := []commact.ExpRec{
		commact.SendRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			ValChnlID:  valueID,
			ValExpCK:   contkey.ADT(7),
			CommID:     commID,
		},
		commact.RecvRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			NewChnlPH:  newChnl,
			ContExp:    commact.CloseSpec{ContChnlPH: contChnl},
		},
		commact.LabRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			ValExpCK:   contkey.ADT(9),
			ValLabQN:   label,
		},
		commact.CaseRec{
			CommChnlPH: commChnl,
			ContChnlID: contID,
			ContExps: map[uniqsym.ADT]commact.ExpSpec{
				label: commact.CloseSpec{ContChnlPH: contChnl},
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
