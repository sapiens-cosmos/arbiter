import { Bech32Address } from "@keplr-wallet/cosmos";

export default function Header({
  bech32Address,
  connectWallet,
  signOut,
}: {
  bech32Address: string;
  connectWallet: () => void;
  signOut: () => void;
}) {
  return (
    <header className="fixed z-30 w-full h-full max-h-header bg-white bg-opacity-70">
      <div className="max-w-default max-h-header w-full h-full m-auto flex justify-between items-center bg-white bg-opacity-70">
        <div className="mr-8 text-3xl font-bold">🛸 Arbiter DAO</div>
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
