// internal/repository/dashboard_repository.go
package repository

import (
	"context"
	"fmt"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// ==================== ADMIN DASHBOARD ====================

func (r *DashboardRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetTotalPosts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM posts").Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetTotalComments(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM comments").Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetTotalLikes(ctx context.Context) (int64, error) {
	// Assuming you have a likes table
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM likes").Scan(&count)
	if err != nil {
		return 0, nil // Return 0 if likes table doesn't exist yet
	}
	return count, nil
}

func (r *DashboardRepository) GetNewUsersLast7Days(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE created_at >= NOW() - INTERVAL '7 days'`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetNewPostsLast7Days(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM posts WHERE created_at >= NOW() - INTERVAL '7 days'`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetNewCommentsLast7Days(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM comments WHERE created_at >= NOW() - INTERVAL '7 days'`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetActiveUsers24h(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(DISTINCT author_id) FROM posts WHERE created_at >= NOW() - INTERVAL '24 hours'`
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetTopPosts(ctx context.Context, limit int) ([]model.TopPost, error) {
	query := `
		SELECT 
			p.id, p.title, u.name as author_name, p.view_count,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count,
			0 as like_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY p.id, p.title, u.name, p.view_count
		ORDER BY p.view_count DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top posts: %w", err)
	}
	defer rows.Close()

	var posts []model.TopPost
	rank := 1
	for rows.Next() {
		var post model.TopPost
		var postID string
		err := rows.Scan(&postID, &post.Title, &post.AuthorName, &post.Views, &post.Comments, &post.Likes)
		if err != nil {
			return nil, err
		}
		post.Rank = rank
		post.PostID = postID
		posts = append(posts, post)
		rank++
	}

	return posts, nil
}

func (r *DashboardRepository) GetTopContributors(ctx context.Context, limit int) ([]model.TopContributor, error) {
	query := `
		SELECT 
			u.id, u.username,
			COUNT(DISTINCT p.id) as post_count,
			COUNT(DISTINCT c.id) as comment_count,
			0 as like_count
		FROM users u
		LEFT JOIN posts p ON u.id = p.author_id
		LEFT JOIN comments c ON u.id = c.author_id
		GROUP BY u.id, u.username
		ORDER BY post_count DESC, comment_count DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top contributors: %w", err)
	}
	defer rows.Close()

	var contributors []model.TopContributor
	rank := 1
	for rows.Next() {
		var contrib model.TopContributor
		var userID string
		err := rows.Scan(&userID, &contrib.Username, &contrib.TotalPosts, &contrib.TotalComments, &contrib.TotalLikes)
		if err != nil {
			return nil, err
		}
		contrib.Rank = rank
		contrib.UserID = userID
		contributors = append(contributors, contrib)
		rank++
	}

	return contributors, nil
}

func (r *DashboardRepository) GetRecentComments(ctx context.Context, limit int) ([]model.RecentComment, error) {
	query := `
		SELECT 
			u.username, c.text, p.title,
			TO_CHAR(c.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		JOIN posts p ON c.post_id = p.id
		ORDER BY c.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent comments: %w", err)
	}
	defer rows.Close()

	var comments []model.RecentComment
	for rows.Next() {
		var comment model.RecentComment
		err := rows.Scan(&comment.Username, &comment.Comment, &comment.PostTitle, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *DashboardRepository) GetPostsByCategory(ctx context.Context) (map[string]int64, error) {
	query := `
		SELECT 
			COALESCE(category, 'Uncategorized') as category,
			COUNT(*) as count
		FROM posts
		GROUP BY category
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by category: %w", err)
	}
	defer rows.Close()

	categoryMap := make(map[string]int64)
	for rows.Next() {
		var category string
		var count int64
		err := rows.Scan(&category, &count)
		if err != nil {
			return nil, err
		}
		categoryMap[category] = count
	}

	return categoryMap, nil
}

// ==================== USER DASHBOARD ====================

func (r *DashboardRepository) GetUserTotalPosts(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE author_id = $1", userID).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetUserTotalViews(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `SELECT COALESCE(SUM(view_count), 0) FROM posts WHERE author_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetUserTotalLikes(ctx context.Context, userID string) (int64, error) {
	// Placeholder - implement when likes table exists
	return 0, nil
}

func (r *DashboardRepository) GetUserTotalComments(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*) FROM comments c
		JOIN posts p ON c.post_id = p.id
		WHERE p.author_id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetUserViewsLast7Days(ctx context.Context, userID string) (int64, error) {
	// This requires a views tracking table - placeholder for now
	return 0, nil
}

func (r *DashboardRepository) GetUserLikesLast7Days(ctx context.Context, userID string) (int64, error) {
	// Placeholder - implement when likes table exists
	return 0, nil
}

func (r *DashboardRepository) GetUserCommentsLast7Days(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*) FROM comments c
		JOIN posts p ON c.post_id = p.id
		WHERE p.author_id = $1 AND c.created_at >= NOW() - INTERVAL '7 days'
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetUserTopPosts(ctx context.Context, userID string, limit int) ([]model.UserTopPost, error) {
	query := `
		SELECT 
			p.id, p.title, p.view_count,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count,
			0 as like_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.author_id = $1
		GROUP BY p.id, p.title, p.view_count
		ORDER BY p.view_count DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user top posts: %w", err)
	}
	defer rows.Close()

	var posts []model.UserTopPost
	rank := 1
	for rows.Next() {
		var post model.UserTopPost
		var postID string
		err := rows.Scan(&postID, &post.Title, &post.Views, &post.Comments, &post.Likes)
		if err != nil {
			return nil, err
		}
		post.Rank = rank
		post.PostID = postID
		posts = append(posts, post)
		rank++
	}

	return posts, nil
}

func (r *DashboardRepository) GetUserRecentPosts(ctx context.Context, userID string, limit int) ([]model.UserRecentPost, error) {
	query := `
		SELECT 
			id, title, is_published, view_count,
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS') as updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY updated_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recent posts: %w", err)
	}
	defer rows.Close()

	var posts []model.UserRecentPost
	for rows.Next() {
		var post model.UserRecentPost
		var postID string
		var isPublished bool
		err := rows.Scan(&postID, &post.Title, &isPublished, &post.Views, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		post.PostID = postID
		if isPublished {
			post.Status = "published"
		} else {
			post.Status = "draft"
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *DashboardRepository) GetUserRecentActivity(ctx context.Context, userID string, limit int) ([]model.UserActivity, error) {
	query := `
		SELECT 
			u.username,
			'commented' as action,
			p.title as post_title,
			c.created_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		JOIN posts p ON c.post_id = p.id
		WHERE p.author_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recent activity: %w", err)
	}
	defer rows.Close()

	var activities []model.UserActivity
	for rows.Next() {
		var activity model.UserActivity
		var createdAt string
		err := rows.Scan(&activity.Username, &activity.Action, &activity.PostTitle, &createdAt)
		if err != nil {
			return nil, err
		}
		activity.TimeAgo = "Just now" // You can implement time ago calculation
		activities = append(activities, activity)
	}

	return activities, nil
}