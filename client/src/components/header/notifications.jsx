import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
let baseURL = "https://ec2-51-20-1-125.eu-north-1.compute.amazonaws.com:8080"

function markNotificationsRead() {
	fetch(baseURL + "readnotifications", {
		credentials: "include",
		method: "POST",
	})
		.then((response) => response.json())
		.then((data) => data)
		.catch((error) =>
			console.error("Error marking notifications as read", error)
		)
}

export function Notifications() {
	const [notifications, setNotifications] = useState([])
	const [showNotifications, setShowNotifications] = useState(false)
	const [notificationsUnread, setNotificationsUnread] = useState(false)
	const [clickedNotificationButton, setClickedNotificationButton] =
		useState(false)
	const { notificationsInterval, setNotificationsInterval } = useContext(
		"notificationsInterval"
	)
	const isLoggin = useContext("isAuthenticated").isAuthenticated
	const navigate = useNavigate()

	useEffect(() => {
		fetchNotifications()
		if (notificationsInterval == 0) {
			setNotificationsInterval(
				setInterval(() => {
					fetchNotifications()
				}, 5000)
			)
		}
		//return () => clearInterval(notificationsInterval) // Doesn't work
	}, [])

	useEffect(() => {
		checkForUnreadNotifications()
	}, [notifications])

	function clickNotifications() {
		setShowNotifications(!showNotifications)
		if (clickedNotificationButton) {
			markNotificationsRead()
		}
		setClickedNotificationButton(true)

		// Change notification button color when window closed
		if (showNotifications) {
			setNotificationsUnread(false)
		}
	}

	function fetchNotifications() {
		fetch(baseURL + "/notifications", {
			credentials: "include",
		})
			.then((response) => response.json())
			.then((data) => setNotifications(data))
			.catch((error) =>
				console.error("Error fetching notifications:", error)
			)
	}

	function checkForUnreadNotifications() {
		if (notifications.length > 0) {
			if (
				notifications.some(
					(notification) => notification.status == "unread"
				)
			) {
				setNotificationsUnread(true)
				return
			}
		}
	}

	return (
		<div className="notifications-container">
			<div
				className={
					notificationsUnread ? "sign__button_unread" : "sign__button"
				}
				onClick={() => {
					clickNotifications()
				}}
			>
				Notifications
			</div>

			{showNotifications ? (
				<div className="notifications-window">
					<p>Your notifications</p>
					<ul className="notifications-window-list">
						{notifications
							.sort(
								(a, b) =>
									new Date(b.creation_date) -
									new Date(a.creation_date)
							)
							.slice(0, 12)
							.map((notification) => (
								<>
									{notification.status == "unread" ? (
										<li
											className="notification"
											onClick={() => {
												clickNotifications()
												navigate(
													`/comments-authorized/${notification.parent_object_id}`
												)
											}}
										>
											{notification.type == "comment"
												? `Your ${notification.related_object_type} has been ${notification.type}ed`
												: `Your ${notification.related_object_type} has been ${notification.type}d`}
										</li>
									) : (
										<li
											className="notification-read"
											onClick={() => {
												clickNotifications()
												navigate(
													`/comments-authorized/${notification.parent_object_id}`
												)
											}}
										>
											{notification.type == "comment"
												? `Your ${notification.related_object_type} has been ${notification.type}ed`
												: `Your ${notification.related_object_type} has been ${notification.type}d`}
										</li>
									)}
								</>
							))}
					</ul>
				</div>
			) : null}
		</div>
	)
}
