import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../Gachi.js/src/core/framework.ts"
import AdminNavBar from "./adminNavBar.jsx"
import ModeratorNavBar from "./moderatorNavBar.jsx"
import UserNavBar from "./userNavBar.jsx"

export default function PersonalNavBar() {
	// const navigateTo = useNavigate()
	const { props, setProps } = useContext("currentProps")
	const { userRole } = useContext("currentUserRole")
	const toggleFilter = (filterType) => {
		if (props !== filterType) {
			setProps("")
			setProps(filterType)
		} else {
			setProps("")
		}
	}

	return (
		<div className="nav__menu_profile">
			<div className="personal-navigator__menu">
				<p>Personal Navigator</p>
			</div>
			<div className="personal-navigator__categories">
				<div
					className={`nav__options_hover ${
						props === "user/posts" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/posts")
					}}
				>
					My Posts
				</div>
				<div
					className={`nav__options_hover ${
						props === "user/liked" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/liked")
					}}
				>
					My Likes
				</div>
				<div
					className={`nav__options_hover ${
						props === "user/disliked" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/disliked")
					}}
				>
					My Dislikes
				</div>
				<div
					className={`nav__options_hover ${
						props === "user/comments" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/comments")
					}}
				>
					My Comments
				</div>
			</div>
			{userRole === "admin" ? (
				<AdminNavBar />
			) : userRole === "moderator" ? (
				<ModeratorNavBar />
			) : (
				<UserNavBar />
			)}
		</div>
	)
}
