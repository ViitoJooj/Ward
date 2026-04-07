"""
Testes de proteção contra SQL Injection.
"""

import pytest
import requests
from config import URL_TOKEN, VALID_USER, TIMEOUT

SQL_PAYLOADS = [
    "' OR '1'='1",
    "' OR '1'='1' --",
    "' OR '1'='1' /*",
    "admin'--",
    "' UNION SELECT * FROM users --",
    "'; DROP TABLE users; --",
    "1' OR '1' = '1",
    "' OR 1=1 --",
    "' OR EXISTS(SELECT * FROM users) --",
    "') OR ('1'='1",
]


class TestSQLInjection:

    @pytest.mark.parametrize("payload", SQL_PAYLOADS)
    def test_sql_injection_no_username(self, payload, usuario_registrado):
        """SQL Injection no campo username não deve resultar em login."""
        resp = requests.post(URL_TOKEN, json={
            "username": payload,
            "password": "qualquer",
        }, timeout=TIMEOUT)
        assert resp.status_code != 200, (
            f"SQL Injection suspeito! Payload '{payload}' retornou 200"
        )

    @pytest.mark.parametrize("payload", SQL_PAYLOADS)
    def test_sql_injection_no_email(self, payload, usuario_registrado):
        """SQL Injection no campo email não deve resultar em login."""
        resp = requests.post(URL_TOKEN, json={
            "email": payload,
            "password": "qualquer",
        }, timeout=TIMEOUT)
        assert resp.status_code != 200

    @pytest.mark.parametrize("payload", SQL_PAYLOADS)
    def test_sql_injection_na_senha(self, payload, usuario_registrado):
        """SQL Injection no campo senha não deve resultar em login."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": payload,
        }, timeout=TIMEOUT)
        assert resp.status_code != 200
