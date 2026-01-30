package main

import (
	"codim/pkg/db"
	"codim/pkg/executors/checkers"
	"codim/pkg/fs"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func createCodeData(fileName string, fileContent string) *json.RawMessage {
	entry := fs.Entry{
		Name:    fileName + ".py",
		Content: fileContent,
	}
	data, _ := json.Marshal(entry)
	result := json.RawMessage(data)
	return &result
}

func createQuizData() *json.RawMessage {
	data := json.RawMessage(`{}`)
	return &data
}

func createTranslationCodeData(text string) *json.RawMessage {
	data := map[string]interface{}{
		"text": text,
	}
	result, _ := json.Marshal(data)
	raw := json.RawMessage(result)
	return &raw
}

func createTranslationQuizData(questions []map[string]interface{}) *json.RawMessage {
	data := map[string]interface{}{
		"questions": questions,
	}
	result, _ := json.Marshal(data)
	raw := json.RawMessage(result)
	return &raw
}

type Translation struct {
	Name        string
	Description string
	Bullets     string
	Content     string
	CodeData    *json.RawMessage
	QuizData    *json.RawMessage
}

type ExerciseSeed struct {
	Type         db.ExerciseType
	Reward       int16
	Data         map[string]interface{}
	Translations map[string]Translation
	CodeChecker  *checkers.CodeChecker
	IoChecker    *checkers.IOChecker
	QuizChecker  *map[string]string
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

func createRawMessage(data []byte) *json.RawMessage {
	result := json.RawMessage(data)
	return &result
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
					"en": {
						Name:        "Introduction to Python",
						Description: "Welcome to Python!",
						Content:     "<div><h2>Welcome to Python!</h2><p>Python is a powerful and versatile programming language. In this lesson, you'll learn the basics of Python programming.</p><p>Let's start with some simple concepts:</p><ul><li>Python is easy to learn</li><li>Python is widely used</li><li>Python has a large community</li></ul></div>",
					},
					"he": {
						Name:        "מבוא לפייתון",
						Description: "ברוכים הבאים לפייתון!",
						Content:     "<div><h2>ברוכים הבאים לפייתון!</h2><p>פייתון היא שפת תכנות חזקה ורב-תכליתית. בשיעור זה תלמדו את היסודות של תכנות בפייתון.</p><p>בואו נתחיל עם כמה מושגים פשוטים:</p><ul><li>פייתון קל ללמידה</li><li>פייתון נמצאת בשימוש נרחב</li><li>לפייתון קהילה גדולה</li></ul></div>",
					},
				},
				Exercises: []ExerciseSeed{
					{
						Type:   db.ExerciseTypeCode,
						Reward: 10,
						Data:   map[string]interface{}{"code": "print('Hello World')", "task": "Print Hello World"},
						Translations: map[string]Translation{
							"en": {
								Name:        "Hello World",
								Description: "Your first program.",
								CodeData:    createRawMessage([]byte(`{"instructions": "<div>Print <code>Hello World</code> to the console.</div>"}`)),
							},
							"he": {
								Name:        "שלום עולם",
								Description: "התוכנית הראשונה שלך.",
								CodeData:    createRawMessage([]byte(`{"instructions": "<div>הדפס <code>Hello World</code> לקונסול.</div>"}`)),
							},
						},
						IoChecker: &checkers.IOChecker{
							Input:          "",
							ExpectedOutput: "Hello World",
						},
					},
					{
						Type:   db.ExerciseTypeQuiz,
						Reward: 5,
						Data:   map[string]interface{}{"question": "Is Python easy?", "options": []string{"Yes", "No"}, "answer": 0},
						Translations: map[string]Translation{
							"en": {
								Name:        "Python Quiz",
								Description: "Simple question.",
								QuizData:    createRawMessage([]byte(`{"1": {"question": "Is Python easy?", "answers": {"1": "Yes", "2": "No"}}, "2": {"question": "Is Python fun?", "answers": {"1": "Yes", "2": "No", "3": "Maybe", "4": "Not sure"}}}`)),
							},
							"he": {
								Name:        "בוחן פייתון",
								Description: "שאלה פשוטה.",
								QuizData:    createRawMessage([]byte(`{"1": {"question": "האם פייתון קל?", "answers": {"1": "כן", "2": "לא"}}, "2": {"question": "האם פייתון כיף?", "answers": {"1": "כן", "2": "לא", "3": "אולי", "4": "לא ברור"}}}`)),
							},
						},
						QuizChecker: &map[string]string{
							"1": "1",
							"2": "4",
						},
					},
				},
			},
			{
				IsPublic: true,
				Translations: map[string]Translation{
					"en": {
						Name:        "Variables and Data Types",
						Description: "Storing information.",
						Content:     "<div><h2>Variables and Data Types</h2><p>Variables are used to store data in Python. Python supports various data types including:</p><ul><li><strong>Integers</strong>: Whole numbers like 5, 10, -3</li><li><strong>Strings</strong>: Text data like 'Hello' or \"World\"</li><li><strong>Floats</strong>: Decimal numbers like 3.14, 2.5</li><li><strong>Booleans</strong>: True or False values</li></ul></div>",
					},
					"he": {
						Name:        "משתנים וסוגי נתונים",
						Description: "אחסון מידע.",
						Content:     "<div><h2>משתנים וסוגי נתונים</h2><p>משתנים משמשים לאחסון נתונים בפייתון. פייתון תומכת בסוגי נתונים שונים כולל:</p><ul><li><strong>מספרים שלמים</strong>: מספרים שלמים כמו 5, 10, -3</li><li><strong>מחרוזות</strong>: נתוני טקסט כמו 'שלום' או \"עולם\"</li><li><strong>מספרים עשרוניים</strong>: מספרים עשרוניים כמו 3.14, 2.5</li><li><strong>בוליאנים</strong>: ערכי True או False</li></ul></div>",
					},
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
						CodeChecker: &checkers.CodeChecker{
							Code:     "from test_utils import TestUtils\ntry:\n    from main import x\n    TestUtils.success(\"Found variable: x\")\n\n    if x != 5:\n        TestUtils.failure(\"Variable x is not 5\")\n    else:\n        TestUtils.success(\"Variable x is set to 5\")\nexcept ImportError:\n    TestUtils.failure(\"Missing variable: x\")",
							FileName: "tests.py",
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
					"en": {
						Name:        "Control Flow",
						Description: "Making decisions.",
						Content:     "<div><h2>Control Flow</h2><p>Control flow statements allow your program to make decisions and execute code conditionally. The main control flow statements in Python are:</p><ul><li><strong>if</strong>: Execute code if a condition is true</li><li><strong>elif</strong>: Check another condition if the previous one was false</li><li><strong>else</strong>: Execute code if all conditions are false</li></ul></div>",
					},
					"he": {
						Name:        "בקרת זרימה",
						Description: "קבלת החלטות.",
						Content:     "<div><h2>בקרת זרימה</h2><p>פקודות בקרת זרימה מאפשרות לתוכנית שלך לקבל החלטות ולבצע קוד בתנאים מסוימים. פקודות בקרת הזרימה העיקריות בפייתון הן:</p><ul><li><strong>if</strong>: בצע קוד אם תנאי הוא נכון</li><li><strong>elif</strong>: בדוק תנאי אחר אם הקודם היה שגוי</li><li><strong>else</strong>: בצע קוד אם כל התנאים שגויים</li></ul></div>",
					},
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
					"en": {
						Name:        "Loops",
						Description: "Repeating actions.",
						Content:     "<div><h2>Loops</h2><p>Loops allow you to repeat code multiple times. Python has two main types of loops:</p><ul><li><strong>for loops</strong>: Iterate over a sequence (like a list or string)</li><li><strong>while loops</strong>: Repeat code while a condition is true</li></ul><p>Loops are essential for processing collections of data efficiently.</p></div>",
					},
					"he": {
						Name:        "לולאות",
						Description: "פעולות חוזרות.",
						Content:     "<div><h2>לולאות</h2><p>לולאות מאפשרות לך לחזור על קוד מספר פעמים. בפייתון יש שני סוגים עיקריים של לולאות:</p><ul><li><strong>לולאות for</strong>: חזור על רצף (כמו רשימה או מחרוזת)</li><li><strong>לולאות while</strong>: חזור על קוד כל עוד תנאי נכון</li></ul><p>לולאות חיוניות לעיבוד אוספי נתונים ביעילות.</p></div>",
					},
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
					"en": {
						Name:        "Functions",
						Description: "Reusable code blocks.",
						Content:     "<div><h2>Functions</h2><p>Functions are reusable blocks of code that perform a specific task. They help organize your code and avoid repetition.</p><p>Key concepts:</p><ul><li><strong>Defining functions</strong>: Use the <code>def</code> keyword</li><li><strong>Parameters</strong>: Pass data to functions</li><li><strong>Return values</strong>: Functions can return results</li><li><strong>Reusability</strong>: Call functions multiple times</li></ul></div>",
					},
					"he": {
						Name:        "פונקציות",
						Description: "בלוקי קוד לשימוש חוזר.",
						Content:     "<div><h2>פונקציות</h2><p>פונקציות הן בלוקי קוד לשימוש חוזר שמבצעים משימה ספציפית. הן עוזרות לארגן את הקוד ולהימנע מחזרה.</p><p>מושגים מרכזיים:</p><ul><li><strong>הגדרת פונקציות</strong>: השתמש במילת המפתח <code>def</code></li><li><strong>פרמטרים</strong>: העבר נתונים לפונקציות</li><li><strong>ערכי החזרה</strong>: פונקציות יכולות להחזיר תוצאות</li><li><strong>שימוש חוזר</strong>: קרא לפונקציות מספר פעמים</li></ul></div>",
					},
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
					"en": {
						Name:        "Lists and Dictionaries",
						Description: "Complex data structures.",
						Content:     "<div><h2>Lists and Dictionaries</h2><p>Python provides powerful data structures for organizing data:</p><ul><li><strong>Lists</strong>: Ordered collections of items, mutable</li><li><strong>Dictionaries</strong>: Key-value pairs, very efficient for lookups</li></ul><p>These structures are fundamental for working with collections of data in Python.</p></div>",
					},
					"he": {
						Name:        "רשימות ומילונים",
						Description: "מבני נתונים מורכבים.",
						Content:     "<div><h2>רשימות ומילונים</h2><p>פייתון מספקת מבני נתונים חזקים לארגון נתונים:</p><ul><li><strong>רשימות</strong>: אוספים מסודרים של פריטים, ניתנים לשינוי</li><li><strong>מילונים</strong>: זוגות מפתח-ערך, יעילים מאוד לחיפושים</li></ul><p>מבנים אלה הם בסיסיים לעבודה עם אוספי נתונים בפייתון.</p></div>",
					},
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
					"en": {
						Name:        "File Handling",
						Description: "Working with files.",
						Content:     "<div><h2>File Handling</h2><p>Python makes it easy to work with files. You can read from and write to files using built-in functions.</p><p>Key operations:</p><ul><li><strong>Opening files</strong>: Use the <code>open()</code> function</li><li><strong>Reading files</strong>: Read entire file or line by line</li><li><strong>Writing files</strong>: Write data to files</li><li><strong>Closing files</strong>: Always close files when done (or use <code>with</code> statement)</li></ul></div>",
					},
					"he": {
						Name:        "טיפול בקבצים",
						Description: "עבודה עם קבצים.",
						Content:     "<div><h2>טיפול בקבצים</h2><p>פייתון מקלה על עבודה עם קבצים. אתה יכול לקרוא מקבצים ולכתוב לקבצים באמצעות פונקציות מובנות.</p><p>פעולות מרכזיות:</p><ul><li><strong>פתיחת קבצים</strong>: השתמש בפונקציה <code>open()</code></li><li><strong>קריאת קבצים</strong>: קרא את כל הקובץ או שורה אחר שורה</li><li><strong>כתיבה לקבצים</strong>: כתוב נתונים לקבצים</li><li><strong>סגירת קבצים</strong>: תמיד סגור קבצים כשסיימת (או השתמש בהצהרת <code>with</code>)</li></ul></div>",
					},
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
				Content:     trans.Content,
			})
			if err != nil {
				log.Fatalf("Failed to create lesson translation (%s): %v", lang, err)
			}
		}

		for j, eSeed := range lSeed.Exercises {
			var codeData *json.RawMessage
			var quizData *json.RawMessage

			if eSeed.Type == db.ExerciseTypeCode {
				// Extract code from Data if available, otherwise use default
				codeContent := "print('Hello World')"
				if code, ok := eSeed.Data["code"].(string); ok {
					codeContent = code
				}
				codeData = createCodeData("main", codeContent)
				quizData = nil
			} else {
				codeData = nil
				quizData = createQuizData()
			}

			params := db.CreateExerciseParams{
				LessonUuid: l.Uuid,
				OrderIndex: int16(j + 1),
				Reward:     eSeed.Reward,
				Type:       eSeed.Type,
				CodeData:   codeData,
				QuizData:   quizData,
			}

			if eSeed.CodeChecker != nil {
				codeCheckerData, _ := json.Marshal(eSeed.CodeChecker)
				params.CodeChecker = createRawMessage(codeCheckerData)
			}
			if eSeed.IoChecker != nil {
				ioCheckerData, _ := json.Marshal(eSeed.IoChecker)
				params.IoChecker = createRawMessage(ioCheckerData)
			}
			if eSeed.QuizChecker != nil {
				quizCheckerData, _ := json.Marshal(eSeed.QuizChecker)
				params.QuizChecker = createRawMessage(quizCheckerData)
			}

			e, err := queries.CreateExercise(ctx, params)
			if err != nil {
				log.Fatalf("Failed to create exercise: %v", err)
			}

			for lang, trans := range eSeed.Translations {
				_, err = queries.CreateExerciseTranslation(ctx, db.CreateExerciseTranslationParams{
					ExerciseUuid: e.Uuid,
					Language:     lang,
					Name:         trans.Name,
					Description:  trans.Description,
					CodeData:     trans.CodeData,
					QuizData:     trans.QuizData,
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
