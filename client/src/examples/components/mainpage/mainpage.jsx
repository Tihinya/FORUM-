import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import CreatePost from "../create-posts/form-input"
import Posts from "../posts/posts"
import ErrorWindow from "../errors/error-window"

export default function MainPage() {
	const { errorMessage, setErrorMessage } = useContext("currentErrorMessage")
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
			<CreatePost />
			<Posts endPointUrl={"posts"} userId={""} />
		</div>
	)
}
