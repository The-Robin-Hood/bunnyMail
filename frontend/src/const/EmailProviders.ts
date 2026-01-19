const emailProviders: Record<
  string,
  {
    imapServer: string;
    imapPort: string;
    imapSecurity: string;
    imapAuth: string;
    smtpServer: string;
    smtpPort: string;
    smtpSecurity: string;
    smtpAuth: string;
  }
> = {
  "gmail.com": {
    imapServer: "imap.gmail.com",
    imapPort: "993",
    imapSecurity: "ssl_tls",
    imapAuth: "OAuth2",
    smtpServer: "smtp.gmail.com",
    smtpPort: "587",
    smtpSecurity: "ssl_tls",
    smtpAuth: "OAuth2",
  },
  "yahoo.com": {
    imapServer: "imap.mail.yahoo.com",
    imapPort: "993",
    imapSecurity: "ssl_tls",
    imapAuth: "OAuth2",
    smtpServer: "smtp.mail.yahoo.com",
    smtpPort: "587",
    smtpSecurity: "ssl_tls",
    smtpAuth: "OAuth2",
  },
  "outlook.com": {
    imapServer: "outlook.office365.com",
    imapPort: "993",
    imapSecurity: "ssl_tls",
    imapAuth: "OAuth2", 
    smtpServer: "smtp.office365.com",
    smtpPort: "587",
    smtpSecurity: "ssl_tls",
    smtpAuth: "OAuth2",
  }
};

export default emailProviders;
