"""
Testes de segurança complementares e edge cases.
"""

import uuid
import pytest
import requests
from config import URL_TOKEN, URL_REGISTER, VALID_USER, TIMEOUT


class TestSegurancaComplementar:

    def test_login_case_sensitive_na_senha(self, usuario_registrado):
        """Senha deve ser case-sensitive."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"].upper(),
        }, timeout=TIMEOUT)
        if VALID_USER["password"] != VALID_USER["password"].upper():
            assert resp.status_code != 200, (
                "Senha NÃO é case-sensitive! Isso é uma falha de segurança."
            )

    def test_login_com_espacos_extras_na_senha(self, usuario_registrado):
        """Espaços extras antes/depois da senha não devem ser aceitos."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": f"  {VALID_USER['password']}  ",
        }, timeout=TIMEOUT)
        assert resp.status_code != 200, (
            "Senha com espaços extras foi aceita! A API pode estar fazendo trim."
        )

    def test_registro_nao_permite_senha_igual_username(self):
        """Registro com senha igual ao username deve ser idealmente rejeitado."""
        uid = uuid.uuid4().hex[:6]
        resp = requests.post(URL_REGISTER, json={
            "username": f"user_{uid}",
            "email": f"user_{uid}@example.com",
            "password": f"user_{uid}",
        }, timeout=TIMEOUT)
        if resp.status_code in (200, 201):
            pytest.xfail(
                "Senha igual ao username foi aceita. "
                "Considere adicionar essa validação."
            )

    def test_registro_nao_permite_senha_igual_email(self):
        """Registro com senha igual ao email deve ser idealmente rejeitado."""
        uid = uuid.uuid4().hex[:6]
        email = f"user_{uid}@example.com"
        resp = requests.post(URL_REGISTER, json={
            "username": f"user_{uid}",
            "email": email,
            "password": email,
        }, timeout=TIMEOUT)
        if resp.status_code in (200, 201):
            pytest.xfail(
                "Senha igual ao email foi aceita. "
                "Considere adicionar essa validação."
            )

    def test_header_server_nao_expoe_versao(self, usuario_registrado):
        """O header 'Server' não deve expor a versão exata do framework."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        server = resp.headers.get("Server", "")
        versoes_perigosas = ["gin/", "echo/", "fiber/", "go/"]
        for v in versoes_perigosas:
            assert v.lower() not in server.lower(), (
                f"Header 'Server' expõe versão: {server}"
            )

    def test_resposta_tem_headers_de_seguranca(self, usuario_registrado):
        """Verifica presença de headers de segurança recomendados."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": VALID_USER["password"],
        }, timeout=TIMEOUT)
        headers = resp.headers
        recomendados = {
            "X-Content-Type-Options": "nosniff",
            "X-Frame-Options": None,
        }
        ausentes = [h for h in recomendados if h not in headers]
        if ausentes:
            pytest.xfail(
                f"Headers de segurança ausentes (recomendado): {ausentes}"
            )
