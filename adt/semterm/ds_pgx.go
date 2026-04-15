package semterm

import (
	"fmt"

	"orglang/go-engine/adt/seqnum"
)

const (
	insertRef = `
		insert into impl_sems (
			impl_id, impl_rn, kind
		) values (
			@impl_id, @impl_rn, @kind
		)`

	insertBind = `
		insert into impl_binds (
			impl_qn, impl_id
		) values (
			@impl_qn, @impl_id
		)`

	selectRefByQN = `
		select
			is.impl_id,
			is.impl_rn
		from impl_sems is
		left join impl_binds ib
			on ib.impl_id = is.impl_id
		where ib.impl_qn = $1`
)

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
