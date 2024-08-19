import pke

extractor = pke.unsupervised.TopicRank()

def keyword_extract(text, n):
    extractor.load_document(input=text, language='ja')
    extractor.candidate_selection(pos={'NOUN', 'PROPN'})
    extractor.candidate_weighting()
    data = extractor.get_n_best(n=n)
    res = []
    for d in data:
        res.append({'tag_name': d[0], 'score': d[1]})
    
    return res