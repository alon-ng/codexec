package main

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Translation struct {
	Name        string
	Description string
	Bullets     string
}

type ExerciseSeed struct {
	Type         db.ExerciseType
	Reward       int16
	Data         map[string]interface{}
	Translations map[string]Translation
}

type LessonSeed struct {
	IsPublic     bool
	Translations map[string]Translation
	Exercises    []ExerciseSeed
}

type CourseSeed struct {
	Subject      string
	Price        int16
	Discount     int16
	Difficulty   int16
	Translations map[string]Translation
	Lessons      []LessonSeed
}

func seedCourse(ctx context.Context, queries *db.Queries) db.Course {
	log.Println("Seeding courses...")

	course := CourseSeed{
		Subject:    "python",
		Price:      0,
		Discount:   0,
		Difficulty: 1,
		Translations: map[string]Translation{
			"en": {
				Name:        "Python for Beginners",
				Description: "Learn the basics of Python programming from scratch.",
				Bullets:     "Variables and Data Types\nControl Flow and Loops\nFunctions and Modules",
			},
			"he": {
				Name:        "פייתון למתחילים",
				Description: "למד את היסודות של תכנות בפייתון מאפס.",
				Bullets:     "משתנים וסוגי נתונים\nבקרת זרימה ולולאות\nפונקציות ומודולים",
			},
		},
		Lessons: []LessonSeed{
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Introduction to Python", Description: "Welcome to Python!"},
					"he": {Name: "מבוא לפייתון", Description: "ברוכים הבאים לפייתון!"},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "print('Hello World')", "task": "Print Hello World"},
						Translations: map[string]Translation{
							"en": {Name: "Hello World", Description: "Your first program."},
							"he": {Name: "שלום עולם", Description: "התוכנית הראשונה שלך."},
						},
					},
					{
						Type:   db.ExerciseTypeQuiz,
						Reward: 5,
						Data:   map[string]interface{}{"question": "Is Python easy?", "options": []string{"Yes", "No"}, "answer": 0},
						Translations: map[string]Translation{
							"en": {Name: "Python Quiz", Description: "Simple question."},
							"he": {Name: "בוחן פייתון", Description: "שאלה פשוטה."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Variables and Data Types", Description: "Storing information."},
					"he": {Name: "משתנים וסוגי נתונים", Description: "אחסון מידע."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "x = 5\nprint(x)", "task": "Create a variable"},
						Translations: map[string]Translation{
							"en": {Name: "Integer Variable", Description: "Store a number."},
							"he": {Name: "משתנה שלם", Description: "שמור מספר."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "name = 'Alice'\nprint(name)", "task": "Create a string"},
						Translations: map[string]Translation{
							"en": {Name: "String Variable", Description: "Store text."},
							"he": {Name: "משתנה מחרוזת", Description: "שמור טקסט."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Control Flow", Description: "Making decisions."},
					"he": {Name: "בקרת זרימה", Description: "קבלת החלטות."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "if True:\n    print('Yes')", "task": "If statement"},
						Translations: map[string]Translation{
							"en": {Name: "If Statement", Description: "Conditional logic."},
							"he": {Name: "פקודת If", Description: "לוגיקה מותנית."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "if False:\n    print('No')\nelse:\n    print('Yes')", "task": "Else statement"},
						Translations: map[string]Translation{
							"en": {Name: "Else Statement", Description: "Alternative logic."},
							"he": {Name: "פקודת Else", Description: "לוגיקה חלופית."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Loops", Description: "Repeating actions."},
					"he": {Name: "לולאות", Description: "פעולות חוזרות."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "for i in range(5):\n    print(i)", "task": "For loop"},
						Translations: map[string]Translation{
							"en": {Name: "For Loop", Description: "Count to 5."},
							"he": {Name: "לולאת For", Description: "ספור עד 5."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "i = 0\nwhile i < 3:\n    print(i)\n    i += 1", "task": "While loop"},
						Translations: map[string]Translation{
							"en": {Name: "While Loop", Description: "While condition is true."},
							"he": {Name: "לולאת While", Description: "כל עוד התנאי מתקיים."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Functions", Description: "Reusable code blocks."},
					"he": {Name: "פונקציות", Description: "בלוקי קוד לשימוש חוזר."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 15,
						Data:   map[string]interface{}{"code": "def greet():\n    print('Hello')", "task": "Define a function"},
						Translations: map[string]Translation{
							"en": {Name: "Define Function", Description: "Create a greet function."},
							"he": {Name: "הגדר פונקציה", Description: "צור פונקציית greet."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 15,
						Data:   map[string]interface{}{"code": "def add(a, b):\n    return a + b", "task": "Function with arguments"},
						Translations: map[string]Translation{
							"en": {Name: "Function Arguments", Description: "Add two numbers."},
							"he": {Name: "ארגומנטים לפונקציה", Description: "חבר שני מספרים."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "Lists and Dictionaries", Description: "Complex data structures."},
					"he": {Name: "רשימות ומילונים", Description: "מבני נתונים מורכבים."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 15,
						Data:   map[string]interface{}{"code": "fruits = ['apple', 'banana']", "task": "Create a list"},
						Translations: map[string]Translation{
							"en": {Name: "Create List", Description: "Store a list of fruits."},
							"he": {Name: "יצירת רשימה", Description: "שמור רשימת פירות."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 15,
						Data:   map[string]interface{}{"code": "person = {'name': 'Alice', 'age': 25}", "task": "Create a dictionary"},
						Translations: map[string]Translation{
							"en": {Name: "Create Dictionary", Description: "Store person info."},
							"he": {Name: "יצירת מילון", Description: "שמור פרטי אדם."},
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {Name: "File Handling", Description: "Working with files."},
					"he": {Name: "טיפול בקבצים", Description: "עבודה עם קבצים."},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 20,
						Data:   map[string]interface{}{"code": "with open('test.txt', 'r') as f:\n    content = f.read()", "task": "Read a file"},
						Translations: map[string]Translation{
							"en": {Name: "Read File", Description: "Read from a file."},
							"he": {Name: "קריאת קובץ", Description: "קרא מקובץ."},
						},
					},
					{
						Type:   db.ExerciseTypeCode,
						Reward: 20,
						Data:   map[string]interface{}{"code": "with open('test.txt', 'w') as f:\n    f.write('Hello')", "task": "Write to a file"},
						Translations: map[string]Translation{
							"en": {Name: "Write File", Description: "Write to a file."},
							"he": {Name: "כתיבה לקובץ", Description: "כתוב לקובץ."},
						},
					},
				},
			},
		},
	}

	// Create course
	c, err := queries.CreateCourse(ctx, db.CreateCourseParams{
		Subject:    course.Subject,
		Price:      course.Price,
		Discount:   course.Discount,
		IsActive:   true,
		Difficulty: course.Difficulty,
	})
	if err != nil {
		log.Fatalf("Failed to create course: %v", err)
	}

	for lang, trans := range course.Translations {
		_, err = queries.CreateCourseTranslation(ctx, db.CreateCourseTranslationParams{
			CourseUuid:  c.Uuid,
			Language:    lang,
			Name:        trans.Name,
			Description: trans.Description,
			Bullets:     trans.Bullets,
		})
		if err != nil {
			log.Fatalf("Failed to create course translation (%s): %v", lang, err)
		}
	}

	for i, lSeed := range course.Lessons {
		l, err := queries.CreateLesson(ctx, db.CreateLessonParams{
			CourseUuid: c.Uuid,
			OrderIndex: int16(i + 1),
			IsPublic:   lSeed.IsPublic,
		})
		if err != nil {
			log.Fatalf("Failed to create lesson: %v", err)
		}

		for lang, trans := range lSeed.Translations {
			_, err = queries.CreateLessonTranslation(ctx, db.CreateLessonTranslationParams{
				LessonUuid:  l.Uuid,
				Language:    lang,
				Name:        trans.Name,
				Description: trans.Description,
			})
			if err != nil {
				log.Fatalf("Failed to create lesson translation (%s): %v", lang, err)
			}
		}

		for j, eSeed := range lSeed.Exercises {
			dataJSON, _ := json.Marshal(eSeed.Data)
			e, err := queries.CreateExercise(ctx, db.CreateExerciseParams{
				LessonUuid: l.Uuid,
				OrderIndex: int16(j + 1),
				Reward:     eSeed.Reward,
				Type:       eSeed.Type,
				Data:       dataJSON,
			})
			if err != nil {
				log.Fatalf("Failed to create exercise: %v", err)
			}

			for lang, trans := range eSeed.Translations {
				_, err = queries.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
					ExerciseUuid: e.Uuid,
					Language:     lang,
					Name:         trans.Name,
					Description:  trans.Description,
				})
				if err != nil {
					log.Fatalf("Failed to create exercise translation (%s): %v", lang, err)
				}
			}
		}
	}

	log.Println("Course seeded successfully!")
	return c
}

func getUUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
