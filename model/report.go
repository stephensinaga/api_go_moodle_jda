package model

type UserCourseQuizData struct {
	UserID                   int
	Fullname                 string
	Email                    string
	City                     string
	Department               string
	CourseID                 int
	CourseName               string
	QuizDNSDHCP              *float64
	QuizRouting              *float64
	QuizIPAddress            *float64
	QuizOSIModel             *float64
	QuizTopologiJaringan     *float64
	QuizAncamanDuniaDigital  *float64
	QuizHukumAturanDuniaDigital *float64
	QuizTujuanMetodeEthicalHacking *float64
	QuizDasarEthicalHackingTipeHacker *float64
	QuizUUPDP                *float64
	QuizUUITE                *float64
	QuizCyberWarfare         *float64
	QuizKerentananCybersecurity *float64
	QuizSumberAncamanCybersecurity *float64
	QuizTipeAncamanCybersecurity *float64
	TestJaringanKomputerDasar *float64
	TestKonsepEthicalHacking *float64
	TestAkhir                *float64

	CourseCompletionDurationHours *float64
}