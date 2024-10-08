import pke
import json

extractor = pke.unsupervised.TopicRank()

def keyword_extract(text, n):
    extractor.load_document(input=text, language='ja')
    extractor.candidate_selection(pos={'NOUN', 'PROPN'})
    extractor.candidate_weighting()
    res = extractor.get_n_best(n=n)
    return json.dumps(res, ensure_ascii=False)

#%%
test = "ハプト藻に関する最古の記載はエーレンベルク（1836）によるものである。彼はバルト海周辺の石灰岩層から微細な円板状の構造物（円石＝coccolith）を発見した。しかし彼は、この構造物を生物由来ではなく、化学的、無機的要因によって生成したものと考えた。その後ハクスリー（1858）が同様の構造物を海底の堆積物の中から発見したが、やはり円石は非生物起源であると考えられた。円石を初めて生物起源であるとしたのは ウォーリッチ（1860）と ソービー（1861）である。彼らは円石が多数結合して中空の球を形成したものを発見し、coccosphere と命名した。現在この語は、円石を持つ細胞全体を、原形質を含めて表す単語として用いられている。しかしながら彼は円石藻という微細藻の存在を提唱したのではなく、coccosphere を有孔虫の生活環の一部と考えるに留まった。1870年代に入ると再び エーレンベルク の円石非生物由来説が支持されるようになった。特に円石の幾何学的な形状から、炭酸カルシウムの凝結、結晶化によると考えられる事が多かった。円石の持ち主を微細藻であると提唱したのは ワイヴィル・トムソン（1874）である。この時初めて円石は単細胞藻の外被であると考えられた。その後、coccosphere の中に色素体があるという報告や、Murray とBlackman（1898）による細胞分裂の描写が為されるに至り、単細胞藻としての円石藻－ハプト藻が認識される事となった。分類上のハプト藻は、体制と光合成色素の類似から、古くは不等毛植物門黄金色藻綱に含められていた経緯がある。ハプト植物門として独立したのは近年（1962）である。"
resp = keyword_extract(test, 10)
print(resp)
print(type(resp))