"""
Testes com entradas malformadas, payloads grandes e métodos HTTP incorretos.
"""

import pytest
import requests
from config import URL_TOKEN, URL_REGISTER, VALID_USER, TIMEOUT


class TestEntradasMalformadas:

    def test_payload_gigante_na_senha(self, usuario_registrado):
        """Enviar uma senha com 100KB deve ser tratado sem crash."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": "A" * 100_000,
        }, timeout=15)
        assert resp.status_code in (400, 401, 413, 422)

    def test_payload_gigante_no_username(self):
        """Enviar um username com 100KB deve ser tratado sem crash."""
        resp = requests.post(URL_TOKEN, json={
            "username": "B" * 100_000,
            "password": "qualquer",
        }, timeout=15)
        assert resp.status_code in (400, 401, 413, 422)

    @pytest.mark.parametrize("caractere", [
        "用户名",                # Chinês
        "المستخدم",             # Árabe
        "🔐🔑💀",                # Emojis
        "\x00\x01\x02",         # Null bytes
        "\n\r\t",               # Caracteres de controle
        "a" * 300,              # Username muito longo
    ])
    def test_caracteres_especiais_no_login(self, caractere):
        """Caracteres especiais/unicode não devem causar erro 500."""
        resp = requests.post(URL_TOKEN, json={
            "username": caractere,
            "password": "qualquer",
        }, timeout=TIMEOUT)
        assert resp.status_code != 500, (
            f"Servidor retornou 500 para caractere: {repr(caractere)}"
        )

    def test_content_type_form_urlencoded(self, usuario_registrado):
        """Enviar form-urlencoded em vez de JSON deve ser tratado."""
        resp = requests.post(URL_TOKEN, data={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        assert resp.status_code != 500

    def test_content_type_text_plain(self):
        """Enviar text/plain deve ser rejeitado."""
        resp = requests.post(
            URL_TOKEN,
            data="username=admin&password=admin",
            headers={"Content-Type": "text/plain"},
            timeout=TIMEOUT,
        )
        assert resp.status_code != 500

    @pytest.mark.parametrize("metodo", ["GET", "PUT", "DELETE", "PATCH"])
    def test_metodo_http_errado_no_login(self, metodo):
        """Apenas POST deve ser aceito no endpoint de login/token."""
        func = getattr(requests, metodo.lower())
        resp = func(URL_TOKEN, timeout=TIMEOUT)
        assert resp.status_code in (404, 405), (
            f"Método {metodo} não foi rejeitado (status {resp.status_code})"
        )

    @pytest.mark.parametrize("metodo", ["GET", "PUT", "DELETE", "PATCH"])
    def test_metodo_http_errado_no_registro(self, metodo):
        """Apenas POST deve ser aceito no endpoint de registro."""
        func = getattr(requests, metodo.lower())
        resp = func(URL_REGISTER, timeout=TIMEOUT)
        assert resp.status_code in (404, 405)
