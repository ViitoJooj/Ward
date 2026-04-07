"""
Testes de proteção contra XSS (Cross-Site Scripting).
"""

import pytest
import requests
from config import URL_TOKEN, URL_REGISTER, _UID, TIMEOUT

XSS_PAYLOADS = [
    "<script>alert('xss')</script>",
    "<img src=x onerror=alert('xss')>",
    "javascript:alert('xss')",
    "';alert('xss');//",
    "<svg/onload=alert('xss')>",
    '"><img src=x onerror=alert(1)>',
]


class TestXSS:

    @pytest.mark.parametrize("payload", XSS_PAYLOADS)
    def test_xss_no_username(self, payload):
        """Payloads XSS no login não devem ser aceitos nem refletidos."""
        resp = requests.post(URL_TOKEN, json={
            "username": payload,
            "password": "qualquer",
        }, timeout=TIMEOUT)
        assert resp.status_code != 200
        assert payload not in resp.text, (
            f"Payload XSS refletido na resposta: {payload}"
        )

    @pytest.mark.parametrize("payload", XSS_PAYLOADS)
    def test_xss_no_registro(self, payload):
        """Payloads XSS no registro não devem ser refletidos."""
        resp = requests.post(URL_REGISTER, json={
            "username": payload,
            "email": f"{_UID}_xss@example.com",
            "password": "Senha123!",
        }, timeout=TIMEOUT)
        assert payload not in resp.text, (
            f"Payload XSS refletido na resposta de registro: {payload}"
        )
