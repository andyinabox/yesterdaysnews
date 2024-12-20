package markov

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// )

// type MultiChainSub struct {
// 	Chain  *BasicChain
// 	Weight float32
// }

// type MultiChain struct {
// 	chains       []MultiChainSub
// 	prefixLength int
// }

// func NewMultiChain(prefixLength int) *MultiChain {
// 	return &MultiChain{
// 		chains:       []MultiChainSub{},
// 		prefixLength: prefixLength,
// 	}
// }

// func (c *MultiChain) Build(chains []MultiChainSub) error {

// 	for _, sub := range chains {
// 		// detect inconsistencies
// 		if c.prefixLength != sub.Chain.PrefixLength() {
// 			return fmt.Errorf("incompatible prefix lengths: %d, %d", c.prefixLength, sub.Chain.prefixLength)
// 		}
// 	}

// 	c.chains = chains

// 	return nil
// }

// func (c *MultiChain) getOrderedChains() []Chain {
// 	chains := make([]Chain, len(c.chains))

// 	for i, sub := range c.chains {
// 		chains[i] = sub.Chain
// 	}

// 	// TODO: some psuedorandom sorting based on weight
// 	rand.Shuffle(len(chains), func(i, j int) {
// 		chains[i], chains[j] = chains[j], chains[i]
// 	})

// 	return chains
// }

// func (c *MultiChain) Start() Prefix {
// 	return c.getOrderedChains()[0].Start()
// }

// func (c *MultiChain) Next(p Prefix) (s string) {
// 	for _, sub := range c.getOrderedChains() {
// 		s = sub.Next(p)
// 		if s != "" {
// 			return
// 		}
// 	}

// 	return
// }

// func (c *MultiChain) End(p Prefix) (s string) {
// 	for _, sub := range c.getOrderedChains() {
// 		s = sub.End(p)
// 		if s != "" {
// 			return
// 		}
// 	}

// 	return
// }

// func (c *MultiChain) PrefixLength() int {
// 	return c.prefixLength
// }

// type multiChainData struct {
// 	PrefixLength int                 `json:"prefixLength"`
// 	Chains       []multiChainDataSub `json:"chains"`
// }

// type multiChainDataSub struct {
// 	Chain  json.RawMessage `json:"chain"`
// 	Weight float32         `json:"weight"`
// }

// func (c *MultiChain) Save() ([]byte, error) {
// 	data := multiChainData{
// 		PrefixLength: c.prefixLength,
// 		Chains:       make([]multiChainDataSub, len(c.chains)),
// 	}

// 	for i, sub := range c.chains {
// 		b, err := sub.Chain.Save()
// 		if err != nil {
// 			return nil, fmt.Errorf("error exporting sub-chain: %w", err)
// 		}
// 		data.Chains[i] = multiChainDataSub{
// 			Chain:  b,
// 			Weight: sub.Weight,
// 		}
// 	}

// 	return json.Marshal(data)
// }

// func (c *MultiChain) Load(b []byte) error {
// 	data := multiChainData{}
// 	err := json.Unmarshal(b, &data)
// 	if err != nil {
// 		return err
// 	}

// 	chains := make([]MultiChainSub, len(data.Chains))

// 	for i, sub := range data.Chains {
// 		chain := NewBasicChain(data.PrefixLength)
// 		err = chain.Load([]byte(sub.Chain))
// 		if err != nil {
// 			return err
// 		}

// 		chains[i] = MultiChainSub{
// 			Chain:  chain,
// 			Weight: sub.Weight,
// 		}
// 	}

// 	c = NewMultiChain(data.PrefixLength)

// 	return c.Build(chains)
// }
