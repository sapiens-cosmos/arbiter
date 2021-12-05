import { useState } from "react";

export default function Stake() {
  const [mode, setMode] = useState<"Stake" | "Unstake">("Stake");

  return (
    <div className="w-full h-full rounded-xl bg-primary pt-8 pb-12 px-16 flex flex-col items-center">
      <div className="text-3xl mb-8">Stake (3, 3)</div>

      <div className="w-full mb-12 flex justify-between">
        <div className="flex flex-col items-center">
          <div className="text-xl">APY</div>
          <div className="text-2xl">87,929.4%</div>
        </div>
        <div className="flex flex-col items-center">
          <div className="text-xl">Total Value Staked</div>
          <div className="text-2xl">$23,311,222</div>
        </div>
        <div className="flex flex-col items-center">
          <div className="text-xl">Current Index</div>
          <div className="text-2xl">23.2 ARB</div>
        </div>
      </div>

      <div className="flex mb-2">
        <div
          className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
            mode === "Stake" ? "border-b-2" : ""
          }`}
          onClick={() => setMode("Stake")}
        >
          Stake
        </div>
        <div
          className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
            mode === "Unstake" ? "border-b-2" : ""
          }`}
          onClick={() => setMode("Unstake")}
        >
          Unstake
        </div>
      </div>

      <div className="w-full flex relative mb-8">
        <input
          type="number"
          className="w-full rounded-md mr-4 py-2 pl-4 pr-14 text-lg"
        />
        <button className="rounded-md p-2.5 absolute right-40 text-lg">
          MAX
        </button>
        <button className="w-button flex-shrink-0 border rounded-md py-2 text-lg">
          {mode} ARB
        </button>
      </div>

      <div className="flex flex-col w-full pt-6 border-t">
        <div className="flex justify-between mb-1">
          <div>Unstaked balance</div>
          <div>231 ARB</div>
        </div>
        <div className="flex justify-between">
          <div>Staked balance</div>
          <div>9542 sARB</div>
        </div>
      </div>
    </div>
  );
}
