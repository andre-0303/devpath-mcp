package models

// Topic is a single learning step in a roadmap.
type Topic struct {
	Name               string   `json:"name"`
	Difficulty         string   `json:"difficulty"` // beginner | intermediate | advanced
	EstimatedMinutes   int      `json:"estimated_minutes"`
	Prerequisites      []string `json:"prerequisites"`
	PracticeChallenges []string `json:"practice_challenges"`
	Hints              []string `json:"hints"`
}

// UserProgress is persisted to data/progress.json.
type UserProgress struct {
	Completed map[string]map[string]bool `json:"completed"`
}

// StudyPlanResponse is returned by GetTodayTopic.
type StudyPlanResponse struct {
	Language      string
	Topic         Topic
	Reason        string
	TimeAvailable int
	FitsInTime    bool
}

// NextTopicResponse is returned by GetNextTopic.
type NextTopicResponse struct {
	Language     string
	CurrentTopic string
	NextTopic    *Topic
	IsComplete   bool
}

// PracticeChallengeResponse is returned by GeneratePractice.
type PracticeChallengeResponse struct {
	Topic     Topic
	Challenge string
	Hint      string
}

// ProgressReviewResponse is returned by ReviewProgress.
type ProgressReviewResponse struct {
	Language        string
	TotalTopics     int
	CompletedTopics []string
	CompletedCount  int
	PercentDone     float64
	NextRecommended *Topic
	WeakAreas       []WeakArea
}

// WeakArea is an uncompleted topic that blocks multiple downstream topics.
type WeakArea struct {
	TopicName      string
	DependentCount int
}
