package owner

import (
	"encoding/json"
	"io/ioutil"
)

type Owner struct {
	Ganache string `json:"ganache"`
	Kovan   string `json:"kovan"`
	Ropsten string `json:"ropsten"`
	Mainnet string `json:"mainnet"`
}

func LoadOwner(path string) (Owner, error) {
	owner := Owner{}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return owner, err
	}
	err = json.Unmarshal(raw, &owner)
	if err != nil {
		return owner, err
	}
	return owner, nil
}
