import Gachi, {useContext, useEffect} from "../../../core/framework"
import { fetchData } from "../../additional-funcitons/api"

export function ContextFetchOwnedPostsIds() {
    const { ownedPostsIds, setOwnedPostsIds } = useContext("currentOwnedPostsIds")
    const { posts, setPosts } = useContext("currentPosts")
    const isLoggin = useContext("isAuthenticated").isAuthenticated

    useEffect(() => {
        if (isLoggin) {
            fetchOwnedPosts()
        }
    }, [posts])
    
    function fetchOwnedPosts() {
        fetchData(null, `user/posts`, "GET").then((resultInJson) => {
            const postIds = resultInJson.map((ownedPost) => ownedPost.id)
            setOwnedPostsIds(postIds)
        })
    }
}

export function ContextFetchOwnedCommentsIds() {
    const { ownedCommentsIds, setOwnedCommentsIds } = useContext("currentOwnedCommentsIds")
    const { comments, setComments } = useContext("currentComments")
    const isLoggin = useContext("isAuthenticated").isAuthenticated

    useEffect(() => {
        if (isLoggin) {
            fetchOwnedComments()
        }
    }, [comments])
    
    function fetchOwnedComments() {
        fetchData(null, `user/createdcomments`, "GET").then((resultInJson) => {
            const commentsIds = resultInJson.map((ownedComment) => ownedComment.id)
            setOwnedCommentsIds(commentsIds)
        })
    }
}
