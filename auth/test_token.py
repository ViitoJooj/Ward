"""
Testes de segurança do Token JWT.
"""

import requests
from config import URL_TOKEN, BASE_URL, VALID_USER, TIMEOUT


class TestToken:

    def test_token_tem_tamanho_adequado(self, token_valido):
        """Token deve ter pelo menos 30 caracteres (JWT típico)."""
        assert len(token_valido) >= 30, "Token muito curto, pode não ser seguro"

    def test_token_eh_formato_jwt(self, token_valido):
        """Token deve ter o formato JWT (3 partes separadas por ponto)."""
        partes = token_valido.split(".")
        assert len(partes) == 3, (
            f"Token não parece ser JWT válido (esperado 3 partes, "
            f"recebido {len(partes)})"
        )

    def test_token_invalido_rejeitado_no_endpoint(self):
        """Tokens forjados devem ser rejeitados."""
        tokens_falsos = [
            "token_falso_12345",
            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.FORJADO.INVALIDO",
            "",
            "null",
            "undefined",
            "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxIn0.",  # alg:none
        ]
        for fake in tokens_falsos:
            headers = {"Authorization": f"Bearer {fake}"}
            resp = requests.get(URL_TOKEN, headers=headers, timeout=TIMEOUT)
            assert resp.status_code in (401, 403, 405), (
                f"Token falso '{fake[:30]}...' não foi rejeitado! "
                f"Status: {resp.status_code}"
            )

    def test_cada_login_gera_token_diferente(self, usuario_registrado):
        """Dois logins consecutivos devem gerar tokens distintos."""
        payload = {
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }
        resp1 = requests.post(URL_TOKEN, json=payload, timeout=TIMEOUT)
        resp2 = requests.post(URL_TOKEN, json=payload, timeout=TIMEOUT)
        token1 = resp1.json().get("token")
        token2 = resp2.json().get("token")
        assert token1 != token2, (
            "Tokens idênticos para logins diferentes — possível token estático!"
        )
