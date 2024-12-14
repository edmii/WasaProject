package api

import "net/http"

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/feed", rt.getFeed)
	rt.router.GET("/profile", rt.getProfile)

	// rt.router.POST("/login", rt.wrap(rt.Login))

	rt.router.GET("/db/:table", rt.wrap(rt.getDB))
	rt.router.GET("/DESTROYDB/sure", rt.wrap(rt.DestroyDB))

	rt.router.POST("/login", rt.wrap(rt.doLogin))
	rt.router.PUT("/changeusername", rt.wrap(rt.setMyUserName))

	rt.router.POST("/banuser", rt.wrap(rt.banUser))
	rt.router.POST("/banuser", rt.wrap(rt.unbanUser))
	rt.router.GET("/getbans", rt.wrap(rt.GetBans))

	rt.router.POST("/followuser", rt.wrap(rt.followUser))
	rt.router.POST("/unfollowuser", rt.wrap(rt.unfollowUser))
	rt.router.GET("/getfollowed", rt.wrap(rt.GetFollowed))
	rt.router.GET("/getfollowers", rt.wrap(rt.GetFollowers))

	rt.router.POST("/createpost", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/deletepost", rt.wrap(rt.DeletePost))
	rt.router.GET("/getuserposts", rt.wrap(rt.GetUserPosts))

	rt.router.POST("/likepost", rt.wrap(rt.LikePost))
	rt.router.GET("/getlikes", rt.wrap(rt.GetLikes))

	rt.router.POST("/commentpost", rt.wrap(rt.CommentPost))
	rt.router.DELETE("/deletecomment", rt.wrap(rt.DeleteComment))
	rt.router.GET("/getcomments", rt.wrap(rt.GetComments))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
