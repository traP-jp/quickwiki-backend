#include <Python.h>
#include <iostream>
#include <vector>
#include <string>
#include <cstring>
#include <cstdlib>

struct Data {
    std::string tag_name;
    double score;
};

// Goとのデータのやり取り用構造体
struct DataArray {
    char** tag_names;
    double* scores;
    size_t size;
};

extern "C"
{

// カレントディレクトリをパスに入れる
void set_path() {
    PyObject *sys = PyImport_ImportModule("sys");
    PyObject *path = PyObject_GetAttrString(sys, "path");
    PyList_Append(path, PyUnicode_DecodeFSDefault("/src/tag"));
}

void finalize_python() {
    Py_Finalize();
}

DataArray extract(const char *text, int num_keywords) {
    static bool is_initialized = false;
    if (!is_initialized) {
        Py_Initialize();
        set_path();
        is_initialized = true;
    }

    std::vector<Data> dataList;
    DataArray dataArray = {nullptr, nullptr, 0};

    // pythonファイル名からモジュールを読み込む
    PyObject *pName = PyUnicode_DecodeFSDefault("keyword_extractor");
    PyObject *pModule = PyImport_Import(pName);
    Py_DECREF(pName);

    if (pModule != nullptr) {
        // モジュールから関数を取得
        PyObject *pFunc = PyObject_GetAttrString(pModule, "keyword_extract");

        if (pFunc && PyCallable_Check(pFunc)) {
            // 引数の設定
            PyObject *pArgs = PyTuple_New(2);

            PyObject *pValue = PyUnicode_DecodeFSDefault(text);
            if (!pValue) {
                Py_DECREF(pArgs);
                Py_DECREF(pModule);
                std::cout << "Cannot convert argument" << std::endl;
                return dataArray;
            }
            PyTuple_SetItem(pArgs, 0, pValue);

            pValue = PyLong_FromLong(num_keywords);
            if (!pValue) {
                Py_DECREF(pArgs);
                Py_DECREF(pModule);
                std::cout << "Cannot convert argument" << std::endl;
                return dataArray;
            }
            PyTuple_SetItem(pArgs, 1, pValue);

            Py_DECREF(pValue);

            // 関数の実行
            PyObject *pList = PyObject_CallObject(pFunc, pArgs);
            Py_DECREF(pArgs);

            if (pList != nullptr && PyList_Check(pList)) {
                // 戻り値の処理
                Py_ssize_t listSize = PyList_Size(pList);
                for (Py_ssize_t i = 0; i < listSize; i++) {
                    PyObject *pDict = PyList_GetItem(pList, i);

                    if (PyDict_Check(pDict)) {
                        PyObject *pTagName = PyDict_GetItemString(pDict, "tag_name");
                        PyObject *pScore = PyDict_GetItemString(pDict, "score");

                        if (PyUnicode_Check(pTagName) && PyFloat_Check(pScore)) {
                            Data data;
                            data.tag_name = PyUnicode_AsUTF8(pTagName);
                            data.score = PyFloat_AsDouble(pScore);
                            dataList.push_back(data);
                        }
                    }
                }

                Py_DECREF(pList);
            } else {
                // 戻り値がなんか変なことになっていたとき
                Py_DECREF(pFunc);
                Py_DECREF(pModule);
                std::cout << "[Error from cpp] Function returned unexpected value" << std::endl;
                PyErr_Print();
                return dataArray;
            }

        } else {
            if(PyErr_Occurred()) {
                std::cout << "[Error from cpp] Failed to get function: keyword_extractor" << std::endl;
                PyErr_Print();
            }
        }
        Py_XDECREF(pFunc);
        Py_DECREF(pModule);
    } else {
        std::cout << "[Error from cpp] Cannot find module: keyword_extractor.py" << std::endl;
    }

    std::cout << "[from cpp] Finish keyword extract" << std::endl;

    // convert vector to DataArray
    dataArray.size = dataList.size();
    dataArray.tag_names = (char**)malloc(dataList.size() * sizeof(char*));
    dataArray.scores = (double*)malloc(dataList.size() * sizeof(double));

    if (dataArray.tag_names == nullptr || dataArray.scores == nullptr) {
        std::cout << "Failed to allocate memory" << std::endl;
        dataArray.size = 0;
        return dataArray;
    }

    for (size_t i = 0; i < dataArray.size; i++) {
        dataArray.tag_names[i] = strdup(dataList[i].tag_name.c_str());
        dataArray.scores[i] = dataList[i].score;
    }

    return dataArray;
}

void free_data_array(DataArray dataArray) {
    if (dataArray.tag_names != nullptr) {
        for (size_t i = 0; i < dataArray.size; i++) {
            free(dataArray.tag_names[i]);
        }
        free(dataArray.tag_names);
    }
    if (dataArray.scores != nullptr) {
        free(dataArray.scores);
    }
}

} // extern "C"