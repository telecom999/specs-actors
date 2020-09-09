package migration

import (
	"context"

	builtin0 "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/specs-actors/v2/actors/builtin"
)

type StateMigration interface {
	// Loads an actor's state from an input store and writes new state to an output store.
	// Returns the new state head CID.
	MigrateState(ctx context.Context, store cbor.IpldStore, head cid.Cid) (cid.Cid, error)
}

type ActorMigration struct {
	InCodeCID      cid.Cid
	OutCodeCID     cid.Cid
	StateMigration StateMigration
}

var migrations = []ActorMigration{ // nolint:varcheck,deadcode,unused
	{
		InCodeCID:      builtin0.AccountActorCodeID,
		OutCodeCID:     builtin.AccountActorCodeID,
		StateMigration: &accountMigrator{},
	},
	{
		InCodeCID:      builtin0.CronActorCodeID,
		OutCodeCID:     builtin.CronActorCodeID,
		StateMigration: &cronMigrator{},
	},
	{
		InCodeCID:      builtin0.InitActorCodeID,
		OutCodeCID:     builtin.InitActorCodeID,
		StateMigration: &initMigrator{},
	},
	{
		InCodeCID:      builtin0.StorageMarketActorCodeID,
		OutCodeCID:     builtin.StorageMarketActorCodeID,
		StateMigration: &marketMigrator{},
	},
}
