import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

import { fetchData } from "../../additional-funcitons/api.js"

export default function LikesAndDislikes({ post, page }) {
	const navigate = useNavigate()
	const likeUrl = "user/liked"
	const disLikeUrl = "user/disliked"

	const { setPosts } = useContext("currentPosts")
	const [likedPosts, setLikedPosts] = useState([])
	const [dislikedPosts, setDislikedPosts] = useState([])

	const fetchLikes = () => {
		fetchData(null, likeUrl, "GET").then((resultInJson) => {
			setLikedPosts(resultInJson)
		})
	}

	const fetchDislikes = () => {
		fetchData(null, disLikeUrl, "GET").then((resultInJson) => {
			setDislikedPosts(resultInJson)
		})
	}

	// console.log(dislikedPosts)

	useEffect(() => {
		fetchLikes()
		fetchDislikes()
	}, [])

	// Create a helper function for making the POST requests
	const handleLike = (type, postId) => {
		const isLiking =
			!likedPosts.some((obj) => obj.id === postId) &&
			!dislikedPosts.includes(postId)
		const endpoint = isLiking ? `${postId}/${type}` : `${postId}/un${type}`

		fetchData(null, `${page}/${endpoint}`, "POST").then((resultInJson) => {
			if (resultInJson.status === "success") {
				fetchData(null, "posts", "GET").then((resultInJson) => {
					setPosts(resultInJson)
					fetchLikes()
					fetchDislikes()
				})
			} else if (resultInJson.status === "error") {
				setErrorMessage(resultInJson.message)
			}
		})
		// if (response.status === "success") {
		// 	// setPosts((prevPosts) => {
		// 	// 	return prevPosts.map((post) => {
		// 	// 		if (post.id === postId) {
		// 	// 			return {
		// 	// 				...post,
		// 	// 				likes:
		// 	// 					type === "like"
		// 	// 						? post.likes + (isLiking ? 1 : -1)
		// 	// 						: post.likes,
		// 	// 				dislikes:
		// 	// 					type === "dislike"
		// 	// 						? post.dislikes + (isLiking ? 1 : -1)
		// 	// 						: post.dislikes,
		// 	// 			}
		// 	// 		}
		// 	// 		return post
		// 	// 	})
		// 	// })
		// 	fetchLikes()
		// 	fetchDislikes()
		// } else {
		// 	console.error(
		// 		`Request failed with status: ${response.status} - ${response.message}`
		// 	)
		// }
		// })
	}

	return (
		<div className="post__likes">
			<img
				onClick={() => handleLike("like", post.id)}
				src="../img/thumbs-up.svg"
			/>
			<p onClick={() => handleLike("like", post.id)}>{post.likes}</p>
			<img
				onClick={() => handleLike("dislike", post.id)}
				src="../img/thumbs-down.svg"
			/>
			<p onClick={() => handleLike("dislike", post.id)}>
				{post.dislikes}
			</p>
		</div>
	)
}
