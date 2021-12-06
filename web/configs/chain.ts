import { ChainInfo } from "@keplr-wallet/types";
import { Bech32Address } from "@keplr-wallet/cosmos";

export const chainInfo: ChainInfo = {
  rpc: "http://localhost:26657",
  rest: "http://localhost:1317",
  chainId: "arbiter-localnet-1",
  chainName: "Arbiter",
  stakeCurrency: {
    coinDenom: "ARB",
    coinMinimalDenom: "uarb",
    coinDecimals: 6,
  },
  bip44: {
    coinType: 118,
  },
  bech32Config: Bech32Address.defaultBech32Config("arbiter"),
  currencies: [
    {
      coinDenom: "ARB",
      coinMinimalDenom: "uarb",
      coinDecimals: 6,
    },
    {
      coinDenom: "GREEN",
      coinMinimalDenom: "ugreen",
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: "ARB",
      coinMinimalDenom: "uarb",
      coinDecimals: 6,
    },
  ],
  features: ["stargate", "ibc-transfer"],
};
