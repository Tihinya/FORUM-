import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

import Header from "../header/header"
import DropdownMenu from "../header/dropdown"
import { NavBar } from "../navbar/navbar"
import Posts from "../posts/posts"
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
