import pke

extractor = pke.unsupervised.TopicRank()

def keyword_extract(text, n):
    extractor.load_document(input=text, language='ja')
    extractor.candidate_selection(pos={'NOUN', 'PROPN'})
    extractor.candidate_weighting()
    return extractor.get_n_best(n=n)