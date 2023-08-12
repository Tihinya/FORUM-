import { PostContainer } from './post.js'
import { PostContainerAuth } from './postAuth.js'
import Gachi from "../core/framework.ts"
import { Router, Route } from "/src/components/router.ts"

const container = document.getElementsByClassName("main__container")[0]

function App() {
	return (
	<Router>
		<Route path="/" element={<PostContainer />} />
		<Route path="/authorized" element={<PostContainerAuth />} />
	</Router>
	)
}

Gachi.render(<App />, container)