package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/foodmngtapp/food-management-apps/src/config"
	"github.com/foodmngtapp/food-management-apps/src/config/database"
	"github.com/foodmngtapp/food-management-apps/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	TOTAL_USERS = 10000
	BATCH_SIZE  = 100
)

// Sample data for random generation
var firstNames = []string{
	"John", "Jane", "Mike", "Sarah", "David", "Emily", "Robert", "Lisa", "James", "Maria",
	"William", "Jennifer", "Richard", "Linda", "Charles", "Elizabeth", "Thomas", "Barbara",
	"Christopher", "Jessica", "Daniel", "Susan", "Matthew", "Karen", "Anthony", "Nancy",
	"Mark", "Lisa", "Donald", "Betty", "Steven", "Helen", "Paul", "Sandra", "Andrew",
	"Donna", "Joshua", "Carol", "Kenneth", "Ruth", "Kevin", "Sharon", "Brian", "Michelle",
	"George", "Laura", "Timothy", "Sarah", "Ronald", "Kimberly", "Jason", "Deborah",
	"Edward", "Dorothy", "Jeffrey", "Lisa", "Ryan", "Nancy", "Jacob", "Karen", "Gary",
	"Betty", "Nicholas", "Helen", "Eric", "Sandra", "Jonathan", "Donna", "Stephen", "Carol",
	"Larry", "Ruth", "Justin", "Sharon", "Scott", "Michelle", "Brandon", "Laura", "Benjamin",
	"Sarah", "Samuel", "Kimberly", "Gregory", "Deborah", "Alexander", "Dorothy", "Frank",
	"Lisa", "Raymond", "Nancy", "Jack", "Karen", "Dennis", "Betty", "Jerry", "Helen",
	"Tyler", "Sandra", "Aaron", "Donna", "Jose", "Carol", "Henry", "Ruth", "Adam", "Sharon",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez",
	"Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor",
	"Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White", "Harris", "Sanchez",
	"Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King", "Wright",
	"Scott", "Torres", "Nguyen", "Hill", "Flores", "Green", "Adams", "Nelson", "Baker",
	"Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts", "Gomez", "Phillips",
	"Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart",
	"Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan",
	"Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox",
	"Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett",
	"Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders",
	"Patel", "Myers", "Long", "Ross", "Foster", "Jimenez", "Powell", "Jenkins", "Perry",
}

var cities = []struct {
	Name      string
	State     string
	Latitude  string
	Longitude string
}{
	{"New York", "NY", "40.7128", "-74.0060"},
	{"Los Angeles", "CA", "34.0522", "-118.2437"},
	{"Chicago", "IL", "41.8781", "-87.6298"},
	{"Houston", "TX", "29.7604", "-95.3698"},
	{"Phoenix", "AZ", "33.4484", "-112.0740"},
	{"Philadelphia", "PA", "39.9526", "-75.1652"},
	{"San Antonio", "TX", "29.4241", "-98.4936"},
	{"San Diego", "CA", "32.7157", "-117.1611"},
	{"Dallas", "TX", "32.7767", "-96.7970"},
	{"San Jose", "CA", "37.3382", "-121.8863"},
	{"Austin", "TX", "30.2672", "-97.7431"},
	{"Jacksonville", "FL", "30.3322", "-81.6557"},
	{"Fort Worth", "TX", "32.7555", "-97.3308"},
	{"Columbus", "OH", "39.9612", "-82.9988"},
	{"San Francisco", "CA", "37.7749", "-122.4194"},
	{"Charlotte", "NC", "35.2271", "-80.8431"},
	{"Indianapolis", "IN", "39.7684", "-86.1581"},
	{"Seattle", "WA", "47.6062", "-122.3321"},
	{"Denver", "CO", "39.7392", "-104.9903"},
	{"Boston", "MA", "42.3601", "-71.0589"},
	{"El Paso", "TX", "31.7619", "-106.4850"},
	{"Detroit", "MI", "42.3314", "-83.0458"},
	{"Nashville", "TN", "36.1627", "-86.7816"},
	{"Portland", "OR", "45.5152", "-122.6784"},
	{"Memphis", "TN", "35.1495", "-90.0490"},
	{"Oklahoma City", "OK", "35.4676", "-97.5164"},
	{"Las Vegas", "NV", "36.1699", "-115.1398"},
	{"Louisville", "KY", "38.2527", "-85.7585"},
	{"Baltimore", "MD", "39.2904", "-76.6122"},
	{"Milwaukee", "WI", "43.0389", "-87.9065"},
}

