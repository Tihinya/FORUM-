import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../Gachi.js/src/core/framework.ts"
import { Router } from "../Gachi.js/src/components/router.ts"
import { importCss } from "../Gachi.js/src/modules/cssLoader.js"
import Login from "./components/login/login.jsx"
import Registration from "./components/registration/registration.jsx"
import ProfilePage from "./components/profile-page/profilePage.jsx"
import { Comments } from "./components/comments/comments.jsx"
import MainPage from "./components/mainpage/mainpage.jsx"
import ErrorPage from "./components/errors/error-page.jsx"
import { RateLimiter } from "./additional-funcitons/ratelimiter.js"
import { ContextFetchOwnedPostsIds } from "./components/context-menu/helpers.jsx"
import { ContextFetchOwnedCommentsIds } from "./components/context-menu/helpers.jsx"
import { UserRole, UserId } from "./components/helpers/helpers.jsx"
import { initializeWebSocket } from "./components/socket/socket.jsx"
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
	const soundUrl = "./sounds/hmm.mp3"

	const [props, setProps] = useState("user/posts")
	Gachi.createContext("currentProps", { props, setProps })

	const [activeSubj, setActiveSubj] = useState("")
	Gachi.createContext("currentCategory", { activeSubj, setActiveSubj })

	const [categories, setCategories] = useState([])
	Gachi.createContext("categories", { categories, setCategories })

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

	const [userId, setUserId] = useState(0)
	Gachi.createContext("currentUserId", {
		userId,
		setUserId,
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

	const [displayNavbar, setDisplayNavbar] = useState(true)
	Gachi.createContext("displayNavbar", {
		displayNavbar,
		setDisplayNavbar,
	})

	const [moderators, setModerators] = useState([])
	Gachi.createContext("currentModerators", {
		moderators,
		setModerators,
	})

	const [selectedModerator, setSelectedModerator] = useState([])
	Gachi.createContext("selectedModerator", {
		selectedModerator,
		setSelectedModerator,
	})

	const [orderedUserList, setOrderedUserList] = useState([])
	Gachi.createContext("orderedUserList", {
		orderedUserList,
		setOrderedUserList,
	})

	const [orderedUsersByLastMessage, setOrderedUsersByLastMessage] = useState(
		[]
	)
	Gachi.createContext("orderedUsersByLastMessage", {
		orderedUsersByLastMessage,
		setOrderedUsersByLastMessage,
	})

	const [user, setUser] = useState([])
	Gachi.createContext("currentUser", {
		user,
		setUser,
	})

	const [allUsers, setAllUsers] = useState([])
	Gachi.createContext("allUsers", {
		allUsers,
		setAllUsers,
	})

	const [onlineUserIdsList, setOnlineUserIdsList] = useState([])
	Gachi.createContext("onlineUserIdsList", {
		onlineUserIdsList,
		setOnlineUserIdsList,
	})

	const [notificationsInterval, setNotificationsInterval] = useState(0)
	Gachi.createContext("notificationsInterval", {
		notificationsInterval,
		setNotificationsInterval,
	})

	const [chatboxRecipient, setChatboxRecipient] = useState(0)
	Gachi.createContext("chatboxRecipient", {
		chatboxRecipient,
		setChatboxRecipient,
	})

	const [notificationSound, setNotificationSound] = useState(undefined)
	Gachi.createContext("notificationSound", {
		notificationSound,
		setNotificationSound,
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

	useEffect(() => {
		if (isAuthenticated) {
			initializeWebSocket()
			document.addEventListener("click", initalizeAudio)
		}
	}, [isAuthenticated])

	function initalizeAudio() {
		const audio = new Audio(soundUrl)
		setNotificationSound(audio)
		document.removeEventListener("click", initalizeAudio)
	}

	ContextFetchOwnedPostsIds()
	ContextFetchOwnedCommentsIds()
	UserRole()
	UserId()

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
