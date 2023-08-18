import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
} from "../core/framework.ts"
import { Router, Route } from "/src/components/router.ts"
import { importCss } from "../modules/cssLoader.js"
import Header from "./components/header/header.jsx"
import Login from "./components/login/login.jsx"
import Registration from "./components/registration/registration.jsx"
import Posts from "./components/posts/posts.jsx"
import ProfilePage from "./components/profile-page/profilePage.jsx"
import { PostsAuth } from "./components/create-posts/postAuth.jsx"
import { CommentAuth } from "./components/comments/commentsAuth.jsx"
import { Comment } from "./components/comments/comments.jsx"
// import MainPage from "./components/pages/mainpage.jsx"
// import DropdownMenu from "./components/header/dropdown.jsx"
importCss("./styles/index.css")

const container = document.getElementById("root")

function Home() {
	return (
		<div>
			<Header />
			<Posts />
		</div>
	)
}

function HomeAuth() {
	return (
		<div>
			<Header />
			<PostsAuth />
		</div>
	)
}

function HomeComment() {
	return (
		<div>
			<Header />
			<Comment />
		</div>
	)
}

function HomeCommentAuth() {
	return (
		<div>
			<Header />
			<CommentAuth />
		</div>
	)
}

function App() {
	return (
		<Router>
			<Route path="/" element={<Home />} />
			<Route path="/authorized" element={<HomeAuth />} />
			<Route path="/login" element={<Login />} />
			<Route path="/registration" element={<Registration />} />
			<Route path="/profile-page" element={<ProfilePage />} />
			<Route path="/internal-error" element={<h1>Error 500</h1>} />
			<Route path="/comments" element={<HomeComment />} />
			<Route path="/comments-authorized" element={<HomeCommentAuth />} />
		</Router>
	)
}

Gachi.render(<App />, container)
