import time
from threading import Thread
from stem import Signal
from stem.control import Controller
import requests
from functools import wraps


class TorProxy:
    def __init__(self, control_port=9051, socks_port=9050, password='proxy'):
        self.control_port = control_port
        self.socks_port = socks_port
        self.password = password
        self._stop = False

    def get_new_ip(self):
        with Controller.from_port(port=self.control_port) as controller:
            controller.authenticate(password=self.password)
            controller.signal(Signal.NEWNYM)

    def get_current_ip(self):
        proxies = {
            'http': f'socks5h://localhost:{self.socks_port}',
            'https': f'socks5h://localhost:{self.socks_port}'
        }
        response = requests.get('http://httpbin.org/ip', proxies=proxies)
        return response.json()['origin']

    def start_changing_ip(self, interval=600):
        def change_ip():
            while not self._stop:
                print("Current IP:", self.get_current_ip())
                self.get_new_ip()
                time.sleep(interval)
        self.thread = Thread(target=change_ip)
        self.thread.start()

    def stop_changing_ip(self):
        self._stop = True
        self.thread.join()



if __name__ == "__main__":
    tor_proxy = TorProxy()
    tor_proxy.start_changing_ip(interval=10) 

