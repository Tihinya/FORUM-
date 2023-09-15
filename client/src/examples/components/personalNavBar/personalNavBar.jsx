import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"

export default function PersonalNavBar() {
	// const navigateTo = useNavigate()
	const { top, setTop } = useContext("currentTop")

	const toggleFilter = (filterType) => {
		if (top !== filterType) {
			setTop("")
			setTop(filterType)
		} else {
			setTop("")
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
						top === "user/posts" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/posts")
					}}
				>
					My Posts
				</div>
				<div
					className={`nav__options_hover ${
						top === "user/liked" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/liked")
					}}
				>
					My Likes
				</div>
				<div
					className={`nav__options_hover ${
						top === "user/comments" ? "nav__options_active" : ""
					}`}
					onClick={() => {
						toggleFilter("user/comments")
					}}
				>
					My Comments
				</div>
			</div>
			<div className="personal-navigator__menu">
				<p>Admin Navigator</p>
			</div>
			<div className="personal-navigator__categories">
				<a>
					<p className="nav__options_hover" id="hover-option">
						Promotions
					</p>
				</a>
				<a>
					<p className="nav__options_hover" id="hover-option">
						Create Category
					</p>
				</a>
			</div>
		</div>
	)
}
