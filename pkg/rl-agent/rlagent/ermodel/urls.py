from django.urls import path

from . import views

urlpatterns = [
    path('get-prediction', views.get_prediction, name='get-prediction'),
]