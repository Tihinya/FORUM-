import { PostContainer } from './post.js'
import Gachi from "../core/framework.ts"
import { Router, Route } from "/src/components/router.ts"

const container = document.getElementsByClassName("main__container")[0]

function App() {
	return (
	<Router>
		<Route path="/" element={<PostContainer />} />
	</Router>
	)
}

Gachi.render(<App />, container)