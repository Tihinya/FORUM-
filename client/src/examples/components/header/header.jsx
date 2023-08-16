import Gachi, {
	useContext,
	useState,
	useNavigate,
} from "../../../core/framework"
import DropdownMenu from "./dropdown"

export default function Header() {
	const navigate = useNavigate()
	const isLogin = false
	return (
		<div className="header">
			<div className="header__logo">
				<p>Cartel Forum</p>
			</div>
			<input className="search__bar" placeholder="Search in progres..." />
			{!isLogin ? (
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
					<div className="sign__button" id="notification-button">
						Notification
					</div>
					<div className="profile-menu">
						<div className="profile-nav">
							<div className="user__info_picture">
								<a href="/src/html/profile-page.html">
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
