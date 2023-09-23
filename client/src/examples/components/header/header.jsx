import Gachi, { useNavigate } from "../../../core/framework"
import DropdownMenu from "./dropdown"
import isLogin from "../../additional-funcitons/isLogin.js"
import { Notifications } from "./notifications"

export default function Header() {
	const navigate = useNavigate()
	const isLoggin = isLogin()

	return (
		<div className="header">
			<div className="header__logo">
				<a onClick={() => navigate("/")}>Cartel Forum</a>
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
					{/* <div
					className="sign__button" 
					id="notification-button"
					onClick={() => console.log(Notifications)}
					> 
						Notification
					</div>*/}
					<Notifications />
					<div className="profile-menu">
						<div className="profile-nav">
							<div className="user__info_picture">
								<a onClick={() => navigate("profile-page")}>
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
