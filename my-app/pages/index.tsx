import { useUser } from "@auth0/nextjs-auth0/client";
import Link from "next/link";

export interface IndexProps {
  // Define your props here if any
}

const Index: React.FC<IndexProps> = () => {
  const { user, error, isLoading } = useUser();
  // Debugging state values
  console.log({ isLoading, error, user });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>{error.message}</div>;

  if (user) {
    return (
      <div>
        <h1>Profile</h1>
        <h2>{user.name}</h2>
        <p>{user.email}</p>
      </div>
    );
  }

  return <Link href="/api/auth/login">Login</Link>;
};

export default Index;
