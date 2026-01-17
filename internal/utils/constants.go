package utils

type IMAPSecurity string

const (
	IMAPSecurityTLS      IMAPSecurity = "ssl_tls"
	IMAPSecuritySTARTTLS IMAPSecurity = "starttls"
	IMAPSecurityNone     IMAPSecurity = "none"
)
