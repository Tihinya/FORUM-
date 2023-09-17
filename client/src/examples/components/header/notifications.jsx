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
	const navigate = useNavigate()

	useEffect(() => {
		fetchNotifications()
	}, [])

	function clickNotifications() {
		setShowNotifications(!showNotifications)
		if (!showNotifications) {
			markNotificationsRead()
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

    return (
        <div className="notifications-container">
			<div
				className="sign__button" 
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
							.map((notification) => (
							<>
								{notification.status == "unread" ? 
									<li 
										className="notification"
										onClick={() => navigate(
											`/comments-authorized/${notification.parent_object_id}`
										)}
									>
										Your {notification.related_object_type} has been {notification.type}d
									</li>
								: 
									<li 
										className="notification-read"
										onClick={() => navigate(
											`/comments-authorized/${notification.parent_object_id}`
										)}
									>
										Your {notification.related_object_type} has been {notification.type}d
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