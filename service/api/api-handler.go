package api

import "net/http"

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/feed", rt.getFeed)
	rt.router.GET("/db/:table", rt.getDB)
	rt.router.GET("/createuser/:username", rt.CreateUser)
	rt.router.GET("/DESTROYDB/sure", rt.DestroyDB)
	rt.router.POST("/createpost/:ownerID", rt.wrap(rt.CreatePost))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
