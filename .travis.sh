# Simulate the environment used by Travis CI so that we can run local tests to
# find and resolve issues that are consistent with the Travis CI environment.
# This is helpful because Travis CI often finds issues that our own local tests
# do not.

#go vet                  `go list ./... | grep -Ev "(vendor)"` && \
#golint -set_exit_status `go list ./... | grep -Ev "(vendor)"` && \

# Test and generate cover profiles
ginkgo --cover core/wallet/status                  \
               core/wallet/swapper                 \
               adapter/binder/btc           \
               adapter/binder/erc20         \
               adapter/binder/eth           \
               adapter/callback             \
               adapter/db                   \
               adapter/fund                 \
               adapter/router               \
               adapter/server               \
               driver/keystore              \
               driver/leveldb               \
               driver/logger                \
               driver/swapperd            &&\

# Merge cover profiles into one root cover profile
covermerge core/wallet/status/status.coverprofile               \
           core/wallet/swapper/swapper.coverprofile             \
           adapter/binder/btc/btc.coverprofile           \
           adapter/binder/erc20/erc20.coverprofile       \
           adapter/binder/eth/eth.coverprofile           \
           adapter/callback/callbak.coverprofile         \
           adapter/db/db.coverprofile                    \
           adapter/fund/fund.coverprofile                \
           adapter/router/router.coverprofile            \
           adapter/server/server.coverprofile            \
           driver/keystore/keystore.coverprofile         \
           driver/leveldb/leveldb.coverprofile           \
           driver/logger/logger.coverprofile             \
           driver/swapperd/swapperd.coverprofile         \
           > swapperd.coverprofile

# Remove auto-generated protobuf files
sed -i '/.pb.go/d' renvm.coverprofile

# Remove marshaling files
sed -i '/marshal.go/d' renvm.coverprofile