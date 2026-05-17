import pytest
from fastapi.testclient import TestClient
from main import app
from src.postback.dependencies import _transaction_repo

@pytest.fixture(autouse=True)
def clean_database():
    # Automatically clear the locked transactions before each test to ensure test isolation
    _transaction_repo._processed_transactions.clear()

@pytest.fixture
def client():
    return TestClient(app)
