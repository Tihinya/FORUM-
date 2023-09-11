import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

export function Notifications() {
	//const {notifications, setNotifications } = useContext("notifications")
	const [notifications, setNotifications] = useState([])
	const [showNotifications, setShowNotifications] = useState(false)

	useEffect(() => {
		fetch("https://localhost:8080/notifications", {
            credentials: "include"
        })
			.then((response) => response.json())
			.then((data) => setNotifications(data))
			.catch((error) => console.error("Error fetching notifications:", error))
	}, [])

    return (
        <div className="notifications-container">
			<div
				className="sign__button" 
				onClick={() => {setShowNotifications(!showNotifications)}}
				>
				Notifications
			</div>

			{ showNotifications && (
                    <div>
                        <p>Your notifications</p>
                    
                    <ul>
                        <li>TEST</li>
                    </ul>
                </div>
            )}

		</div>
    )
}