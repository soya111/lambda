import React from "react";
import { InferGetServerSidePropsType } from "next";
import { withPageAuthRequired, getSession } from "@auth0/nextjs-auth0";
import { GetServerSideProps } from "next";

export default function Profile({
  user,
  data,
  error,
}: InferGetServerSidePropsType<typeof getServerSideProps>): React.ReactElement {
  console.log(data);
  if (error) console.error(error);
  return (
    <>
      <h1>Profile (server rendered)</h1>
      <pre data-testid="profile">{JSON.stringify(user, null, 2)}</pre>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = withPageAuthRequired({
  async getServerSideProps(ctx) {
    const session = await getSession(ctx.req, ctx.res);
    console.log(session);
    let data = null;
    let error = null;

    try {
      const response = await fetch("http://localhost:8080/api/private", {
        headers: {
          Authorization: `Bearer ${session?.accessToken}`,
        },
      });

      data = await response.json();
    } catch (e: any) {
      // エラーが発生した場合、エラーメッセージをキャッチします。
      console.error("Error during API call:", e.message);
      error = e.message;
    }

    return { props: { data, error } };
  },
});
