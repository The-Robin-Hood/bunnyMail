from PyQt5.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QPushButton,
    QLabel,
    QListWidget,
    QListWidgetItem,
    QHBoxLayout,
)
from PyQt5.QtCore import Qt,QThread

from mail import EmailFetcher, MailItem, MailFetcherWorker
from utils import Utils


class HomeScreen(QWidget):
    def __init__(self, stack, fetcher: EmailFetcher):
        super().__init__()

        self.current_index = 0
        self.stack = stack
        self.fetcher = fetcher
        self.mails = fetcher.fetched_mail_list

        self.setWindowTitle("Bunny Mail")

        self.header_label = QLabel("Inbox")
        self.header_label.setObjectName("headerLabel")

        self.list_widget = QListWidget()
        self.list_widget.setObjectName("mailList")
        self.populate_list()


        self.prev_button = QPushButton("Prev Page")
        self.prev_button.setObjectName("prevButton")
        self.prev_button.clicked.connect(self.load_previous)
        self.next_button = QPushButton("Next Page")
        self.next_button.setObjectName("nextButton")
        self.next_button.clicked.connect(self.load_next)

        pagination_button_layout = QHBoxLayout()
        pagination_button_layout.addStretch()
        pagination_button_layout.addWidget(self.prev_button)
        pagination_button_layout.addWidget(self.next_button)

        layout = QVBoxLayout()
        layout.addWidget(self.header_label)
        layout.addWidget(self.list_widget)
        layout.addLayout(pagination_button_layout)
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

            #prevButton,
            #nextButton {
                background-color: #4CAF50;
                color: white;
                border: none;
                padding: 8px 16px;
                border-radius: 6px;
                font-weight: 500;
            }
                           
            #prevButton:hover,
            #nextButton:hover {
                background-color: #45a049;
            }
        """)
        
    
    def formatted_label_text(self, mail: MailItem):
        sender = mail.from_addr
        subject = Utils.trimmed_text(mail.subject, max_length=50)
        date = mail.date
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

    def mails_loaded(self, mails):
        self.mails = mails
        self.stack.widget(2).stop_loading()
        self.populate_list()
        self.stack.setCurrentIndex(0)
        print("Fetched more mails.")
        
    def load_next(self):
        if self.current_index + 10 >= len(self.mails):

            self.stack.widget(2).start_loading()
            self.stack.setCurrentIndex(2)

            self.thread = QThread()
            self.worker = MailFetcherWorker(self.fetcher)
            self.worker.moveToThread(self.thread)
            self.thread.started.connect(self.worker.run)
            self.worker.finished.connect(self.mails_loaded)
            self.worker.finished.connect(self.thread.quit)
            self.worker.finished.connect(self.worker.deleteLater)
            self.thread.finished.connect(self.thread.deleteLater)
            self.thread.start()

        self.current_index += 10
        self.populate_list()
    
    def load_previous(self):
        if self.current_index - 10 >= 0:
            self.current_index -= 10
            self.populate_list()

    def open_mail(self, selected_mail):
        self.stack.widget(1).current_index = self.mails.index(selected_mail)
        self.stack.widget(1).update_view(selected_mail)
        self.stack.widget(1).update_buttons()
        self.stack.setCurrentIndex(1)
