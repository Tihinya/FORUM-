import Gachi, { useNavigate, useContext } from "../../../core/framework"
import DropdownMenu from "./dropdown"
import { Notifications } from "./notifications"

export default function Header() {
	const navigate = useNavigate()
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	return (
		<div className="header">
			<div className="header__logo">
				<a
					onClick={() => {
						navigate("/")
						// window.location.reload() // Reload the page
					}}
				>
					Cartel Forum
				</a>
			</div>
			<input className="search__bar" placeholder="Search in progres..." />
			{!isLoggin ? (
				<>
					<a
						className="sign__button"
						onClick={() => navigate("/login")}
					>
						Sign In
					</a>
					<a
						className="sign__button"
						onClick={() => navigate("/registration")}
					>
						Sign Up
					</a>
				</>
			) : (
				<>
					<Notifications />
					<div className="profile-menu">
						<div className="profile-nav">
							<div className="user__info_picture">
								<a onClick={() => navigate("/profile-page")}>
									<img src="../img/avatarka.jpeg" />
								</a>
							</div>
							<DropdownMenu />
						</div>
					</div>
				</>
			)}
		</div>
	)
}
