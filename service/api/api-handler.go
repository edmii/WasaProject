package api

import "net/http"

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/feed", rt.getFeed)

	rt.router.GET("/db/:table", rt.wrap(rt.getDB))
	rt.router.GET("/DESTROYDB/sure", rt.wrap(rt.DestroyDB))

	rt.router.GET("/createuser/:username", rt.wrap(rt.CreateUser))
	rt.router.GET("/banuser/:ownerID/:prayID", rt.wrap(rt.BanUser))
	rt.router.GET("/unbanuser/:ownerID/:prayID", rt.wrap(rt.UnbanUser))
	rt.router.GET("/followuser/:ownerID/:followedID", rt.wrap(rt.Followuser))

	rt.router.POST("/createpost/:ownerID", rt.wrap(rt.CreatePost))
	rt.router.POST("/likepost/:postID/:ownerID", rt.wrap(rt.LikePost))
	rt.router.POST("/commentpost/:postID/:ownerID/:content", rt.wrap(rt.CommentPost))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
