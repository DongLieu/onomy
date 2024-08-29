#!/bin/bash
set -xeu

# always returns true so set -e doesn't exit if it is not running.
killall onomyd || true
rm -rf $HOME/.onomyd/

mkdir $HOME/.onomyd
cd $HOME/.onomyd/
mkdir $HOME/.onomyd/validator1
mkdir $HOME/.onomyd/validator2
mkdir $HOME/.onomyd/validator3

# init all three validators
onomyd init --chain-id=testing-1 validator1 --home=$HOME/.onomyd/validator1
onomyd init --chain-id=testing-1 validator2 --home=$HOME/.onomyd/validator2
onomyd init --chain-id=testing-1 validator3 --home=$HOME/.onomyd/validator3

# create keys for all three validators
# onomyd keys add validator1 --keyring-backend=test --home=$HOME/.onomyd/validator1
# onomyd keys add validator2 --keyring-backend=test --home=$HOME/.onomyd/validator2
# onomyd keys add validator3 --keyring-backend=test --home=$HOME/.onomyd/validator3
echo $(cat /Users/donglieu/script/keys/mnemonic1)| onomyd keys add validator1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1
# cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg
echo $(cat /Users/donglieu/script/keys/mnemonic2)| onomyd keys add validator2 --recover --keyring-backend=test --home=$HOME/.onomyd/validator2
# cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc
echo $(cat /Users/donglieu/script/keys/mnemonic3)| onomyd keys add validator3 --recover --keyring-backend=test --home=$HOME/.onomyd/validator3

# create validator node with tokens to transfer to the three other nodes
onomyd add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator1 
onomyd add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator1 
onomyd add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator1 
onomyd add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator2 
onomyd add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator2 
onomyd add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator2 
onomyd add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator3 
onomyd add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator3 
onomyd add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000usdt --home=$HOME/.onomyd/validator3 
onomyd gentx validator1 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator1 --chain-id=testing-1
onomyd gentx validator2 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator2 --chain-id=testing-1
onomyd gentx validator3 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator3 --chain-id=testing-1

cp validator2/config/gentx/*.json $HOME/.onomyd/validator1/config/gentx/
cp validator3/config/gentx/*.json $HOME/.onomyd/validator1/config/gentx/
onomyd collect-gentxs --home=$HOME/.onomyd/validator1 

# cp validator1/config/genesis.json $HOME/.onomyd/validator2/config/genesis.json
# cp validator1/config/genesis.json $HOME/.onomyd/validator3/config/genesis.json


# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.onomyd/validator1/config/app.toml
VALIDATOR2_APP_TOML=$HOME/.onomyd/validator2/config/app.toml
VALIDATOR3_APP_TOML=$HOME/.onomyd/validator3/config/app.toml

# validator1
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9050|g' $VALIDATOR1_APP_TOML
sed -i -E 's|127.0.0.1:9090|127.0.0.1:9050|g' $VALIDATOR1_APP_TOML

# validator2
sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1316|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|127.0.0.1:9090|127.0.0.1:9088|g' $VALIDATOR2_APP_TOML
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $VALIDATOR2_APP_TOML

# validator3
sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1315|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|127.0.0.1:9090|127.0.0.1:9086|g' $VALIDATOR3_APP_TOML
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $VALIDATOR3_APP_TOML

# change config.toml values
VALIDATOR1_CONFIG=$HOME/.onomyd/validator1/config/config.toml
VALIDATOR2_CONFIG=$HOME/.onomyd/validator2/config/config.toml
VALIDATOR3_CONFIG=$HOME/.onomyd/validator3/config/config.toml


# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR1_CONFIG


# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $VALIDATOR2_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26630"|g' $VALIDATOR2_CONFIG

# validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $VALIDATOR3_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26620"|g' $VALIDATOR3_CONFIG

# copy validator1 genesis file to validator2-3
# update
update_test_genesis () {
    # EX: update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/.onomyd/validator1/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.onomyd/validator1/config/genesis.json
}

update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"] = "15s"'

cp $HOME/.onomyd/validator1/config/genesis.json $HOME/.onomyd/validator2/config/genesis.json
cp $HOME/.onomyd/validator1/config/genesis.json $HOME/.onomyd/validator3/config/genesis.json

# copy tendermint node id of validator1 to persistent peers of validator2-3
node1=$(onomyd tendermint show-node-id --home=$HOME/.onomyd/validator1)
node2=$(onomyd tendermint show-node-id --home=$HOME/.onomyd/validator2)
node3=$(onomyd tendermint show-node-id --home=$HOME/.onomyd/validator3)
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26656,$node3@localhost:26656\"|g" $HOME/.onomyd/validator1/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26656,$node3@localhost:26656\"|g" $HOME/.onomyd/validator2/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26656,$node3@localhost:26656\"|g" $HOME/.onomyd/validator3/config/config.toml


# # start all three validators/
# onomyd start --home=$HOME/.onomyd/validator1
screen -S onomy1 -t onomy1 -d -m onomyd start --home=$HOME/.onomyd/validator1
screen -S onomy2 -t onomy2 -d -m onomyd start --home=$HOME/.onomyd/validator2
screen -S onomy3 -t onomy3 -d -m onomyd start --home=$HOME/.onomyd/validator3
# onomyd start --home=$HOME/.onomyd/validator3

sleep 7

onomyd tx bank send $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 100000stake --keyring-backend=test --chain-id=testing-1 -y --home=$HOME/.onomyd/validator1 --fees 100000000000000osmo

onomyd q staking validators
onomyd keys list --keyring-backend=test --home=$HOME/.onomyd/validator1
onomyd keys list --keyring-backend=test --home=$HOME/.onomyd/validator2
onomyd keys list --keyring-backend=test --home=$HOME/.onomyd/validator3

sleep 7
# killall onomyd || true
# onomyd testnet testing-1  onomyvaloper1wa3u4knw74r598quvzydvca42qsmk6jrya79zd --accounts-to-fund="onomy1wa3u4knw74r598quvzydvca42qsmk6jrc6uj7m,onomy1w7f3xx7e75p4l7qdym5msqem9rd4dyc4y47xsd,onomy1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn2ygn94a" --home=$HOME/.onomyd/validator1 --skip-confirmation

onomyd tx gov submit-proposal add-stable-coin "d" "d" "usdt" "100000000000000000000000000000" "1" "0.001" "0.001" onomy1wa3u4knw74r598quvzydvca42qsmk6jrc6uj7m  10000000000000000000stake --keyring-backend=test  --home=$HOME/.onomyd/validator1 --from onomy1wa3u4knw74r598quvzydvca42qsmk6jrc6uj7m -y --chain-id testing-1

sleep 7
onomyd tx gov vote 1 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y 
onomyd tx gov vote 1 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y 
onomyd tx gov vote 1 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y 

sleep 15
echo "==================="

onomyd q psm stablecoin usdt

onomyd q bank balances onomy1wa3u4knw74r598quvzydvca42qsmk6jrc6uj7m
echo "==================="
onomyd tx psm swap-to-ist 100000000000000000000000000000usdt --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y

sleep 7

onomyd tx psm swap-to-stablecoin usdt 1000IST --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y

# modulenema: 3CF27F7408755A5E12B08812B22AD88CB17322B5
# onomy18ne87aqgw4d9uy4s3qfty2kc3jchxg44ec3vm5