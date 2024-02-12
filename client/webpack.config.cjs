const path = require("path")
const HtmlWebpackPlugin = require("html-webpack-plugin")

module.exports = {
	mode: "production",
	entry: "./src/index.js", // Entry point of the application
	output: {
		path: path.resolve(__dirname, "dist"), // Output directory path
		filename: "[name].[contenthash].js", // Output bundle file name
		clean: true,
		publicPath: "/",
	},
	module: {
		rules: [
			{
				test: /\.(js|jsx|ts|tsx)$/,
				exclude: /node_modules/, // Exclude node_modules folder from processing
				use: {
					loader: "babel-loader", // Use Babel to transpile JavaScript
					options: {
						presets: [
							[
								"@babel/preset-typescript",
								{
									jsxPragma: "Gachi.createElement",
									jsxPragmaFrag: "Gachi.Fragment",
								},
							],
						],
						plugins: [
							[
								"@babel/plugin-transform-react-jsx",
								{
									pragma: "Gachi.createElement",
									pragmaFrag: "Gachi.Fragment",
								},
							],
						],
					},
				},
			},
			{
				test: /\.css$/i,
				use: ["style-loader", "css-loader"],
			},
			{
				test: /\.(png|svg|jpg|jpeg|gif)$/i,
				type: "asset/resource",
			},
		],
	},
	resolve: {
		extensions: [".js", ".jsx", ".ts", ".tsx"], // Add support for resolving .js and .mjs extensions
	},
	devServer: {
		https: {
			ca: "./ssl/server.pem",
			key: "./ssl/server.key",
			cert: "./ssl/server.crt",
			passphrase: "webpack-dev-server",
			requestCert: false,
		},
		server: "https",
		static: path.join(__dirname, "./src"),
		historyApiFallback: true,
		port: 3000,
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: path.resolve(__dirname, "src/index.html"), // Path to your HTML template
			filename: "index.html", // Output HTML file name
		}),
	],
}
