import Gachi, { useState, useNavigate } from "../../../core/framework"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import CreatePost from "../create-posts/form-input"
import Posts from "../posts/posts"
export default function MainPage() {
	return (
		<div>
			<Header />
			<NavBar />
			<CreatePost />
			<Posts endPointUrl={"posts"} userId={""} />
		</div>
	)
}
