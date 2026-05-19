package roadmap

import (
	"strings"

	"github.com/andre-0303/devpath-mcp/internal/models"
)

var GolangRoadmap = []models.Topic{
	{
		Name:             "Variables",
		Difficulty:       "beginner",
		EstimatedMinutes: 20,
		Prerequisites:    []string{},
		PracticeChallenges: []string{
			"Declare variables of 5 different types (int, float64, string, bool, byte) and print each one.",
			"Swap two integer variables without using a temporary variable.",
			"Use short declaration (:=) inside a function and explain why it cannot be used at package level.",
		},
		Hints: []string{
			"Use fmt.Printf with %T to print the type of each variable.",
			"XOR swap: a ^= b; b ^= a; a ^= b — or try the tuple assignment: a, b = b, a.",
			"Short declaration is syntax sugar for var + assignment; the compiler infers the type.",
		},
	},
	{
		Name:             "Constants",
		Difficulty:       "beginner",
		EstimatedMinutes: 15,
		Prerequisites:    []string{"Variables"},
		PracticeChallenges: []string{
			"Create an iota-based enum for days of the week (Sunday=0 through Saturday=6).",
			"Define a typed constant for pi and show that it cannot be assigned at runtime.",
			"Use iota with bit-shifting to create permission flags (Read=1, Write=2, Execute=4).",
		},
		Hints: []string{
			"const ( Sunday = iota; Monday; ... ) — iota resets to 0 at each const block.",
			"const Pi = 3.14159 — try Pi = 3.0 and observe the compile error.",
			"const ( Read = 1 << iota; Write; Execute ) — iota increments by 1, shift doubles the value.",
		},
	},
	{
		Name:             "Functions",
		Difficulty:       "beginner",
		EstimatedMinutes: 30,
		Prerequisites:    []string{"Variables", "Constants"},
		PracticeChallenges: []string{
			"Write a function that returns both the quotient and remainder of integer division (multiple return values).",
			"Implement a variadic function sum(...int) int and call it with 0, 1, and many arguments.",
			"Write a higher-order function that accepts a func(int) int and applies it to each element of a slice.",
		},
		Hints: []string{
			"func divide(a, b int) (int, int) { return a / b, a % b }",
			"func sum(nums ...int) int — inside the function, nums is a []int.",
			"func apply(nums []int, f func(int) int) []int — range over nums, call f on each element.",
		},
	},
	{
		Name:             "Conditionals",
		Difficulty:       "beginner",
		EstimatedMinutes: 20,
		Prerequisites:    []string{"Variables", "Functions"},
		PracticeChallenges: []string{
			"Write a function that classifies a number as negative, zero, positive, or large (>1000) using if/else if/else.",
			"Use a switch statement to map a month number (1-12) to its season (summer, autumn, winter, spring).",
			"Implement FizzBuzz for numbers 1-100 using a switch with no condition (switch { case ... }).",
		},
		Hints: []string{
			"Go if-else requires braces even for single-line bodies. No parentheses around condition.",
			"switch month { case 12, 1, 2: return 'winter' ... } — multiple values per case.",
			"switch { case n%15==0: ... case n%3==0: ... } — evaluated top to bottom, first match wins.",
		},
	},
	{
		Name:             "Loops",
		Difficulty:       "beginner",
		EstimatedMinutes: 25,
		Prerequisites:    []string{"Variables", "Conditionals"},
		PracticeChallenges: []string{
			"Use a for loop to compute the factorial of n (handle n=0 as edge case).",
			"Iterate over a string with range and count how many times each character appears (use a map).",
			"Implement a while-style loop that reads from a channel until it receives a sentinel value (use range on channel).",
		},
		Hints: []string{
			"Go has only for — no while or do-while. for i := 1; i <= n; i++ { result *= i }",
			"for i, ch := range str — ch is a rune (Unicode code point), not a byte.",
			"for val := range ch — exits when channel is closed. Close it with close(ch).",
		},
	},
	{
		Name:             "Arrays",
		Difficulty:       "beginner",
		EstimatedMinutes: 20,
		Prerequisites:    []string{"Variables", "Loops"},
		PracticeChallenges: []string{
			"Declare an array of 5 integers, fill it with squares (1,4,9,16,25), and print it.",
			"Write a function that reverses an array [5]int in place without extra memory.",
			"Find the maximum and minimum values in an array of 10 random integers using a single pass.",
		},
		Hints: []string{
			"var squares [5]int — array size is part of the type; [5]int and [6]int are different types.",
			"Swap arr[i] and arr[len(arr)-1-i] while i < len(arr)/2.",
			"Initialize max = arr[0], min = arr[0], then iterate from index 1.",
		},
	},
	{
		Name:             "Slices",
		Difficulty:       "beginner",
		EstimatedMinutes: 35,
		Prerequisites:    []string{"Arrays"},
		PracticeChallenges: []string{
			"Write a function that removes all duplicate integers from a slice, preserving order.",
			"Implement a stack (push, pop, peek) using a slice as the underlying storage.",
			"Write a function that flattens a [][]int into a single []int.",
		},
		Hints: []string{
			"Use a map[int]bool as a seen-set; append to result only if !seen[v].",
			"Push: stack = append(stack, v). Pop: v, stack = stack[len-1], stack[:len-1]. Peek: stack[len-1].",
			"Pre-allocate with make([]int, 0, totalLen) and append each inner slice with append(result, inner...).",
		},
	},
	{
		Name:             "Maps",
		Difficulty:       "beginner",
		EstimatedMinutes: 30,
		Prerequisites:    []string{"Variables", "Loops"},
		PracticeChallenges: []string{
			"Count word frequencies in a string using a map[string]int.",
			"Invert a map[string]string (swap keys and values) and handle duplicate values.",
			"Implement a simple in-memory cache with get/set/delete operations using a map.",
		},
		Hints: []string{
			"Split string with strings.Fields, then freq[word]++ for each word.",
			"For duplicate values, the last key wins — or collect collisions in a []string.",
			"Always check the ok value: val, ok := cache[key] — ok is false if key doesn't exist.",
		},
	},
	{
		Name:             "Structs",
		Difficulty:       "intermediate",
		EstimatedMinutes: 40,
		Prerequisites:    []string{"Variables", "Functions"},
		PracticeChallenges: []string{
			"Create a BankAccount struct with owner, balance, and accountNumber fields. Add a constructor NewBankAccount.",
			"Implement an embedded struct: Person with name/age, and Employee that embeds Person and adds role/salary.",
			"Write a function that deep-copies a struct containing a slice field (shallow copy won't work).",
		},
		Hints: []string{
			"Use pointer receiver for methods that modify state: func (a *BankAccount) Deposit(amount float64).",
			"type Employee struct { Person; Role string; Salary float64 } — access Person fields directly: e.Name.",
			"Copy the struct value, then make a new slice and copy elements: copy(newSlice, original.Items).",
		},
	},
	{
		Name:             "Methods",
		Difficulty:       "intermediate",
		EstimatedMinutes: 35,
		Prerequisites:    []string{"Structs"},
		PracticeChallenges: []string{
			"Add Deposit, Withdraw, and Balance methods to BankAccount; Withdraw must return an error if balance is insufficient.",
			"Implement the String() string method on a custom type so fmt.Println formats it nicely.",
			"Create a Rectangle and Circle type; add Area() and Perimeter() float64 methods to each.",
		},
		Hints: []string{
			"func (a *BankAccount) Withdraw(amount float64) error — return fmt.Errorf('insufficient funds') when amount > a.balance.",
			"func (b BankAccount) String() string — this makes BankAccount implement fmt.Stringer automatically.",
			"Value receiver for read-only methods is idiomatic when the struct is small and doesn't need mutation.",
		},
	},
	{
		Name:             "Interfaces",
		Difficulty:       "intermediate",
		EstimatedMinutes: 45,
		Prerequisites:    []string{"Methods"},
		PracticeChallenges: []string{
			"Define a Shape interface with Area() and Perimeter() methods; implement it for Rectangle and Circle.",
			"Write a function PrintShapeInfo(s Shape) that prints area and perimeter — call it with both types.",
			"Use a type switch to detect whether an interface value holds a Rectangle, Circle, or unknown type.",
		},
		Hints: []string{
			"Interfaces are satisfied implicitly — no 'implements' keyword needed.",
			"func PrintShapeInfo(s Shape) — pass a Rectangle or Circle value directly; Go handles the dispatch.",
			"switch v := s.(type) { case Rectangle: ... case Circle: ... default: ... }",
		},
	},
	{
		Name:             "Error Handling",
		Difficulty:       "intermediate",
		EstimatedMinutes: 40,
		Prerequisites:    []string{"Functions"},
		PracticeChallenges: []string{
			"Write a SafeDivide(a, b float64) (float64, error) that returns an error for division by zero.",
			"Create a custom error type ValidationError with a Field and Message; implement the error interface.",
			"Wrap errors with fmt.Errorf('context: %w', err) and unwrap them with errors.As and errors.Is.",
		},
		Hints: []string{
			"The error interface has one method: Error() string. Return nil for success, errors.New(...) for failure.",
			"type ValidationError struct { Field, Message string }; func (e *ValidationError) Error() string { ... }",
			"errors.As(err, &target) fills target if err's chain contains a *ValidationError. errors.Is checks identity.",
		},
	},
	{
		Name:             "Packages",
		Difficulty:       "intermediate",
		EstimatedMinutes: 30,
		Prerequisites:    []string{"Functions", "Error Handling"},
		PracticeChallenges: []string{
			"Split your BankAccount code into a separate package 'bank'; export only the types and methods users need.",
			"Create an internal helper function (unexported) inside the package and verify it's not accessible from outside.",
			"Write a package-level init() function that validates required configuration on startup.",
		},
		Hints: []string{
			"Exported = starts with uppercase (BankAccount). Unexported = lowercase (validateAmount).",
			"Unexported identifiers cause a compile error when accessed from another package — no runtime check needed.",
			"func init() runs automatically before main(); use it for one-time setup, not business logic.",
		},
	},
	{
		Name:             "Concurrency",
		Difficulty:       "advanced",
		EstimatedMinutes: 60,
		Prerequisites:    []string{"Functions", "Loops"},
		PracticeChallenges: []string{
			"Launch 5 goroutines, each printing its index; use a WaitGroup to wait for all to finish.",
			"Identify the data race in a counter incremented by multiple goroutines; fix it with sync.Mutex.",
			"Use sync.Once to initialize a shared resource exactly once across multiple goroutines.",
		},
		Hints: []string{
			"var wg sync.WaitGroup; wg.Add(5); go func() { defer wg.Done(); ... }(); wg.Wait()",
			"go run -race detects data races. Fix: var mu sync.Mutex; mu.Lock(); counter++; mu.Unlock()",
			"var once sync.Once; once.Do(func() { resource = expensiveInit() }) — safe for concurrent calls.",
		},
	},
	{
		Name:             "Channels",
		Difficulty:       "advanced",
		EstimatedMinutes: 50,
		Prerequisites:    []string{"Concurrency"},
		PracticeChallenges: []string{
			"Create a pipeline: generator goroutine → square goroutine → printer goroutine, connected by channels.",
			"Implement a fan-out pattern: one input channel, multiple worker goroutines reading from it.",
			"Use select with a timeout channel (time.After) to abort a slow operation after 2 seconds.",
		},
		Hints: []string{
			"Each stage reads from its input channel and sends to its output channel; close output when input is closed.",
			"Multiple goroutines can read from the same channel — Go distributes sends fairly among receivers.",
			"select { case result := <-work: ...; case <-time.After(2*time.Second): return errors.New('timeout') }",
		},
	},
	{
		Name:             "Context",
		Difficulty:       "advanced",
		EstimatedMinutes: 45,
		Prerequisites:    []string{"Channels"},
		PracticeChallenges: []string{
			"Pass a context.WithTimeout to a function that simulates work; cancel it early and observe the behavior.",
			"Propagate a request ID through a call chain using context.WithValue; retrieve it in a deep function.",
			"Write a function that respects ctx.Done() by selecting on it alongside its work channel.",
		},
		Hints: []string{
			"ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()",
			"Use a typed key (not a string) for context values to avoid collisions: type ctxKey string; const requestID ctxKey = 'requestID'",
			"select { case <-ctx.Done(): return ctx.Err(); case result := <-workCh: return result, nil }",
		},
	},
	{
		Name:             "Testing",
		Difficulty:       "intermediate",
		EstimatedMinutes: 40,
		Prerequisites:    []string{"Functions", "Packages"},
		PracticeChallenges: []string{
			"Write table-driven tests for SafeDivide covering: normal case, division by zero, negative numbers.",
			"Use t.Run to group subtests; add t.Parallel() to run them concurrently.",
			"Write a benchmark (BenchmarkSafeDiv) and compare performance with and without the zero check.",
		},
		Hints: []string{
			"tests := []struct{ a, b float64; wantErr bool }{ {10,2,false}, {5,0,true} }; for _, tt := range tests { t.Run(...) }",
			"t.Parallel() must be called at the start of each subtest; the outer test must NOT call t.Parallel().",
			"func BenchmarkSafeDiv(b *testing.B) { for i := 0; i < b.N; i++ { SafeDivide(10, 3) } }; go test -bench=.",
		},
	},
	{
		Name:             "REST API",
		Difficulty:       "advanced",
		EstimatedMinutes: 90,
		Prerequisites:    []string{"Structs", "Error Handling", "Packages", "Context"},
		PracticeChallenges: []string{
			"Build a CRUD API for a Todo resource using net/http (no framework): POST /todos, GET /todos, GET /todos/{id}, DELETE /todos/{id}.",
			"Add JSON request validation: return 400 Bad Request with a structured error body when required fields are missing.",
			"Implement a middleware chain: logging middleware wraps auth middleware wraps your handler; compose them manually.",
		},
		Hints: []string{
			"Use http.NewServeMux() and mux.HandleFunc('/todos/', handler). Parse the ID from r.URL.Path manually or use strings.TrimPrefix.",
			"json.NewDecoder(r.Body).Decode(&req) returns an error if body is malformed; check required fields before processing.",
			"type Middleware func(http.Handler) http.Handler — wrap handlers: handler = logging(auth(realHandler)). Call next.ServeHTTP(w, r).",
		},
	},
}

// Registry maps language names to their roadmap slices.
var Registry = map[string][]models.Topic{
	"golang": GolangRoadmap,
}

// Get returns the roadmap for a language (case-insensitive). Returns nil if not found.
func Get(language string) []models.Topic {
	roadmap, ok := Registry[strings.ToLower(language)]
	if !ok {
		return nil
	}
	return roadmap
}

// TopicByName finds a topic by name (case-insensitive). Returns nil if not found.
func TopicByName(roadmap []models.Topic, name string) *models.Topic {
	lower := strings.ToLower(name)
	for i := range roadmap {
		if strings.ToLower(roadmap[i].Name) == lower {
			return &roadmap[i]
		}
	}
	return nil
}

// TopicIndex returns the index of a topic by name (case-insensitive). Returns -1 if not found.
func TopicIndex(roadmap []models.Topic, name string) int {
	lower := strings.ToLower(name)
	for i, t := range roadmap {
		if strings.ToLower(t.Name) == lower {
			return i
		}
	}
	return -1
}

// SupportedLanguages returns a sorted list of supported language names.
func SupportedLanguages() []string {
	langs := make([]string, 0, len(Registry))
	for k := range Registry {
		langs = append(langs, k)
	}
	return langs
}
