import Gachi from "../core/framework.ts"
import { Router } from "/src/components/router.ts"
import { importCss } from "../modules/cssLoader.js"
import Header from "./components/header/header.jsx"
import Login from "./components/login/login.jsx"
import Registration from "./components/registration/registration.jsx"
import ProfilePage from "./components/profile-page/profilePage.jsx"
import { CommentAuth } from "./components/comments/commentsAuth.jsx"
import MainPage from "./components/mainpage/mainpage.jsx"
importCss("/styles/index.css")

const container = document.getElementById("root")

function Home() {
	return (
		<div>
			<MainPage />
		</div>
	)
}

function HomeComment({ params }) {
	return (
		<div>
			<Header />
			<CommentAuth postId={params.postId} />
		</div>
	)
}

// function HomeCommentAuth() {
// 	return (
// 		<div>
// 			<Header />
// 			<CommentAuth />
// 		</div>
// 	)
// }

function App() {
	return (
		<Router
			routes={[
				{ path: "/", element: <Home /> },
				{ path: "/login", element: <Login /> },
				{ path: "/registration", element: <Registration /> },
				{ path: "/profile-page", element: <ProfilePage /> },
				{ path: "/internal-error", element: <h1>Error 500</h1> },
				{
					path: "/comments-authorized/:postId",
					element: <HomeComment />,
				},
			]}
		/>

		/* <Route path="/" element={<Home />} />
			<Route path="/authorized" element={<HomeAuth />} />
			<Route path="/login" element={<Login />} />
			<Route path="/registration" element={<Registration />} />
			<Route path="/profile-page" element={<ProfilePage />} />
			<Route path="/internal-error" element={<h1>Error 500</h1>} />
			<Route path="/comments" element={<HomeComment />} />
			<Route path="/comments-authorized" element={<HomeCommentAuth />} /> */
	)
}

Gachi.render(<App />, container)
