import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api"
import ConfirmationWindow from "./confirmationwindow"

export default function PostContextMenu( {post} ) {
    const [ showButtonContent, setShowButtonContent ] = useState(false)
    const [ showConfirmationWindow, setShowConfirmationWindow ] = useState(false)
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


    function editPost(post) {
        // do stuff lmao
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
            
            <div 
                className="edit-button"
                onClick={() => setShowButtonContent(!showButtonContent)}
            >
                <button
                    className={showButtonContent ? "edit-button-content-edit" : "edit-button-content-hidden"}
                    onClick={() => editPost(post)}
                >
                    Edit post
                </button>
                <button
                    className={showButtonContent ? "edit-button-content-delete" : "edit-button-content-hidden"}
                    onClick={() => setShowConfirmationWindow(!showConfirmationWindow)}
                >
                    Delete post
                </button>
            </div>
        </div>
    )
}