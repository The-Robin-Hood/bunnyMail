from mail import EmailFetcher, MailItem
from gui import MailViewerScreen,MainScreen
from PyQt5.QtWidgets import QApplication,  QStackedWidget
from dotenv import load_dotenv
import os
import sys
import json


def load_emails_from_json(file_path: str) -> list[MailItem]:
    with open(file_path, "r") as f:
        emails_json = json.load(f)
    return [MailItem(**mail) for mail in emails_json]

if __name__ == "__main__":
    load_dotenv()
    IMAP_SERVER = os.getenv("IMAP_SERVER")
    EMAIL = os.getenv("EMAIL")
    APP_PASSWORD = os.getenv("APP_PASSWORD")
    fetcher = EmailFetcher(EMAIL, APP_PASSWORD, IMAP_SERVER)
    emails_data = fetcher.fetch_next(10)  
    
    # emails_data = load_emails_from_json("sample_emails.json")

    app = QApplication([])
    stack = QStackedWidget()
    screen2 = MailViewerScreen(stack, emails_data)
    screen1 = MainScreen(stack, emails_data)


    stack.addWidget(screen1)  # index 0
    stack.addWidget(screen2)  # index 1
    stack.setCurrentIndex(0)  # show first screen
    stack.resize(900, 700)
    stack.show()
    sys.exit(app.exec_())
