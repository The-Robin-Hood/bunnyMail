import imaplib
import email
from email.header import decode_header
from pyexpat.errors import messages
from utils import Utils
from PyQt5.QtCore import QObject, pyqtSignal

class MailItem:
    def __init__(self,id, date, from_addr, to_addr, subject, body, seen):
        self.id = id
        self.date = Utils.format_mail_date(date)
        self.from_addr = from_addr
        self.to_addr = to_addr
        self.subject = subject
        self.body = body
        self.seen = seen

class MailFetcherWorker(QObject):
    finished = pyqtSignal(list)

    def __init__(self, fetcher):
        super().__init__()
        self.fetcher = fetcher

    def run(self):
        self.fetcher.fetch_mails()
        self.finished.emit(self.fetcher.fetched_mail_list)

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
        self.entire_mail_count = 0
        self.fetched_mail_list = []
        self.fetch_mails()

    def connect(self):
        self.imap = imaplib.IMAP4_SSL(self.imap_server)
        self.imap.login(self.email_address, self.app_password)
        self.imap.select(self.folder)
        status, messages = self.imap.search(None, "ALL")
        self.email_ids = messages[0].split()
        self.email_ids.reverse()
        self.entire_mail_count = len(self.email_ids)

    def close(self):
        if self.imap:
            self.imap.close()
            self.imap.logout()
            self.imap = None
    
    def get_mail(self,index:int)->MailItem:
        if 0 <= index < len(self.fetched_mail_list):
            return self.fetched_mail_list[index]
        return None

    def fetch_mails(self, count=20) -> None:
        print("Fetching mails...")
        if not self.imap:
            self.connect()

        next_mails = []
        next_index = self.current_index + count

        id_list = ','.join(eid.decode() if isinstance(eid, bytes) else str(eid) 
                        for eid in self.email_ids[self.current_index : next_index])
        
        status, messages = self.imap.fetch(id_list, "(BODY.PEEK[] FLAGS)")
        messages = [m for m in messages if isinstance(m, tuple)]
        messages.reverse()

        for flags_data,msg_data in messages:
            flags = flags_data.decode()
            seen = False
            if('\\Seen' in flags):
                seen = True
            if status != "OK" or not msg_data:
                continue

            msg = email.message_from_bytes(msg_data)

            subject, encoding = decode_header(msg.get("Subject"))[0]
            if isinstance(subject, bytes):
                subject = subject.decode(encoding or "utf-8")

            body = self._get_body(msg)
            next_mails.append(
                MailItem(
                    id=flags_data.decode().split()[0],
                    date=msg.get("Date"),
                    from_addr=msg.get("From"),
                    to_addr=msg.get("To"),
                    subject=subject,
                    body=body,
                    seen=seen,
                )
            )
        
        self.current_index = next_index
        self.fetched_mail_list.extend(next_mails)
    
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
