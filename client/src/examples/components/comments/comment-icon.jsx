import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"

export default function CommentsIcon({ post }) {
	const navigate = useNavigate()

	return (
		<div className="post__info">
			<a
				onClick={() => {
					navigate(`/comments-authorized/${post.id}`)
				}}
			>
				<img src="../img/message-square.svg" />
			</a>
			<p
				onClick={() => {
					navigate(`/comments-authorized/${post.id}`)
				}}
			>
				{post.comment_count}
			</p>
		</div>
	)
}
