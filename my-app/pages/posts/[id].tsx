import axios from 'axios';
import { GetServerSideProps } from 'next';

type Post = {
    userId: number;
    id: number;
    title: string;
    body: string;
};

type PostDetailProps = {
    post: Post;
};

const PostDetail: React.FC<PostDetailProps> = ({ post }) => {
    return (
        <div className="container">
            <h1>{post.title}</h1>
            <p>{post.body}</p>
        </div>
    );
};

export const getServerSideProps: GetServerSideProps = async (context) => {
    const id = context.params?.id;
    let post: Post | null = null;

    try {
        const response = await axios.get(`https://jsonplaceholder.typicode.com/posts/${id}`);
        post = response.data;
    } catch (error) {
        console.error("Error fetching post:", error);
    }

    return {
        props: { post }
    };
};

export default PostDetail;
