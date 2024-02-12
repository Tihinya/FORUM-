import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import { NavBar } from "../navbar/navbar.jsx"
import Posts from "../posts/posts"
import Header from "../header/header"
import CreateComment from "./create-comment"
import { UserListBar } from "../userlistbar/userlistbar.jsx"
import ChatBoxComponent from "../chat/chatBoxComponent.jsx"

export function Comments({ postId: navigatePostId }) {
	const { chatboxRecipient } = useContext("chatboxRecipient")

	return (
		<div>
			<Header />
			<UserListBar />
			{chatboxRecipient == 0 ? (
				<div>
					<div>
						<Posts endPointUrl={"post"} userId={navigatePostId} />
					</div>
					<CreateComment
						endPointUrl={"comment"}
						userId={navigatePostId}
					/>
					<Posts endPointUrl={"comments"} userId={navigatePostId} />
				</div> // Extra div because <></> doesn't work with chatbox
			) : (
				<ChatBoxComponent recipientId={chatboxRecipient} />
			)}
		</div>
	)
}
