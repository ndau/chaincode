package main

import (
	"sort"

	"github.com/oneiro-ndev/chaincode/pkg/chain"
	"github.com/oneiro-ndev/ndau/pkg/ndau"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
)

type index struct {
	Name  string
	Value byte
}

func getNdauIndexMap() map[string]byte {
	indices := make(map[string]byte)
	objects := []interface{}{
		backing.AccountData{},
		backing.Lock{},
		ndau.Transfer{},
		ndau.ChangeValidation{},
		ndau.ReleaseFromEndowment{},
		ndau.ChangeSettlementPeriod{},
		ndau.Delegate{},
		ndau.CreditEAI{},
		ndau.Lock{},
		ndau.Notify{},
		ndau.SetRewardsDestination{},
		ndau.ClaimAccount{},
		ndau.TransferAndLock{},
	}

	for _, o := range objects {
		ks, _ := chain.ExtractConstants(o)
		for k, v := range ks {
			indices[k] = v
		}
	}
	return indices
}

func getNdauIndices() []index {
	indices := getNdauIndexMap()
	out := make([]index, 0)
	for k, v := range indices {
		out = append(out, index{Name: k, Value: v})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Value < out[j].Value
	})
	return out
}
