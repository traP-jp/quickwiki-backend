#include <Python.h>
#include <iostream>
#include <vector>
#include <string>
#include <cstring>
#include <cstdlib>

extern "C"
{

// カレントディレクトリをパスに入れる
void set_path() {
    PyObject *sys = PyImport_ImportModule("sys");
    PyObject *path = PyObject_GetAttrString(sys, "path");
    PyList_Append(path, PyUnicode_DecodeFSDefault("/src/tag"));
    PyList_Append(path, PyUnicode_DecodeFSDefault("."));  // debug
}

void finalize_python() {
    Py_Finalize();
}

void initialize_python() {
    Py_Initialize();
    set_path();
}

std::string result = "";

// return: "tagname:score,tagname:score..."
const char* extract(const char *text, int num_keywords) {

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
                return result.c_str();
            }
            PyTuple_SetItem(pArgs, 0, pValue);

            pValue = PyLong_FromLong(num_keywords);
            if (!pValue) {
                Py_DECREF(pArgs);
                Py_DECREF(pModule);
                std::cout << "Cannot convert argument" << std::endl;
                return result.c_str();
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
                            std::string tag_name = PyUnicode_AsUTF8(pTagName);
                            double score = PyFloat_AsDouble(pScore);
                            result += tag_name + ":" + std::to_string(score) + ",";
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
                return result.c_str();
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

    result = result.substr(0, result.size() - 1);  // 最後の","を削除
    std::cout << "[from cpp] Finish keyword extract" << std::endl;
    printf("%s\n", result.c_str());

    return result.c_str();
}

} // extern "C"