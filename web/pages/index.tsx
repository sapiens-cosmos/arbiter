import type { NextPage } from "next";
import Header from "components/header";
import Stake from "components/stake";
import Bond from "components/bond";
import Head from "next/head";
import { getKeplrFromWindow } from "@keplr-wallet/stores";
import { Keplr } from "@keplr-wallet/types";
import { useEffect, useState } from "react";
import { chainInfo } from "configs/chain";

const KeyAccountAutoConnect = "account_auto_connect";

const Home: NextPage = () => {
  const [keplr, setKeplr] = useState<Keplr | null>(null);
  const [bech32Address, setBech32Address] = useState<string>("");

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
    setBech32Address("");
  };

  useEffect(() => {
    const shouldAutoConnectAccount =
      localStorage?.getItem(KeyAccountAutoConnect) != null;

    const geyKeySetAccountInfo = async () => {
      if (keplr) {
        const key = await keplr.getKey(chainInfo.chainId);
        setBech32Address(key.bech32Address);
      }
    };

    if (shouldAutoConnectAccount) {
      connectWallet();
    }
    geyKeySetAccountInfo();
  }, [keplr]);

  return (
    <>
      <Head>
        <title>Arbiter DAO</title>
        <meta name="description" content="Arbiter DAO" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header
        bech32Address={bech32Address}
        connectWallet={connectWallet}
        signOut={signOut}
      />
      <main className="pt-24 pb-10 max-w-default w-full mx-auto">
        <div className="mb-12">
          <Bond keplr={keplr} bech32Address={bech32Address} />
        </div>

        <Stake keplr={keplr} bech32Address={bech32Address} />
      </main>
    </>
  );
};

export default Home;
