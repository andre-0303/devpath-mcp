package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/andre-0303/devpath-mcp/internal/models"
	"github.com/andre-0303/devpath-mcp/internal/roadmap"
)

// Service holds business logic. No I/O — saveFn closure handles persistence.
type Service struct {
	progress *models.UserProgress
	saveFn   func(*models.UserProgress) error
}

// New creates a Service with the loaded progress and a save function.
func New(progress *models.UserProgress, saveFn func(*models.UserProgress) error) *Service {
	return &Service{progress: progress, saveFn: saveFn}
}

// GetTodayTopic returns the ideal topic to study today given a language and available minutes.
func (s *Service) GetTodayTopic(language string, timeAvailable int) (*models.StudyPlanResponse, error) {
	rm := roadmap.Get(language)
	if rm == nil {
		return nil, fmt.Errorf("language %q not found. Supported: %s", language, strings.Join(roadmap.SupportedLanguages(), ", "))
	}

	completed := s.completedFor(language)

	// Find first uncompleted topic with all prerequisites satisfied.
	var candidate *models.Topic
	for i := range rm {
		t := &rm[i]
		if completed[t.Name] {
			continue
		}
		if !s.prerequisitesMet(t, completed) {
			continue
		}
		candidate = t
		break
	}

	if candidate == nil {
		return nil, fmt.Errorf("congratulations: you have completed all topics in the %s roadmap", language)
	}

	reason := s.buildReason(candidate, completed)
	fitsInTime := candidate.EstimatedMinutes <= timeAvailable

	return &models.StudyPlanResponse{
		Language:      language,
		Topic:         *candidate,
		Reason:        reason,
		TimeAvailable: timeAvailable,
		FitsInTime:    fitsInTime,
	}, nil
}

// GetNextTopic returns the next topic after currentTopic in the roadmap.
func (s *Service) GetNextTopic(language string, currentTopic string) (*models.NextTopicResponse, error) {
	rm := roadmap.Get(language)
	if rm == nil {
		return nil, fmt.Errorf("language %q not found. Supported: %s", language, strings.Join(roadmap.SupportedLanguages(), ", "))
	}

	idx := roadmap.TopicIndex(rm, currentTopic)
	if idx == -1 {
		return nil, fmt.Errorf("topic %q not found in %s roadmap", currentTopic, language)
	}

	// Canonical name for the response.
	canonical := rm[idx].Name

	if idx == len(rm)-1 {
		return &models.NextTopicResponse{
			Language:     language,
			CurrentTopic: canonical,
			IsComplete:   true,
		}, nil
	}

	completed := s.completedFor(language)

	// Scan forward from idx+1 for the first topic with prerequisites met.
	for i := idx + 1; i < len(rm); i++ {
		t := &rm[i]
		if s.prerequisitesMet(t, completed) {
			return &models.NextTopicResponse{
				Language:     language,
				CurrentTopic: canonical,
				NextTopic:    t,
			}, nil
		}
	}

	// All remaining topics have unmet prerequisites — return the immediate next anyway.
	next := &rm[idx+1]
	return &models.NextTopicResponse{
		Language:     language,
		CurrentTopic: canonical,
		NextTopic:    next,
	}, nil
}

// GeneratePractice returns a daily-rotated practice challenge for the given topic.
func (s *Service) GeneratePractice(topic string, language string) (*models.PracticeChallengeResponse, error) {
	if language == "" {
		language = "golang"
	}
	rm := roadmap.Get(language)
	if rm == nil {
		return nil, fmt.Errorf("language %q not found. Supported: %s", language, strings.Join(roadmap.SupportedLanguages(), ", "))
	}

	t := roadmap.TopicByName(rm, topic)
	if t == nil {
		return nil, fmt.Errorf("topic %q not found in %s roadmap", topic, language)
	}

	if len(t.PracticeChallenges) == 0 {
		return nil, fmt.Errorf("no practice challenges defined for topic %q", t.Name)
	}

	idx := time.Now().YearDay() % len(t.PracticeChallenges)
	challenge := t.PracticeChallenges[idx]
	hint := ""
	if idx < len(t.Hints) {
		hint = t.Hints[idx]
	}

	return &models.PracticeChallengeResponse{
		Topic:     *t,
		Challenge: challenge,
		Hint:      hint,
	}, nil
}

