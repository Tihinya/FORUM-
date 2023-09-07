import Gachi, { useNavigate } from "../../../core/framework"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import Posts from "../posts/posts"

export default function MainPage() {
	const navigate = useNavigate()

	return (
		<div>
			<Header />
			<NavBar />
			<Posts />
		</div>
	)
}
