import numpy as np

def create_struct_array():
    # 構造体のdtypeを定義
    dtype = np.dtype([('name', 'U10'), ('value', np.float64)])
    
    # サンプルデータを作成
    data = np.array([('item1', 1.0), ('item2', 2.5), ('item3', 3.75)], dtype=dtype)
    
    return data
