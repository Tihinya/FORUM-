import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"
import { fetchData } from "../../additional-funcitons/api"
import { Event } from "../socket/events"
import { getWebSocket } from "../socket/socket"

export function UserListBar() {
	const [isConnected, setIsConnected] = useState(false)
	const { user, setUser } = useContext("currentUser")
	const { allUsers, setAllUsers } = useContext("allUsers")
	const { onlineUserIdsList, setOnlineUserIdsList } =
		useContext("onlineUserIdsList")
	const { orderedUserList, setOrderedUserList } =
		useContext("orderedUserList")
	const { orderedUsersByLastMessage, setOrderedUsersByLastMessage } =
		useContext("orderedUsersByLastMessage")

	const { chatboxRecipient, setChatboxRecipient } =
		useContext("chatboxRecipient")
	const { notificationSound } = useContext("notificationSound")
	const isLoggin = useContext("isAuthenticated").isAuthenticated

	const fetchUser = `user/me`
	const fetchAllUsers = `users/get`

	let ws

	useEffect(() => {
		if (isLoggin && chatboxRecipient == 0) {
			fetchData(null, fetchUser, "GET").then((resultInJson) => {
				setUser(resultInJson)
			})

			fetchData(null, fetchAllUsers, "GET").then((resultInJson) => {
				setAllUsers(resultInJson)
			})

			ws = getWebSocket()

			ws.onopen = () => {
				setIsConnected(true)
			}

			ws.onmessage = (incomingEvent) => {
				const eventData = JSON.parse(incomingEvent.data)
				const event = Object.assign(new Event(), eventData)

				switch (event.type) {
					case "online_users_list":
						setOnlineUserIdsList(event.payload.list)
						break

					case "receive_users_by_last_message":
						setOrderedUsersByLastMessage(event.payload.list)
						break

					case "receive_message":
						if (notificationSound == undefined) {
							break
						}
						notificationSound.play()
						break
				}
			}
		}
	}, [isLoggin, chatboxRecipient, notificationSound])

	useEffect(() => {
		if (isLoggin) {
			setOrderedUserList(
				orderUserList(user, allUsers, orderedUsersByLastMessage)
			)
		}
	}, [onlineUserIdsList, allUsers, user, orderedUsersByLastMessage])

	if (isLoggin) {
		return (
			<div className="online-users-main-container">
				<div className="online-users-container">
					<div className="users-list-menu-name">
						<p1>Users</p1>
					</div>
					{orderedUserList.map((user) => {
						return (
							<>
								<div
									className="users"
									onClick={() => {
										if (isConnected) {
											setChatboxRecipient(user.id)
										}
									}}
								>
									<div className="user-online-status">
										<div className="users-avatar-online-status-container">
											<div
												className={`users-avatar-online-status${
													onlineUserIdsList.includes(
														user.id
													)
														? `-true`
														: ""
												}`}
											></div>
										</div>
									</div>
									<div className="users-user-info">
										{user.username}
									</div>
								</div>
							</>
						)
					})}
				</div>
			</div>
		)
	}
}

function orderUserList(user, allUsers, orderedUsersByLastMessage) {
	const allUsersCopy = [...allUsers]

	// Remove logged in user from user list
	const indexToRemove = allUsersCopy.findIndex((u) => u.id === user.id)
	if (indexToRemove > -1) {
		allUsersCopy.splice(indexToRemove, 1)
	}

	// First create a new map of user objects with
	// the order of orderedUsersByLastMessage
	const usersByLastMessage = orderedUsersByLastMessage
		.map((userId) => allUsersCopy.find((user) => user.id === userId))
		.filter((user) => user !== undefined)

	// Filter out all users that have no conversations
	// with logged in user
	const remainingUsers = allUsersCopy.filter(
		(user) => !orderedUsersByLastMessage.includes(user.id)
	)

	// Sort remaining users alphabetically that
	const sortedRemainingUsers = remainingUsers.sort((a, b) =>
		a.username.localeCompare(b.username)
	)

	const orderedUsers = usersByLastMessage.concat(sortedRemainingUsers)

	return orderedUsers
}
