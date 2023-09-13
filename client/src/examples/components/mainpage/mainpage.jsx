import Gachi, { useState, useNavigate } from "../../../core/framework"

import Header from "../header/header"
import { NavBar } from "../navbar/navbar"
import { PostsAuth } from "../create-posts/postAuth"
import PostList from "../post-list/postList"

export default function MainPage() {
	return (
		<div>
			<Header />
			<NavBar />
			<PostList />

			{/* <PostsAuth /> */}
		</div>
	)
}
