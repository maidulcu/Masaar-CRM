package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// UAE demo data — realistic contacts and leads for client demos
var users = []struct {
	Name     string
	Email    string
	Password string
	Role     string
	Lang     string
}{
	{"Ahmed Al Mansoori", "ahmed@masaar.local", "Demo@1234", "admin", "ar"},
	{"Sarah Johnson", "sarah@masaar.local", "Demo@1234", "agent", "en"},
	{"Mohammed Al Rashidi", "mohammed@masaar.local", "Demo@1234", "agent", "ar"},
}

var contacts = []struct {
	PhoneWA  string
	FullName string
	Email    string
	Language string
	Score    int
}{
	{"+971501234567", "خالد العامري", "khaled@example.ae", "ar", 85},
	{"+971502345678", "Rania Boutros", "rania@boutros.ae", "en", 72},
	{"+971503456789", "يوسف المهيري", "yousuf@alheri.ae", "ar", 91},
	{"+971504567890", "David Chen", "david@sinobizdubai.com", "en", 60},
	{"+971505678901", "فاطمة الزيدي", "fatima@zaidi.ae", "ar", 45},
	{"+971506789012", "Priya Sharma", "priya@dubaifintech.io", "en", 78},
	{"+971507890123", "عبدالله بن زايد", "abdulla@bz.ae", "ar", 95},
	{"+971508901234", "Lara Hadid", "lara@realestateuae.com", "en", 67},
}

var leadStages = []string{"new", "contacted", "qualified", "proposal", "won", "lost"}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://masaar:masaar@localhost:5432/masaar?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	fmt.Println("🌱 Seeding Masaar CRM...")

	// ── Users ────────────────────────────────────────────────────────────────
	userIDs := make([]uuid.UUID, len(users))
	for i, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		id := uuid.New()
		_, err := pool.Exec(ctx, `
			INSERT INTO users (id, name, email, password_hash, role, lang_pref)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (email) DO NOTHING
		`, id, u.Name, u.Email, string(hash), u.Role, u.Lang)
		if err != nil {
			log.Fatalf("insert user %s: %v", u.Email, err)
		}
		userIDs[i] = id
		fmt.Printf("  ✓ User: %s (%s)\n", u.Name, u.Email)
	}

	// ── Contacts ─────────────────────────────────────────────────────────────
	contactIDs := make([]uuid.UUID, len(contacts))
	for i, c := range contacts {
		id := uuid.New()
		assigned := userIDs[i%len(userIDs)]
		_, err := pool.Exec(ctx, `
			INSERT INTO contacts (id, phone_wa, full_name, email, language, lead_score, assigned_to)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT (phone_wa) DO NOTHING
		`, id, c.PhoneWA, c.FullName, c.Email, c.Language, c.Score, assigned)
		if err != nil {
			log.Fatalf("insert contact %s: %v", c.FullName, err)
		}
		contactIDs[i] = id
		fmt.Printf("  ✓ Contact: %s (%s)\n", c.FullName, c.PhoneWA)
	}

	// ── Leads ─────────────────────────────────────────────────────────────────
	dealValues := []float64{85000, 42000, 210000, 15000, 320000, 78000, 550000, 29000}
	sources := []string{"whatsapp", "web", "referral", "event", "whatsapp", "web", "referral", "whatsapp"}

	for i, cID := range contactIDs {
		stage := leadStages[i%len(leadStages)]
		_, err := pool.Exec(ctx, `
			INSERT INTO leads (id, contact_id, stage, source, deal_value, currency, notes)
			VALUES ($1,$2,$3,$4,$5,'AED',$6)
		`, uuid.New(), cID, stage, sources[i], dealValues[i],
			fmt.Sprintf("Lead from %s — follow up required", sources[i]))
		if err != nil {
			log.Printf("  ⚠ insert lead for contact %s: %v", cID, err)
			continue
		}
		fmt.Printf("  ✓ Lead: stage=%s value=AED %.0f\n", stage, dealValues[i])
	}

	// ── WhatsApp threads + messages ───────────────────────────────────────────
	waMessages := [][]struct{ dir, body string }{
		{
			{"inbound", "السلام عليكم، أبحث عن شقة في دبي مارينا"},
			{"outbound", "وعليكم السلام! يسعدنا مساعدتك. ما هي ميزانيتك؟"},
			{"inbound", "ميزانيتي حوالي 100 ألف درهم سنوياً"},
		},
		{
			{"inbound", "Hi, I'm interested in your office spaces in DIFC"},
			{"outbound", "Hello! Great choice. DIFC is prime. What size are you looking for?"},
			{"inbound", "Around 1500 sqft for a team of 10"},
			{"outbound", "Perfect, I have 3 options ready for you. Shall we schedule a viewing?"},
		},
		{
			{"inbound", "مرحبا، هل لديكم عروض على الوحدات التجارية؟"},
			{"outbound", "أهلاً! نعم لدينا عروض ممتازة. متى يمكنك الحضور للمعاينة؟"},
		},
	}

	for i := 0; i < 3 && i < len(contactIDs); i++ {
		threadID := uuid.New()
		_, err := pool.Exec(ctx, `
			INSERT INTO whatsapp_threads (id, contact_id, wa_account_id, thread_status)
			VALUES ($1,$2,'15551234567','open')
		`, threadID, contactIDs[i])
		if err != nil {
			log.Printf("  ⚠ insert thread: %v", err)
			continue
		}

		for _, msg := range waMessages[i] {
			_, _ = pool.Exec(ctx, `
				INSERT INTO whatsapp_messages (id, thread_id, direction, body, wa_message_id)
				VALUES ($1,$2,$3,$4,$5)
			`, uuid.New(), threadID, msg.dir, msg.body, uuid.New().String())
		}

		_, _ = pool.Exec(ctx, `
			UPDATE whatsapp_threads
			SET last_message_at=NOW(), message_count=$1 WHERE id=$2
		`, len(waMessages[i]), threadID)

		fmt.Printf("  ✓ Thread with %d messages for contact %d\n", len(waMessages[i]), i+1)
	}

	fmt.Println("\n✅ Seed complete!")
	fmt.Println("   Login: ahmed@masaar.local / Demo@1234  (admin, Arabic)")
	fmt.Println("   Login: sarah@masaar.local  / Demo@1234  (agent, English)")
}
