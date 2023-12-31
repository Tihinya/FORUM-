import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../core/framework.ts"
import { Router } from "/src/components/router.ts"
import { importCss } from "../modules/cssLoader.js"
import Header from "./components/header/header.jsx"
import Login from "./components/login/login.jsx"
import Registration from "./components/registration/registration.jsx"
import ProfilePage from "./components/profile-page/profilePage.jsx"
import { Comments } from "./components/comments/comments.jsx"
import MainPage from "./components/mainpage/mainpage.jsx"
import ErrorPage from "./components/errors/error-page.jsx"
import { RateLimiter } from "./additional-funcitons/ratelimiter.js"
import { ContextFetchOwnedPostsIds } from "./components/context-menu/helpers.jsx"
import { ContextFetchOwnedCommentsIds } from "./components/context-menu/helpers.jsx"
import { UserRole } from "./components/context-menu/helpers.jsx"
importCss("/styles/index.css")

const container = document.getElementById("root")

function HomeComment({ params }) {
	return (
		<div>
			<Comments postId={params.postId} />
		</div>
	)
}

const ErrorNotFound = {
	message: "Page Not Found",
	status: "404",
}

const ErrorBadRequest = {
	message: "Bad Request",
	status: "400",
}

const ErrorInternalError = {
	message: "Internal Server Error",
	status: "500",
}

export function App() {
	const [top, setTop] = useState("user/posts")
	Gachi.createContext("currentTop", { top, setTop })

	const [activeSubj, setActiveSubj] = useState("")
	Gachi.createContext("currentCategory", { activeSubj, setActiveSubj })

	const [posts, setPosts] = useState([])
	Gachi.createContext("currentPosts", { posts, setPosts })

	const [comments, setComments] = useState([])
	Gachi.createContext("currentComments", { comments, setComments })

	const [errorMessage, setErrorMessage] = useState("")
	Gachi.createContext("currentErrorMessage", {
		errorMessage,
		setErrorMessage,
	})

	const [ownedPostsIds, setOwnedPostsIds] = useState("")
	Gachi.createContext("currentOwnedPostsIds", {
		ownedPostsIds,
		setOwnedPostsIds,
	})

	const [userRole, setUserRole] = useState("")
	Gachi.createContext("currentUserRole", {
		userRole,
		setUserRole,
	})

	const [ownedCommentsIds, setOwnedCommentsIds] = useState("")
	Gachi.createContext("currentOwnedCommentsIds", {
		ownedCommentsIds,
		setOwnedCommentsIds,
	})

	const [isAuthenticated, setIsAuthenticated] = useState(false)
	Gachi.createContext("isAuthenticated", {
		isAuthenticated,
		setIsAuthenticated,
	})

	// Check if the user is authenticated on page load
	useEffect(() => {
		fetch("https://localhost:8080/authorized", {
			credentials: "include",
		})
			.then((response) => {
				if (response.ok) {
					setIsAuthenticated(true)
				} else if (response.status === 401) {
					setIsAuthenticated(false)
				}
			})
			.catch(() => setIsAuthenticated(false))
	}, [])

	ContextFetchOwnedPostsIds()
	ContextFetchOwnedCommentsIds()
	UserRole()

	return (
		<Router
			routes={[
				{ path: "/", element: <MainPage /> },
				{ path: "/login", element: <Login /> },
				{ path: "/registration", element: <Registration /> },
				{ path: "/profile-page", element: <ProfilePage /> },
				{
					path: "/comments-authorized/:postId",
					element: <HomeComment />,
				},
				{
					path: "serverded",
					element: <ErrorPage error={ErrorInternalError} />,
				},
				{
					path: "bad",
					element: <ErrorPage error={ErrorBadRequest} />,
				},
				{
					path: "*",
					element: <ErrorPage error={ErrorNotFound} />,
				},
			]}
		/>
	)
}

Gachi.render(<RateLimiter />, container)
