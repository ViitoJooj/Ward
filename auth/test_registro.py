"""
Testes para o endpoint POST /register.
"""

import requests
from config import URL_REGISTER, VALID_USER, TIMEOUT


class TestRegistro:

    def test_registro_com_dados_validos(self, usuario_registrado):
        """Registro com username, email e senha válidos deve retornar sucesso."""
        assert usuario_registrado["success"] is True
        assert "data" in usuario_registrado
        user = usuario_registrado["data"]
        assert user["username"] == VALID_USER["username"]
        assert user["email"] == VALID_USER["email"]
        assert "id" in user
        assert "created_at" in user

    def test_registro_usuario_duplicado(self, usuario_registrado):
        """Registrar o mesmo username/email deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json=VALID_USER, timeout=TIMEOUT)
        assert resp.status_code in (400, 409, 422)
        data = resp.json()
        assert data["success"] is False

    def test_registro_sem_username(self):
        """Registro sem username deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={
            "email": "nouser@example.com",
            "password": "Senha123!",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_sem_email(self):
        """Registro sem email deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={
            "username": "sem_email_user",
            "password": "Senha123!",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_sem_senha(self):
        """Registro sem senha deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={
            "username": "sem_senha_user",
            "email": "semsena@example.com",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_body_vazio(self):
        """Registro com body vazio deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={}, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_email_invalido(self):
        """Registro com email mal formatado deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={
            "username": "email_invalido_user",
            "email": "isso-nao-e-email",
            "password": "Senha123!",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_senha_fraca(self):
        """Registro com senha muito curta/fraca deve ser rejeitado."""
        resp = requests.post(URL_REGISTER, json={
            "username": "senha_fraca_user",
            "email": "fraca@example.com",
            "password": "123",
        }, timeout=TIMEOUT)
        assert resp.status_code in (400, 422)

    def test_registro_nao_retorna_senha(self, usuario_registrado):
        """A resposta do registro NUNCA deve conter a senha."""
        raw = str(usuario_registrado)
        assert VALID_USER["password"] not in raw, (
            "A senha do usuário aparece na resposta do registro!"
        )
