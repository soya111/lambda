import axios from "axios";
import { GetServerSideProps } from "next";
import Link from "next/link";

type Post = {
  userId: number;
  id: number;
  title: string;
  body: string;
};

type HomeProps = {
  posts: Post[];
};

const Home: React.FC<HomeProps> = ({ posts }) => {
  return (
    <div className="container">
      <h1>Posts from JSONPlaceholder:</h1>
      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            {/* 以下の部分を変更 */}
            <Link href={`/posts/${post.id}`}>
              <h2>{post.title}</h2>
            </Link>
            <p>{post.body}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export const getServerSideProps: GetServerSideProps = async () => {
  let posts: Post[] = [];

  try {
    const response = await axios.get(
      "https://jsonplaceholder.typicode.com/posts"
    );
    posts = response.data;
  } catch (error) {
    console.error("Error fetching data:", error);
  }

  return {
    props: { posts },
  };
};

export default Home;
