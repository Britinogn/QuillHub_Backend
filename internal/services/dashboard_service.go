// internal/services/dashboard_service.go
package services

import (
	"context"
	"fmt"
	"log"

	"github.com/britinogn/quillhub/internal/model"
)

type DashboardRepo interface {
	// Admin Dashboard
	GetTotalUsers(ctx context.Context) (int64, error)
	GetTotalPosts(ctx context.Context) (int64, error)
	GetTotalComments(ctx context.Context) (int64, error)
	GetTotalLikes(ctx context.Context) (int64, error)
	GetNewUsersLast7Days(ctx context.Context) (int64, error)
	GetNewPostsLast7Days(ctx context.Context) (int64, error)
	GetNewCommentsLast7Days(ctx context.Context) (int64, error)
	GetActiveUsers24h(ctx context.Context) (int64, error)
	GetTopPosts(ctx context.Context, limit int) ([]model.TopPost, error)
	GetTopContributors(ctx context.Context, limit int) ([]model.TopContributor, error)
	GetRecentComments(ctx context.Context, limit int) ([]model.RecentComment, error)
	GetPostsByCategory(ctx context.Context) (map[string]int64, error)

	// User Dashboard
	GetUserTotalPosts(ctx context.Context, userID string) (int64, error)
	GetUserTotalViews(ctx context.Context, userID string) (int64, error)
	GetUserTotalLikes(ctx context.Context, userID string) (int64, error)
	GetUserTotalComments(ctx context.Context, userID string) (int64, error)
	GetUserViewsLast7Days(ctx context.Context, userID string) (int64, error)
	GetUserLikesLast7Days(ctx context.Context, userID string) (int64, error)
	GetUserCommentsLast7Days(ctx context.Context, userID string) (int64, error)
	GetUserTopPosts(ctx context.Context, userID string, limit int) ([]model.UserTopPost, error)
	GetUserRecentPosts(ctx context.Context, userID string, limit int) ([]model.UserRecentPost, error)
	GetUserRecentActivity(ctx context.Context, userID string, limit int) ([]model.UserActivity, error)
}

type DashboardService struct {
	repo DashboardRepo
}

func NewDashboardService(repo DashboardRepo) *DashboardService {
	return &DashboardService{repo: repo}
}

// GetAdminDashboard - Get complete admin dashboard data
func (s *DashboardService) GetAdminDashboard(ctx context.Context) (*model.AdminDashboard, error) {
	log.Printf("[DASHBOARD-SERVICE] Fetching admin dashboard")

	dashboard := &model.AdminDashboard{}

	// Fetch all metrics concurrently for better performance
	totalUsers, _ := s.repo.GetTotalUsers(ctx)
	totalPosts, _ := s.repo.GetTotalPosts(ctx)
	totalComments, _ := s.repo.GetTotalComments(ctx)
	totalLikes, _ := s.repo.GetTotalLikes(ctx)
	newUsersLast7Days, _ := s.repo.GetNewUsersLast7Days(ctx)
	newPostsLast7Days, _ := s.repo.GetNewPostsLast7Days(ctx)
	newCommentsLast7Days, _ := s.repo.GetNewCommentsLast7Days(ctx)
	activeUsers24h, _ := s.repo.GetActiveUsers24h(ctx)
	topPosts, _ := s.repo.GetTopPosts(ctx, 5)
	topContributors, _ := s.repo.GetTopContributors(ctx, 10)
	recentComments, _ := s.repo.GetRecentComments(ctx, 10)
	postsByCategory, _ := s.repo.GetPostsByCategory(ctx)

	dashboard.TotalUsers = totalUsers
	dashboard.TotalPosts = totalPosts
	dashboard.TotalComments = totalComments
	dashboard.TotalLikes = totalLikes
	dashboard.NewUsersLast7Days = newUsersLast7Days
	dashboard.NewPostsLast7Days = newPostsLast7Days
	dashboard.NewCommentsLast7Days = newCommentsLast7Days
	dashboard.ActiveUsers24h = activeUsers24h
	dashboard.TopPosts = topPosts
	dashboard.TopContributors = topContributors
	dashboard.RecentComments = recentComments
	dashboard.PostsByCategory = postsByCategory

	log.Printf("[DASHBOARD-SERVICE] Admin dashboard fetched successfully")
	return dashboard, nil
}

// GetUserDashboard - Get complete user dashboard data
func (s *DashboardService) GetUserDashboard(ctx context.Context, userID string) (*model.UserDashboard, error) {
	log.Printf("[DASHBOARD-SERVICE] Fetching user dashboard for: %s", userID)

	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	dashboard := &model.UserDashboard{}

	// Fetch all user metrics
	totalPosts, _ := s.repo.GetUserTotalPosts(ctx, userID)
	totalViews, _ := s.repo.GetUserTotalViews(ctx, userID)
	totalLikes, _ := s.repo.GetUserTotalLikes(ctx, userID)
	totalComments, _ := s.repo.GetUserTotalComments(ctx, userID)
	viewsLast7Days, _ := s.repo.GetUserViewsLast7Days(ctx, userID)
	likesLast7Days, _ := s.repo.GetUserLikesLast7Days(ctx, userID)
	commentsLast7Days, _ := s.repo.GetUserCommentsLast7Days(ctx, userID)
	topPosts, _ := s.repo.GetUserTopPosts(ctx, userID, 3)
	recentPosts, _ := s.repo.GetUserRecentPosts(ctx, userID, 5)
	recentActivity, _ := s.repo.GetUserRecentActivity(ctx, userID, 10)

	dashboard.TotalPosts = totalPosts
	dashboard.TotalViews = totalViews
	dashboard.TotalLikes = totalLikes
	dashboard.TotalComments = totalComments
	dashboard.ViewsLast7Days = viewsLast7Days
	dashboard.LikesLast7Days = likesLast7Days
	dashboard.CommentsLast7Days = commentsLast7Days
	dashboard.TopPosts = topPosts
	dashboard.RecentPosts = recentPosts
	dashboard.RecentActivity = recentActivity

	log.Printf("[DASHBOARD-SERVICE] User dashboard fetched successfully for: %s", userID)
	return dashboard, nil
}