var streets = []string{
	"Main St", "Oak Ave", "Pine St", "Elm St", "Cedar Ave", "Maple St", "Park Ave", "First St",
	"Second St", "Third St", "Fourth St", "Fifth St", "Washington St", "Lincoln Ave", "Jefferson St",
	"Madison Ave", "Jackson St", "Monroe St", "Adams St", "Wilson Ave", "Johnson St", "Brown St",
	"Davis Ave", "Miller St", "Wilson St", "Moore Ave", "Taylor St", "Anderson Ave", "Thomas St",
	"Jackson Ave", "White St", "Harris Ave", "Martin St", "Thompson Ave", "Garcia St", "Martinez Ave",
}

var localities = []string{
	"Downtown", "Midtown", "Uptown", "Old Town", "Historic District", "Business District",
	"Shopping Center", "Residential Area", "Suburbs", "City Center", "Financial District",
	"Arts District", "University Area", "Industrial District", "Waterfront", "Heights",
	"Village", "Plaza", "Gardens", "Hills", "Valley", "Ridge", "Park", "Square", "Commons",
}

var profileTypes = []models.ProfileType{
	models.UserType,
	models.RiderType,
	models.RestaurantType,
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg, err := config.LoadConfig("../configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	client, err := database.InitializeDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userCollection := database.OpenCollection(client, "users")

	// Hash the common password
	password := "Test@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	hashedPasswordStr := string(hashedPassword)

	fmt.Printf("🚀 Starting batch insert of %d users with batch size %d\n", TOTAL_USERS, BATCH_SIZE)
	fmt.Println("📊 This will create batches of:", TOTAL_USERS/BATCH_SIZE)

	totalBatches := (TOTAL_USERS + BATCH_SIZE - 1) / BATCH_SIZE
	var totalInserted int
	var profileCounts = make(map[models.ProfileType]int)

	startTime := time.Now()

	for batchNum := 0; batchNum < totalBatches; batchNum++ {
		batchStartTime := time.Now()
		
		// Calculate how many users to create in this batch
		remainingUsers := TOTAL_USERS - (batchNum * BATCH_SIZE)
		currentBatchSize := BATCH_SIZE
		if remainingUsers < BATCH_SIZE {
			currentBatchSize = remainingUsers
		}

		// Generate batch of users
		batch := make([]interface{}, currentBatchSize)
		for i := 0; i < currentBatchSize; i++ {
			user := generateRandomUser(hashedPasswordStr)
			batch[i] = user
			profileCounts[*user.ProfileType]++
		}

		// Insert batch
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		
		result, err := userCollection.InsertMany(ctx, batch)
		cancel()

		if err != nil {
			log.Printf("❌ Failed to insert batch %d: %v", batchNum+1, err)
			continue
		}

		totalInserted += len(result.InsertedIDs)
		batchDuration := time.Since(batchStartTime)
		
		// Progress reporting
		progress := float64(batchNum+1) / float64(totalBatches) * 100
		fmt.Printf("✅ Batch %d/%d completed - Inserted %d users (%.1f%%) - Duration: %v\n", 
			batchNum+1, totalBatches, len(result.InsertedIDs), progress, batchDuration)

		// Optional: Add a small delay between batches to reduce database load
		if batchNum < totalBatches-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	totalDuration := time.Since(startTime)

	// Final report
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🎉 BATCH INSERT COMPLETED!\n")
	fmt.Printf("📊 Total users inserted: %d/%d\n", totalInserted, TOTAL_USERS)
	fmt.Printf("⏱️  Total time: %v\n", totalDuration)
	fmt.Printf("🚀 Average rate: %.0f users/second\n", float64(totalInserted)/totalDuration.Seconds())
	fmt.Printf("📦 Batch size: %d users per batch\n", BATCH_SIZE)
	fmt.Printf("🔐 Password for all users: %s\n", password)
	
	fmt.Println("\n📈 Profile Type Distribution:")
	for profileType, count := range profileCounts {
		percentage := float64(count) / float64(totalInserted) * 100
		fmt.Printf("- %s: %d users (%.1f%%)\n", profileType, count, percentage)
	}

	fmt.Println("\n📍 Users distributed across", len(cities), "major US cities")
	fmt.Println("✨ All users have realistic names, locations, and timestamps")
}

func generateRandomUser(hashedPassword string) models.User {
	now := time.Now()
	
	// Random profile type (exclude admin for bulk users)
	profileType := profileTypes[rand.Intn(len(profileTypes))]
	
	// Random name
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	
	// Generate email based on profile type
	var emailDomain string
	switch profileType {
	case models.UserType:
		emailDomain = "@example.com"
	case models.RiderType:
		emailDomain = "@rider.com"
	case models.RestaurantType:
		emailDomain = "@restaurant.com"
	}
	
	email := fmt.Sprintf("%s.%s%d%s", 
		toLowerCase(firstName), 
		toLowerCase(lastName), 
		rand.Intn(1000), 
		emailDomain)
	
	// Random city
	city := cities[rand.Intn(len(cities))]
	
	// Random street
	streetNumber := rand.Intn(9999) + 1
	streetName := streets[rand.Intn(len(streets))]
	street := fmt.Sprintf("%d %s", streetNumber, streetName)
	
	// Random locality
	locality := localities[rand.Intn(len(localities))]
	
	// Random phone number
	phone := fmt.Sprintf("+1%d%d%d%d%d%d%d%d%d%d", 
		rand.Intn(9)+1, rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
	
	// Random boolean values with realistic distributions
	isVerify := rand.Float32() < 0.85     // 85% verified
	isComplete := rand.Float32() < 0.90   // 90% complete
	
	// Random last active time (within last week)
	maxHoursAgo := 24 * 7 // 1 week
	hoursAgo := rand.Intn(maxHoursAgo)
	lastActive := now.Add(-time.Duration(hoursAgo) * time.Hour)
	
	// Slightly vary coordinates within city bounds
	latVariation := (rand.Float64() - 0.5) * 0.1  // ±0.05 degrees
	lonVariation := (rand.Float64() - 0.5) * 0.1  // ±0.05 degrees
	
	latitude := fmt.Sprintf("%.4f", parseFloat(city.Latitude)+latVariation)
	longitude := fmt.Sprintf("%.4f", parseFloat(city.Longitude)+lonVariation)

	userID := primitive.NewObjectID()
	
	return models.User{
		ID:           userID,
		UserID:       userID.Hex(),
		FirstName:    &firstName,
		LastName:     &lastName,
		Email:        &email,
		Password:     &hashedPassword,
		ProfileType:  &profileType,
		MobileNumber: &phone,
		IsAdmin:      false,
		IsVerify:     isVerify,
		IsComplete:   isComplete,
		CreatedAt:    now.Add(-time.Duration(rand.Intn(30*24)) * time.Hour), // Created within last 30 days
		UpdatedAt:    now,
		LastActive:   lastActive,
		Location: models.Location{
			Street:          street,
			LocalityVerbose: locality,
			City:            city.Name,
			State:           city.State,
			Coordinates: struct {
				Latitude  string `bson:"latitude" json:"latitude"`
				Longitude string `bson:"longitude" json:"longitude"`
			}{
				Latitude:  latitude,
				Longitude: longitude,
			},
		},
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func profileTypePtr(pt models.ProfileType) *models.ProfileType {
	return &pt
}

func toLowerCase(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func parseFloat(s string) float64 {
	var result float64
	var sign float64 = 1
	var decimal bool
	var decimalPlace float64 = 0.1
	
	for i, r := range s {
		if i == 0 && r == '-' {
			sign = -1
			continue
		}
		if r == '.' {
			decimal = true
			continue
		}
		if r >= '0' && r <= '9' {
			digit := float64(r - '0')
			if decimal {
				result += digit * decimalPlace
				decimalPlace *= 0.1
			} else {
				result = result*10 + digit
			}
		}
	}
	return result * sign
}
