import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api"
import ConfirmationWindow from "./confirmation-window"
import ErrorWindow from "../errors/error-window"

export default function PostContextMenu( {post} ) {
    const deletePostUrl = `post/${post.id}`
    const editPostUrl = `post/${post.id}`
    const postsUrl = `posts`
    const [ showButtonContent, setShowButtonContent ] = useState(false)
    const [ showConfirmationWindow, setShowConfirmationWindow ] = useState(false)
    const [ showEditInput, setShowEditInput ] = useState(false)
    const [ errorMessage, setErrorMessage ] = useState(false)
    const { ownedPostsIds, setOwnedPostsIds } = useContext("currentOwnedPostsIds")
    const { posts, setPosts } = useContext("currentPosts")

    if (!ownedPostsIds.includes(post.id)) {
        return
    }

    function deletePost() {
        fetchData(null, deletePostUrl, "DELETE").then((responseInJson) => {
            if (responseInJson.status !== "success") {
                setErrorMessage("Post deletion failed")
                return
            }

            fetchData(null, postsUrl, "GET").then((resultInJson) => {
                setPosts(resultInJson)
            })
        })
        
        setShowButtonContent(false)
        setShowConfirmationWindow(false)
    }

    function dismissDeletion() {
        setShowButtonContent(false)
        setShowConfirmationWindow(false)
    }

    function editPost(e) {
        e.preventDefault()
        const form = e.target
		const formData = new FormData(form)
		const formJson = Object.fromEntries(formData.entries())

        fetchData(formJson, editPostUrl, "PATCH").then((responseInJson) => {
            if (responseInJson.status !== "success") {
                setErrorMessage("Post editing failed")
                return
            }

            fetchData(null, postsUrl, "GET").then((resultInJson) => {
                setPosts(resultInJson)
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
            {showConfirmationWindow ? 
                <ConfirmationWindow 
                    message={"Are you sure you want to delete your post?"}
                    onYes={() => deletePost()}
                    onNo={() => dismissDeletion()}
                />
            : null }
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
                    className={
                        `context-button-content-edit ${!showButtonContent ? "hidden" : ""}`
                    }
                    onClick={() => setShowEditInput(!showEditInput)}
                >
                    Edit post
                </button>
                <button
                    className={
                        `context-button-content-delete ${!showButtonContent ? "hidden" : ""}`
                    }
                    onClick={() => setShowConfirmationWindow(!showConfirmationWindow)}
                >
                    Delete post
                </button>
            </div>

            <form onSubmit={editPost}>
                <input
                    className={`edit-button-title-window ${!showEditInput ? "hidden" : ""}`}
                    name="title"
                    id="titleValue"
                    defaultValue={post.title}
                />
                <textarea
                    className={`edit-button-content-window ${!showEditInput ? "hidden" : ""}`}
                    name="content"
                    id="contentValue"
                    defaultValue={post.content}
                />
                
                <button
                    className={`edit-button-publish ${!showEditInput ? "hidden" : ""}`}
                    type="submit"
                >
                    Publish
                </button>
                <button
                    className={`edit-button-cancel ${!showEditInput ? "hidden" : ""}`}
                    type="button"
                    onClick={() => dismissEdit()}
                >
                    Cancel
                </button>
            </form>
        </div>
    )
}