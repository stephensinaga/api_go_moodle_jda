package constant

const (
	// ganti nama moodle url dan token dengan benar
	MoodleURL         = "http://192.168.95.68/moodle/webservice/rest/server.php" // di ganti menyesuaikan url domain dari app moodle nya
	MoodleToken       = "6b5893a337b538f2660bfe5d58733db3" // ganti sesuai token yang di beri dari moodle
	WSFunctionCreate  = "core_user_create_users"  // tidak boleh diganti
	WSFunctionEnrol   = "enrol_manual_enrol_users"  // tidak boleh diganti
	RoleID			  = 5 // tidak boleh diganti
	CourseID	      = 3  // ganti sesuai dengan id dari course moodle
	HOST              = "192.168.95.114:4500"   // ini sesuai dengan host yang ingin anda buat
)