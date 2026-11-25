package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	compexp "orglang/go-engine/internal/proc/compexec/exp"
)

func TestExpRecRoundTripPreservesConnectionIDs(t *testing.T) {
	connID := identity.New()
	contID := identity.New()
	valueID := identity.New()
	connChnl := symbol.ADT("conn")
	contChnl := symbol.ADT("cont")
	newChnl := symbol.ADT("new")
	label := uniqsym.New(symbol.ADT("label"))

	tests := []compexp.ExpRec{
		compexp.SendRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValChnlID:  valueID,
			ValExpCK:   contkey.ADT(7),
			ConnID:     connID,
		},
		compexp.RecvRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			NewChnlPH:  newChnl,
			ContExp:    compexp.CloseSpec{ContChnlPH: contChnl},
		},
		compexp.LabRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ValExpCK:   contkey.ADT(9),
			ValLabQN:   label,
		},
		compexp.CaseRec{
			ConnChnlPH: connChnl,
			ContChnlID: contID,
			ContExps: map[uniqsym.ADT]compexp.ExpSpec{
				label: compexp.CloseSpec{ContChnlPH: contChnl},
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
