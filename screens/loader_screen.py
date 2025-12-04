from PyQt5.QtWidgets import QStackedWidget, QLabel
from PyQt5.QtGui import QMovie
from PyQt5.QtCore import Qt, QSize

class LoaderScreen(QStackedWidget):
    def __init__(self, parent=None):
        super().__init__(parent)

        # Create loader widget ONCE
        self.loader = QLabel(self)
        self.loader.setAlignment(Qt.AlignCenter)

        self.movie = QMovie("assets/loader.gif")
        self.movie.setScaledSize(QSize(75, 75))
        self.loader.setMovie(self.movie)

        # Add only one screen to this QStackedWidget
        self.addWidget(self.loader)
        self.setCurrentWidget(self.loader)

    def start_loading(self):
        """Show loading animation"""
        self.movie.start()
        # self.setCurrentWidget(self.loader)

    def stop_loading(self):
        """Stop animation"""
        self.movie.stop()
