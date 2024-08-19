#define NPY_NO_DEPRECATED_API NPY_1_7_API_VERSION
#include </usr/include/python3.10/Python.h>
#include <numpy/arrayobject.h>
#include <stdio.h>

void set_path() {
    PyObject *sys = PyImport_ImportModule("sys");
    PyObject *path = PyObject_GetAttrString(sys, "path");
    PyList_Append(path, PyUnicode_DecodeFSDefault("."));
}

void call_python_function() {
    char name_str[16];
    // Pythonインタープリタを初期化
    Py_Initialize();
    set_path();
    
    import_array();  // NumPy C APIの初期化

    // Pythonモジュールをインポート
    PyObject *pName = PyUnicode_DecodeFSDefault("np");
    PyObject *pModule = PyImport_Import(pName);
    Py_DECREF(pName);

    if (pModule != NULL) {
        printf("Module loaded\n");
        // Python関数を取得
        PyObject *pFunc = PyObject_GetAttrString(pModule, "create_struct_array");
        printf("Function loaded\n");

        if (pFunc && PyCallable_Check(pFunc)) {
            // Python関数を呼び出し
            PyObject *pValue = PyObject_CallObject(pFunc, NULL);
            printf("Function called\n");

            if (pValue != NULL && PyArray_Check(pValue)) {
                // NumPy配列としてデータを取得
                PyArrayObject *np_array = (PyArrayObject*)pValue;
                int num_elements = PyArray_SIZE(np_array);
                printf("Num elements: %d\n", num_elements);
                
                // 各要素のアクセス
                for (int i = 0; i < num_elements; i++) {
                    // nameとvalueにアクセス
                    PyObject *item = PyArray_GETITEM(np_array, PyArray_GETPTR1(np_array, i));
                    printf("successfuly get item\n");
                    PyObject *name = PyObject_GetAttrString(item, "name");
                    PyObject *value = PyObject_GetAttrString(item, "value");
                    printf("successfuly get name and value\n");

                    double value_double = PyFloat_AsDouble(value);
                    printf("successfuly get value_double\n");
                    name_str = PyUnicode_AsUTF8(name);
                    printf("successfuly get name_str\n");

                    printf("Item %d: name = %s, value = %f\n", i, name_str, value_double);

                    Py_DECREF(name);
                    Py_DECREF(value);
                }

                Py_DECREF(pValue);
            } else {
                PyErr_Print();
            }

            Py_XDECREF(pFunc);
        } else {
            PyErr_Print();
        }

        Py_DECREF(pModule);
    } else {
        PyErr_Print();
    }

    // Pythonインタープリタを終了
    Py_Finalize();
}

int main() {
    call_python_function();
    return 0;
}
