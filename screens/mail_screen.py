from PyQt5.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QHBoxLayout,
    QPushButton,
    QLabel,
    QGroupBox,
    QStyle,
)
from PyQt5.QtWebEngineWidgets import QWebEngineView
from PyQt5.QtCore import Qt
from mail import MailItem

class MailScreen(QWidget):
    def __init__(self, stack,fetcher):
        super().__init__()

        self.stack = stack
        self.fetcher = fetcher
        self.current_index = 0
        self.subject_label = QLabel()
        self.from_label = QLabel()
        self.to_label = QLabel()
        self.date_label = QLabel()

        self.subject_label.setTextInteractionFlags(Qt.TextSelectableByMouse)
        self.from_label.setTextInteractionFlags(Qt.TextSelectableByMouse)
        self.to_label.setTextInteractionFlags(Qt.TextSelectableByMouse)
        self.date_label.setTextInteractionFlags(Qt.TextSelectableByMouse)

        info_layout = QVBoxLayout()
        info_layout.addWidget(self.date_label)
        info_layout.addWidget(self.subject_label)
        info_layout.addWidget(self.from_label)
        info_layout.addWidget(self.to_label)

        info_group = QGroupBox("")
        info_group.setLayout(info_layout)

        self.view = QWebEngineView()

        self.back_button = QPushButton("Back")

        pixmapi = QStyle.SP_ArrowBack
        icon = self.style().standardIcon(pixmapi)

        self.back_button.setIcon(icon)
        self.prev_button = QPushButton("Previous")
        self.next_button = QPushButton("Next")
        self.back_button.clicked.connect(lambda: self.stack.setCurrentIndex(0))
        self.prev_button.clicked.connect(self.go_previous)
        self.next_button.clicked.connect(self.go_next)

        button_layout = QHBoxLayout()
        button_layout.addWidget(self.back_button)
        button_layout.addStretch()
        button_layout.addWidget(self.prev_button)
        button_layout.addWidget(self.next_button)

        main_layout = QVBoxLayout()
        main_layout.addLayout(button_layout)
        main_layout.addWidget(info_group)
        main_layout.addWidget(self.view, stretch=1)

        self.setLayout(main_layout)
        self.setWindowTitle("Bunny Mail")

        self.setStyleSheet(
            """
            QWidget {
                font-family: Arial;
                font-size: 14px;
            }
            QGroupBox {
                font-weight: bold;
                border: 1px solid #ddd;
                border-radius: 5px;
                margin-top: 10px;
            }
            QGroupBox::title {
                subcontrol-origin: margin;
                subcontrol-position: top left;
                padding: 5px 10px;
            }
            QPushButton {
                padding: 6px 12px;
                border-radius: 5px;
                background-color: #4CAF50;
                color: white;
            }
            QPushButton:disabled {
                background-color: #ccc;
            }
        """
        )

    def update_view(self,mail:MailItem):
        self.subject_label.setText(f"<b>Subject:</b> {mail.subject}")
        self.from_label.setText(f"<b>From:</b> {mail.from_addr}")
        self.to_label.setText(f"<b>To:</b> {mail.to_addr}")
        self.date_label.setText(f"<b>Date:</b> {mail.date}")
        self.view.setHtml(mail.body)

    def go_previous(self):
        mail = self.fetcher.get_mail(self.current_index - 1)
        if mail:
            self.current_index -= 1
            self.update_view(mail)
        self.update_buttons()

    def go_next(self):
        mail = self.fetcher.get_mail(self.current_index + 1)
        if mail:
            self.current_index += 1
            self.update_view(mail)
        elif self.current_index + 1 < self.fetcher.entire_mail_count:
            self.fetcher.fetch_mails()
            mail = self.fetcher.get_mail(self.current_index + 1)
            if mail:
                self.current_index += 1
                self.update_view(mail)
        self.update_buttons()

    def update_buttons(self):
        has_prev = self.current_index > 0
        has_next = self.current_index < self.fetcher.entire_mail_count - 1
        self.prev_button.setEnabled(has_prev)
        self.next_button.setEnabled(has_next)