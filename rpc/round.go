package rpc

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/mixin/kernel"
	"github.com/MixinNetwork/mixin/storage"
)

func getRoundByNumber(store storage.Store, params []interface{}) (map[string]interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("invalid params count")
	}
	node, err := crypto.HashFromString(fmt.Sprint(params[0]))
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseUint(fmt.Sprint(params[1]), 10, 64)
	if err != nil {
		return nil, err
	}
	head, err := store.ReadRound(node)
	if err != nil {
		return nil, err
	}
	hash := node
	start := head.Timestamp
	end := head.Timestamp
	references := head.References

	snapshots, err := store.ReadSnapshotsForNodeRound(node, number)
	if err != nil {
		return nil, err
	}
	if head.Number == number {
	} else if len(snapshots) > 0 {
		rawSnapshots := make([]*common.Snapshot, len(snapshots))
		for i, s := range snapshots {
			rawSnapshots[i] = &s.Snapshot
		}
		start, end, hash = kernel.ComputeRoundHash(node, number, rawSnapshots)
		round, err := store.ReadRound(hash)
		if err != nil {
			return nil, err
		}
		if round == nil {
			return nil, fmt.Errorf("round malformed %s:%d:%d %s", node, number, start, hash)
		}
		if round.NodeId != node || round.Number != number || round.Timestamp != start {
			return nil, fmt.Errorf("round malformed %s:%d:%d %s:%d:%d", node, number, start, round.NodeId, round.Number, round.Timestamp)
		}
		references = round.References
	} else {
		return nil, fmt.Errorf("round not found")
	}
	return map[string]interface{}{
		"node":       node,
		"hash":       hash,
		"start":      start,
		"end":        end,
		"number":     number,
		"references": references,
		"snapshots":  snapshotsToMap(snapshots, nil, false),
	}, nil
}

func getRoundByHash(store storage.Store, params []interface{}) (map[string]interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("invalid params count")
	}
	hash, err := crypto.HashFromString(fmt.Sprint(params[0]))
	if err != nil {
		return nil, err
	}
	round, err := store.ReadRound(hash)
	if err != nil {
		return nil, err
	}
	if round == nil {
		return nil, errors.New("round not found")
	}
	start := round.Timestamp
	end := round.Timestamp

	snapshots, err := store.ReadSnapshotsForNodeRound(round.NodeId, round.Number)
	if err != nil {
		return nil, err
	}
	if round.NodeId == hash {
	} else if len(snapshots) > 0 {
		rawSnapshots := make([]*common.Snapshot, len(snapshots))
		for i, s := range snapshots {
			rawSnapshots[i] = &s.Snapshot
		}
		s, e, h := kernel.ComputeRoundHash(round.NodeId, round.Number, rawSnapshots)
		if h != hash {
			return nil, fmt.Errorf("round malformed %s:%d:%d:%s %s", round.NodeId, round.Number, round.Timestamp, hash, h)
		}
		start, end = s, e
	} else {
		return nil, fmt.Errorf("round malformed %s:%d", round.NodeId, round.Number)
	}
	return map[string]interface{}{
		"node":       round.NodeId,
		"hash":       hash,
		"start":      start,
		"end":        end,
		"number":     round.Number,
		"references": round.References,
		"snapshots":  snapshotsToMap(snapshots, nil, false),
	}, nil
}