import { useState } from "react";
import { GrClose } from "react-icons/gr";
import useSWR from "swr";
import { fetcher } from "utils/api";
import InfoPlaceholder from "components/infoPlaceholder";
import { toPrettyCoin } from "utils/coin";
import { Dec, DecUtils } from "@keplr-wallet/unit";
import { chainInfo } from "configs/chain";
import { Keplr } from "@keplr-wallet/types";
import {
  BroadcastMode,
  makeSignDoc,
  makeStdTx,
  Msg,
  StdFee,
} from "@cosmjs/launchpad";

function BondModal({
  closeModal,
  keplr,
  bech32Address,
}: {
  closeModal: () => void;
  keplr: Keplr | null;
  bech32Address: string;
}) {
  const [mode, setMode] = useState<"Bond" | "Redeem">("Bond");
  const [inputValue, setInputValue] = useState("");

  const { data: accountData } = useSWR(
    `/auth/accounts/${bech32Address}`,
    fetcher
  );
  const { data: arbiterBondData } = useSWR(
    "/arbiter/bond/v1beta1/bond_info/ugreen",
    fetcher
  );
  const { data: arbiterRedeemData } = useSWR(
    `/arbiter/bond/v1beta1/redeemable/${bech32Address}`,
    fetcher
  );
  const { data: cosmosBankData } = useSWR(
    `/cosmos/bank/v1beta1/balances/${bech32Address}/ugreen`,
    fetcher
  );

  const yourBalance = cosmosBankData
    ? toPrettyCoin(
        cosmosBankData.balance ? cosmosBankData.balance.amount : "0",
        "uarb"
      )
        .trim(true)
        .hideDenom(true)
    : null;
  const youWillGet = arbiterBondData
    ? toPrettyCoin(inputValue || "0", "ugreen")
        .mul(
          DecUtils.getPrecisionDec(
            chainInfo.currencies.find(
              (currency) => currency.coinMinimalDenom === "uarb"
            )!.coinDecimals
          )
        )
        .mul(new Dec((1 / arbiterBondData.executing_price).toString() || "0"))
        .locale(true)
        .trim(true)
        .hideDenom(true)
    : null;
  const maxYouCanBuy =
    arbiterBondData && yourBalance
      ? yourBalance.mul(
          new Dec((1 / arbiterBondData.executing_price).toString() || "0")
        )
      : null;
  const claimableRewards = arbiterRedeemData
    ? toPrettyCoin(
        arbiterRedeemData.coin ? arbiterRedeemData.coin.amount : "0",
        "uarb"
      )
        .trim(true)
        .hideDenom(true)
    : null;

  const bondOrRedeem = async () => {
    if (!keplr) {
      throw Error("Keplr isn't connected");
    }

    if (mode === "Bond" && (!inputValue || inputValue === "0")) {
      return;
    }

    const accountNumber =
      (accountData && accountData.result.value.account_number) || "0";
    const sequence = (accountData && accountData.result.value.sequence) || "0";

    const aminoMsgs: Msg[] = [
      {
        type: mode === "Bond" ? "arbiter/MsgBondIn" : "arbiter/MsgRedeem",
        value: {
          bonder: bech32Address,
          ...(mode === "Bond" && {
            coin: {
              amount: new Dec(inputValue)
                .mul(
                  DecUtils.getPrecisionDec(
                    chainInfo.currencies.find(
                      (currency) => currency.coinMinimalDenom === "ugreen"
                    )!.coinDecimals
                  )
                )
                .truncate()
                .toString(),
              denom: "ugreen",
            },
          }),
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
  };

  return (
    <div
      className="fixed top-0 z-10 w-full h-full bg-black bg-opacity-70 flex justify-center items-center"
      onClick={(e) => {
        if (e.target !== e.currentTarget) {
          return;
        }

        closeModal();
      }}
    >
      <div className="rounded-xl w-full max-w-modal pt-8 pb-12 px-16 bg-secondary relative">
        <div
          className="absolute top-6 right-6 cursor-pointer"
          onClick={closeModal}
        >
          <GrClose />
        </div>

        <div className="w-full mb-6 flex justify-around text-3xl">GREEN</div>

        <div className="w-full mb-4 flex justify-around">
          <div className="flex flex-col items-center">
            <div className="text-lg">Bond Price(1 ARB)</div>
            {arbiterBondData ? (
              <div className="text-xl">
                {parseFloat(
                  parseFloat(arbiterBondData.executing_price).toFixed(6)
                )}{" "}
                GREEN
              </div>
            ) : (
              <InfoPlaceholder className="w-24 h-5 mt-1" />
            )}
          </div>
          <div className="flex flex-col items-center">
            <div className="text-lg">Market Price(1 ARB)</div>
            <div className="text-xl">2 GREEN</div>
          </div>
        </div>

        <div className="flex w-full justify-center mb-2">
          <div
            className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
              mode === "Bond" ? "border-b-2" : ""
            }`}
            onClick={() => {
              setMode("Bond");
              setInputValue("");
            }}
          >
            Bond
          </div>
          <div
            className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
              mode === "Redeem" ? "border-b-2" : ""
            }`}
            onClick={() => {
              setMode("Redeem");
              setInputValue("");
            }}
          >
            Redeem
          </div>
        </div>

        <div className="w-full flex relative mb-4">
          {mode === "Bond" && (
            <input
              type="number"
              className="w-full rounded-md mr-4 py-2 pl-4 pr-14 text-lg"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
            />
          )}
          {mode === "Bond" && (
            <button
              className="rounded-md p-2.5 absolute right-40 text-lg"
              onClick={() => {
                setInputValue(
                  mode === "Bond"
                    ? (yourBalance && yourBalance.locale(false).toString()) ||
                        "0"
                    : claimableRewards!.locale(false).toString()
                );
              }}
            >
              MAX
            </button>
          )}
          <button
            className={`${
              mode === "Bond" ? "w-button" : "w-full"
            } flex-shrink-0 border rounded-md py-2 text-lg`}
            onClick={bondOrRedeem}
          >
            {mode === "Bond" ? "Bond GREEN" : "Claim ARB"}
          </button>
        </div>

        {mode === "Bond" ? (
          <div className="flex flex-col w-full pb-3">
            <div className="flex justify-between mb-1">
              <div>Your balance</div>
              {yourBalance ? (
                <div>{yourBalance.toString()} GREEN</div>
              ) : (
                <InfoPlaceholder className="w-24 h-5" />
              )}
            </div>
            <div className="flex justify-between mb-1">
              <div>You will get</div>
              {youWillGet ? (
                <div>{youWillGet.locale(true).toString()} ARB</div>
              ) : (
                <InfoPlaceholder className="w-24 h-5" />
              )}
            </div>
            <div className="flex justify-between mb-1">
              <div>Max you can buy</div>
              {maxYouCanBuy ? (
                <div>{maxYouCanBuy.toString()} ARB</div>
              ) : (
                <InfoPlaceholder className="w-24 h-5" />
              )}
            </div>
          </div>
        ) : (
          <div className="flex flex-col w-full pb-3">
            <div className="flex justify-between mb-1">
              <div>Cliamable Rewards</div>
              {claimableRewards ? (
                <div>{claimableRewards!.locale(true).toString()} ARB</div>
              ) : (
                <InfoPlaceholder className="w-24 h-5" />
              )}
            </div>
            <div className="flex justify-between mb-1">
              <div>Pending Rewards</div>
              <div>0 ARB</div>
            </div>
          </div>
        )}
        {mode === "Bond" ? (
          <div className="flex flex-col w-full pt-3 border-t">
            <div className="flex justify-between">
              <div>Vesting term</div>
              <div>5 days</div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col w-full pt-6 border-t">
            <div className="flex justify-between mb-1">
              <div>Time until fully vested</div>
              <div>0 days</div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function Bond({
  keplr,
  bech32Address,
}: {
  keplr: Keplr | null;
  bech32Address: string;
}) {
  const [isOpenModal, setIsOpenModal] = useState(false);
  const openModal = () => setIsOpenModal(true);
  const closeModal = () => setIsOpenModal(false);

  const { data: arbiterBondData } = useSWR("/arbiter/bond/v1beta1/bond_info/ugreen", fetcher);
  const { data: arbiterStakeData } = useSWR("/arbiter/stake/v1beta1/total_reserve", fetcher);

  const treasuryBalance = arbiterStakeData ? toPrettyCoin(arbiterStakeData.totalReserve || "0", "ugreen")
  .mul(
    DecUtils.getPrecisionDec(
      chainInfo.currencies.find(
        (currency) => currency.coinMinimalDenom === "ugreen"
      )!.coinDecimals
    )
  ).trim(true).hideDenom(true) : null;

  return (
    <div className="w-full h-full rounded-xl bg-secondary pt-8 pb-12 px-16 flex flex-col items-center">
      {isOpenModal && (
        <BondModal
          keplr={keplr}
          bech32Address={bech32Address}
          closeModal={closeModal}
        />
      )}
      <div className="text-3xl mb-8">Bond (1, 1)</div>

      <div className="w-full mb-12 flex justify-around">
        <div className="flex flex-col items-center">
          <div className="text-lg">Treasury Balance</div>
          {treasuryBalance ? <div className="text-xl">{treasuryBalance.toString()} GREEN (${treasuryBalance.mul(new Dec(10)).trim(true).toString()})</div> : <InfoPlaceholder className="w-24 h-5" />} 
        </div>
        <div className="flex flex-col items-center">
          <div className="text-xl">ARB Price</div>
          <div className="text-xl">2 GREEN ($20) </div>
        </div>
      </div>

      <div className="flex flex-col w-full">
        <div
          className="flex justify-between items-center mb-4 py-4 px-6 border rounded-lg cursor-pointer"
          onClick={openModal}
        >
          <div className="text-xl">GREEN (eco-credit)</div>
          {arbiterBondData ? (
            <div className="text-xl">
              {parseFloat(
                parseFloat(`${parseFloat(arbiterBondData.executing_price) * 50}`).toFixed(
                  2
                )
              )}
              % D/C
            </div>
          ) : (
            <InfoPlaceholder className="w-24 h-6" />
          )}
        </div>
      </div>
    </div>
  );
}
