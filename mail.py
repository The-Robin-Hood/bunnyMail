import imaplib
import email
from email.header import decode_header


class MailItem:
    def __init__(self, date, from_addr, to_addr, subject, body, seen):
        self.date = date
        self.from_addr = from_addr
        self.to_addr = to_addr
        self.subject = subject
        self.body = body
        self.seen = seen


class EmailFetcher:
    def __init__(
        self, email_address, app_password, imap_server="imap.gmail.com", folder="INBOX"
    ):
        self.email_address = email_address
        self.app_password = app_password
        self.imap_server = imap_server
        self.folder = folder
        self.imap = None
        self.email_ids = []
        self.current_index = 0

    def connect(self):
        self.imap = imaplib.IMAP4_SSL(self.imap_server)
        self.imap.login(self.email_address, self.app_password)
        self.imap.select(self.folder)
        status, messages = self.imap.search(None, "ALL")
        self.email_ids = messages[0].split()
        self.email_ids.reverse()

    def close(self):
        if self.imap:
            self.imap.close()
            self.imap.logout()
            self.imap = None

    def fetch_next(self, count=10) -> list[MailItem]:
        if not self.imap:
            self.connect()

        mails: list[MailItem] = []
        next_index = self.current_index + count
        for eid in self.email_ids[self.current_index : next_index]:
            status, msg_data = self.imap.fetch(eid, "(BODY.PEEK[])")
            s,flags_data =  self.imap.fetch(eid, "(FLAGS)")
            flags = flags_data[0].decode()
            seen = False
            if('\\Seen' in flags):
                seen = True
            if status != "OK" or not msg_data:
                print(f"Failed to fetch email id {eid}")
                continue

            raw_email = None
            if (
                isinstance(msg_data[0], tuple)
                and len(msg_data[0]) > 1
                and isinstance(msg_data[0][1], (bytes, bytearray))
            ):
                raw_email = msg_data[0][1]
            if not raw_email:
                continue

            msg = email.message_from_bytes(raw_email)

            subject, encoding = decode_header(msg.get("Subject"))[0]
            if isinstance(subject, bytes):
                subject = subject.decode(encoding or "utf-8")

            body = self._get_body(msg)
            mails.append(
                MailItem(
                    date=msg.get("Date"),
                    from_addr=msg.get("From"),
                    to_addr=msg.get("To"),
                    subject=subject,
                    body=body,
                    seen=seen,
                )
            )

        self.current_index = next_index
        return mails

    def _get_body(self, msg):
        body = ""
        if msg.is_multipart():
            for part in msg.walk():
                content_type = part.get_content_type()
                if content_type == "text/html":
                    part_bytes = part.get_payload(decode=True)
                    body += part_bytes.decode(errors="ignore")
        else:
            part_bytes = msg.get_payload(decode=True)
            if part_bytes:
                body += part_bytes.decode(errors="ignore")
        return body.strip()
