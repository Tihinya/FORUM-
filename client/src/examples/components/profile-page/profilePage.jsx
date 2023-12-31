import Gachi, { useNavigate, useContext } from "../../../core/framework"

import { NavBar } from "../navbar/navbar"
import Posts from "../posts/posts"
import Header from "../header/header"
import PersonalNavBar from "../personalNavBar/personalNavBar"

export default function ProfilePage() {
	// const defaultposts = createContext("user/posts")
	const { top } = useContext("currentTop")

	return (
		<div>
			<Header />
			<NavBar />
			<PersonalNavBar />
			<Posts endPointUrl={top} userId={""} />
		</div>
	)
}
