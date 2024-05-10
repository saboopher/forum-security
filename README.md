## Web Forum Project

This project involves creating a web forum where users can communicate with each other, associate categories with posts, like and dislike posts and comments, and filter posts. SQLite is used as the database library to store data such as users, posts, and comments.

### Features

1. **Authentication**
   - Users can register with their email, username, and password.
   - Registration checks for duplicate emails and ensures password encryption.
   - Registered users can log in with their credentials.
   - Session management is implemented using cookies with expiration dates.

2. **Communication**
   - Registered users can create posts and comments.
   - Posts can be associated with one or more categories.
   - Posts and comments are visible to all users.
   - Non-registered users can view posts and comments.

3. **Likes and Dislikes**
   - Registered users can like or dislike posts and comments.
   - The number of likes and dislikes is visible to all users.

4. **Filtering**
   - Users can filter posts by categories, created posts, and liked posts.
   - Filtering by categories allows users to view posts in specific subforums.

5. **Docker**
   - Docker is used for containerizing the application.
   - Docker basics are implemented for image creation and application containerization.

### Technologies and Packages Used

- Go
- SQLite
- bcrypt (for password encryption)
- UUID
- Docker

### Setup Instructions

1. Clone the repository to your local machine.
2. Ensure you have Go and Docker installed.
3. Navigate to the project root directory.
4. Docker way
   1. Run `docker build -t forum-app .` to build the Docker image.
   2. Run `docker run -p 8989:8989 forum-app` to start the application.
5. Traditional way
   1. `go run ./cmd/api`
6. Access the forum at `http://localhost:8989` in your web browser.


### Usage

- Create categories in the forum to organize posts.
- Register as a user to participate in discussions and interactions.
- Create posts and associate them with relevant categories.
- Like or dislike posts and comments to express opinions.
- Use the filter mechanism to navigate and find relevant content.

