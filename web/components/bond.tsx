import { useState } from "react";
import { GrClose } from "react-icons/gr";
import useSWR from "swr";
import { fetcher } from "utils/api";
import InfoPlaceholder from "components/infoPlaceholder";

function BondModal({ closeModal }: { closeModal: () => void }) {
  const [mode, setMode] = useState<"Bond" | "Redeem">("Bond");
  const { data, error } = useSWR(
    "/arbiter/bond/v1beta1/bond_info/ugreen",
    fetcher
  );

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
            {data ? (
              <div className="text-xl">
                {parseFloat(parseFloat(data.executing_price).toFixed(2))} GREEN
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

        <div className="w-full flex relative mb-4">
          <input
            type="number"
            className="w-full rounded-md mr-4 py-2 pl-4 pr-14 text-lg"
          />
          <button className="rounded-md p-2.5 absolute right-40 text-lg">
            MAX
          </button>
          <button className="w-button flex-shrink-0 border rounded-md py-2 text-lg">
            {mode === "Bond" ? "Bond GREEN" : "Claim ARB"}
          </button>
        </div>

        {mode === "Bond" ? (
          <div className="flex flex-col w-full pb-3">
            <div className="flex justify-between mb-1">
              <div>Your balance</div>
              <div>950 GREEN</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>You will get</div>
              <div>10 ARB</div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col w-full pb-3">
            <div className="flex justify-between mb-1">
              <div>Cliamable Rewards</div>
              <div>1 ARB</div>
            </div>
            <div className="flex justify-between mb-1">
              <div>Pending Rewards</div>
              <div>0 ARB</div>
            </div>
          </div>
        )}
        {mode === "Bond" ? (
          <div className="flex flex-col w-full pt-3 border-t">
            <div className="flex justify-between mb-1">
              <div>Max you can buy</div>
              <div>10 ARB</div>
            </div>
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

export default function Bond() {
  const [isOpenModal, setIsOpenModal] = useState(false);
  const openModal = () => setIsOpenModal(true);
  const closeModal = () => setIsOpenModal(false);

  const { data, error } = useSWR(
    "/arbiter/bond/v1beta1/bond_info/ugreen",
    fetcher
  );

  return (
    <div className="w-full h-full rounded-xl bg-secondary pt-8 pb-12 px-16 flex flex-col items-center">
      {isOpenModal && <BondModal closeModal={closeModal} />}
      <div className="text-3xl mb-8">Bond (1, 1)</div>

      <div className="w-full mb-12 flex justify-around">
        <div className="flex flex-col items-center">
          <div className="text-lg">Treasury Balance</div>
          <div className="text-xl">30000 GREEN ($300,000)</div>
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
          {data ? (
            <div className="text-xl">
              {parseFloat(
                parseFloat(`${parseFloat(data.executing_price) * 50}`).toFixed(
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
