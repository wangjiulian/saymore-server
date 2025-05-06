package repository

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

// MockRepositories 模拟仓库实例
type MockRepositories struct {
	// 模拟数据
	Students          []db.Student
	Teachers          []db.Teacher
	Courses           []db.Course
	StudentPackages   []db.StudentPackage
	Subjects          []db.Subject
	Banners           []db.Banner
	CourseEvaluations []db.CourseEvaluation
	
	// 配置
	Config *config.Config
}

// NewMockRepositories 创建新的模拟仓库
func NewMockRepositories() *MockRepositories {
	return &MockRepositories{
		Students:          []db.Student{},
		Teachers:          []db.Teacher{},
		Courses:           []db.Course{},
		StudentPackages:   []db.StudentPackage{},
		Subjects:          []db.Subject{},
		Banners:           []db.Banner{},
		CourseEvaluations: []db.CourseEvaluation{},
		Config: &config.Config{
			Course: &config.Course{
				CancelInterval: 24,
				CancelRefund:   1,
				CancelRule:     "课程开始前24小时可以免费取消，之后不予退款",
				CourseUnit:     60,
			},
		},
	}
}

// InitMockDB 初始化模拟数据库数据
func (m *MockRepositories) InitMockDB() {
	// 添加学科数据
	m.addSubjects()
	
	// 添加学生数据
	m.addStudents()
	
	// 添加教师数据
	m.addTeachers()
	
	// 添加课程数据
	m.addCourses()
	
	// 添加Banner数据
	m.addBanners()
	
	// 设置仓库
	repository.Repos = &repository.Repositories{
		Config: m.Config,
	}
}

// addSubjects 添加学科数据
func (m *MockRepositories) addSubjects() {
	m.Subjects = []db.Subject{
		{
			Base:      db.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:      "英语",
			Code:      "english",
			ParentId:  0,
			SortOrder: 1,
		},
		{
			Base:      db.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:      "口语",
			Code:      "speaking",
			ParentId:  1,
			SortOrder: 1,
		},
		{
			Base:      db.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Name:      "听力",
			Code:      "listening",
			ParentId:  1,
			SortOrder: 2,
		},
	}
}

// addStudents 添加学生数据
func (m *MockRepositories) addStudents() {
	m.Students = []db.Student{
		{
			Base:            db.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Phone:           "13800138000",
			AvatarUrl:       "https://example.com/avatar1.jpg",
			Nickname:        "测试学生1",
			Gender:          1,
			BirthDate:       time.Now().AddDate(-20, 0, 0),
			StudentType:     2,
			LearningPurpose: 3,
			EnglishLevel:    4,
			IsActive:        1,
		},
		{
			Base:            db.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Phone:           "13800138001",
			AvatarUrl:       "https://example.com/avatar2.jpg",
			Nickname:        "测试学生2",
			Gender:          2,
			BirthDate:       time.Now().AddDate(-18, 0, 0),
			StudentType:     3,
			LearningPurpose: 1,
			EnglishLevel:    2,
			IsActive:        1,
		},
	}
}

// addTeachers 添加教师数据
func (m *MockRepositories) addTeachers() {
	m.Teachers = []db.Teacher{
		{
			Base:                  db.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Phone:                 "13900139000",
			Name:                  "王老师",
			Nickname:              "Teacher Wang",
			Gender:                1,
			CourseHours:           decimal.NewFromFloat(120),
			AvatarUrl:             "https://example.com/teacher1.jpg",
			Background:            "5年教学经验",
			VideoUrl:              "https://example.com/video1.mp4",
			EducationSchool:       "北京大学",
			EducationLevel:        2,
			TeachingStartDate:     time.Now().AddDate(-5, 0, 0),
			Notes:                 "擅长口语教学",
			TeachingExperience:    "曾在多家知名英语培训机构任教",
			SuccessCases:          "帮助多名学生提升英语口语能力",
			TeachingAchievements:  "获得年度优秀教师奖",
			IsActive:              1,
			IsRecommend:           1,
			Evaluation:            decimal.NewFromFloat(4.8),
		},
		{
			Base:                  db.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Phone:                 "13900139001",
			Name:                  "李老师",
			Nickname:              "Teacher Li",
			Gender:                2,
			CourseHours:           decimal.NewFromFloat(80),
			AvatarUrl:             "https://example.com/teacher2.jpg",
			Background:            "3年教学经验",
			VideoUrl:              "https://example.com/video2.mp4",
			EducationSchool:       "清华大学",
			EducationLevel:        3,
			TeachingStartDate:     time.Now().AddDate(-3, 0, 0),
			Notes:                 "擅长听力训练",
			TeachingExperience:    "专注于雅思托福备考",
			SuccessCases:          "多名学生成功通过雅思考试",
			TeachingAchievements:  "编写了系列英语听力教材",
			IsActive:              1,
			IsRecommend:           0,
			Evaluation:            decimal.NewFromFloat(4.5),
		},
	}
}

