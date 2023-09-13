import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"
import CreatePost from "../create-posts/form-input"
import Posts from "../posts/posts1"

export default function PostList() {
	return (
		<div>
			<CreatePost />
			<Posts endPointUrl={"posts"} userId={""} />
		</div>
	)
}
