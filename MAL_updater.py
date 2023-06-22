import requests
import json
import Levenshtein  # I think this is a better comparer than Sequence Matcher

CLIENT_ID = 'e4894e28af942615b14299e7e411eabf'
SIMILARITY_CUTOFF = 0.8  # half arbitrarily chosen


def get_access_token() -> str:
    with open('MAL_token.txt', 'r') as f:
        # refresh_token = f.readlines()[1]
        refresh_token = f.read()
    f.close()
    response = json.loads(requests.post('https://myanimelist.net/v1/oauth2/token',
                                        data={'client_id': CLIENT_ID, 'grant_type': 'refresh_token',
                                              'refresh_token': refresh_token}).text)
    with open('MAL_token.txt', 'w') as f:
        f.write(response['refresh_token'])
    f.close()
    return response['access_token']


def search_closest(json_data: dict, name: str) -> [float, int]:
    max_similarity = [0.0, -1]
    for fake_node in json_data['data']:
        node = fake_node['node']
        if Levenshtein.ratio(name, node['title']) > max_similarity[0]:
            max_similarity = [Levenshtein.ratio(name, node['title']), node['id']]
        for alternative_title in node['alternative_titles'].values():
            if type(alternative_title) is str:
                if Levenshtein.ratio(name, alternative_title) > max_similarity[0]:
                    max_similarity = [Levenshtein.ratio(name, alternative_title), node['id']]
            elif type(alternative_title) is list:
                for title in alternative_title:
                    if Levenshtein.ratio(name, title) > max_similarity[0]:
                        max_similarity = [Levenshtein.ratio(name, title), node['id']]
    return max_similarity


if __name__ == '__main__':

    input_list = [line[:-1] for line in open('send_to_tracker.txt', 'r')]
    assert len(input_list) / 2 == len(input_list) // 2

    for i in range(0, len(input_list), 2):
        response = json.loads(
            requests.get(f'https://api.myanimelist.net/v2/anime?q={input_list[i]}&limit=4&fields=alternative_titles,id',
                         headers={'X-MAL-CLIENT-ID': CLIENT_ID}).text)
        closest = search_closest(response, input_list[i])
        if closest[0] >= SIMILARITY_CUTOFF:
            access_token = get_access_token()
            updated = requests.patch(f'https://api.myanimelist.net/v2/anime/{closest[1]}/my_list_status',data={'num_watched_episodes': input_list[1]}, headers={'Content-Type': 'application/x-www-form-urlencoded','Authorization': f'Bearer {access_token}'})