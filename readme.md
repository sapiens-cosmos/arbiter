# Arbiter DAO

**Arbiter DAO** is a blockchain built using Cosmos SDK and Tendermint.

Arbiter DAO is a Olympus DAO and Klima DAO fork brought to the cosmos ecosystem combined with the carbon market along with [Regen](https://github.com/regen-network/regen-ledger) to secure carbon assets.

With Arbiter DAO enforcing price rise in the carbon credit market by sucking in eco-credits and carbon credits into the protocol's treasury, Arbiter DAO aims to start a movement for Earth within the cosmos ecosystem. Each ARB token would be backed by any verified eco-credit token with the equivalancy of  1 tonne of carbon reduction.

## What is Arbiter?

Arbiter DAO is named after a spaceship from the game **starcraft**. An arbiter in starcraft would use a skill called "stasis field" to freeze any units in the field. Arbiter DAO serves the same purpose of freezing carbon credits within the protocol to inflate the price of carbon credits and contribute to building a healthier earth.


## Launching Arbiter DAO in Local Environment

To support easier testing in local environments, we have created a docker file in `web/localnet/Dockerfile`. 

Docker installation is necessary to run local environment that has been pre-set for testing puposes. Dockerfile runs in amd64 cpu, also supporting arm64 (tested with M1 Mac).

The dockerfile initializes node with params and genesis suitable for testing purposes (ex. change default epoch from 72,000 blocks to 40 blocks) then creates and provides three different accounts to test with. The project is also setted with a unique bech32 prefix, allowing to test the accounts with the mnemonics provided with Keplr.

To start the node in your local environment along with integration of front end, run the following commands.
```
cd web
yarn install
yarn localnet
yarn dev
```


## Structure of Arbiter DAO

Arbiter DAO has to main modules. The `bond` module handling the bonding moduels and the `stake` module which manages staking and works as a treasury for the protocol.

### Bond
The debt ratio is calculated using the equation `bonds_outstanding / base_supply`. 

Premium is calculated by the equation  `premium = 1 + (debt_ratio * control_variable)` where control_varaible is controlled with params and governance.

Executing price of ARB token is calculated from the following equation: `executing_price = risk_free_price * premium (premium >= 1) `.


### Stake
Stake module param contains reward_ratio, which can be adjusted with governance in accordance to total supply and runway of ARB.

APY and rewards for the stakers goes through the process of every 8 hours, where APY is calculated with the equation `APY = (1 + reward_yield)^1095`.


Reward yield is detirmined using the following equation: `rewardYield = ARB_distributed / ARB_total_staked`.

The number of ARB distributed is calculated from ARB total supply with the following equation: `ARB_distributed = ARB_total_supply * reward_rate`

## Significance of Merging Olympus DAO with Cosmos

One of the shortcomings of Olympus DAO or Klima DAO was that the tokens treasury could manage would be strictly restricted to chain. With Web3.0 flourishing, more and more demand to put carbon credits on chain would increase as time goes on.


With the power of IBC,  Arbiter DAO would not have to be restricted by chain to manage carbon assets from all different chains. Any carbon assets from any chain within the cosmos ecosystem could come together as one in Arbiter DAO, along with the support of [Regen](https://github.com/regen-network/regen-ledger)'s eco credit module to bring more and more eco credit into the cosmos ecosystem.