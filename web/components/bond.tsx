import { useState } from "react";

function BondModal({ closeModal }: { closeModal: () => void }) {
  const [mode, setMode] = useState<"Bond" | "Redeem">("Bond");

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
      <div className="rounded-xl w-full max-w-modal pt-8 pb-12 px-16 bg-secondary">
        <div className="w-full mb-4 flex justify-around">
          <div className="flex flex-col items-center">
            <div className="text-xl">Bond Price</div>
            <div className="text-2xl">$612.3</div>
          </div>
          <div className="flex flex-col items-center">
            <div className="text-xl">Market Price</div>
            <div className="text-2xl">$655.3</div>
          </div>
        </div>

        <div className="flex w-full justify-center mb-2">
          <div
            className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
              mode === "Bond" ? "border-b-2" : ""
            }`}
            onClick={() => setMode("Bond")}
          >
            Bond
          </div>
          <div
            className={`text-lg py-0.5 mx-4 my-2 cursor-pointer ${
              mode === "Redeem" ? "border-b-2" : ""
            }`}
            onClick={() => setMode("Redeem")}
          >
            Redeem
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
            {mode === "Bond" ? "Bond BCT" : "Claim ARB"}
          </button>
        </div>

        {mode === "Bond" ? (
          <div className="flex flex-col w-full pt-6 border-t">
            <div className="flex justify-between mb-1">
              <div>Your balance</div>
              <div>231 BCT</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>You will get</div>
              <div>542 ARB</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>Max you can buy</div>
              <div>1,524 ARB</div>
            </div>
            <div className="flex justify-between">
              <div>Vesting term</div>
              <div>5 days</div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col w-full pt-6 border-t">
            <div className="flex justify-between mb-1">
              <div>Pending Rewards</div>
              <div>0 ARB</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>Claimable Rewards</div>
              <div>22 ARB</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>Time until fully vested</div>
              <div>2 days</div>
            </div>
            <div className="flex justify-between">
              <div>Vesting term</div>
              <div>5 days</div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function Bond() {
  const [isOpenModal, setIsOpenModal] = useState(false);
  const openModal = () => setIsOpenModal(true);
  const closeModal = () => setIsOpenModal(false);

  return (
    <div className="w-full h-full rounded-xl bg-secondary pt-8 pb-12 px-16 flex flex-col items-center">
      {isOpenModal && <BondModal closeModal={closeModal} />}
      <div className="text-3xl mb-8">Bond (1, 1)</div>

      <div className="w-full mb-12 flex justify-around">
        <div className="flex flex-col items-center">
          <div className="text-xl">Treasury Balance</div>
          <div className="text-2xl">$54,223,242</div>
        </div>
        <div className="flex flex-col items-center">
          <div className="text-xl">ARB Price</div>
          <div className="text-2xl">$655.3</div>
        </div>
      </div>

      <div className="flex flex-col w-full">
        <div
          className="flex justify-between items-center mb-4 py-4 px-6 border rounded-lg cursor-pointer"
          onClick={openModal}
        >
          <div className="text-xl">BCT(Base Carbon Tons)</div>
          <div className="text-xl">5% D/C</div>
        </div>
      </div>
    </div>
  );
}
