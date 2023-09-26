import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"
import { NavBar } from "../navbar/navbar.jsx"
import Posts from "../posts/posts"
import Header from "../header/header"
import CreateComment from "./create-comment"

export function Comments({ postId: navigatePostId }) {
	return (
		<div>
			<Header />
			<div>
				<Posts endPointUrl={"post"} userId={navigatePostId} />
			</div>

			<CreateComment endPointUrl={"comment"} userId={navigatePostId} />
			<Posts endPointUrl={"comments"} userId={navigatePostId} />
		</div>
	)
}