// MarkTopicComplete marks a topic as completed and persists progress.
func (s *Service) MarkTopicComplete(language string, topicName string) error {
	rm := roadmap.Get(language)
	if rm == nil {
		return fmt.Errorf("language %q not found. Supported: %s", language, strings.Join(roadmap.SupportedLanguages(), ", "))
	}

	t := roadmap.TopicByName(rm, topicName)
	if t == nil {
		return fmt.Errorf("topic %q not found in %s roadmap", topicName, language)
	}

	if s.progress.Completed[language] == nil {
		s.progress.Completed[language] = make(map[string]bool)
	}
	s.progress.Completed[language][t.Name] = true

	if err := s.saveFn(s.progress); err != nil {
		return fmt.Errorf("failed to save progress: %w", err)
	}
	return nil
}

// IsTopicCompleted reports whether a topic is already marked complete.
func (s *Service) IsTopicCompleted(language string, topicName string) bool {
	rm := roadmap.Get(language)
	if rm == nil {
		return false
	}
	t := roadmap.TopicByName(rm, topicName)
	if t == nil {
		return false
	}
	return s.completedFor(language)[t.Name]
}

// ReviewProgress returns a full progress summary for a language.
func (s *Service) ReviewProgress(language string) (*models.ProgressReviewResponse, error) {
	rm := roadmap.Get(language)
	if rm == nil {
		return nil, fmt.Errorf("language %q not found. Supported: %s", language, strings.Join(roadmap.SupportedLanguages(), ", "))
	}

	completed := s.completedFor(language)

	var completedNames []string
	for _, t := range rm {
		if completed[t.Name] {
			completedNames = append(completedNames, t.Name)
		}
	}

	count := len(completedNames)
	total := len(rm)
	pct := 0.0
	if total > 0 {
		pct = float64(count) / float64(total) * 100
	}

	// Next recommended: first uncompleted with prerequisites met.
	var nextRec *models.Topic
	for i := range rm {
		t := &rm[i]
		if !completed[t.Name] && s.prerequisitesMet(t, completed) {
			nextRec = t
			break
		}
	}

	weakAreas := s.computeWeakAreas(rm, completed)

	return &models.ProgressReviewResponse{
		Language:        language,
		TotalTopics:     total,
		CompletedTopics: completedNames,
		CompletedCount:  count,
		PercentDone:     pct,
		NextRecommended: nextRec,
		WeakAreas:       weakAreas,
	}, nil
}

// CompletedCountFor returns the number of completed topics and the total for a language.
func (s *Service) CompletedCountFor(language string) (done, total int) {
	rm := roadmap.Get(language)
	if rm == nil {
		return 0, 0
	}
	completed := s.completedFor(language)
	for _, t := range rm {
		total++
		if completed[t.Name] {
			done++
		}
	}
	return
}

// --- helpers ---

func (s *Service) completedFor(language string) map[string]bool {
	if m, ok := s.progress.Completed[language]; ok {
		return m
	}
	return map[string]bool{}
}

func (s *Service) prerequisitesMet(t *models.Topic, completed map[string]bool) bool {
	for _, req := range t.Prerequisites {
		if !completed[req] {
			return false
		}
	}
	return true
}

func (s *Service) buildReason(t *models.Topic, completed map[string]bool) string {
	if len(t.Prerequisites) == 0 {
		return "This is the first topic in the roadmap."
	}
	done := make([]string, 0, len(t.Prerequisites))
	for _, req := range t.Prerequisites {
		if completed[req] {
			done = append(done, req)
		}
	}
	if len(done) > 0 {
		return fmt.Sprintf("You've completed %s — %s is the natural next step.", strings.Join(done, ", "), t.Name)
	}
	return fmt.Sprintf("%s is next in the %s roadmap.", t.Name, t.Name)
}

func (s *Service) computeWeakAreas(rm []models.Topic, completed map[string]bool) []models.WeakArea {
	// Count how many topics depend on each topic name.
	dependentCount := make(map[string]int, len(rm))
	for _, t := range rm {
		for _, req := range t.Prerequisites {
			dependentCount[req]++
		}
	}

	var weak []models.WeakArea
	for _, t := range rm {
		if completed[t.Name] {
			continue
		}
		if dependentCount[t.Name] >= 2 {
			weak = append(weak, models.WeakArea{
				TopicName:      t.Name,
				DependentCount: dependentCount[t.Name],
			})
		}
	}

	// Sort descending by dependent count (simple insertion sort — 18 elements max).
	for i := 1; i < len(weak); i++ {
		for j := i; j > 0 && weak[j].DependentCount > weak[j-1].DependentCount; j-- {
			weak[j], weak[j-1] = weak[j-1], weak[j]
		}
	}

	if len(weak) > 3 {
		weak = weak[:3]
	}
	return weak
}
