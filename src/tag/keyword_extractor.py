import pke

extractor = pke.unsupervised.MultipartiteRank()
print('[from python]Keyword extractor loaded.', flush=True)

def keyword_extract(text, n):
    print('[from python]Extracting keywords...', flush=True)
    extractor.load_document(input=text, language='ja', normalization=None)
    extractor.candidate_selection(pos={'NOUN', 'PROPN'})
    extractor.candidate_weighting(threshold=0.74, method='average', alpha=1.1)
    data = extractor.get_n_best(n=n)
    res = []
    for d in data:
        res.append({'tag_name': d[0], 'score': d[1]})
        print(f'[from python] {d[0]}: {d[1]}', flush=True)

    print('[from python]Keywords extracted.', flush=True)
    return res