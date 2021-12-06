import { CoinPretty, Dec } from "@keplr-wallet/unit";
import { chainInfo } from "configs/chain";

export const toPrettyCoin = (
  amount: string | Dec,
  denom: string = chainInfo.currencies[0].coinMinimalDenom
): CoinPretty => {
  return new CoinPretty(
    chainInfo.currencies.find(
      (currency) => currency.coinMinimalDenom === denom
    ) || chainInfo.currencies[0],
    new Dec(amount.toString())
  );
};
