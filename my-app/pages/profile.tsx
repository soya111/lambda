import React from "react";
import { InferGetServerSidePropsType } from "next";
import { withPageAuthRequired, getSession } from "@auth0/nextjs-auth0";

export default function Profile({
  user,
  session,
}: InferGetServerSidePropsType<typeof getServerSideProps>): React.ReactElement {
  return (
    <>
      <h1>Profile (server rendered)</h1>
      <pre data-testid="profile">{JSON.stringify(user, null, 2)}</pre>
    </>
  );
}

export const getServerSideProps = withPageAuthRequired({
  // async getServerSideProps(ctx) {
  //   const session = await getSession(ctx.req, ctx.res);
  //   const serializableSession = JSON.parse(JSON.stringify(session));
  //   console.log(serializableSession);
  //   return { props: { user: session?.user, session: serializableSession } };
  // },
});
