import Gachi, {
	useContext,
	useEffect,
	useNavigate,
	useState,
} from "../core/framework.ts"
import { importCss } from "../modules/cssLoader.js"
import Button from "./button.jsx"
importCss("./index.css")

const container = document.getElementsByClassName("post__tags")[0]

const text = <p className="tag">testingaaa</p>

Gachi.render(text, container)