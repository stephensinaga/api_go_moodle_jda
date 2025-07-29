package service

import (
	"moodleinix/database"
	"moodleinix/model"

)

func GetUserCourseQuizData() ([]model.UserCourseQuizData, error) {
	query := `
	SELECT
	  u.id AS user_id,
	  CONCAT(u.firstname, ' ', u.lastname) AS fullname,
	  u.email,
	  u.city,
	  u.department,
	  c.id AS course_id,
	  c.fullname AS course_name,

	  MAX(CASE WHEN q.name = 'Quiz DNS dan DHCP' THEN qg.grade END) AS Quiz_DNS_DHCP,
	  MAX(CASE WHEN q.name = 'Quiz Routing' THEN qg.grade END) AS Quiz_Routing,
	  MAX(CASE WHEN q.name = 'Quiz IP Address' THEN qg.grade END) AS Quiz_IP_Address,
	  MAX(CASE WHEN q.name = 'Quiz OSI Model' THEN qg.grade END) AS Quiz_OSI_Model,
	  MAX(CASE WHEN q.name = 'Quiz Topologi Jaringan' THEN qg.grade END) AS Quiz_Topologi_Jaringan,
	  MAX(CASE WHEN q.name = 'Quiz Ancaman Dunia Digital' THEN qg.grade END) AS Quiz_Ancaman_Dunia_Digital,
	  MAX(CASE WHEN q.name = 'Quiz Hukum/Aturan dalam Dunia Digital' THEN qg.grade END) AS Quiz_HukumAturan_dalam_Dunia_Digital,
	  MAX(CASE WHEN q.name = 'Quiz Tujuan dan Metode Ethical Hacking' THEN qg.grade END) AS Quiz_Tujuan_dan_Metode_Ethical_Hacking,
	  MAX(CASE WHEN q.name = 'Quiz Dasar Ethical Hacking dan Tipe Hacker' THEN qg.grade END) AS Quiz_Dasar_Ethical_Hacking_dan_Tipe_Hacker,
	  MAX(CASE WHEN q.name = 'Quiz UU PDP' THEN qg.grade END) AS Quiz_UU_PDP,
	  MAX(CASE WHEN q.name = 'Quiz UU ITE' THEN qg.grade END) AS Quiz_UU_ITE,
	  MAX(CASE WHEN q.name = 'Quiz Cyber Warfare' THEN qg.grade END) AS Quiz_Cyber_Warfare,
	  MAX(CASE WHEN q.name = 'Quiz Kerentanan Cybersecurity' THEN qg.grade END) AS Quiz_Kerentanan_Cybersecurity,
	  MAX(CASE WHEN q.name = 'Quiz Sumber Ancaman Cybersecurity' THEN qg.grade END) AS Quiz_Sumber_Ancaman_Cybersecurity,
	  MAX(CASE WHEN q.name = 'Quiz Tipe dan Ancaman Cybersecurity' THEN qg.grade END) AS Quiz_Tipe_dan_Ancaman_Cybersecurity,
	  MAX(CASE WHEN q.name = 'Test Jaringan Komputer Dasar' THEN qg.grade END) AS Test_Jaringan_Komputer_Dasar,
	  MAX(CASE WHEN q.name = 'Test Konsep Ethical Hacking' THEN qg.grade END) AS Test_Konsep_Ethical_Hacking,
	  MAX(CASE WHEN q.name = 'Test Akhir' THEN qg.grade END) AS Test_Akhir,

	  ROUND((cc.timecompleted - cc.timestarted) / 3600, 2) AS Course_Completion_Duration_Hours,
	  cc.timecompleted AS End_Date

	FROM mdl_user u
	JOIN mdl_user_enrolments ue ON ue.userid = u.id
	JOIN mdl_enrol e ON e.id = ue.enrolid
	JOIN mdl_course c ON c.id = e.courseid
	LEFT JOIN mdl_quiz q ON q.course = c.id
	LEFT JOIN mdl_quiz_grades qg ON qg.quiz = q.id AND qg.userid = u.id
	LEFT JOIN mdl_course_completions cc ON cc.course = c.id AND cc.userid = u.id

	GROUP BY u.id, u.firstname, u.lastname, u.email, u.city, u.department,
	         c.id, c.fullname, cc.timestarted, cc.timecompleted;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.UserCourseQuizData

	for rows.Next() {
		var ucqd model.UserCourseQuizData
		err := rows.Scan(
			&ucqd.UserID,
			&ucqd.Fullname,
			&ucqd.Email,
			&ucqd.City,
			&ucqd.Department,
			&ucqd.CourseID,
			&ucqd.CourseName,
			&ucqd.QuizDNSDHCP,
			&ucqd.QuizRouting,
			&ucqd.QuizIPAddress,
			&ucqd.QuizOSIModel,
			&ucqd.QuizTopologiJaringan,
			&ucqd.QuizAncamanDuniaDigital,
			&ucqd.QuizHukumAturanDuniaDigital,
			&ucqd.QuizTujuanMetodeEthicalHacking,
			&ucqd.QuizDasarEthicalHackingTipeHacker,
			&ucqd.QuizUUPDP,
			&ucqd.QuizUUITE,
			&ucqd.QuizCyberWarfare,
			&ucqd.QuizKerentananCybersecurity,
			&ucqd.QuizSumberAncamanCybersecurity,
			&ucqd.QuizTipeAncamanCybersecurity,
			&ucqd.TestJaringanKomputerDasar,
			&ucqd.TestKonsepEthicalHacking,
			&ucqd.TestAkhir,
			&ucqd.CourseCompletionDurationHours,
			&ucqd.EndDate,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, ucqd)
	}
	return results, nil
}