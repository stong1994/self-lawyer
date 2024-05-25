import Head from "next/head";
import { Inter } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { SearchDialog } from "@/components/SearchDialog";
import Image from "next/image";
import Link from "next/link";
import Script from "next/script";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  return (
    <>
      <Head>
        <title> 私人律师</title>
        <meta name="description" content="私人律师" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Script
        src="https://www.googletagmanager.com/gtag/js?id=GTM-NC7WFP48"
        strategy="afterInteractive"
      />
      <Script id="google-analytics" strategy="afterInteractive">
        {`
          window.dataLayer = window.dataLayer || [];
          function gtag(){window.dataLayer.push(arguments);}
          gtag('js', new Date());

          gtag('config', 'GTM-NC7WFP48');
        `}
      </Script>
      <main className={styles.main}>
        <h1 className="text-slate-700 font-bold text-2xl mb-12 flex items-center gap-3 dark:text-slate-400">
          <Image
            src={"/lawyer.png"}
            width="100"
            height="100"
            alt="MagickPen logo"
          />
          私人律师
        </h1>
        <div className={styles.center}>
          <SearchDialog />
        </div>

        <footer className="w-full flex justify-center items-center  mt-auto py-4">
          <div className="opacity-75 transition hover:opacity-100 cursor-pointer">
            <Link
              href="https://github.com/stong1994/self-lawyer"
              className="flex items-center justify-center"
            >
              <Image
                src={"/github.svg"}
                width="24"
                height="24"
                alt="Github logo"
              />
            </Link>
          </div>
        </footer>
      </main>
    </>
  );
}
