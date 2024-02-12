import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import CreatePost from "../create-posts/form-input"
import Posts from "../posts/posts"
import ErrorWindow from "../errors/error-window"
import ChatBoxComponent from "../chat/chatBoxComponent"
import { UserListBar } from "../userlistbar/userlistbar"

export default function MainPage() {
	const { errorMessage, setErrorMessage } = useContext("currentErrorMessage")
	const { chatboxRecipient } = useContext("chatboxRecipient")

	return (
		<div>
			{errorMessage != "" ? (
				<ErrorWindow
					errorMessage={errorMessage}
					onClose={() => setErrorMessage("")}
				/>
			) : (
				""
			)}

			<Header />
			<NavBar />
			<UserListBar />
			{chatboxRecipient == 0 ? <CreatePost /> : null}
			{chatboxRecipient != 0 ? (
				<ChatBoxComponent recipientId={chatboxRecipient} />
			) : (
				// For some reason it breaks when CreatePost component
				// is included here with empty fragments
				<Posts endPointUrl={"posts"} userId={""} />
			)}
		</div>
	)
}
