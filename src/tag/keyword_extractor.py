import pke
extractor = pke.unsupervised.MultipartiteRank()
print('[from python]Keyword extractor loaded.', flush=True)

print('[from python]Extracting keywords...', flush=True)

num = 0
texts = []

with open('/src/tag/tmp.txt', 'r') as f:
    for line in f:
        if num == 0:
            num = int(line)
        else:
            texts.append(line)

results = []

for text in texts:
    try:
        print('[from python]Loading document...', flush=True)
        extractor.load_document(input=text, language='ja', normalization=None, stoplist='stoplist_ja.txt')
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

    res = ""
    for d in data:
        res += f'{d[0]}:{d[1]},'
        print(f'[from python] {d[0]}: {d[1]}', flush=True)

    results.append(res)

with open('/src/tag/tmp.txt', 'w') as f:
    f.truncate(0)
    for r in results:
        f.write(r + '\n')

print('[from python]Keywords extraction completed.', flush=True)