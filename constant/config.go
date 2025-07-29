package constant

const (
	// ganti nama moodle url dan token dengan benar
	MoodleURL                = "https://jida.inixindobdg.co.id/webservice/rest/server.php" // di ganti menyesuaikan url domain dari app moodle nya
	MoodleToken              = "246d4b3ab5e32f9f64bc8cf4d373bfe8"                          // ganti sesuai token yang di beri dari moodle
	WSFunctionCreate         = "core_user_create_users"
	WSFunctionEnrol          = "enrol_manual_enrol_users"
	WSFunctionGetUsers       = "core_user_get_users"
	WSFunctionGetNameGroup   = "core_group_get_course_groups"
	WSFunctionAddUsertoGroup = "core_group_add_group_members"
	HOST                     = "192.168.95.114:8080" // ini sesuai dengan host yang ingin anda buat
	RoleID                   = 5                // tidak boleh diganti
	BasicAuthUserName        = "inixapijida"
	BasicAuthUPassword       = "inixjayasentosa"
	API_TOKEN 				 = "c84f3e2b-5a91-44d3-8f52-1b7d2f946eca"

	// Production
	PORT                     = ":8443"
	EncryptFullChain         = "/etc/letsencrypt/live/jida.inixindobdg.co.id/fullchain.pem"
	EncryptPriveKey          = "/etc/letsencrypt/live/jida.inixindobdg.co.id/privkey.pem"
)
