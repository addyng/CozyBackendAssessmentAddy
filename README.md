# CozyBackendAssessmentAddy

## Within routes.go contains code for a RESTful API. The code has 4 endpoints with the following possible urls:
---
### /posts
Example URL:

http://localhost:3000/posts

http://localhost:3000/posts?page=1&limit=5

http://localhost:3000/posts?page=1&limit=5&user=10

### /posts/:id
Example URL:

http://localhost:3000/posts/1

http://localhost:3000/posts/1?user=10

### /posts/:id/likes
Example URL:

http://localhost:3000/posts/10/likes

http://localhost:3000/posts/10/likes?page=1&limit=5

### /users/:id
Example URL:

http://localhost:3000/users/11

#### NOTE:

http://localhost:3000/posts?user=10

should work as intended, but it overloads postgres
"pq: sorry, too many clients already"
