import { PostContainer } from './post.js'
import Gachi from "../core/framework.ts"

const container = document.getElementsByClassName("main__container")[0]

function App() {
	return (
		<PostContainer />
	)
}

Gachi.render(<App />, container)


