import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api"
import ConfirmationWindow from "./confirmation-window"
import ErrorWindow from "../errors/error-window"
import { UserRole } from "./helpers.jsx"

export default function ContextMenu({ obj, endpoint }) {
	const [contextType, setContextType] = useState(undefined)

	const [deleteUrl, setDeleteUrl] = useState("")
	const [editUrl, setEditUrl] = useState("")
	const [fetchUrl, setFetchUrl] = useState("")

	const [showButtonContent, setShowButtonContent] = useState(false)
	const [showConfirmationWindow, setShowConfirmationWindow] = useState(false)
	const [showEditInput, setShowEditInput] = useState(false)
	const [errorMessage, setErrorMessage] = useState(false)
	const { ownedPostsIds, setOwnedPostsIds } = useContext(
		"currentOwnedPostsIds"
	)
	const { ownedCommentsIds, setOwnedCommentsIds } = useContext(
		"currentOwnedCommentsIds"
	)
	const { posts, setPosts } = useContext("currentPosts")
	const { comments, setComments } = useContext("currentComments")
	const { userRole } = useContext("currentUserRole")

	if (obj.hasOwnProperty("comment_count")) {
		setContextType("post")
	} else if (obj.hasOwnProperty("post_id")) {
		setContextType("comment")
	}

	// Duplicate useEffects cuz multi-dependancy useEffect brokie
	useEffect(() => {
		if (contextType == "post") {
			if (userRole == "moderator") {
				setDeleteUrl(`post/${obj.id}/mod`)
			} else {
				setDeleteUrl(`post/${obj.id}`)
				setEditUrl(`post/${obj.id}`)
			}
			if (endpoint == "post") {
				setFetchUrl(`post/${obj.id}`)
			} else if (endpoint == "posts") {
				setFetchUrl(`posts`)
			}
		} else if (contextType == "comment") {
			setDeleteUrl(`comment/${obj.id}`)
			setEditUrl(`comment/${obj.id}`)
			setFetchUrl(`comments/${obj.post_id}`)
		}
	}, [contextType])

	useEffect(() => {
		if (contextType == "post") {
			if (userRole == "moderator") {
				setDeleteUrl(`post/${obj.id}/mod`)
			} else {
				setDeleteUrl(`post/${obj.id}`)
				setEditUrl(`post/${obj.id}`)
			}
			if (endpoint == "post") {
				setFetchUrl(`post/${obj.id}`)
			} else if (endpoint == "posts") {
				setFetchUrl(`posts`)
			}
		} else if (contextType == "comment") {
			setDeleteUrl(`comment/${obj.id}`)
			setEditUrl(`comment/${obj.id}`)
			setFetchUrl(`comments/${obj.post_id}`)
		}
	}, [posts])

	useEffect(() => {
		if (contextType == "post") {
			if (userRole == "moderator") {
				setDeleteUrl(`post/${obj.id}/mod`)
			} else {
				setDeleteUrl(`post/${obj.id}`)
				setEditUrl(`post/${obj.id}`)
			}
			if (endpoint == "post") {
				setFetchUrl(`post/${obj.id}`)
			} else if (endpoint == "posts") {
				setFetchUrl(`posts`)
			}
		} else if (contextType == "comment") {
			setDeleteUrl(`comment/${obj.id}`)
			setEditUrl(`comment/${obj.id}`)
			setFetchUrl(`comments/${obj.post_id}`)
		}
	}, [comments])

	if (userRole !== "moderator") {
		if (!ownedPostsIds.includes(obj.id) && contextType == "post") {
			return
		}
	}

	if (!ownedCommentsIds.includes(obj.id) && contextType == "comment") {
		return
	}

	function deleteObj() {
		fetchData(null, deleteUrl, "DELETE").then((responseInJson) => {
			if (responseInJson.status !== "success") {
				setErrorMessage(`${contextType} deletion failed`)
				return
			}

			fetchData(null, fetchUrl, "GET").then((resultInJson) => {
				if (contextType == "post") {
					setPosts(resultInJson)
				} else {
					setComments(resultInJson)
				}
			})
		})

		setShowButtonContent(false)
		setShowConfirmationWindow(false)
	}

	function dismissDeletion() {
		setShowButtonContent(false)
		setShowConfirmationWindow(false)
	}

	function editObj(e) {
		e.preventDefault()
		const form = e.target
		const formData = new FormData(form)
		const formJson = Object.fromEntries(formData.entries())

		if (contextType == "post") {
			formJson.categories = obj.categories
		}

		fetchData(formJson, editUrl, "PATCH").then((responseInJson) => {
			if (responseInJson.status !== "success") {
				setErrorMessage(`${contextType} editing failed`)
				return
			}

			fetchData(null, fetchUrl, "GET").then((resultInJson) => {
				if (contextType == "post") {
					setPosts(resultInJson)
				} else {
					setComments(resultInJson)
				}
			})
		})

		setShowButtonContent(false)
		setShowEditInput(false)
	}

	function dismissEdit() {
		setShowButtonContent(false)
		setShowEditInput(false)
	}

	function handleErrorMessageClose() {
		setErrorMessage("")
	}

	return (
		<div>
			{showConfirmationWindow ? (
				<ConfirmationWindow
					message={"Are you sure you want to delete your post?"}
					onYes={() => deleteObj()}
					onNo={() => dismissDeletion()}
				/>
			) : null}
			{errorMessage != "" ? (
				<ErrorWindow
					errorMessage={errorMessage}
					onClose={handleErrorMessageClose}
				/>
			) : (
				""
			)}

			<div
				className="context-button"
				onClick={() => setShowButtonContent(!showButtonContent)}
			>
				<button
					className={`context-button-content-edit ${
						!showButtonContent ? "hidden" : ""
					}`}
					onClick={() => setShowEditInput(!showEditInput)}
				>
					Edit post
				</button>
				<button
					className={`context-button-content-delete ${
						!showButtonContent ? "hidden" : ""
					}`}
					onClick={() =>
						setShowConfirmationWindow(!showConfirmationWindow)
					}
				>
					Delete post
				</button>
			</div>

			<form onSubmit={editObj}>
				{contextType === "post" ? (
					<input
						className={`edit-button-title-window ${
							!showEditInput ? "hidden" : ""
						}`}
						name="title"
						id="titleValue"
						defaultValue={obj.title}
					/>
				) : null}

				<textarea
					className={`edit-button-content-window 
                    ${contextType == "comment" ? "comment" : ""}
                    ${!showEditInput ? "hidden" : ""}`}
					name="content"
					id="contentValue"
					defaultValue={obj.content}
				/>

				<button
					className={`edit-button-publish 
                    ${contextType == "comment" ? "comment" : ""}
                    ${!showEditInput ? "hidden" : ""}`}
					type="submit"
				>
					Publish
				</button>
				<button
					className={`edit-button-cancel 
                    ${contextType == "comment" ? "comment" : ""}
                    ${!showEditInput ? "hidden" : ""}`}
					type="button"
					onClick={() => dismissEdit()}
				>
					Cancel
				</button>
			</form>
		</div>
	)
}
