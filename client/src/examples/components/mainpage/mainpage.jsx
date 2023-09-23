import Gachi, { useState, useNavigate } from "../../../core/framework"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import { PostsAuth } from "../create-posts/postAuth"

export default function MainPage() {
	const [activeSubj, setActiveSubj] = useState("")
	Gachi.createContext("currentCategory", { activeSubj, setActiveSubj })

	const navigate = useNavigate()

	return (
		<div>
			<Header />
			<NavBar />
			<PostsAuth />
		</div>
	)
}
