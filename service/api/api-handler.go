package api

import "net/http"

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/feed", rt.getFeed)

	rt.router.GET("/db/:table", rt.wrap(rt.getDB))
	rt.router.GET("/DESTROYDB/sure", rt.wrap(rt.DestroyDB))

	rt.router.GET("/createuser", rt.wrap(rt.CreateUser))
	rt.router.GET("/banuser", rt.wrap(rt.BanUser))
	rt.router.GET("/followuser", rt.wrap(rt.FollowUser))

	rt.router.POST("/createpost", rt.wrap(rt.CreatePost))
	rt.router.POST("/likepost", rt.wrap(rt.LikePost))
	rt.router.POST("/commentpost", rt.wrap(rt.CommentPost))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
