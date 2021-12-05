import { getKeplrFromWindow } from "@keplr-wallet/stores";
import { Keplr, ChainInfo } from "@keplr-wallet/types";
import { Bech32Address } from "@keplr-wallet/cosmos";
import { useEffect, useState } from "react";

const chainInfo: ChainInfo = {
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
const KeyAccountAutoConnect = "account_auto_connect";

export default function Header() {
  const [keplr, setKeplr] = useState<Keplr | null>(null);
  const [bech32Address, setBech32Address] = useState<string | null>(null);

  const connectWallet = async () => {
    try {
      const newKeplr = await getKeplrFromWindow();
      if (!newKeplr) {
        throw new Error("Keplr extension not found");
      }

      await newKeplr.experimentalSuggestChain(chainInfo);
      await newKeplr.enable(chainInfo.chainId);

      localStorage?.setItem(KeyAccountAutoConnect, "true");
      setKeplr(newKeplr);
    } catch (e) {
      alert(e);
    }
  };

  const signOut = () => {
    localStorage?.removeItem(KeyAccountAutoConnect);
    setKeplr(null);
    setBech32Address(null);
  };

  useEffect(() => {
    const shouldAutoConnectAccount =
      localStorage?.getItem(KeyAccountAutoConnect) != null;

    const geyKeySetAddress = async () => {
      if (keplr) {
        const key = await keplr.getKey(chainInfo.chainId);
        setBech32Address(key.bech32Address);
      }
    };

    if (shouldAutoConnectAccount) {
      connectWallet();
    }
    geyKeySetAddress();
  }, [keplr, bech32Address]);

  return (
    <header className="fixed w-full h-full max-h-header bg-white bg-opacity-95">
      <div className="max-w-default max-h-header w-full h-full m-auto flex justify-between items-center">
        <div className="mr-8 text-xl">Arbiter DAO</div>
        <div>
          <span className="mr-3">
            {bech32Address && Bech32Address.shortenAddress(bech32Address, 18)}
          </span>
          <button
            className={`rounded-lg py-2 px-4 ${
              bech32Address
                ? "bg-white text-black border"
                : "bg-black text-white"
            }`}
            onClick={() => (bech32Address ? signOut() : connectWallet())}
          >
            {bech32Address ? "Sign Out" : "Connect Wallet"}
          </button>
        </div>
      </div>
    </header>
  );
}
