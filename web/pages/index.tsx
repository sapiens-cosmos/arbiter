import type { NextPage } from "next";
import Header from "components/header";
import Stake from "components/stake";
import Bond from "components/bond";
import Head from "next/head";

const Home: NextPage = () => {
  return (
    <>
      <Head>
        <title>Arbiter DAO</title>
        <meta name="description" content="Arbiter DAO" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />
      <main className="pt-24 pb-10 max-w-default w-full mx-auto">
        <div className="mb-12">
          <Bond />
        </div>

        <Stake />
      </main>
    </>
  );
};

export default Home;
