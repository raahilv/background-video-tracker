import requests
import webbrowser
import secrets

CLIENT_ID = 'e4894e28af942615b14299e7e411eabf'

if __name__ == '__main__':
    code_challenge = secrets.token_urlsafe(95)
    print('Authencate client from browser')
    webbrowser.open(f'https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id={CLIENT_ID}&code_challenge={code_challenge}')
    auth_code = input('Enter the auth code from the URL here:')
    response = requests.post('https://myanimelist.net/v1/oauth2/token', data={'client_id': CLIENT_ID, 'grant_type': 'authorization_code','code': auth_code, 'code_verifier': code_challenge})
    with open('MAL_token.txt', 'w') as file:
        # file.write(code_challenge + '\n')
        file.write(response.text.split(':')[-1][1:-2])
