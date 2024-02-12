import Gachi, {
	useNavigate,
	useContext,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import DropdownMenu from "./dropdown"
import { Notifications } from "./notifications"
import { fetchData } from "../../additional-funcitons/api"

import ProfilePicture from "../../img/avatarka.jpeg"

export default function Header() {
	const navigate = useNavigate()
	const isLoggin = useContext("isAuthenticated").isAuthenticated
	const { notificationsInterval } = useContext("notificationsInterval")
	const { setChatboxRecipient } = useContext("chatboxRecipient")

	// Workaround for notifications component useEffect cleanup not working
	useEffect(() => {
		if (!isLoggin) {
			clearInterval(notificationsInterval)
		}
	}, [isLoggin])

	return (
		<div className="header">
			<div className="header__logo">
				<a
					onClick={() => {
						setChatboxRecipient(0)
						navigate("/")
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
									<img src={ProfilePicture} />
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
