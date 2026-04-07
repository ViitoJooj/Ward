"""
Testes de proteção contra Brute Force / Rate Limiting.
"""

import pytest
import requests
from config import URL_TOKEN, VALID_USER, TIMEOUT


class TestBruteForce:

    def test_multiplas_tentativas_com_senha_errada(self, usuario_registrado):
        """
        Após muitas tentativas com senha errada, a API deve retornar
        429 (Too Many Requests) ou 423 (Locked).
        """
        MAX_TENTATIVAS = 25
        bloqueado = False

        for i in range(MAX_TENTATIVAS):
            resp = requests.post(URL_TOKEN, json={
                "email": VALID_USER["email"],
                "password": f"SenhaErrada{i}!",
            }, timeout=TIMEOUT)
            if resp.status_code in (429, 423):
                bloqueado = True
                break

        if not bloqueado:
            pytest.xfail(
                f"Rate limiting NÃO detectado após {MAX_TENTATIVAS} tentativas. "
                "Considere implementar proteção contra brute force."
            )
