package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cozy-software/interview-test/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/* Struct for /users/:id */
type UsersId struct {
	UserId   int               `json:"userid"`
	Name     string            `json:"name"`
	Birthday int               `json:"birthday"`
	Avatar   string            `json:"avatar"`
	Posts    []PostsForUsersId `json:"postsforusersid"`
}

/* Struct for /users/:id */
type PostsForUsersId struct {
	PostId   int    `json:"postid"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	PostDate int    `json:"postdate"`
	NumLikes int    `json:"numlikes"`
}

/* Struct for /posts and /posts/:id */
type Posts struct {
	PostId         int    `json:"postid"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	PostDate       int    `json:"postdate"`
	AuthorId       int    `json:"authorid"`
	AuthorName     string `json:"authorname"`
	AuthorBirthday int    `json:"authorbirthday"`
	AuthorAvatar   string `json:"avatar"`
	NumLikes       int    `json:"numlikes"`
}

/* Struct for optional user query, /posts and /posts/:id */
type PostsLikedByUser struct {
	PostId          int    `json:"postid"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	PostDate        int    `json:"postdate"`
	AuthorId        int    `json:"authorid"`
	AuthorName      string `json:"authorname"`
	AuthorBirthday  int    `json:"authorbirthday"`
	AuthorAvatar    string `json:"avatar"`
	NumLikes        int    `json:"numlikes"`
	PostLikedByUser string `json:"postlikedbyuser"`
}

/* Struct for /posts/:id/likes */
type PostsIdLikes struct {
	LikeDate int    `json:"likedate"`
	UserId   int    `json:"userid"`
	Name     string `json:"name"`
	Birthday int    `json:"birthday"`
	Avatar   string `json:"avatar"`
}

func Mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Subroute for /posts and /users
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", getPosts)
		r.Get("/{id}", getPostsWithID)
		r.Get("/{id}/likes", getPostsWithIDLikes)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", getUsers)
	})

	return r
}

/* API get endpoint for /posts */
func getPosts(w http.ResponseWriter, r *http.Request) {
	// Query for pagination
	page := r.URL.Query().Get("page")
	pageint, err := strconv.Atoi(page)
	checkErr(err)
	limit := r.URL.Query().Get("limit")
	limitint, err := strconv.Atoi(limit)
	checkErr(err)
	offset := (pageint - 1) * limitint

	// Optional query for user
	user := r.URL.Query().Get("user")

	// Query for posts
	rows, err := database.DB.Query("SELECT posts.id, posts.title, posts.content, posts.post_date, posts.author_id, users.name, users.birthday, users.avatar, COUNT(likes.user_id) FROM posts INNER JOIN users ON posts.author_id = users.id LEFT JOIN likes ON posts.id = likes.post_id GROUP BY posts.id, users.id LIMIT $1 OFFSET $2", limitint, offset)
	checkErr(err)

	// Slice of posts
	allPosts := []Posts{}
	allPostsLikedByUser := []PostsLikedByUser{}

	for rows.Next() {
		var aPost Posts

		err := rows.Scan(
			&aPost.PostId,
			&aPost.Title,
			&aPost.Content,
			&aPost.PostDate,
			&aPost.AuthorId,
			&aPost.AuthorName,
			&aPost.AuthorBirthday,
			&aPost.AuthorAvatar,
			&aPost.NumLikes,
		)

		checkErr(err)

		if user != "" {
			userInt, err := strconv.Atoi(user)
			checkErr(err)
			liked, err := database.DB.Query("SELECT COUNT(likes.user_id) FROM likes WHERE $1 = likes.user_id AND $2 = likes.post_id;", userInt, aPost.PostId)
			checkErr(err)
			var postLiked int
			if liked.Next() {
				err := liked.Scan(&postLiked)
				checkErr(err)
			}
			var likedYesOrNo string
			if postLiked > 0 {
				likedYesOrNo = "Yes"
			} else {
				likedYesOrNo = "No"
			}

			var WithLikePost PostsLikedByUser
			WithLikePost.PostId = aPost.PostId
			WithLikePost.Title = aPost.Title
			WithLikePost.Content = aPost.Content
			WithLikePost.PostDate = aPost.PostDate
			WithLikePost.AuthorId = aPost.AuthorId
			WithLikePost.AuthorName = aPost.AuthorName
			WithLikePost.AuthorBirthday = aPost.AuthorBirthday
			WithLikePost.AuthorAvatar = aPost.AuthorAvatar
			WithLikePost.NumLikes = aPost.NumLikes
			WithLikePost.PostLikedByUser = likedYesOrNo
			allPostsLikedByUser = append(allPostsLikedByUser, WithLikePost)
		}

		allPosts = append(allPosts, aPost)
	}
	if user != "" {
		json.NewEncoder(w).Encode(allPostsLikedByUser)
	} else {
		json.NewEncoder(w).Encode(allPosts)
	}
}

/* API get endpoint for /posts/:id */
func getPostsWithID(w http.ResponseWriter, r *http.Request) {
	// Optional query for user
	user := r.URL.Query().Get("user")

	// Parameter for PostID
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	checkErr(err)

	// Query for post
	row, err := database.DB.Query("SELECT posts.id, posts.title, posts.content, posts.post_date, posts.author_id, users.name, users.birthday, users.avatar, COUNT(likes.user_id) FROM posts INNER JOIN users ON posts.author_id = users.id LEFT JOIN likes ON posts.id = likes.post_id WHERE posts.id = $1 GROUP BY posts.id, users.id", idInt)
	checkErr(err)
	onePost := []Posts{}
	onePostLikedByUser := []PostsLikedByUser{}

	for row.Next() {
		var aPost Posts

		err := row.Scan(
			&aPost.PostId,
			&aPost.Title,
			&aPost.Content,
			&aPost.PostDate,
			&aPost.AuthorId,
			&aPost.AuthorName,
			&aPost.AuthorBirthday,
			&aPost.AuthorAvatar,
			&aPost.NumLikes,
		)

		checkErr(err)

		if user != "" {
			userInt, err := strconv.Atoi(user)
			checkErr(err)
			liked, err := database.DB.Query("SELECT COUNT(likes.user_id) FROM likes WHERE $1 = likes.user_id AND $2 = likes.post_id;", userInt, &aPost.PostId)
			checkErr(err)
			var postLiked int
			if liked.Next() {
				err := liked.Scan(&postLiked)
				checkErr(err)
			}
			var likedYesOrNo string
			if postLiked > 0 {
				likedYesOrNo = "Yes"
			} else {
				likedYesOrNo = "No"
			}

			var WithLikePost PostsLikedByUser
			WithLikePost.PostId = aPost.PostId
			WithLikePost.Title = aPost.Title
			WithLikePost.Content = aPost.Content
			WithLikePost.PostDate = aPost.PostDate
			WithLikePost.AuthorId = aPost.AuthorId
			WithLikePost.AuthorName = aPost.AuthorName
			WithLikePost.AuthorBirthday = aPost.AuthorBirthday
			WithLikePost.AuthorAvatar = aPost.AuthorAvatar
			WithLikePost.NumLikes = aPost.NumLikes
			WithLikePost.PostLikedByUser = likedYesOrNo
			onePostLikedByUser = append(onePostLikedByUser, WithLikePost)
		}

		onePost = append(onePost, aPost)
	}
	if user != "" {
		json.NewEncoder(w).Encode(onePostLikedByUser)
	} else {
		json.NewEncoder(w).Encode(onePost)
	}
}

/* API get endpoint for /posts/:id/likes */
func getPostsWithIDLikes(w http.ResponseWriter, r *http.Request) {
	// Query for pagination
	page := r.URL.Query().Get("page")
	pageint, err := strconv.Atoi(page)
	checkErr(err)
	limit := r.URL.Query().Get("limit")
	limitint, err := strconv.Atoi(limit)
	checkErr(err)
	offset := (pageint - 1) * limitint

	// Parameter for Post ID
	postId := chi.URLParam(r, "id")
	postIdInt, err := strconv.Atoi(postId)
	checkErr(err)

	// Query for posts
	rows, err := database.DB.Query("SELECT likes.like_date, users.id, users.name, users.birthday, users.avatar FROM likes INNER JOIN users ON likes.user_id = users.id WHERE likes.post_id = $1 LIMIT $2 OFFSET $3", postIdInt, limitint, offset)
	checkErr(err)
	ListOfUsers := []PostsIdLikes{}

	for rows.Next() {
		var aUser PostsIdLikes

		err := rows.Scan(
			&aUser.LikeDate,
			&aUser.UserId,
			&aUser.Name,
			&aUser.Birthday,
			&aUser.Avatar,
		)
		checkErr(err)
		ListOfUsers = append(ListOfUsers, aUser)
	}
	json.NewEncoder(w).Encode(ListOfUsers)
}

/* API get endpoint for /users/:id */
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Parameter for UserID
	userId := chi.URLParam(r, "id")
	userIdInt, err := strconv.Atoi(userId)
	checkErr(err)

	// Query for desired user
	rows, err := database.DB.Query("SELECT users.id, users.name, users.birthday, users.avatar FROM users WHERE users.id = $1", userIdInt)
	checkErr(err)
	// Query for 5 latest posts ordered by post_date in descending order
	postsRows, err := database.DB.Query("SELECT posts.id, posts.title, posts.content, posts.post_date, COUNT(likes.user_id) FROM posts LEFT JOIN likes ON posts.id = likes.post_id WHERE posts.author_id = $1 GROUP BY posts.id ORDER BY post_date DESC LIMIT 5", userIdInt)
	checkErr(err)

	oneUser := []UsersId{}
	postsSlice := []PostsForUsersId{}

	for rows.Next() {
		var aUser UsersId

		err := rows.Scan(
			&aUser.UserId,
			&aUser.Name,
			&aUser.Birthday,
			&aUser.Avatar,
		)
		checkErr(err)

		for postsRows.Next() {
			var aPost PostsForUsersId

			err := postsRows.Scan(
				&aPost.PostId,
				&aPost.Title,
				&aPost.Content,
				&aPost.PostDate,
				&aPost.NumLikes,
			)

			checkErr(err)

			postsSlice = append(postsSlice, aPost)
		}
		aUser.Posts = postsSlice
		oneUser = append(oneUser, aUser)
	}
	json.NewEncoder(w).Encode(oneUser)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
