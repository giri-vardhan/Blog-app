# Blog-app

Blog-api is a simple Twitter kind of an app where we have functionalites like create post , fetching post , creating comment & fetching comment by PostID.
It is a backend project developed using gin package of go lang for routing , database - postgres sql, postman for api requests.

we used GET for retrieve data from database.
POST for adding or creating a data in database through api.

Initially we created a connectDB() function with a variable DB pointing to that database. also database constants like host, port, username, dbname and password are saved in main package.

In the main function we initialized our gin router and run it on port "localhost:8080".
Created the functions for above methods(GET,POST) with some routing paths.
  router.GET("/posts", getPost)
	router.POST("/posts", CreatePost)
	router.GET("/comment/:id", getCommentByPostID)
	router.POST("/comment", createComment)
  
 
