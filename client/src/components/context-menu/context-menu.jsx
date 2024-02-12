import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../Gachi.js/src/core/framework.ts"

import { fetchData } from "../../additional-funcitons/api"
import ConfirmationWindow from "./confirmation-window"
import ErrorWindow from "../errors/error-window"

export default function ContextMenu({ obj, endpoint }) {
	const [contextType, setContextType] = useState(undefined)

	const [deleteUrl, setDeleteUrl] = useState("")
	const [editUrl, setEditUrl] = useState("")
	const [fetchUrl, setFetchUrl] = useState("")

	const [showButtonContent, setShowButtonContent] = useState(false)
	const [showConfirmationWindow, setShowConfirmationWindow] = useState(false)
	const [showEditInput, setShowEditInput] = useState(false)
	const [errorMessage1, setErrorMessage1] = useState(false)
	const { setErrorMessage } = useContext("currentErrorMessage")
	const { ownedPostsIds, setOwnedPostsIds } = useContext(
		"currentOwnedPostsIds"
	)
	const { ownedCommentsIds, setOwnedCommentsIds } = useContext(
		"currentOwnedCommentsIds"
	)
	const { posts, setPosts } = useContext("currentPosts")
	const { comments, setComments } = useContext("currentComments")
	const { userRole } = useContext("currentUserRole")
	const [reportPostWindow, setReportPostWindow] = useState("")

	if (obj.hasOwnProperty("comment_count")) {
		setContextType("post")
	} else if (obj.hasOwnProperty("post_id")) {
		setContextType("comment")
	}

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
	}, [contextType, userRole, comments, posts])

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
				setErrorMessage1(`${contextType} deletion failed`)
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
				setErrorMessage1(`${contextType} editing failed`)
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

	function reportAnswerButton(formData) {
		fetchData(formData, "postreport/create", "POST").then(
			(resultInJson) => {
				if (resultInJson.status === "error") {
					setErrorMessage(resultInJson.message)
				}
				setReportPostWindow("")
			}
		)
	}

	function handleErrorMessageClose() {
		setErrorMessage1("")
	}
	const toggleFilter = (filterType) => {
		if (reportPostWindow !== filterType) {
			setReportPostWindow("")
			setReportPostWindow(filterType)
		} else {
			setReportPostWindow("")
		}
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
			{errorMessage1 != "" ? (
				<ErrorWindow
					errorMessage={errorMessage1}
					onClose={handleErrorMessageClose}
				/>
			) : (
				""
			)}

			<div
				className="context-button"
				onClick={() => setShowButtonContent(!showButtonContent)}
			></div>
			<div
				className={`context_container ${
					!showButtonContent ? "hidden" : ""
				}`}
			>
				<a
					className={"context__button_hover"}
					onClick={() => {
						setShowEditInput(!showEditInput)
						setShowButtonContent(false)
					}}
				>
					Edit post
				</a>
				<a
					className={"context__button_hover"}
					onClick={() =>
						setShowConfirmationWindow(!showConfirmationWindow)
					}
				>
					Delete post
				</a>
				<a
					className={`context__button_hover ${
						reportPostWindow === "reportPost"
							? "nav__options_active"
							: ""
					}`}
					onClick={() => {
						toggleFilter("reportPost")
					}}
				>
					Report Post
				</a>
				{reportPostWindow === "reportPost" ? (
					<div className={"context_container report_box"}>
						<a
							className={"context__button_hover"}
							onClick={() => {
								const formData = {
									post_id: obj.id,
									message: "Irrelevant post",
								}
								reportAnswerButton(formData)
								setShowButtonContent(false)
							}}
							style={"cursor: pointer"}
						>
							Irrelevant Post
						</a>
						<a
							className={"context__button_hover"}
							onClick={() => {
								const formData = {
									post_id: obj.id,
									message: "Obscene Post",
								}
								reportAnswerButton(formData)
								setShowButtonContent(false)
							}}
							style={"cursor: pointer"}
						>
							Obscene Post
						</a>
						<a
							className={"context__button_hover"}
							onClick={() => {
								const formData = {
									post_id: obj.id,
									message: "Illegal Post",
								}
								reportAnswerButton(formData)
								setShowButtonContent(false)
							}}
							style={"cursor: pointer"}
						>
							Illegal Post
						</a>
						<a
							className={"context__button_hover"}
							onClick={() => {
								const formData = {
									post_id: obj.id,
									message: "Irrelevant post",
								}
								reportAnswerButton(formData)
								setShowButtonContent(false)
							}}
							style={"cursor: pointer"}
						>
							Insulting Post
						</a>
					</div>
				) : null}
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
