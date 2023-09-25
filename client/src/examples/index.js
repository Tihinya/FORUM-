import Gachi, { useState, useEffect } from "../core/framework.ts"
import { Router } from "/src/components/router.ts"
import { importCss } from "../modules/cssLoader.js"
import Header from "./components/header/header.jsx"
import Login from "./components/login/login.jsx"
import Registration from "./components/registration/registration.jsx"
import ProfilePage from "./components/profile-page/profilePage.jsx"
import { CommentAuth } from "./components/comments/commentsAuth.jsx"
import MainPage from "./components/mainpage/mainpage.jsx"
import ErrorPage from "./components/errors/error-page.jsx"
import { RateLimiter } from "./additional-funcitons/ratelimiter.js"
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
	const [isAuthenticated, setIsAuthenticated] = useState(false);

	Gachi.createContext("isAuthenticated", { isAuthenticated, setIsAuthenticated })

	// Check if the user is authenticated on page load
	useEffect(() => {
		fetch('https://localhost:8080/authorized', {
			credentials: "include",
		})
			.then((response) => {
				if (response.ok) {
			  		setIsAuthenticated(true);
				} else if (response.status === 401) {
			  		setIsAuthenticated(false);
				}
		  	})
			.catch(() => setIsAuthenticated(false));
	}, []);

	return (
		<Router
			routes={[
				{ path: "/", element: <Home /> },
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