"""
🔐 Menu Interativo — Suite de Testes de Segurança

Execute:  python menu.py
"""

import subprocess
import sys
import os

# Diretório onde estão os testes
TEST_DIR = os.path.dirname(os.path.abspath(__file__))

TESTES = [
    ("test_registro.py",              "📝 Registro",                  "Registrar com dados válidos, duplicados, campos ausentes, email inválido, senha fraca"),
    ("test_login.py",                 "🔑 Login",                     "Login com email/username + senha, campos ausentes, body vazio, strings vazias"),
    ("test_sql_injection.py",         "💉 SQL Injection",             "10 payloads de SQL injection em username, email e senha (30 testes)"),
    ("test_xss.py",                   "🌐 XSS",                      "6 payloads de Cross-Site Scripting no login e registro"),
    ("test_brute_force.py",           "🔨 Brute Force",              "25 tentativas consecutivas — verifica rate limiting (429/423)"),
    ("test_token.py",                 "🎟️  Token JWT",                "Formato, tamanho, tokens falsos rejeitados, unicidade por login"),
    ("test_entradas_malformadas.py",  "🧪 Entradas Malformadas",     "Payloads de 100KB, unicode, null bytes, Content-Type errado, métodos HTTP"),
    ("test_vazamento_info.py",        "🔍 Vazamento de Informação",  "Stack traces, enumeração de usuários, hash de senha na resposta"),
    ("test_timing_attack.py",         "⏱️  Timing Attack",            "Diferença de tempo entre user inexistente vs senha errada"),
    ("test_seguranca_complementar.py","🛡️  Segurança Complementar",  "Case-sensitivity, espaços na senha, senha=username, headers de segurança"),
]


def limpar_tela():
    os.system("cls" if os.name == "nt" else "clear")


def exibir_menu():
    limpar_tela()
    print()
    print("  ╔══════════════════════════════════════════════════════════════╗")
    print("  ║          🔐  SUITE DE TESTES DE SEGURANÇA — AUTH  🔐        ║")
    print("  ╠══════════════════════════════════════════════════════════════╣")
    print("  ║                                                              ║")

    for i, (_, nome, descricao) in enumerate(TESTES, 1):
        num = f" {i}" if i < 10 else f"{i}"
        print(f"  ║   [{num}]  {nome:<30}                       ║")

    print("  ║                                                              ║")
    print("  ║  ─────────────────────────────────────────────────────────── ║")
    print("  ║   [ 0]  🚀 Executar TODOS os testes                         ║")
    print("  ║   [ Q]  ❌ Sair                                              ║")
    print("  ║                                                              ║")
    print("  ╚══════════════════════════════════════════════════════════════╝")
    print()


def exibir_descricao(indice: int):
    _, nome, descricao = TESTES[indice]
    print(f"  {nome}")
    print(f"  └─ {descricao}")
    print()


def rodar_teste(arquivo: str, verbose: bool = True):
    """Executa um arquivo de teste com pytest."""
    caminho = os.path.join(TEST_DIR, arquivo)
    cmd = [sys.executable, "-m", "pytest", caminho, "-v", "--tb=short"]
    if not verbose:
        cmd.append("-q")
    
    print(f"\n  ▶ Executando: pytest {arquivo} -v --tb=short")
    print(f"  {'─' * 55}\n")
    
    resultado = subprocess.run(cmd, cwd=TEST_DIR)
    return resultado.returncode


def rodar_todos():
    """Executa todos os testes de uma vez."""
    print(f"\n  ▶ Executando: pytest (todos os módulos) -v --tb=short")
    print(f"  {'─' * 55}\n")
    
    cmd = [sys.executable, "-m", "pytest", TEST_DIR, "-v", "--tb=short"]
    resultado = subprocess.run(cmd, cwd=TEST_DIR)
    return resultado.returncode


def main():
    while True:
        exibir_menu()
        escolha = input("  Escolha uma opção: ").strip().lower()

        if escolha == "q":
            print("\n  👋 Até logo!\n")
            break

        if escolha == "0":
            rodar_todos()
            input("\n  Pressione ENTER para voltar ao menu...")
            continue

        try:
            indice = int(escolha) - 1
            if 0 <= indice < len(TESTES):
                exibir_descricao(indice)
                arquivo = TESTES[indice][0]
                rodar_teste(arquivo)
                input("\n  Pressione ENTER para voltar ao menu...")
            else:
                input("  ⚠️  Opção inválida. Pressione ENTER...")
        except ValueError:
            input("  ⚠️  Digite um número válido. Pressione ENTER...")


if __name__ == "__main__":
    main()
