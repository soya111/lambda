import Link from "next/link";

const Header: React.FC = () => {
  return (
    <header>
      <div className="container">
        <Link href="/">Home</Link>
        <Link href="/api/auth/login">Login</Link>
        <Link href="/api/auth/logout">Logout</Link>
        <Link href="/profile">Profile</Link>
      </div>
      <style jsx>{`
        header {
          background-color: cornflowerblue;
          padding: 10px 0;
        }
        .container {
          max-width: 1200px;
          margin: 0 auto;
          padding: 0 15px;
          display: flex;
          justify-content: space-between;
        }
        a {
          color: white;
          text-decoration: none;
          margin-right: 20px;
        }
      `}</style>
    </header>
  );
};

export default Header;
