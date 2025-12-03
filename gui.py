from PyQt5.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QHBoxLayout,
    QPushButton,
    QLabel,
    QGroupBox,
    QListWidget,
    QStyle,
    QListWidgetItem,
)
from PyQt5.QtWebEngineWidgets import QWebEngineView
from PyQt5.QtCore import Qt
from mail import MailItem
from email.utils import parsedate_to_datetime
from datetime import datetime


class MailViewerScreen(QWidget):
    def __init__(self, stack, mails: list[MailItem] = None):
        super().__init__()

        self.stack = stack

        if mails is None or len(mails) == 0:    
            mails = [MailItem("<h1>No Content</h1>", "", "", "", "")]
        self.mails = mails
        self.current_index = 0

        # Info labels
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

        # Web view for email body
        self.view = QWebEngineView()

        # Navigation buttons
        self.back_button = QPushButton("Back")

        # change back button icon to a left arrow and set color of icon to red

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

        # Main layout
        main_layout = QVBoxLayout()
        main_layout.addLayout(button_layout)
        main_layout.addWidget(info_group)
        main_layout.addWidget(self.view, stretch=1)

        self.setLayout(main_layout)
        self.setWindowTitle("Bunny Mail")

        # Apply initial content
        self.update_view()
        self.update_buttons()

        # Apply some modern styling
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

    def update_view(self):
        mail = self.mails[self.current_index]
        self.subject_label.setText(f"<b>Subject:</b> {mail.subject}")
        self.from_label.setText(f"<b>From:</b> {mail.from_addr}")
        self.to_label.setText(f"<b>To:</b> {mail.to_addr}")
        self.date_label.setText(f"<b>Date:</b> {mail.date}")
        self.view.setHtml(mail.body)

    def go_previous(self):
        if self.current_index > 0:
            self.current_index -= 1
            self.update_view()
        self.update_buttons()

    def go_next(self):
        if self.current_index < len(self.mails) - 1:
            self.current_index += 1
            self.update_view()
        self.update_buttons()

    def update_buttons(self):
        count = len(self.mails)
        has_prev = self.current_index > 0
        has_next = self.current_index < count - 1

        self.prev_button.setEnabled(has_prev)
        self.next_button.setEnabled(has_next)
        self.prev_button.setVisible(has_prev)
        self.next_button.setVisible(has_next)

class MainScreen(QWidget):
    def __init__(self, stack, mails: list[MailItem] = None):
        super().__init__()

        self.stack = stack

        if mails is None or len(mails) == 0:
            mails = [MailItem("<h1>No Content</h1>", "", "", "", "")]
        self.mails = mails
        self.current_index = 0

        self.setWindowTitle("Bunny Mail")

        self.header_label = QLabel("Inbox")
        self.header_label.setObjectName("headerLabel")

        self.list_widget = QListWidget()
        self.list_widget.setObjectName("mailList")
        self.populate_list()
        self.list_widget.itemClicked.connect(
            lambda item: self.open_mail(item.data(Qt.UserRole), self.mails)
        )

        self.next_button = QPushButton("Next 10")
        self.next_button.setObjectName("nextButton")
        self.next_button.clicked.connect(self.load_next)

        layout = QVBoxLayout()
        layout.addWidget(self.header_label)
        layout.addWidget(self.list_widget)
        layout.addWidget(self.next_button, alignment=Qt.AlignRight)
        self.setLayout(layout)


        # Apply modern styles
        self.setStyleSheet("""
            QWidget {
                font-family: 'Segoe UI', 'Arial', sans-serif;
                background-color: transparent;
                color: #111;
            }
            #headerLabel {
                font-size: 24px;
                font-weight: 600;
                padding: 15px 10px;
                           color: white;
            }
            #mailList {
                background-color: #fff;
                border: 1px solid #e0e0e0;
                border-radius: 8px;
                padding: 5px;
            }
            QListWidget::item {
                min-height: 50px;
                border-bottom: 1px solid #f0f0f0;
            }

            QListWidget::item:selected {
                background-color: #e6f7ff;
            }
            QListWidget::item:hover {
                background-color: #f5f5f5;
            }

            #nextButton {
                background-color: #4CAF50;
                color: white;
                border: none;
                padding: 8px 16px;
                border-radius: 6px;
                font-weight: 500;
            }
            #nextButton:hover {
                background-color: #45a049;
            }
        """)
    
    def trimmed_text(self, text, max_length=80):
        if len(text) <= max_length:
            return text
        return text[: max_length - 3] + "..."
    

    def format_mail_date(self,date_str):
        try:
            dt = parsedate_to_datetime(date_str)
            local_tz = datetime.now().astimezone().tzinfo
            dt_local = dt.astimezone(local_tz)
            now = datetime.now().astimezone()
            if dt_local.date() == now.date():
                return dt_local.strftime("%-I:%M %p")
            if dt_local.year == now.year:
                return dt_local.strftime("%b %d")
            return dt_local.strftime("%d/%m/%y")

        except Exception:
            return date_str
    
    def formatted_label_text(self, mail: MailItem):
        sender = mail.from_addr
        subject = self.trimmed_text(mail.subject)
        date = self.format_mail_date(mail.date)
        html = f"""
            <table width="100%" style="font-size:14px; border-collapse: collapse;">
                <tr>
                    <td width="25%" style="font-weight:500; {mail.seen and 'color:#777;' or 'color:#333;'}">
                    {'<span style="color:#4CAF50;font-size:10px">&#9679; </span>' if not mail.seen else ''}{sender}
                    </td>
                    <td width="65%" style="padding-left:20px;color:#555;">
                        {subject}
                    </td>
                    <td width="10%" style="color:#999; font-size:12px;" align="right">
                        {date}
                    </td>
                </tr>
            </table>
        """

        return html



    def populate_list(self):
        self.list_widget.clear()
        for mail in self.mails[self.current_index:self.current_index + 10]:
            item = QListWidgetItem()
            item.setData(Qt.UserRole, mail)

            label = QLabel()
            label.setTextFormat(Qt.RichText)  # Enable HTML rendering
            label.setText(self.formatted_label_text(mail))
            label.setContentsMargins(10, 5, 10, 5)
            label.setCursor(Qt.PointingHandCursor)
            label.setStyleSheet("""
                QLabel {
                    font-size: 14px;
                    background-color: transparent;
                }
            """)
            label.mouseDoubleClickEvent = lambda event, m=mail: self.open_mail(m)

            self.list_widget.addItem(item)
            self.list_widget.setItemWidget(item, label)

    def load_next(self):
        if self.current_index + 10 < len(self.mails):
            self.current_index += 10
            self.populate_list()

    def open_mail(self, selected_mail):
        self.stack.widget(1).current_index = self.mails.index(selected_mail)
        self.stack.widget(1).update_view()
        self.stack.widget(1).update_buttons()
        self.stack.setCurrentIndex(1)  