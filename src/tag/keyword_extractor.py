import pke
extractor = pke.unsupervised.MultipartiteRank()
print('[from python]Keyword extractor loaded.')

print('[from python]Extracting keywords...')

num = 0
texts = []

with open('/src/tag/tmp.txt', 'r') as f:
    for line in f:
        if num == 0:
            num = int(line)
        else:
            # print(f'[from python]Text: {line}', flush=True)
            strs = line.split(',', maxsplit=1)
            texts.append({'id': strs[0], 'text': strs[1]})

print(f'[from python]Number of keywords to extract: {num}')
results = []

for text in texts:
    try:
        print('[from python]Loading document...', flush=True)
        print(f'[from python]Text: {text["text"]}', flush=True)
        extractor.load_document(input=text['text'], language='ja', normalization=None, stoplist='stoplist_ja.txt')
        print('[from python]Document loaded.', flush=True)
    except Exception as e:
        print(f'[from python]Error loading document: {e}', flush=True)

    try:
        extractor.candidate_selection(pos={'NOUN', 'PROPN'})
        print('[from python]Candidates selected.', flush=True)
        extractor.candidate_weighting(threshold=0.74, method='average', alpha=1.1)
        print('[from python]Candidates weighted.', flush=True)
        data = extractor.get_n_best(n=num)
        print('[from python]Keywords extracted.', flush=True)
    except Exception as e:
        print(f'[from python]Error during extraction: {e}', flush=True)

    # res = []
    # for d in data:
    #     res.append({'tag_name': d[0], 'score': d[1]})
    #     print(f'[from python] {d[0]}: {d[1]}', flush=True)

    res = f'{text["id"]}|'
    for d in data:
        res += f'{d[0]}:{d[1]},'
        print(f'[from python] {d[0]}: {d[1]}')

    results.append(res)

with open('/src/tag/tmp.txt', 'w') as f:
    f.truncate(0)
    for r in results:
        f.write(r + '\n')

print('[from python]Keywords extraction completed.')