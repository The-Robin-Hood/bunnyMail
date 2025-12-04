from email.utils import parsedate_to_datetime
from datetime import datetime

class Utils:
    @staticmethod
    def format_mail_date(date_str):
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
        
    
    @staticmethod
    def trimmed_text(text, max_length=80):
        if len(text) <= max_length:
            return text
        return text[: max_length - 3] + "..."