from mail import EmailFetcher, MailItem
from PyQt5.QtWidgets import QApplication, QStackedWidget, QLabel
from PyQt5.QtCore import Qt, QCoreApplication
from dotenv import load_dotenv
import os
import sys
import json
from screens import HomeScreen, MailScreen, LoaderScreen

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
    
    # emails_data = load_emails_from_json("test/testEmails.json")
    QCoreApplication.setAttribute(Qt.AA_UseSoftwareOpenGL)
    app = QApplication([])
    stack = QStackedWidget()
    screen1 = HomeScreen(stack, fetcher)
    screen2 = MailScreen(stack, fetcher)
    screen3 = LoaderScreen(stack)

    stack.addWidget(screen1)  
    stack.addWidget(screen2)  
    stack.addWidget(screen3)
    stack.setCurrentIndex(0)
    stack.resize(900, 700)
    stack.show()
    sys.exit(app.exec_())
