package rpc

import (
	"github.com/MixinNetwork/mixin/kernel"
	"github.com/MixinNetwork/mixin/storage"
)

func listAllNodes(store storage.Store, node *kernel.Node) ([]map[string]interface{}, error) {
	nodes := node.SortAllNodesByTimestampAndId()
	result := make([]map[string]interface{}, len(nodes))
	for i, n := range nodes {
		item := map[string]interface{}{
			"id":          n.IdForNetwork(node.NetworkId()),
			"signer":      n.Signer,
			"payee":       n.Payee,
			"transaction": n.Transaction,
			"timestamp":   n.Timestamp,
			"state":       n.State,
		}
		result[i] = item
	}
	return result, nil
}
