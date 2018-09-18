package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/onsi/ginkgo"
)

// LocalContext allows you to mark a ginkgo context as being local-only.
// It won't run if the CI environment variable is true.
func LocalContext(description string, f func()) {
	var local bool

	ciEnv := os.Getenv("CI")
	ci, err := strconv.ParseBool(ciEnv)
	if err != nil {
		ci = false
	}

	// Assume tests are running locally if CI environment variable is not defined
	local = !ci

	if local {
		ginkgo.Context(description, f)
	} else {
		ginkgo.PContext(description, func() {
			ginkgo.It("SKIPPING LOCAL TESTS", func() {})
		})
	}
}

type TestKeystores struct {
	Alice TestKeystore `json:"alice"`
	Bob   TestKeystore `json:"bob"`
}

type TestKeystore struct {
	Ethereum string `json:"ethereum"`
	Bitcoin  string `json:"bitcoin"`
}

func LoadTestKeys(loc string) TestKeystores {
	keystores := TestKeystores{}
	data, err := ioutil.ReadFile(loc)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &keystores); err != nil {
		panic(err)
	}
	return keystores
}
