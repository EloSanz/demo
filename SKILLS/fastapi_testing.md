# FastAPI Testing Best Practices Skill Guide

This guide details industry-standard patterns for testing robust FastAPI applications using `pytest`, `TestClient`, and dependency injection overrides.

---

## 💡 Core Principles

### 1. Leverage Dependency Injection Overrides (No Hard Mocking)
FastAPI provides a built-in mechanism to dynamically swap components (databases, repositories, queue adapters) using `app.dependency_overrides`. Avoid hard mocking libraries like `unittest.mock` when testing endpoints; instead, inject test-specific adapters.

**Pro-Tip**: Always clear overrides in a `yield` fixture to prevent state leaks between tests.

```python
# Clear overrides after test run
app.dependency_overrides.clear()
```

### 2. Isolate Database and Adapter State
Every test must start with a clean state. Use `pytest` fixtures with a `function` scope to reset databases or in-memory tables before each execution.

### 3. Stick to a Single TestClient Approach
Use a standard synchronous `TestClient` for regular requests. Keep your code clean by utilizing a reusable fixture inside `conftest.py`.

---

## 🛠️ Practical Implementation Examples

### 1. Setting Up `conftest.py`
Place your global fixtures inside `conftest.py` at the root of your test directory.

```python
import pytest
from fastapi.testclient import TestClient
from main import app
from src.postback.dependencies import get_queue, get_transaction_repo
from src.infrastructure.queue.mock_queue import MockQueueAdapter
from src.infrastructure.repository.in_memory_transaction_repository import InMemoryTransactionRepository

@pytest.fixture
def test_queue():
    # Return a clean Mock Queue Adapter
    return MockQueueAdapter()

@pytest.fixture
def test_repo():
    # Return a clean In-Memory Repository
    return InMemoryTransactionRepository()

@pytest.fixture
def client(test_queue, test_repo):
    # Override real dependencies with isolated ones
    app.dependency_overrides[get_queue] = lambda: test_queue
    app.dependency_overrides[get_transaction_repo] = lambda: test_repo
    
    with TestClient(app) as c:
        yield c
        
    # Crucial: Reset overrides to keep tests fully isolated!
    app.dependency_overrides.clear()
```

### 2. Testing Endpoints and Business Rules
Write clean, descriptive test functions targeting different use cases (success, validation failure, declarative locking).

```python
def test_postback_success(client):
    # Standard successful deterministic/probabilistic request
    payload = {
        "transaction_id": "tx_123",
        "campaign_id": "camp_abc",
        "os": "ios",
        "device_id": "idfa-456",
        "revenue": 12.50
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 202
    assert response.json()["status"] == "success"

def test_postback_validation_error(client):
    # Missing required field (transaction_id)
    payload = {
        "campaign_id": "camp_abc",
        "os": "ios"
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 422  # FastAPI default validation status code
```

### 3. Testing Concurrency / Duplicate Transactions
Verify that the transaction lock acts as a shield against processing duplicate transactions simultaneously.

```python
def test_postback_duplicate_transaction_error(client):
    payload = {
        "transaction_id": "duplicate_tx_999",
        "campaign_id": "camp_abc",
        "os": "android",
        "revenue": 5.0
    }
    
    # 1. First execution succeeds
    response1 = client.post("/conversions", json=payload)
    assert response1.status_code == 202
    
    # 2. Second execution with the same transaction_id fails with a 400 Bad Request
    response2 = client.post("/conversions", json=payload)
    assert response2.status_code == 400
    assert response2.json()["error"] == "DuplicateTransactionError"
```
