def test_eac_android_success(client):
    payload = {
        "event_id": "evt_andr_ok",
        "event_type": "purchase",
        "os": "android",
        "device_id": "gaid-1234",
        "revenue": 19.99
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["event_id"] == "evt_andr_ok"
    assert data["status"] == "accepted"
    assert data["attribution_method"] == "deterministic"


def test_eac_android_missing_device_id_fails(client):
    payload = {
        "event_id": "evt_andr_fail",
        "event_type": "purchase",
        "os": "android",
        "device_id": None,
        "revenue": 19.99
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 400
    data = response.json()
    assert "detail" in data
    assert data["detail"] == "device_id is required for Android deterministic matching."


def test_eac_ios_deterministic_success(client):
    payload = {
        "event_id": "evt_ios_det",
        "event_type": "purchase",
        "os": "ios",
        "device_id": "idfa-5678",
        "revenue": 25.50
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["event_id"] == "evt_ios_det"
    assert data["status"] == "accepted"
    assert data["attribution_method"] == "deterministic"


def test_eac_ios_probabilistic_success(client):
    payload = {
        "event_id": "evt_ios_prob",
        "event_type": "purchase",
        "os": "ios",
        "device_id": None,
        "ip_address": "181.45.67.89",
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X)...",
        "revenue": 25.50
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["event_id"] == "evt_ios_prob"
    assert data["status"] == "accepted"
    assert data["attribution_method"] == "probabilistic"


def test_eac_ios_probabilistic_missing_ip_fails(client):
    payload = {
        "event_id": "evt_ios_fail_ip",
        "event_type": "purchase",
        "os": "ios",
        "device_id": None,
        "ip_address": None,
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X)...",
        "revenue": 25.50
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 400
    data = response.json()
    assert "detail" in data
    assert data["detail"] == "ip_address and user_agent are required for iOS probabilistic matching when device_id is missing."


def test_eac_ios_probabilistic_missing_ua_fails(client):
    payload = {
        "event_id": "evt_ios_fail_ua",
        "event_type": "purchase",
        "os": "ios",
        "device_id": None,
        "ip_address": "181.45.67.89",
        "user_agent": "",
        "revenue": 25.50
    }
    response = client.post("/v1/eac/events", json=payload)
    assert response.status_code == 400
    data = response.json()
    assert "detail" in data
    assert data["detail"] == "ip_address and user_agent are required for iOS probabilistic matching when device_id is missing."
