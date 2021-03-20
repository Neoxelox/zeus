from locust import task
from locust.contrib.fasthttp import FastHttpUser


class LoadTest(FastHttpUser):
    @task
    def health(self):
        pass
