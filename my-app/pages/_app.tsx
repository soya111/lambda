import Header from "../components/Header";
import { UserProvider } from "@auth0/nextjs-auth0/client";
import type { AppProps } from "next/app";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <UserProvider>
      <Header />
      <Component {...pageProps} />
    </UserProvider>
  );
}

export default MyApp;
