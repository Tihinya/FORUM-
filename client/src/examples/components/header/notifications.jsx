import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

function markNotificationsRead() {
	fetch("https://localhost:8080/readnotifications", {
		credentials: "include",
		method: "POST"
	})
		.then((response) => response.json())
		.then((data) => data)
		.catch((error) => console.error("Error marking notifications as read", error))
}

export function Notifications() {
	const [notifications, setNotifications] = useState([])
	const [showNotifications, setShowNotifications] = useState(false)
	const [notificationsUnread, setNotificationsUnread] = useState(false)
	const [clickedNotificationButton, setClickedNotificationButton] = useState(false)
	const navigate = useNavigate()

	useEffect(() => {
		fetchNotifications()

		const interval = setInterval(() => {
			// Future reference: React useRef()
			fetchNotifications()
		}, 5000)
	
		return () => clearInterval(interval)
	}, [])

	useEffect(() => {
		checkForUnreadNotifications();
	}, [notifications]);

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
		fetch("https://localhost:8080/notifications", {
			credentials: "include"
		})
			.then((response) => response.json())
			.then((data) => setNotifications(data))
			.catch((error) => console.error("Error fetching notifications:", error))
	}

	function checkForUnreadNotifications() {
		if (notifications.length > 0) {
			if (notifications.some((notification) => notification.status == "unread")) {
				setNotificationsUnread(true)
				return
			}
		}
	}

    return (
        <div className="notifications-container">
			<div
				className={notificationsUnread ? "sign__button_unread" : "sign__button"}
				onClick={() => {
					clickNotifications()
				}}
				>
				Notifications
			</div>

			{ showNotifications ? (
				<div className="notifications-window">
					<p>Your notifications</p>
					<ul className="notifications-window-list">
						{notifications
							.sort(
								(a, b) =>
									new Date(b.creation_date) - new Date(a.creation_date)
							)
							.slice(0, 12)
							.map((notification) => (
							<>
								{notification.status == "unread" ? 
									<li 
										className="notification"
										onClick={() => {
											clickNotifications()
											navigate(
											`/comments-authorized/${notification.parent_object_id}`
											)
										}}
									>
										{notification.type == "comment" ? 
											`Your ${notification.related_object_type} has been ${notification.type}ed` 
										: 
											`Your ${notification.related_object_type} has been ${notification.type}d`
										}
									</li>
								: 
									<li 
										className="notification-read"
										onClick={() => {
											clickNotifications()
											navigate(
											`/comments-authorized/${notification.parent_object_id}`
											)
										}}
									>
										{notification.type == "comment" ? 
											`Your ${notification.related_object_type} has been ${notification.type}ed` 
										: 
											`Your ${notification.related_object_type} has been ${notification.type}d`
										}
									</li>
								}
							</>
						))}
					</ul>
				</div>
			) : null }

		</div>
    )
}