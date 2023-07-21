#!/bin/python3
if __name__ == '__main__':
    input_list = [line for line in open('nearing_completion.txt', 'r')]

    # remove extra info
    for i in range(0, len(input_list)):
        if input_list[i][0] == '[':
            input_list[i] = input_list[i].split(']', maxsplit=1)[1]
        input_list[i] = input_list[i].rsplit('[')[0]

    for i in range(0, len(input_list)):
        input_list[i] = input_list[i].rsplit(' - ', maxsplit=1)
        for j in range(1,5):
            try:
                num = int(input_list[i][1][0:j])
            except ValueError:
                break
        input_list[i][1] = num
    with open('send_to_tracker.txt', 'w') as file:
        for item in input_list:
            file.write(item[0] + '\n')
            file.write(str(item[1]) + '\n')

    exec(open('MAL_updater.py').read())
