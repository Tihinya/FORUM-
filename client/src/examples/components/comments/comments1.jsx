import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"
import { NavBar } from "../navbar/navbar.jsx"
import Posts from "../posts/posts1"
import Header from "../header/header"
import CreateComment from "../comments/create-comment"

export function Comments({ postId: navigatePostId }) {
	return (
		<div>
			<Header />
			{/* <NavBar /> */}
			<Posts endPointUrl={"post"} userId={navigatePostId} />
			<CreateComment endPointUrl={"comment"} userId={navigatePostId} />
			<Posts endPointUrl={"comments"} userId={navigatePostId} />
		</div>
	)
}
