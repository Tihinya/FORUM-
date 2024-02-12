import Gachi, {
	useContext,
	useState,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { getWebSocket } from "../socket/socket"
import {
	SendMessageEvent,
	Event,
	RequestMessageHistoryEvent,
	IsTypingEvent,
} from "../socket/events"
import { convertTime } from "../../additional-funcitons/post"

export default function ChatBoxComponent({ recipientId }) {
	const [messages, setMessages] = useState([])
	const [timesScrolled, setTimesScrolled] = useState(1)
	const [previousRecipientId, setPreviousRecipientId] = useState(undefined)
	const [isTyping, setIsTyping] = useState(false)
	const isLoggin = useContext("isAuthenticated").isAuthenticated
	const { user } = useContext("currentUser")
	const { allUsers } = useContext("allUsers")
	const { onlineUserIdsList, setOnlineUserIdsList } =
		useContext("onlineUserIdsList")
	const { setOrderedUsersByLastMessage } = useContext(
		"orderedUsersByLastMessage"
	)

	const recipientUser = allUsers.find((user) => user.id === recipientId)
	const ws = getWebSocket()

	useEffect(() => {
		if (previousRecipientId != undefined) {
			let event = new RequestMessageHistoryEvent(recipientId)
			event = new Event("read_messages_history", event)
			ws.send(JSON.stringify(event))

			setTimeout(() => {
				const chatBox = document.querySelector(".chat-box")
				if (chatBox) {
					chatBox.scrollTop = chatBox.scrollHeight
				}
			}, 0)
		}
		setPreviousRecipientId(recipientId)
		setIsTyping(false)
	}, [recipientId])

	useEffect(() => {
		if (isLoggin) {
			ws.onmessage = (incomingEvent) => {
				const eventData = JSON.parse(incomingEvent.data)
				const event = Object.assign(new Event(), eventData)

				switch (event.type) {
					case "receive_message":
						if (
							event.payload.sender_id == user.id ||
							(event.payload.receiver_id == user.id &&
								event.payload.sender_id == recipientId)
						) {
							setMessages((prevMessages) => {
								return [event.payload, ...prevMessages]
							})
						}
						break

					case "receive_messages_history":
						setMessages(
							event.payload.slice(-10 * timesScrolled).reverse()
						)
						break

					case "online_users_list":
						setOnlineUserIdsList(event.payload.list)
						break

					case "receive_users_by_last_message":
						setOrderedUsersByLastMessage(event.payload.list)
						break

					case "typing_status":
						if (event.payload.sender_id == recipientId) {
							setIsTyping(event.payload.typing_status)
						}
						break
				}
			}
		}
	}, [isLoggin, timesScrolled, recipientId])

	useEffect(() => {
		if (isLoggin) {
			let event = new RequestMessageHistoryEvent(recipientId)
			event = new Event("read_messages_history", event)
			ws.send(JSON.stringify(event))
		}
	}, [isLoggin])

	const debouncedScrollEvent = debounce(
		(setTimesScrolled, recipientId, ws) => {
			setTimesScrolled(timesScrolled + 1)
			let event = new RequestMessageHistoryEvent(recipientId)
			event = new Event("read_messages_history", event)
			ws.send(JSON.stringify(event))
		}
	)

	const debounceTyping = debounce(() => {
		sendTypingEvent(false)
	}, 1000)

	const handleScroll = (event) => {
		const div = event.target
		const topPosition = div.scrollHeight - div.clientHeight
		const atTop = -div.scrollTop == topPosition

		if (atTop) {
			debouncedScrollEvent(setTimesScrolled, recipientId, ws)
		}
	}

	const handleTyping = () => {
		if (!isTyping) {
			sendTypingEvent(true)
		}
		debounceTyping()
	}

	const handleSubmit = (event) => {
		event.preventDefault()

		const form = event.target
		const formDatafied = new FormData(form)
		const formJson = Object.fromEntries(formDatafied.entries())

		if (formJson.message != null) {
			let event = new SendMessageEvent(
				formJson.message,
				user.id,
				recipientId
			)
			event = new Event("send_message", event)
			ws.send(JSON.stringify(event))
		}
		form.reset()
	}

	const sendTypingEvent = (typingStatus) => {
		let event = new IsTypingEvent(user.id, recipientId, typingStatus)
		event = new Event("typing_status", event)
		ws.send(JSON.stringify(event))
	}

	return (
		<div className="main__chatbox-container">
			<div className="chat-box-window">
				<div className="users">
					<div className="users-avatar-container">
						<div className="users-avatar">
							<div className="users-avatar-online-status-container">
								<div
									className={`users-avatar-online-status${
										onlineUserIdsList.includes(recipientId)
											? `-true`
											: ""
									}`}
								></div>
							</div>
						</div>
					</div>
					<div className="users-user-info">
						{recipientUser.username}
					</div>
					<div className="users-user-type-in-progress">
						{isTyping ? "User is typing..." : ""}
					</div>
				</div>
				<div onScroll={handleScroll} className="chat-box">
					{messages.map((message) => {
						if (message.sender_id == user.id) {
							return (
								<div className="message-box-reciever">
									<div className="message-box-for-reciever">
										<div className="message-box-additional-info">
											<div className="message-box-user-info">
												<p>{user.username}</p>
												<p>
													{convertTime(
														message.sent_time
													)}
												</p>
											</div>
											<div className="message-content">
												<p>{message.message}</p>
											</div>
										</div>
										<div className="users-avatar"></div>
									</div>
								</div>
							)
						} else {
							return (
								<div className="message-box-sender">
									<div className="message-box">
										<div className="users-avatar"></div>
										<div className="message-box-additional-info">
											<div className="message-box-user-info">
												<p>{recipientUser.username}</p>
												<p>
													{convertTime(
														message.sent_time
													)}
												</p>
											</div>
											<div className="message-content">
												<p>{message.message}</p>
											</div>
										</div>
									</div>
								</div>
							)
						}
					})}
				</div>
				<form
					className="message-send-container"
					onSubmit={handleSubmit}
				>
					<textarea
						name="message"
						className="message-input"
						rows={10}
						placeholder="Write here your message"
						onInput={handleTyping}
					/>
					<button className="send-button" type="submit"></button>
				</form>
			</div>
		</div>
	)
}

function debounce(func, timeout = 300) {
	let timer
	return (...args) => {
		clearTimeout(timer)
		timer = setTimeout(() => {
			func.apply(this, args)
		}, timeout)
	}
}
