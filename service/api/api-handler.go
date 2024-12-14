package api

import "net/http"

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/feed", rt.getMyStream)
	rt.router.GET("/profile", rt.getUserProfile)

	// rt.router.POST("/login", rt.wrap(rt.Login))

	rt.router.GET("/db/:table", rt.wrap(rt.getDB))
	rt.router.GET("/DESTROYDB/sure", rt.wrap(rt.destroyDB))

	rt.router.POST("/login", rt.wrap(rt.doLogin))
	rt.router.PUT("/changeusername", rt.wrap(rt.setMyUserName))

	rt.router.POST("/banuser", rt.wrap(rt.banUser))
	rt.router.POST("/unbanuser", rt.wrap(rt.unbanUser))
	rt.router.GET("/getbans", rt.wrap(rt.getBans))

	rt.router.POST("/followuser", rt.wrap(rt.followUser))
	rt.router.POST("/unfollowuser", rt.wrap(rt.unfollowUser))
	rt.router.GET("/getfollowed", rt.wrap(rt.getFollowed))
	rt.router.GET("/getfollowers", rt.wrap(rt.getFollowers))

	rt.router.POST("/createpost", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/deletepost", rt.wrap(rt.deletePhoto))
	rt.router.GET("/getuserposts", rt.wrap(rt.getUserPosts))

	rt.router.POST("/likepost", rt.wrap(rt.likePhoto))
	rt.router.POST("/unlikepost", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/getlikes", rt.wrap(rt.getLikes))

	rt.router.POST("/commentpost", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/deletecomment", rt.wrap(rt.uncommentPhoto))
	rt.router.GET("/getcomments", rt.wrap(rt.getComments))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
