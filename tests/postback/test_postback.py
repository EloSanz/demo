
def test_postback_ios_deterministic_success(client):
    # iOS with device_id -> Deterministic attribution
    payload = {
        "transaction_id": "tx_ios_det_1",
        "campaign_id": "camp_abc",
        "os": "ios",
        "device_id": "idfa-999",
        "revenue": 10.5
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 202
    data = response.json()
    assert data["status"] == "success"
    assert data["transaction_id"] == "tx_ios_det_1"
    assert data["attribution_type"] == "deterministic"
    assert data["internal_status"] == "attributed"

def test_postback_ios_probabilistic_success(client):
    # iOS without device_id -> Probabilistic attribution & queued
    payload = {
        "transaction_id": "tx_ios_prob_1",
        "campaign_id": "camp_abc",
        "os": "ios",
        "device_id": None,
        "revenue": 20.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 202
    data = response.json()
    assert data["status"] == "success"
    assert data["transaction_id"] == "tx_ios_prob_1"
    assert data["attribution_type"] == "probabilistic"
    assert data["internal_status"] == "queued_for_probabilistic_attribution"

def test_postback_android_deterministic_success(client):
    # Android with device_id -> Probabilistic attribution with internal_status=attributed
    payload = {
        "transaction_id": "tx_andr_det_1",
        "campaign_id": "camp_abc",
        "os": "android",
        "device_id": "gps-ad-id-123",
        "revenue": 5.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 202
    data = response.json()
    assert data["status"] == "success"
    assert data["transaction_id"] == "tx_andr_det_1"
    assert data["attribution_type"] == "probabilistic"
    assert data["internal_status"] == "attributed"

def test_postback_android_missing_device_id_fails(client):
    # Android without device_id -> Raises AndroidMissingIDError (400 Bad Request)
    payload = {
        "transaction_id": "tx_andr_fail_1",
        "campaign_id": "camp_abc",
        "os": "android",
        "device_id": None,
        "revenue": 5.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 400
    data = response.json()
    assert data["error"] == "Bad Request"
    assert "device_id is mandatory" in data["message"]

def test_campaign_missing_id_exception_handler(client):
    # postback without campaign -> Raises CampaignMissingIDError (400 Bad Request)
    payload = {
        "transaction_id": "tx_andr_fail_1",
        "campaign_id": None,
        "os": "ios",
        "device_id": "idfa-555",
        "revenue": 5.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 400
    data = response.json()
    assert data["error"] == "Bad Request"

def test_postback_duplicate_transaction_locked(client):
    # Double submission -> Raises DuplicateTransactionError (423 Locked)
    from src.postback.dependencies import _transaction_repo
    
    payload = {
        "transaction_id": "tx_duplicate_lock_test",
        "campaign_id": "camp_abc",
        "os": "ios",
        "device_id": "idfa-555",
        "revenue": 15.0
    }
    
    # Pre-lock the transaction to simulate a concurrent request or an already processed transaction
    _transaction_repo.lock("tx_duplicate_lock_test")
    
    # Sending the request with the locked transaction_id must fail with 423
    response = client.post("/conversions", json=payload)
    assert response.status_code == 423
    data = response.json()
    assert data["error"] == "Locked"
    assert "already processed" in data["message"]

def test_postback_validation_negative_revenue_fails(client):
    # Revenue must be greater than 0.0 -> Pydantic Validation Error (422)
    payload = {
        "transaction_id": "tx_val_fail_1",
        "campaign_id": "camp_abc",
        "os": "ios",
        "revenue": -1.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 422

def test_postback_validation_empty_transaction_id_fails(client):
    # transaction_id must be at least 1 character -> Pydantic Validation Error (422)
    payload = {
        "transaction_id": "",
        "campaign_id": "camp_abc",
        "os": "ios",
        "revenue": 10.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 422

def test_postback_validation_invalid_os_fails(client):
    # OS must be 'ios' or 'android' -> Pydantic Validation Error (422)
    payload = {
        "transaction_id": "tx_val_fail_3",
        "campaign_id": "camp_abc",
        "os": "windows",
        "revenue": 10.0
    }
    response = client.post("/conversions", json=payload)
    assert response.status_code == 422
