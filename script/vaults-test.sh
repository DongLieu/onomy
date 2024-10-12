#!/bin/bash
set -xeu

# always returns true so set -e doesn't exit if it is not running.
killall onomyd || true
rm -rf $HOME/.onomyd/

mkdir $HOME/.onomyd
mkdir $HOME/.onomyd/validator1
mkdir $HOME/.onomyd/validator2
mkdir $HOME/.onomyd/validator3

# init all three validators
onomyd init --chain-id=testing-1 validator1 --home=$HOME/.onomyd/validator1
onomyd init --chain-id=testing-1 validator2 --home=$HOME/.onomyd/validator2
onomyd init --chain-id=testing-1 validator3 --home=$HOME/.onomyd/validator3

# create keys for all three validators
mnemonic1="ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy"
mnemonic2="soap step crash ceiling path virtual this armor accident pond share track spice woman vault discover share holiday inquiry oak shine scrub bulb arrive"
mnemonic3="travel jelly basic visa apart kidney piano lumber elevator fat unknown guard matter used high drastic umbrella humble crush stock banner enlist mule unique"

echo $mnemonic1 | onomyd keys add validator1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1
# cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg
echo $mnemonic2 | onomyd keys add validator2 --recover --keyring-backend=test --home=$HOME/.onomyd/validator2
# cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc
echo $mnemonic3| onomyd keys add validator3 --recover --keyring-backend=test --home=$HOME/.onomyd/validator3

# create validator node with tokens to transfer to the three other nodes
onomyd genesis add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator1 
onomyd genesis add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator1 
onomyd genesis add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator1 
onomyd genesis add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator2 
onomyd genesis add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator2 
onomyd genesis add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator2 
onomyd genesis add-genesis-account $(onomyd keys show validator1 -a --keyring-backend=test --home=$HOME/.onomyd/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator3 
onomyd genesis add-genesis-account $(onomyd keys show validator2 -a --keyring-backend=test --home=$HOME/.onomyd/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator3 
onomyd genesis add-genesis-account $(onomyd keys show validator3 -a --keyring-backend=test --home=$HOME/.onomyd/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000nomUSD,10000000000000000000000000000000atom --home=$HOME/.onomyd/validator3 
onomyd genesis gentx validator1 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator1 --chain-id=testing-1
onomyd genesis gentx validator2 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator2 --chain-id=testing-1
onomyd genesis gentx validator3 1000000000000000000000stake --keyring-backend=test --home=$HOME/.onomyd/validator3 --chain-id=testing-1

cp $HOME/.onomyd/validator2/config/gentx/*.json $HOME/.onomyd/validator1/config/gentx/
cp $HOME/.onomyd/validator3/config/gentx/*.json $HOME/.onomyd/validator1/config/gentx/
onomyd genesis collect-gentxs --home=$HOME/.onomyd/validator1 

# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.onomyd/validator1/config/app.toml
VALIDATOR2_APP_TOML=$HOME/.onomyd/validator2/config/app.toml
VALIDATOR3_APP_TOML=$HOME/.onomyd/validator3/config/app.toml

# validator1
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9050|g' $VALIDATOR1_APP_TOML
sed -i -E 's|127.0.0.1:9090|127.0.0.1:9050|g' $VALIDATOR1_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR1_APP_TOML

# validator2
sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1316|g' $VALIDATOR2_APP_TOML
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $VALIDATOR2_APP_TOML
sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR2_APP_TOML

# validator3
sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1315|g' $VALIDATOR3_APP_TOML
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $VALIDATOR3_APP_TOML
sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR3_APP_TOML

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
# copy, update validator1 genesis file to validator2-3
update_test_genesis () {
    cat $HOME/.onomyd/validator1/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.onomyd/validator1/config/genesis.json
}

update_test_genesis '.app_state["gov"]["params"]["voting_period"] = "15s"'
update_test_genesis '.app_state["gov"]["params"]["expedited_voting_period"] = "10s"'

# aution param
update_test_genesis '.app_state["aution"]["params"]["auction_periods"] = "10s"'
update_test_genesis '.app_state["aution"]["params"]["reduce_step"] = "5s"'

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
# set price
sleep 7
onomyd tx oracle set-price nomUSD 1 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price fet 1.3 --home=$HOME/.onomyd/validator2  --from validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y
onomyd tx oracle set-price atom 8.13 --home=$HOME/.onomyd/validator3  --from validator3 --keyring-backend test --fees 20stake --chain-id testing-1 -y
sleep 7
onomyd q oracle  get-price fet
onomyd tx gov submit-proposal ./script/proposal-1.json --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# # vote
sleep 7
onomyd tx gov vote 1 yes  --from validator1 --keyring-backend test --home ~/.onomyd/validator1 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator2 --keyring-backend test --home ~/.onomyd/validator2 --chain-id testing-1 -y --fees 20stake
onomyd tx gov vote 1 yes  --from validator3 --keyring-backend test --home ~/.onomyd/validator3 --chain-id testing-1 -y --fees 20stake

# wait voting_perio=15s
echo "========sleep=========="
sleep 7
test1="today auto lazy finger shoulder abstract oppose south sunny glass similar great feature rhythm raise evil owner orange auction absurd half mail ice glory"
echo $test1 | onomyd keys add test1 --recover --keyring-backend=test --home=$HOME/.onomyd/validator1
# onomyd keys add test1 --home=$HOME/.onomyd/validator1 --keyring-backend test

onomyd tx bank send $( onomyd keys show validator1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) $( onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a) 30000000nomUSD,10000stake --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 8
onomyd q gov proposals
onomyd tx vaults create-vault 10000000atom 20000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y
# onomyd tx vaults create-vault 10000000atom 20000000nomUSD --from validator2 --home=$HOME/.onomyd/validator2 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults mint 0 20000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7 

# onomyd tx vaults deposit 0 20000000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults withdraw 0 1000atom --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

# sleep 7

# onomyd tx vaults repay 0 40000000nomUSD --from validator1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 7

onomyd q bank balances $(onomyd keys show validator1 -a --keyring-backend test --home /Users/donglieu/.onomyd/validator1)

onomyd tx oracle set-price atom 2 --home=$HOME/.onomyd/validator1  --from validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

sleep 31
onomyd q bank balances $(onomyd keys show test1 --home=$HOME/.onomyd/validator1  --keyring-backend test -a)

onomyd tx auction bid 0 20000000nomUSD 2.2 --from test1 --home=$HOME/.onomyd/validator1 --keyring-backend test --fees 20stake --chain-id testing-1 -y

echo "wating long time, query auction ratecurrent = 1.1...liquidate"
# onomyd tx aution 