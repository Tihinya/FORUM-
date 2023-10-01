import Gachi, {
	useContext,
	useState,
	useNavigate,
	useEffect,
} from "../../../core/framework"
export default function Categories({ post }) {
	return (
		<div className="post__tags">
			{post.categories.map((categories) => (
				<p className="tag">{categories}</p>
			))}
		</div>
	)
}
