package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"math/rand"

	"github.com/DenysBahachuk/gopher_social/internal/store"
)

var names = []string{
	"Aiden", "Liam", "Emma", "Noah", "Olivia",
	"Mason", "Sophia", "Logan", "Isla", "Lucas",
	"Mia", "Ethan", "Ava", "James", "Ella",
	"Jack", "Chloe", "Henry", "Grace", "Leo",
	"Zoe", "Owen", "Lily", "Caleb", "Nina",
	"Ryan", "Ruby", "Finn", "Maya", "Jade",
	"Luke", "Ivy", "Jake", "Anna", "Max",
	"Tina", "Seth", "Kate", "Nora", "Hugo",
	"Luna", "Rory", "Eli", "Cora", "Sam",
	"Troy", "Faye", "Milo", "Dara", "Zara",
	"Drew", "Gina", "Sage", "Beau", "Tess",
	"Reed", "Lila", "Cody", "Juno", "Bree",
	"Rita", "Kira", "Jace", "Lana", "Vera",
	"Wade", "Kara", "Dawn", "Zane", "Jett",
	"Ayla", "Nate", "Luca", "Skye", "Mira",
	"Rami", "Yara", "Joni", "Omar", "Zuri",
	"Liam", "Rami", "Cleo", "Milo", "Nico",
	"Tali", "Zane", "Cami", "Suki", "Demi",
	"Kian", "Rhea", "Tori", "Cora", "Gale",
}

var titles = []string{
	"The Art of Mindfulness",
	"10 Tips for Healthy Living",
	"Exploring the Cosmos",
	"The Future of Technology",
	"Traveling on a Budget",
	"Delicious Plant-Based Recipes",
	"Mastering Time Management",
	"The Power of Positive Thinking",
	"Sustainable Fashion Trends",
	"Fitness Routines for Busy People",
	"Understanding Mental Health",
	"DIY Home Decor Ideas",
	"The Benefits of Meditation",
	"Innovative Startups to Watch",
	"Gardening for Beginners",
	"Essential Skills for Remote Work",
	"Cultural Festivals Around the World",
	"The Science of Happiness",
	"Creative Writing Prompts",
	"Navigating Career Changes",
}

var contents = []string{
	"Embrace the magic of today! What are you grateful for? #Gratitude",
	"Ready to conquer your goals? Take one small step today! #Motivation",
	"Just finished an amazing book! What’s on your reading list? #BookLovers",
	"Creativity is contagious! Share your latest project with us! #ArtCommunity",
	"Autumn leaves and cozy vibes. What’s your favorite fall activity? #FallFun",
	"Let's make a difference! How do you contribute to your community? #GiveBack",
	"Music is the soundtrack of our lives. What’s your current favorite song? #MusicLovers",
	"Innovation starts with a single idea. What inspires you to create? #Innovation",
	"Every day is a new beginning. What will you do differently today? #NewBeginnings",
	"Fitness journey update: What’s your favorite way to stay active? #FitnessGoals",
	"Cooking up something delicious! Share your favorite recipe with us! #Foodie",
	"Nature is calling! What's your favorite outdoor adventure? #NatureLovers",
	"Take a deep breath and relax. How do you unwind after a long day? #SelfCare",
	"Let’s talk about dreams! What’s one dream you’re working towards? #DreamBig",
	"Celebrating small wins today! What’s something you’re proud of? #CelebrateSuccess",
	"Capture the moment! Share your favorite photo from this week. #Photography",
	"Collaboration over competition! Who inspires you in your field? #Networking",
	"Time management tip: Prioritize tasks for maximum productivity! #WorkSmart",
	"Express yourself through art! What medium do you love the most? #CreativeExpression",
	"Remember, every expert was once a beginner. Keep learning and growing! #LifelongLearning",
}

var tags = []string{
	"inspiration",
	"motivation",
	"wellness",
	"travel",
	"foodie",
	"lifestyle",
	"fitness",
	"photography",
	"nature",
	"adventure",
	"mindfulness",
	"creativity",
	"selfcare",
	"books",
	"technology",
	"fashion",
	"art",
	"music",
	"entrepreneurship",
	"success",
	"mentalhealth",
	"quotes",
	"family",
	"pets",
	"sustainability",
	"community",
	"events",
	"DIY",
	"tutorials",
	"reviews",
}

var commentsContent = []string{
	"Great post! Love this!",
	"Absolutely inspiring!",
	"This made my day!",
	"So true! Couldn't agree more.",
	"Fantastic insight, thanks for sharing!",
	"Wow, this is amazing!",
	"I needed this today, thank you!",
	"Such a beautiful perspective.",
	"Love the positivity here!",
	"This is spot on!",
	"Thanks for the motivation!",
	"Your creativity shines through!",
	"This is so relatable!",
	"I admire your passion!",
	"Keep up the great work!",
	"This is pure gold!",
	"Such an interesting take!",
	"You always know how to inspire.",
	"This made me smile!",
	"Incredible content as always!",
	"You have a way with words.",
	"Such a thoughtful post!",
	"I appreciate your insights.",
	"This is exactly what I needed to hear.",
	"Your posts always brighten my day.",
	"Love your enthusiasm!",
	"This topic is so important.",
	"You nailed it with this one!",
	"So insightful! Thank you for sharing.",
	"Your perspective is refreshing.",
	"This is a game changer!",
	"I love how you think!",
	"Such a well-written post!",
	"You have a gift for storytelling.",
	"This resonates with me deeply.",
	"Thanks for sparking this conversation!",
	"Your work is always top-notch!",
	"I can relate to this so much.",
	"Such a positive vibe here!",
	"You inspire me to do better.",
	"This is beautifully expressed.",
	"Thanks for sharing your thoughts!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)

	tx, _ := db.BeginTx(ctx, nil)

	for _, u := range users {
		if err := store.Users.Create(ctx, tx, u); err != nil {
			tx.Rollback()
			log.Printf("err:%v, data:%v", err, u)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, p := range posts {
		if err := store.Posts.Create(ctx, p); err != nil {
			log.Printf("err:%v, data:%v", err, p)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, c := range comments {
		if err := store.Comments.Create(ctx, c); err != nil {
			log.Printf("err:%v, data:%v", err, c)
			return
		}
	}
	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		randomIndex := r.Intn(len(names))

		users[i] = &store.User{
			Username: names[randomIndex] + fmt.Sprintf("_%d", i),
			Email:    names[randomIndex] + fmt.Sprintf("_%d@example.com", i),
			//Password: "123123",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		posts[i] = &store.Post{
			UserID:  users[r.Intn(len(users))].ID,
			Content: contents[r.Intn(len(contents))],
			Title:   titles[r.Intn(len(titles))],
			Tags: []string{
				tags[r.Intn(len(tags))],
				tags[r.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < num; i++ {
		randUser := users[r.Intn(len(users))]
		comments[i] = &store.Comment{
			PostId:  posts[r.Intn(len(posts))].ID,
			UserId:  randUser.ID,
			Content: commentsContent[r.Intn(len(commentsContent))],
			User:    *randUser,
		}
	}

	return comments
}
