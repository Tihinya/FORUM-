import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api"
import ConfirmationWindow from "./confirmationwindow"
import ErrorWindow from "../errors/error-window"

export default function PostContextMenu( {post} ) {
    const [ showButtonContent, setShowButtonContent ] = useState(false)
    const [ showConfirmationWindow, setShowConfirmationWindow ] = useState(false)
    const [ showEditInput, setShowEditInput ] = useState(false)
    const [ errorMessage, setErrorMessage ] = useState("")
    const { posts, setPosts } = useContext("currentPosts")

    function deletePost() {
        fetchData(null, `post/${post.id}`, "DELETE").then((responseInJson) => {
            if (!responseInJson.status === "success") {
                return
            }

            fetchData(null, "posts", "GET").then((resultInJson) => {
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

        fetchData(formJson, `post/${post.id}`, "PATCH").then((responseInJson) => {
            console.log(responseInJson)
            if (responseInJson.status !== "success") {
                setErrorMessage("Post editing failed")
                return
            }
            fetchData(null, "posts", "GET").then((resultInJson) => {
                setPosts(resultInJson)
            })

        })

        setShowButtonContent(!showButtonContent)
        setShowEditInput(false)
    }

    function dismissEdit() {
        setErrorMessage("test")
        setShowButtonContent(false)
        setShowEditInput(false)
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
					onClose={setErrorMessage("")}
				/>
			) : (
                ""
            )}

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
                    onClick={() => dismissEdit()}
                >
                    Cancel
                </button>
            </form>
            
            <div 
                className="context-button"
                onClick={() => setShowButtonContent(!showButtonContent)}
            >
                <button
                    className={showButtonContent ? 
                        "context-button-content-edit" 
                    : 
                        "context-button-content-hidden"}
                    onClick={() => setShowEditInput(!showEditInput)}
                >
                    Edit post
                </button>
                <button
                    className={showButtonContent ? 
                        "context-button-content-delete" 
                    : 
                        "context-button-content-hidden"}
                    onClick={() => setShowConfirmationWindow(!showConfirmationWindow)}
                >
                    Delete post
                </button>
            </div>
        </div>
    )
}