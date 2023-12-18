module.exports = {
	globDirectory: '/home/mrhumble/Documents/projects/gamifiedev',
	globPatterns: [
		'**/*.{css,json,md,key,png,jpeg,html,js,webmanifest}'
	],
	swDest: '/home/mrhumble/Documents/projects/gamifiedev/sw.js',
	ignoreURLParametersMatching: [
		/^utm_/,
		/^fbclid$/,
		/^version/
	]
};