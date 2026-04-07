"""
Testes de vazamento de informações sensíveis nas respostas da API.
"""

import requests
from config import URL_TOKEN, VALID_USER, TIMEOUT


class TestVazamentoDeInformacao:

    TERMOS_PERIGOSOS = [
        "stack", "traceback", "exception", "panic", "runtime error",
        "sql", "query", "database", "table", "column",
        "internal server error", "goroutine",
    ]

    def test_erro_login_nao_vaza_detalhes(self, usuario_registrado):
        """Resposta de login inválido não deve conter info do servidor."""
        resp = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": "SenhaErrada!",
        }, timeout=TIMEOUT)
        body = resp.text.lower()
        for termo in self.TERMOS_PERIGOSOS:
            assert termo not in body, (
                f"Termo sensível '{termo}' encontrado na resposta de erro!"
            )

    def test_erro_usuario_inexistente_nao_vaza_detalhes(self):
        """Resposta para usuário inexistente não deve conter info interna."""
        resp = requests.post(URL_TOKEN, json={
            "username": "inexistente_xyz_99",
            "password": "QualquerSenha!",
        }, timeout=TIMEOUT)
        body = resp.text.lower()
        for termo in self.TERMOS_PERIGOSOS:
            assert termo not in body

    def test_mensagem_generica_impede_enumeracao_de_usuarios(
        self, usuario_registrado
    ):
        """
        A mensagem de erro para 'usuário inexistente' e 'senha errada'
        deve ser IGUAL, impedindo enumeração de usuários.
        """
        resp1 = requests.post(URL_TOKEN, json={
            "username": "fantasma_xyz_99",
            "password": "QualquerSenha!",
        }, timeout=TIMEOUT)

        resp2 = requests.post(URL_TOKEN, json={
            "email": VALID_USER["email"],
            "password": "SenhaErradaXYZ!",
        }, timeout=TIMEOUT)

        assert resp1.status_code == resp2.status_code, (
            f"Status codes diferentes permitem enumeração de usuários! "
            f"Inexistente: {resp1.status_code}, Senha errada: {resp2.status_code}"
        )

    def test_registro_nao_retorna_hash_da_senha(self, usuario_registrado):
        """A resposta do registro não deve conter password hash."""
        raw = str(usuario_registrado).lower()
        assert "bcrypt" not in raw
        assert "$2a$" not in raw
        assert "$2b$" not in raw
        assert VALID_USER["password"] not in str(usuario_registrado)
