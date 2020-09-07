package migration

import (
	"context"

	market0 "github.com/filecoin-project/specs-actors/actors/builtin/market"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	market2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/market"
)

type marketMigrator struct {
}

func (*marketMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, head cid.Cid) (cid.Cid, error) {
	var inState market0.State
	if err := store.Get(ctx, head, &inState); err != nil {
		return cid.Undef, err
	}

	proposalsRoot, err := migrateProposals(ctx, store, inState.Proposals)
	if err != nil {
		return cid.Undef, err
	}

	statesRoot, err := migrateStates(ctx, store, inState.States)
	if err != nil {
		return cid.Undef, err
	}

	pendingRoot, err := migratePendingProposals(ctx, store, inState.PendingProposals)
	if err != nil {
		return cid.Undef, err
	}

	escrowRoot, err := migrateBalanceTable(ctx, store, inState.EscrowTable)
	if err != nil {
		return cid.Undef, err
	}

	lockedRoot, err := migrateBalanceTable(ctx, store, inState.LockedTable)
	if err != nil {
		return cid.Undef, err
	}

	dealOpsRoot, err := migrateDealOps(ctx, store, inState.DealOpsByEpoch)
	if err != nil {
		return cid.Undef, err
	}

	outState := market2.State{
		Proposals:                     proposalsRoot,
		States:                        statesRoot,
		PendingProposals:              pendingRoot,
		EscrowTable:                   escrowRoot,
		LockedTable:                   lockedRoot,
		NextID:                        inState.NextID,
		DealOpsByEpoch:                dealOpsRoot,
		LastCron:                      inState.LastCron,
		TotalClientLockedCollateral:   inState.TotalClientLockedCollateral,
		TotalProviderLockedCollateral: inState.TotalProviderLockedCollateral,
		TotalClientStorageFee:         inState.TotalClientStorageFee,
	}
	return store.Put(ctx, &outState)
}

func migrateProposals(ctx context.Context, store cbor.IpldStore, root cid.Cid) (cid.Cid, error) {
	// AMT and both the key and value type unchanged between v0 and v2.
	// Verify that the value type is identical.
	var _ = market0.DealProposal(market2.DealProposal{})

	// Just return the old root.
	return root, nil
}

func migrateStates(ctx context.Context, store cbor.IpldStore, root cid.Cid) (cid.Cid, error) {
	// AMT and both the key and value type unchanged between v0 and v2.
	// Verify that the value type is identical.
	var _ = market0.DealState(market2.DealState{})

	// Just return the old root.
	return root, nil
}

func migratePendingProposals(ctx context.Context, store cbor.IpldStore, root cid.Cid) (cid.Cid, error) {
	// The HAMT has changed, but the value type is identical.
	var _ = market0.DealProposal(market2.DealProposal{})

	return migrateHAMTRaw(ctx, store, root)

	//inMap, err := adt0.AsMap(adt0.WrapStore(ctx, store), root)
	//if err != nil {
	//	return cid.Undef, err
	//}
	//outMap := adt2.MakeEmptyMap(adt2.WrapStore(ctx, store))
	//
	//var inProposal market0.DealProposal
	//if err = inMap.ForEach(&inProposal, func(key int64) error {
	//	outProposal := market2.DealProposal(inProposal) // Identical
	//	return outMap.Set(uint64(key), &outProposal)
	//}); err != nil {
	//	return cid.Undef, err
	//}
	//
	//return outMap.Root()
}

func migrateBalanceTable(ctx context.Context, store cbor.IpldStore, root cid.Cid) (cid.Cid, error) {
	// The HAMT has changed, but the value type (abi.TokenAmount) is identical.
	return migrateHAMTRaw(ctx, store, root)
}

func migrateDealOps(ctx context.Context, store cbor.IpldStore, root cid.Cid) (cid.Cid, error) {
	// The HAMT has changed, at each level, but the final value type (abi.DealID) is identical.
	return migrateHAMTHAMTRaw(ctx, store, root)
}
