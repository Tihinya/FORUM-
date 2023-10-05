import Gachi, {useContext, useEffect} from "../../../core/framework"
import { fetchData } from "../../additional-funcitons/api"

export function ContextFetchOwnedPostIds() {
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

