package migration

import (
	"context"

	init0 "github.com/filecoin-project/specs-actors/actors/builtin/init"
	adt0 "github.com/filecoin-project/specs-actors/actors/util/adt"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"

	init2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/init"
	adt2 "github.com/filecoin-project/specs-actors/v2/actors/util/adt"
)

type initMigrator struct {
}

func (*initMigrator) MigrateState(ctx context.Context, store cbor.IpldStore, head cid.Cid) (cid.Cid, error) {
	var inState init0.State
	if err := store.Get(ctx, head, &inState); err != nil {
		return cid.Undef, err
	}

	// Migrate address resolution map
	inAddrMap, err := adt0.AsMap(adt0.WrapStore(ctx, store), inState.AddressMap)
	if err != nil {
		return cid.Undef, err
	}

	outAddrMap := adt2.MakeEmptyMap(adt2.WrapStore(ctx, store))

	var actorId cbg.CborInt
	if err = inAddrMap.ForEach(&actorId, func(key string) error {
		return outAddrMap.Put(StringKey(key), &actorId)
	}); err != nil {
		return cid.Undef, err
	}

	addrMapRoot, err := outAddrMap.Root()
	if err != nil {
		return cid.Undef, err
	}

	outState := init2.State{
		AddressMap:  addrMapRoot,
		NextID:      inState.NextID,
		NetworkName: inState.NetworkName,
	}
	return store.Put(ctx, &outState)
}

