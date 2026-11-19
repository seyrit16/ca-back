package room

import (
	"ca/internal/model/enum"
	"crypto/rand"
	"math/big"
)

type SecretReward struct {
	pool []SecretResource

	Cache map[string]*SecretResource
}

type SecretResource struct {
	Resource enum.ItemTemplateId
	Weight   int
}

func buildStandardSecretReward() *SecretReward {
	return &SecretReward{
		pool: []SecretResource{
			{Resource: enum.ItemTemplateIdHeal, Weight: 15},
			{Resource: enum.ItemTemplateIdBuffDebuff, Weight: 20},
			{Resource: enum.ItemTemplateIdProvocation, Weight: 10},
			{Resource: enum.ItemTemplateIdEgoStrike, Weight: 10},
			{Resource: enum.ItemTemplateIdTeleport, Weight: 10},
			{Resource: enum.ItemTemplateIdCamp, Weight: 10},
			{Resource: enum.ItemTemplateIdReturnStone, Weight: 5},
			{Resource: enum.ItemTemplateIdBuff, Weight: 10},
			{Resource: enum.ItemTemplateIdDebuff, Weight: 10},
		},
		Cache: make(map[string]*SecretResource),
	}
}

func (sr *SecretReward) PickSecret() SecretResource {
	total := 0
	for _, e := range sr.pool {
		total += e.Weight
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(total)))
	roll := int(n.Int64())

	for _, e := range sr.pool {
		roll -= e.Weight
		if roll < 0 {
			return e
		}
	}
	return sr.pool[len(sr.pool)-1]
}
