// internal/model/dashboard.go
package model

// AdminDashboard - System-wide metrics for admins
type AdminDashboard struct {
	TotalUsers          int64                `json:"total_users"`
	TotalPosts          int64                `json:"total_posts"`
	TotalComments       int64                `json:"total_comments"`
	TotalLikes          int64                `json:"total_likes"`
	NewUsersLast7Days   int64                `json:"new_users_last_7_days"`
	NewPostsLast7Days   int64                `json:"new_posts_last_7_days"`
	NewCommentsLast7Days int64               `json:"new_comments_last_7_days"`
	ActiveUsers24h      int64                `json:"active_users_24h"`
	TopPosts            []TopPost            `json:"top_posts"`
	TopContributors     []TopContributor     `json:"top_contributors"`
	RecentComments      []RecentComment      `json:"recent_comments"`
	PostsByCategory     map[string]int64     `json:"posts_by_category"`
}

// UserDashboard - Personal stats for individual users
type UserDashboard struct {
	TotalPosts          int64           `json:"total_posts"`
	TotalViews          int64           `json:"total_views"`
	TotalLikes          int64           `json:"total_likes"`
	TotalComments       int64           `json:"total_comments"`
	ViewsLast7Days      int64           `json:"views_last_7_days"`
	LikesLast7Days      int64           `json:"likes_last_7_days"`
	CommentsLast7Days   int64           `json:"comments_last_7_days"`
	TopPosts            []UserTopPost   `json:"top_posts"`
	RecentPosts         []UserRecentPost `json:"recent_posts"`
	RecentActivity      []UserActivity  `json:"recent_activity"`
}

// TopPost - For admin dashboard top posts
type TopPost struct {
	Rank       int    `json:"rank"`
	PostID     string `json:"post_id"`
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
	Views      int64  `json:"views"`
	Likes      int64  `json:"likes"`
	Comments   int64  `json:"comments"`
}

// TopContributor - For admin dashboard top users
type TopContributor struct {
	Rank       int    `json:"rank"`
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	TotalPosts int64  `json:"total_posts"`
	TotalComments int64 `json:"total_comments"`
	TotalLikes int64 `json:"total_likes"`
}

// RecentComment - For admin dashboard recent comments
type RecentComment struct {
	Username  string `json:"username"`
	Comment   string `json:"comment"`
	PostTitle string `json:"post_title"`
	CreatedAt string `json:"created_at"`
}

// UserTopPost - For user dashboard top posts
type UserTopPost struct {
	Rank     int    `json:"rank"`
	PostID   string `json:"post_id"`
	Title    string `json:"title"`
	Views    int64  `json:"views"`
	Likes    int64  `json:"likes"`
	Comments int64  `json:"comments"`
}

// UserRecentPost - For user dashboard recent posts
type UserRecentPost struct {
	PostID      string `json:"post_id"`
	Title       string `json:"title"`
	Status      string `json:"status"` // "published" or "draft"
	Views       int64  `json:"views"`
	UpdatedAt   string `json:"updated_at"`
}

// UserActivity - For user dashboard recent activity
type UserActivity struct {
	Username  string `json:"username"`
	Action    string `json:"action"` // "liked", "commented", "viewed"
	PostTitle string `json:"post_title"`
	TimeAgo   string `json:"time_ago"`
}