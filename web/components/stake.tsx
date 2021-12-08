import { useState } from "react";
import useSWR, { useSWRConfig } from "swr";
import { fetcher } from "utils/api";
import { toPrettyCoin } from "utils/coin";
import InfoPlaceholder from "components/infoPlaceholder";
import { Keplr } from "@keplr-wallet/types";
import {
  BroadcastMode,
  makeSignDoc,
  makeStdTx,
  Msg,
  StdFee,
} from "@cosmjs/launchpad";
import { chainInfo } from "configs/chain";
import { Dec, DecUtils } from "@keplr-wallet/unit";

export default function Stake({
  keplr,
  bech32Address,
}: {
  keplr: Keplr | null;
  bech32Address: string;
}) {
  const [mode, setMode] = useState<"Stake" | "Unstake">("Stake");
  const [inputValue, setInputValue] = useState("");

  const { mutate } = useSWRConfig();
  const { data: accountData } = useSWR(
    `/auth/accounts/${bech32Address}`,
    fetcher
  );
  const { data: arbiterStakeData } = useSWR(
    `/arbiter/stake/v1beta1/stake_info/${bech32Address}`,
    fetcher
  );

  const stakeOrUnstake = async () => {
    if (!keplr) {
      throw Error("Keplr isn't connected");
    }

    if (!inputValue || inputValue === "0") {
      return;
    }

    const accountNumber =
      (accountData && accountData.result.value.account_number) || "0";
    const sequence = (accountData && accountData.result.value.sequence) || "0";

    const aminoMsgs: Msg[] = [
      {
        type:
          mode === "Stake" ? "arbiter/stake/join-stake" : "arbiter/stake/claim",
        value: {
          sender: bech32Address,
          token_in:
            mode === "Stake"
              ? {
                  amount: new Dec(inputValue)
                    .mul(
                      DecUtils.getPrecisionDec(
                        chainInfo.currencies.find(
                          (currency) => currency.coinMinimalDenom === "uarb"
                        )!.coinDecimals
                      )
                    )
                    .truncate()
                    .toString(),
                  denom: "uarb",
                }
              : new Dec(inputValue)
                  .mul(
                    DecUtils.getPrecisionDec(
                      chainInfo.currencies.find(
                        (currency) => currency.coinMinimalDenom === "uarb"
                      )!.coinDecimals
                    )
                  )
                  .truncate()
                  .toString(),
        },
      },
    ];
    const fee: StdFee = {
      gas: "200000",
      amount: [
        {
          amount: "0",
          denom: "uarb",
        },
      ],
    };
    const broadCastMode = "async";

    const signDoc = makeSignDoc(
      aminoMsgs,
      fee,
      chainInfo.chainId,
      "",
      accountNumber.toString(),
      sequence.toString()
    );
    const signResponse = await keplr.signAmino(
      chainInfo.chainId,
      bech32Address,
      signDoc,
      undefined
    );
    const signedTx = makeStdTx(signResponse.signed, signResponse.signature);

    await keplr.sendTx(
      chainInfo.chainId,
      signedTx,
      broadCastMode as BroadcastMode
    );

    setInputValue("");
    const timer = window.setTimeout(() => {
      mutate(`/auth/accounts/${bech32Address}`);
      mutate(`/arbiter/stake/v1beta1/stake_info/${bech32Address}`);
      clearTimeout(timer);
    }, 4000);
  };

  const setMax = () => {
    setInputValue(
      mode === "Stake"
        ? (arbiterStakeData &&
            arbiterStakeData.balance &&
            toPrettyCoin(
              (parseFloat(arbiterStakeData.balance) - 5000).toString(),
              "uarb"
            )
              .trim(true)
              .hideDenom(true)
              .locale(false)
              .toString()) ||
            "0"
        : (arbiterStakeData &&
            arbiterStakeData.staked &&
            toPrettyCoin(arbiterStakeData.staked, "uarb")
              .trim(true)
              .hideDenom(true)
              .toString()) ||
            "0"
    );
  };

  const unstakedBalance = arbiterStakeData
    ? arbiterStakeData.balance
      ? toPrettyCoin(arbiterStakeData.balance, "uarb")
          .trim(true)
          .hideDenom(true)
          .toString()
      : "0"
    : null;
  const stakedBalance = arbiterStakeData
    ? arbiterStakeData.staked
      ? toPrettyCoin(arbiterStakeData.staked, "uarb")
          .trim(true)
          .hideDenom(true)
          .toString()
      : "0"
    : null;
  const totalValueStaked = arbiterStakeData
    ? toPrettyCoin(arbiterStakeData.totalStaked || "0", "uarb")
        .trim(true)
        .hideDenom(true)
    : null;
  const apy = arbiterStakeData
    ? parseFloat(arbiterStakeData.rewardYield) === 0
      ? "0"
      : parseFloat(
          (
            Math.pow(1 + parseFloat(arbiterStakeData.rewardYield), 1095) * 100
          ).toFixed(2)
        )
    : null;

  return (
    <div className="w-full h-full rounded-xl bg-primary pt-8 pb-12 px-16 flex flex-col items-center">
      <div className="text-3xl mb-8 font-semibold">Stake (3, 3) 🧊</div>

      <div className="w-full mb-12 flex justify-around">
        <div className="flex flex-col items-center">
          <div className="text-lg">APY</div>
          {apy ? (
            <div className="text-xl">{apy}%</div>
          ) : (
            <InfoPlaceholder className="w-20 h-5 mt-1" />
          )}
        </div>
        <div className="flex flex-col items-center">
          <div className="text-lg">Total Value Staked</div>
          {totalValueStaked ? (
            <div className="text-xl">
              {totalValueStaked.toString()} ($
              {totalValueStaked.mul(new Dec(10)).toString()})
            </div>
          ) : (
            <InfoPlaceholder className="w-20 h-5 mt-1" />
          )}
        </div>
      </div>

      <div className="flex mb-2">
        <div
          className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
            mode === "Stake" ? "border-b-2" : ""
          }`}
          onClick={() => {
            setMode("Stake");
            setInputValue("");
          }}
        >
          Stake
        </div>
        <div
          className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
            mode === "Unstake" ? "border-b-2" : ""
          }`}
          onClick={() => {
            setMode("Unstake");
            setInputValue("");
          }}
        >
          Unstake
        </div>
      </div>

      <div className="w-full flex relative mb-8">
        <input
          type="number"
          className="w-full rounded-md mr-4 py-2 pl-4 pr-14 text-lg"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
        />
        <button
          className="rounded-md p-2.5 absolute right-40 text-lg"
          onClick={setMax}
        >
          MAX
        </button>
        <button
          className="w-button flex-shrink-0 border rounded-md py-2 text-lg"
          onClick={() => stakeOrUnstake()}
        >
          {mode} ARB
        </button>
      </div>

      <div className="flex flex-col w-full pt-6 border-t">
        <div className="flex justify-between mb-1">
          <div>Unstaked balance</div>
          {unstakedBalance ? (
            <div>{unstakedBalance} ARB</div>
          ) : (
            <InfoPlaceholder className="w-20 h-4" />
          )}
        </div>
        <div className="flex justify-between">
          <div>Staked balance</div>
          {stakedBalance ? (
            <div>
              {stakedBalance}
              &nbsp;sARB
            </div>
          ) : (
            <InfoPlaceholder className="w-20 h-4" />
          )}
        </div>
      </div>
    </div>
  );
}