// addCourses 添加课程数据
func (m *MockRepositories) addCourses() {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	
	m.Courses = []db.Course{
		{
			Base:         db.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			CourseType:   db.CourseTypeRegular,
			Name:         "英语口语提高课程",
			SubjectId:    2,
			TeacherId:    1,
			StudentId:    1,
			Status:       db.CourseStatusFinished,
			IsEvaluated:  db.CourseIsEvaluated,
			CancelReason: "",
			StartTime:    uint64(yesterday.Add(-2 * time.Hour).Unix()),
			EndTime:      uint64(yesterday.Add(-1 * time.Hour).Unix()),
		},
		{
			Base:         db.Base{ID: 2, CreatedAt: now, UpdatedAt: now},
			CourseType:   db.CourseTypeRegular,
			Name:         "英语听力训练课程",
			SubjectId:    3,
			TeacherId:    2,
			StudentId:    1,
			Status:       db.CourseStatusWaiting,
			IsEvaluated:  db.CourseIsNotEvaluated,
			CancelReason: "",
			StartTime:    uint64(tomorrow.Add(2 * time.Hour).Unix()),
			EndTime:      uint64(tomorrow.Add(3 * time.Hour).Unix()),
		},
		{
			Base:         db.Base{ID: 3, CreatedAt: now, UpdatedAt: now},
			CourseType:   db.CourseTypeTrial,
			Name:         "英语口语试听课",
			SubjectId:    2,
			TeacherId:    1,
			StudentId:    2,
			Status:       db.CourseStatusCanceled,
			IsEvaluated:  db.CourseIsNotEvaluated,
			CancelReason: "学生临时有事",
			StartTime:    uint64(yesterday.Unix()),
			EndTime:      uint64(yesterday.Add(time.Hour).Unix()),
		},
	}
	
	// 添加课程评价
	m.CourseEvaluations = []db.CourseEvaluation{
		{
			Base:      db.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			CourseId:  1,
			StudentId: 1,
			TeacherId: 1,
			Score:     5,
			Content:   "老师讲课非常清晰，很满意这节课",
			IsPublic:  1,
		},
	}
}

// addBanners 添加Banner数据
func (m *MockRepositories) addBanners() {
	m.Banners = []db.Banner{
		{
			Base:      db.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Title:     "精品课程推荐",
			ImageUrl:  "https://example.com/banner1.jpg",
			TargetUrl: "https://example.com/courses",
			Sort:      1,
		},
		{
			Base:      db.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Title:     "名师在线",
			ImageUrl:  "https://example.com/banner2.jpg",
			TargetUrl: "https://example.com/teachers",
			Sort:      2,
		},
	}
}

// MockDB 模拟 GORM DB 接口的简单实现
// 这个实现只是一个骨架，需要根据实际需求进一步完善
type MockDB struct {
	mockRepos *MockRepositories
}

// WithContext 实现 gorm.DB 接口
func (m *MockDB) WithContext(ctx interface{}) *gorm.DB {
	return &gorm.DB{}
}

// 在实际使用中，需要实现更多的 gorm.DB 接口方法
// 根据测试需要，模拟数据库查询结果 